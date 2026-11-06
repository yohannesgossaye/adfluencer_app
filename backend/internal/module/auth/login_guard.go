package auth

import (
	"context"
	"sync"
	"time"

	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/ratelimit"
)

// loginGuard counts wrong passwords in three layers (see ratelimit.LoginPolicy)
// and pauses logins when one is exceeded. Only failures are counted. Errors
// from the store are ignored: a broken counter must never stop people logging in.
type loginGuard struct {
	store ratelimit.Store
}

func loginPolicy(admin bool) ratelimit.LoginPolicy {
	if admin {
		return ratelimit.AdminLogin
	}
	return ratelimit.UserLogin
}

type loginKeys struct {
	ip, pair, pairLock, pairLocks, account, accountLock string
}

func keysFor(email, ip string) loginKeys {
	if ip == "" {
		ip = "unknown"
	}
	pair := email + "|" + ip
	return loginKeys{
		ip:          "rl:login:ip:" + ip,
		pair:        "rl:login:pair:" + pair,
		pairLock:    "rl:login:pairlock:" + pair,
		pairLocks:   "rl:login:pairlocks:" + pair,
		account:     "rl:login:acct:" + email,
		accountLock: "rl:login:acctlock:" + email,
	}
}

// checkNetwork refuses when this IP has made too many wrong attempts lately.
func (g loginGuard) checkNetwork(ctx context.Context, p ratelimit.LoginPolicy, k loginKeys) error {
	if n, left, _ := g.store.Get(ctx, k.ip); n >= int64(p.IPMaxFailures) && left > 0 {
		return customErr.Wait(customErr.ErrTooManyAttempts, left)
	}
	return nil
}

// checkAccount refuses while this email is locked, from this IP or everywhere.
func (g loginGuard) checkAccount(ctx context.Context, k loginKeys) error {
	_, pairLeft, _ := g.store.Get(ctx, k.pairLock)
	_, accountLeft, _ := g.store.Get(ctx, k.accountLock)
	if left := max(pairLeft, accountLeft); left > 0 {
		return customErr.Wait(customErr.ErrLoginLocked, left)
	}
	return nil
}

// failed records a wrong password. It returns a wait error when this attempt
// started a pause, and accountLocked when the email-wide lock just began (so
// the owner can be warned once).
func (g loginGuard) failed(ctx context.Context, p ratelimit.LoginPolicy, k loginKeys) (wait error, accountLocked bool) {
	ipN, _ := g.store.Incr(ctx, k.ip, p.IPWindow)
	pairN, _ := g.store.Incr(ctx, k.pair, p.PairWindow)
	accountN, _ := g.store.Incr(ctx, k.account, p.AccountWindow)

	var lock time.Duration
	if pairN >= int64(p.PairMaxFailures) {
		locks, _ := g.store.Incr(ctx, k.pairLocks, p.PairLockMemory)
		lock = p.PairLock(locks)
		_ = g.store.Set(ctx, k.pairLock, 1, lock)
		_ = g.store.Del(ctx, k.pair)
	}
	if accountN >= int64(p.AccountMaxFailures) {
		_ = g.store.Set(ctx, k.accountLock, 1, p.AccountLock)
		_ = g.store.Del(ctx, k.account)
		accountLocked = true
		lock = max(lock, p.AccountLock)
	}
	if lock > 0 {
		return customErr.Wait(customErr.ErrLoginLocked, lock), accountLocked
	}
	if ipN >= int64(p.IPMaxFailures) {
		_, left, _ := g.store.Get(ctx, k.ip)
		return customErr.Wait(customErr.ErrTooManyAttempts, max(left, time.Second)), false
	}
	return nil, false
}

// succeeded forgets this email's failures from this network: it was the owner.
// The network and email-wide counts stay, so an attacker who also knows one
// password cannot use it to reset them.
func (g loginGuard) succeeded(ctx context.Context, k loginKeys) {
	_ = g.store.Del(ctx, k.pair, k.pairLocks)
}

// allow counts one use of a limited action (sign-up, reset request) for key
// and refuses once limit is passed within window.
func (g loginGuard) allow(ctx context.Context, key string, limit int, window time.Duration, limitErr error) error {
	n, err := g.store.Incr(ctx, key, window)
	if err != nil || n <= int64(limit) {
		return nil
	}
	_, left, _ := g.store.Get(ctx, key)
	return customErr.Wait(limitErr, max(left, time.Second))
}

var (
	dummyHashOnce sync.Once
	dummyHash     string
)

// checkDummyPassword spends the same time as a real password check, so the
// response time does not reveal whether an email has an account.
func checkDummyPassword(password string) {
	dummyHashOnce.Do(func() {
		dummyHash, _ = jwt.HashPassword("adfluence-no-such-account")
	})
	if dummyHash != "" {
		_, _ = jwt.VerifyPassword(password, dummyHash)
	}
}
