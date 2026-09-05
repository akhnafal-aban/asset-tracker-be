package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akhnafal-aban/asset_tracker_be/internal/model"
	"github.com/akhnafal-aban/asset_tracker_be/internal/service"
)

type AssetHandler struct {
	service service.AssetService
}

func NewAssetHandler(service service.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

const defaultUserID int64 = 1

func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	asset, err := h.service.CreateAsset(r.Context(), defaultUserID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, asset)
}

func (h *AssetHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	// Go 1.22+ PathValue
	idStr := r.PathValue("id")
	assetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	asset, err := h.service.GetAsset(r.Context(), defaultUserID, assetID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, asset)
}

func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := h.service.ListAssets(r.Context(), defaultUserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if assets == nil {
		assets = make([]*model.Asset, 0)
	}

	respondJSON(w, http.StatusOK, map[string]any{"data": assets})
}

func (h *AssetHandler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	assetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var req service.UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	asset, err := h.service.UpdateAsset(r.Context(), defaultUserID, assetID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, asset)
}

func (h *AssetHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	assetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	err = h.service.DeleteAsset(r.Context(), defaultUserID, assetID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AssetHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	if categories == nil {
		respondJSON(w, http.StatusOK, map[string]any{"data": []any{}})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"data": categories})
}
