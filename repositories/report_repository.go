package repositories

import (
	"context"
	"database/sql"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetDailyReport(ctx context.Context, start time.Time, end time.Time) (*models.DailyReport, error) {
	var report models.DailyReport

	// 1. total revenue & total transaksi
	q1 := `
		SELECT
			COALESCE(SUM(subtotal), 0),
			COUNT(DISTINCT transaction_id)
		FROM transaction_details
		WHERE created_at >= $1
		AND created_at < $2
	`
	if err := r.db.QueryRowContext(ctx, q1, start, end).
		Scan(&report.TotalRevenue, &report.TotalTransaction); err != nil {
		return nil, err
	}

	// 2. best seller (slice)
	q2 := `
		SELECT
			p.name,
			SUM(td.quantity)
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		WHERE td.created_at >= $1
		AND td.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`

	rows, err := r.db.QueryContext(ctx, q2, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bs models.BestSeller
		if err := rows.Scan(&bs.ProductName, &bs.QuantitySold); err != nil {
			return nil, err
		}
		report.BestSeller = append(report.BestSeller, bs)
	}

	// kalau tidak ada transaksi → array kosong []
	if report.BestSeller == nil {
		report.BestSeller = []models.BestSeller{}
	}

	return &report, nil

}

func (r *ReportRepository) GetReportByRange(ctx context.Context, start time.Time, end time.Time) (*models.DailyReport, error) {
	// ambil data total revenue & transaksi
	var report models.DailyReport
	q1 := `
		SELECT
			COALESCE(SUM(subtotal), 0),
			COUNT(DISTINCT transaction_id)
		FROM transaction_details
		WHERE created_at >= $1
		AND created_at < $2
	`
	if err := r.db.QueryRowContext(ctx, q1, start, end).
		Scan(&report.TotalRevenue, &report.TotalTransaction); err != nil {
		return nil, err
	}

	// ambil data best seller 5 produk
	q2 := `
		SELECT
			p.name,
			SUM(td.quantity)
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		WHERE td.created_at >= $1
		AND td.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY SUM(td.quantity) DESC
		LIMIT 5
	`

	rows, err := r.db.QueryContext(ctx, q2, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bs models.BestSeller
		if err := rows.Scan(&bs.ProductName, &bs.QuantitySold); err != nil {
			return nil, err
		}
		report.BestSeller = append(report.BestSeller, bs)
	}

	if report.BestSeller == nil {
		report.BestSeller = []models.BestSeller{}
	}

	return &report, nil
}
