package services

import (
	"context"
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) DailyReport(ctx context.Context, date time.Time) (*models.DailyReport, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	return s.repo.GetDailyReport(ctx, start, end)
}

func (s *ReportService) ReportByRange(
	ctx context.Context,
	start time.Time,
	end time.Time,
) (*models.DailyReport, error) {

	if end.Before(start) {
		return nil, errors.New("end_date tidak boleh sebelum start_date")
	}

	return s.repo.GetReportByRange(ctx, start, end)
}
