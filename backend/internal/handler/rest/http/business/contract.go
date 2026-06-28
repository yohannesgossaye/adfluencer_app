package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const maxContractBodySize = 1 << 20

func (h *businessHandler) CreateContract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxContractBodySize)

	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	var req businessDto.CreateContractRequest
	if decodeJSON(w, r, &req) {
		return
	}

	res, err := h.service.CreateContract(ctx, userID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "contract created successfully", res, nil)
}

func (h *businessHandler) ListContracts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	query := contractDto.ListQuery{
		Status: r.URL.Query().Get("status"),
		Page:   parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:  parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListContracts(ctx, userID, query)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "contracts retrieved successfully", res, nil)
}

func (h *businessHandler) GetContract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	contractID := chi.URLParam(r, "contract_id")
	res, err := h.service.GetContract(ctx, userID, contractID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "contract retrieved successfully", res, nil)
}

func (h *businessHandler) PayContract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	contractID := chi.URLParam(r, "contract_id")
	res, err := h.service.PayContract(ctx, userID, contractID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "contract payment recorded (escrow held)", res, nil)
}

func (h *businessHandler) ApproveWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	contractID := chi.URLParam(r, "contract_id")
	res, err := h.service.ApproveWork(ctx, userID, contractID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "work approved — payout queued", res, nil)
}

func (h *businessHandler) CancelContract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	contractID := chi.URLParam(r, "contract_id")
	res, err := h.service.CancelContract(ctx, userID, contractID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "contract cancelled successfully", res, nil)
}
