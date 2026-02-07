package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(name string) ([]models.Product, error) {
	//query := "SELECT id, name, price, stock FROM products"
	query := "SELECT p.id, p.name, p.price, p.stock, c.id, c.name, c.description FROM products AS p LEFT JOIN categories AS c ON p.category_id = c.id"
	var args []interface{}
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		//err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category.ID, &p.Category.Name, &p.Category.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(data models.Product) error {
	query := "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRow(query, data.Name, data.Price, data.Stock).Scan(&data.ID)
	return err
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	/*
		query := "SELECT id, name, price, stock FROM products WHERE id = $1"
		row := r.db.QueryRow(query, id)

		var p models.Product
		err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	*/
	query := "SELECT p.id, p.name, p.price, p.stock, c.id, c.name, c.description FROM products p JOIN categories c ON p.category_id = c.id WHERE p.id = $1"
	row := r.db.QueryRow(query, id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category.ID, &p.Category.Name, &p.Category.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *ProductRepository) UpdatebyID(id int, data models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := r.db.Exec(query, data.Name, data.Price, data.Stock, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r *ProductRepository) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("Product not found")
	}
	return err
}
