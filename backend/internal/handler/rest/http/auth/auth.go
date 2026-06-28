package auth

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/utils"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type authHandler struct {
	service module.AuthService
	logger  utils.Logger
}

func InitAuthHandler(s module.AuthService, log utils.Logger) *authHandler {
	return &authHandler{
		service: s,
		logger:  log,
	}
}

const MaxRequestBodySize = 1 << 20 // 1 MB

func sessionMetaFromRequest(r *http.Request) authDto.SessionMeta {
	return authDto.SessionMeta{
		IPAddress: utils.ClientIPFromRequest(r),
		UserAgent: r.UserAgent(),
	}
}

func userIDFromClaims(claims *jwt.AccessClaims) string {
	if claims.UserID != "" {
		return claims.UserID
	}
	return claims.Subject
}

func (h *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "signup rejected: invalid request body")
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	if err := req.ValidateSignUpRequest(); err != nil {
		h.logger.Warn(ctx, "signup rejected: validation failed",
			zap.String("reason", customErr.ClientErrorCode(err)),
			zap.String("email", req.Email),
		)
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.SignUp(ctx, req, sessionMetaFromRequest(r))
	if err != nil {
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "signup rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.String("email", req.Email),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "signup failed",
			zap.String("email", req.Email),
			zap.String("reason", err.Error()),
		)
		response.SendErrorResponse(w, http.StatusInternalServerError, "signup_failed", err, "failed to signup user")
		return
	}

	h.logger.Info(ctx, "signup completed",
		zap.String("user_id", res.User.ID),
		zap.String("email", res.User.Email),
		zap.String("role", res.User.Role),
	)
	response.SendSuccessResponse(w, r, http.StatusCreated, "signup successful", res, nil)
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "login rejected: invalid request body")
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	if err := req.ValidateLoginRequest(); err != nil {
		h.logger.Warn(ctx, "login rejected: validation failed",
			zap.String("reason", customErr.ClientErrorCode(err)),
			zap.String("email", req.Email),
		)
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.Login(ctx, req, sessionMetaFromRequest(r))
	if err != nil {
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "login rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.String("email", req.Email),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "login failed",
			zap.String("email", req.Email),
			zap.String("reason", err.Error()),
		)
		response.SendErrorResponse(w, http.StatusInternalServerError, "login_failed", err, "failed to login")
		return
	}

	h.logger.Info(ctx, "login completed",
		zap.String("user_id", res.User.ID),
		zap.String("email", res.User.Email),
		zap.String("role", res.User.Role),
	)
	response.SendSuccessResponse(w, r, http.StatusOK, "login successful", res, nil)
}

func (h *authHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := middleware.GetAccessClaims(ctx)
	if !ok {
		h.logger.Warn(ctx, "get me rejected: missing token claims")
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return
	}

	userID := userIDFromClaims(claims)

	res, err := h.service.GetMe(ctx, userID)
	if err != nil {
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "get me rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.String("user_id", userID),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "get me failed",
			zap.String("user_id", userID),
			zap.String("reason", err.Error()),
		)
		response.SendErrorResponse(w, http.StatusInternalServerError, "get_me_failed", err, "failed to get current user")
		return
	}

	h.logger.Info(ctx, "get me completed",
		zap.String("user_id", res.User.ID),
		zap.String("role", res.User.Role),
	)
	response.SendSuccessResponse(w, r, http.StatusOK, "ok", res, nil)
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	claims, ok := middleware.GetAccessClaims(ctx)
	if !ok {
		h.logger.Warn(ctx, "logout rejected: missing token claims")
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return
	}

	var req authDto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "logout rejected: invalid request body")
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	if err := req.ValidateLogoutRequest(); err != nil {
		h.logger.Warn(ctx, "logout rejected: validation failed",
			zap.String("reason", customErr.ClientErrorCode(err)),
		)
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	err := h.service.Logout(ctx, userIDFromClaims(claims), req, sessionMetaFromRequest(r))
	if err != nil {
		if errors.Is(err, customErr.ErrInvalidRefreshToken) {
			h.logger.Warn(ctx, "logout rejected", zap.String("reason", "invalid_refresh_token"))
			response.SendErrorResponse(w, http.StatusUnauthorized, "invalid_refresh_token", err, err.Error())
			return
		}
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "logout rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "logout failed", zap.String("reason", err.Error()))
		response.SendErrorResponse(w, http.StatusInternalServerError, "logout_failed", err, "failed to logout")
		return
	}

	h.logger.Info(ctx, "logout completed", zap.String("user_id", userIDFromClaims(claims)))
	response.SendSuccessResponse(w, r, http.StatusOK, "logged out successfully", nil, nil)
}

func (h *authHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "forgot password rejected: invalid request body")
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	if err := req.ValidateForgotPasswordRequest(); err != nil {
		h.logger.Warn(ctx, "forgot password rejected: validation failed",
			zap.String("reason", customErr.ClientErrorCode(err)),
			zap.String("email", req.Email),
		)
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	if err := h.service.ForgotPassword(ctx, req); err != nil {
		h.logger.Error(ctx, "forgot password failed",
			zap.String("email", req.Email),
			zap.String("reason", err.Error()),
		)
		response.SendErrorResponse(w, http.StatusInternalServerError, "forgot_password_failed", err, "failed to process forgot password request")
		return
	}

	h.logger.Info(ctx, "forgot password processed", zap.String("email", req.Email))
	response.SendSuccessResponse(w, r, http.StatusOK, "if an account exists for this email, a reset link has been sent", nil, nil)
}

