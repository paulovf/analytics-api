package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/paulovf/analytics-api/internal/domain"
)

type AnalyticsHandler struct {
	repo domain.AnalyticsRepository
}

func NewAnalyticsHandler(repo domain.AnalyticsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

func (h *AnalyticsHandler) GetDailyStats(w http.ResponseWriter, r *http.Request) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	stats, err := h.repo.GetDailyStatusStats(r.Context(), startDate, endDate)
	if err != nil {
		http.Error(w, "Error on search analytics metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
