package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"log"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.DailyReport(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ReportHandler) DailyReport(w http.ResponseWriter, r *http.Request) {
	dateParam := r.PathValue("date")

	if dateParam == "" || dateParam == "today" {
		dateParam = time.Now().Format("2006-01-02")
	}

	// Validasi format tanggal
	parsedDate, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		http.Error(w, "Format tanggal harus YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	report, err := h.service.DailyReport(r.Context(), parsedDate)
	if err != nil {
		log.Printf("daily report error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(report); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func (h *ReportHandler) ReportByRange(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	startStr := q.Get("start_date")
	endStr := q.Get("end_date")

	if startStr == "" {
		http.Error(w, "start_date wajib diisi (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		http.Error(w, "Format start_date harus YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	var endDate time.Time
	if endStr == "" {
		endDate = startDate
	} else {
		endDate, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			http.Error(w, "Format end_date harus YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	// normalisasi ke range [start, end)
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.Local).
		Add(24 * time.Hour)

	report, err := h.service.ReportByRange(r.Context(), start, end)
	if err != nil {
		http.Error(w, "Gagal mengambil report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