func (h *authHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "reset password rejected: invalid request body")
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	if err := req.ValidateResetPasswordRequest(); err != nil {
		h.logger.Warn(ctx, "reset password rejected: validation failed",
			zap.String("reason", customErr.ClientErrorCode(err)),
		)
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.ResetPassword(ctx, req)
	if err != nil {
		if errors.Is(err, authDto.ErrInvalidResetToken) {
			h.logger.Warn(ctx, "reset password rejected", zap.String("reason", "invalid_reset_token"))
			response.SendErrorResponse(w, http.StatusBadRequest, "invalid_reset_token", err, err.Error())
			return
		}
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "reset password rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "reset password failed", zap.String("reason", err.Error()))
		response.SendErrorResponse(w, http.StatusInternalServerError, "reset_password_failed", err, "failed to reset password")
		return
	}

	h.logger.Info(ctx, "reset password completed")
	response.SendSuccessResponse(w, r, http.StatusOK, "password updated successfully", res, nil)
}

func (h *authHandler) SendVerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := middleware.GetAccessClaims(ctx)
	if !ok {
		h.logger.Warn(ctx, "send verify email rejected: missing token claims")
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return
	}

	userID := userIDFromClaims(claims)
	if err := h.service.SendVerifyEmail(ctx, userID); err != nil {
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "send verify email rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.String("user_id", userID),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "send verify email failed",
			zap.String("user_id", userID),
			zap.String("reason", err.Error()),
		)
		response.SendErrorResponse(w, http.StatusInternalServerError, "send_verify_email_failed", err, "failed to send verification email")
		return
	}

	h.logger.Info(ctx, "verification email sent", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "verification code sent", nil, nil)
}

func (h *authHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	claims, ok := middleware.GetAccessClaims(ctx)
	if !ok {
		h.logger.Warn(ctx, "verify email rejected: missing token claims")
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return
	}

	var req authDto.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "verify email rejected: invalid request body")
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	if err := req.ValidateVerifyEmailRequest(); err != nil {
		h.logger.Warn(ctx, "verify email rejected: validation failed",
			zap.String("reason", customErr.ClientErrorCode(err)),
		)
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	userID := userIDFromClaims(claims)
	res, err := h.service.VerifyEmail(ctx, userID, req)
	if err != nil {
		if errors.Is(err, authDto.ErrInvalidVerifyOTP) {
			h.logger.Warn(ctx, "verify email rejected", zap.String("reason", "invalid_verify_otp"))
			response.SendErrorResponse(w, http.StatusBadRequest, "invalid_verify_otp", err, err.Error())
			return
		}
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "verify email rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.String("user_id", userID),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "verify email failed",
			zap.String("user_id", userID),
			zap.String("reason", err.Error()),
		)
		response.SendErrorResponse(w, http.StatusInternalServerError, "verify_email_failed", err, "failed to verify email")
		return
	}

	h.logger.Info(ctx, "email verified", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "email verified successfully", res, nil)
}

func (h *authHandler) SendVerifyPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := middleware.GetAccessClaims(ctx)
	if !ok {
		h.logger.Warn(ctx, "send verify phone rejected: missing token claims")
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return
	}

	userID := userIDFromClaims(claims)
	if err := h.service.SendVerifyPhone(ctx, userID); err != nil {
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "send verify phone rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.String("user_id", userID),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "send verify phone failed",
			zap.String("user_id", userID),
			zap.String("reason", err.Error()),
		)
		response.SendErrorResponse(w, http.StatusInternalServerError, "send_verify_phone_failed", err, "failed to send verification sms")
		return
	}

	h.logger.Info(ctx, "verification sms sent", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "verification code sent", nil, nil)
}

func (h *authHandler) VerifyPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	claims, ok := middleware.GetAccessClaims(ctx)
	if !ok {
		h.logger.Warn(ctx, "verify phone rejected: missing token claims")
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return
	}

	var req authDto.VerifyPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "verify phone rejected: invalid request body")
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	if err := req.ValidateVerifyPhoneRequest(); err != nil {
		h.logger.Warn(ctx, "verify phone rejected: validation failed",
			zap.String("reason", customErr.ClientErrorCode(err)),
		)
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	userID := userIDFromClaims(claims)
	res, err := h.service.VerifyPhone(ctx, userID, req)
	if err != nil {
		if errors.Is(err, authDto.ErrInvalidVerifyOTP) {
			h.logger.Warn(ctx, "verify phone rejected", zap.String("reason", "invalid_verify_otp"))
			response.SendErrorResponse(w, http.StatusBadRequest, "invalid_verify_otp", err, err.Error())
			return
		}
		if status, ok := customErr.HTTPStatusFromError(err); ok {
			h.logger.Warn(ctx, "verify phone rejected",
				zap.String("reason", customErr.ClientErrorCode(err)),
				zap.String("user_id", userID),
				zap.Int("status", status),
			)
			response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
			return
		}
		h.logger.Error(ctx, "verify phone failed",
			zap.String("user_id", userID),
			zap.String("reason", err.Error()),
		)
		response.SendErrorResponse(w, http.StatusInternalServerError, "verify_phone_failed", err, "failed to verify phone")
		return
	}

	h.logger.Info(ctx, "phone verified", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "phone verified successfully", res, nil)
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "refresh rejected: invalid request body")
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	res, err := h.service.RefreshToken(ctx, req, sessionMetaFromRequest(r))
	if err != nil {
		h.logger.Warn(ctx, "refresh rejected", zap.String("reason", "invalid_refresh_token"))
		response.SendErrorResponse(w, http.StatusUnauthorized, "invalid_refresh_token", err, "invalid or expired refresh token")
		return
	}

	h.logger.Info(ctx, "token refreshed", zap.String("redirect_to", res.RedirectTo))
	response.SendSuccessResponse(w, r, http.StatusOK, "token refreshed", res, nil)
}
