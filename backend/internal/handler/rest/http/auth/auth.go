package auth

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/internal/handler/rest/httputil"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/utils"
	"encoding/json"
	"net/http"
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

func (h *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.ValidateSignUpRequest(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.SignUp(ctx, req, sessionMetaFromRequest(r))
	if err != nil {
		httputil.WriteServiceErrorOr(w, err, "signup_failed", "failed to signup user")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusCreated, "signup successful", res, nil)
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.ValidateLoginRequest(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.Login(ctx, req, sessionMetaFromRequest(r))
	if err != nil {
		httputil.WriteServiceErrorOr(w, err, "login_failed", "failed to login")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "login successful", res, nil)
}

func (h *authHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := middleware.ExtractUserInfo(ctx)
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return
	}

	res, err := h.service.GetMe(ctx, user.UserID)
	if err != nil {
		httputil.WriteServiceErrorOr(w, err, "get_me_failed", "failed to get current user")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "ok", res, nil)
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	user, err := middleware.ExtractUserInfo(ctx)
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return
	}

	var req authDto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.ValidateLogoutRequest(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	if err := h.service.Logout(ctx, user.UserID, req, sessionMetaFromRequest(r)); err != nil {
		httputil.WriteServiceErrorOr(w, err, "logout_failed", "failed to logout",
			httputil.ErrorMapping{Err: customErr.ErrInvalidRefreshToken, Status: http.StatusUnauthorized, Code: "invalid_refresh_token"},
		)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "logged out successfully", nil, nil)
}

func (h *authHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.ValidateForgotPasswordRequest(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	if err := h.service.ForgotPassword(ctx, req); err != nil {
		httputil.WriteServiceErrorOr(w, err, "forgot_password_failed", "failed to process forgot password request")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "if an account exists for this email, a reset link has been sent", nil, nil)
}

func (h *authHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.ValidateResetPasswordRequest(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.ResetPassword(ctx, req)
	if err != nil {
		httputil.WriteServiceErrorOr(w, err, "reset_password_failed", "failed to reset password",
			httputil.ErrorMapping{Err: authDto.ErrInvalidResetToken, Status: http.StatusBadRequest, Code: "invalid_reset_token"},
		)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "password updated successfully", res, nil)
}

func (h *authHandler) SendVerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.ExtractUserInfo(ctx)
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return
	}

	if err := h.service.SendVerifyEmail(ctx, user.UserID); err != nil {
		httputil.WriteServiceErrorOr(w, err, "send_verify_email_failed", "failed to send verification email")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "verification code sent", nil, nil)
}

func (h *authHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	user, err := middleware.ExtractUserInfo(ctx)
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return
	}

	var req authDto.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.ValidateVerifyEmailRequest(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.VerifyEmail(ctx, user.UserID, req)
	if err != nil {
		httputil.WriteServiceErrorOr(w, err, "verify_email_failed", "failed to verify email",
			httputil.ErrorMapping{Err: authDto.ErrInvalidVerifyOTP, Status: http.StatusBadRequest, Code: "invalid_verify_otp"},
		)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "email verified successfully", res, nil)
}

func (h *authHandler) SendVerifyPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.ExtractUserInfo(ctx)
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return
	}

	if err := h.service.SendVerifyPhone(ctx, user.UserID); err != nil {
		httputil.WriteServiceErrorOr(w, err, "send_verify_phone_failed", "failed to send verification sms")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "verification code sent", nil, nil)
}

func (h *authHandler) VerifyPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	user, err := middleware.ExtractUserInfo(ctx)
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return
	}

	var req authDto.VerifyPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.ValidateVerifyPhoneRequest(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.VerifyPhone(ctx, user.UserID, req)
	if err != nil {
		httputil.WriteServiceErrorOr(w, err, "verify_phone_failed", "failed to verify phone",
			httputil.ErrorMapping{Err: authDto.ErrInvalidVerifyOTP, Status: http.StatusBadRequest, Code: "invalid_verify_otp"},
		)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "phone verified successfully", res, nil)
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req authDto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	res, err := h.service.RefreshToken(ctx, req, sessionMetaFromRequest(r))
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "invalid_refresh_token", err, "invalid or expired refresh token")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "token refreshed", res, nil)
}
