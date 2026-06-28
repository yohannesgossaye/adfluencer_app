package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *adminHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	res, err := h.service.ListCategories(ctx)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "categories retrieved successfully", res, nil)
}

func (h *adminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.CreateCategory(ctx, adminID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "category created successfully", res, nil)
}

func (h *adminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	categoryID := chi.URLParam(r, "category_id")
	res, err := h.service.UpdateCategory(ctx, adminID, categoryID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "category updated successfully", res, nil)
}

func (h *adminHandler) DisableCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	categoryID := chi.URLParam(r, "category_id")
	res, err := h.service.DisableCategory(ctx, adminID, categoryID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "category disabled successfully", res, nil)
}

func (h *adminHandler) EnableCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	categoryID := chi.URLParam(r, "category_id")
	res, err := h.service.EnableCategory(ctx, adminID, categoryID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "category enabled successfully", res, nil)
}
