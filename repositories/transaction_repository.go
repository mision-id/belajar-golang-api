package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItems) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	//Inisialisasi subtotal -> jumlah transaksi keseluruhan
	totalAmount := 0
	//Inisialisasi Transaction details -> akan kita import ke DB
	details := make([]models.TransactionDetail, 0)
	// Loop setiap items

	for _, item := range items {
		//get product dapet pricing
		var ProductName string
		var ProductID, price, stock int
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&ProductID, &ProductName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product ID %d not Found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}
		//hitung current total = quantity * pricing
		subotal := item.Quantity * price
		//masukan kedalam subtotal
		totalAmount += subotal

		//Kurangi jumlah stok nya
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
		//item nya masukan ke transaction details
		details = append(details, models.TransactionDetail{
			ProductID:    ProductID,
			ProductName:  ProductName,
			ProductPrice: price,
			Quantity:     item.Quantity,
			Subtotal:     subotal,
		})
	}

	//insert transcation
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}
	//insert transaction details
	for i, detail := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1,$2,$3,$4)", transactionID,
			detail.ProductID, detail.Quantity, detail.Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}
