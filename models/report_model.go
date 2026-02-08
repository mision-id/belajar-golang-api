package models

type BestSeller struct {
	ProductName  string `json:"name"`
	QuantitySold int    `json:"quantity_sold"`
}

type DailyReport struct {
	TotalRevenue     int          `json:"total_revenue"`
	TotalTransaction int          `json:"total_transaction"`
	BestSeller       []BestSeller `json:"best_seller"`
}
