package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"carbon_intensity/models"
)

var (
	minimumDuration = 30
	maximumDuration = 1440
)

type Processor interface {
	GetSlots(ctx context.Context, requiredDuration int, isContinuous bool) (models.Response, error)
}

type Handler struct {
	Processor Processor
}

func (h *Handler) GetSlotsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		duration   int
		continuous bool
	)
	ctx := r.Context()
	query := r.URL.Query()

	durationStr := query.Get("duration")
	if durationStr != "" {
		d, err := strconv.Atoi(durationStr)
		if err != nil || d < minimumDuration || d > maximumDuration {
			http.Error(w, "invalid duration", http.StatusBadRequest)
			return
		}
		duration = d
	}

	continuousStr := query.Get("continuous")
	if continuousStr == "true" {
		continuous = true
	}

	slots, err := h.Processor.GetSlots(ctx, duration, continuous)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}
