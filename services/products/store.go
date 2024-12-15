package products

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hammad177/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateProduct(product types.CreateProductPayload) error {
	_, err := s.db.Exec("INSERT INTO products (name, price, image, description, quantity) VALUES (?, ?, ?, ?, ?)", product.Name, product.Price, product.Image, product.Description, product.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateProductData(productID int, product types.UpdateProductPayload) error {
	if product.Name == nil && product.Price == nil && product.Image == nil && product.Description == nil && product.Quantity == nil {
		return fmt.Errorf("no fields to update")
	}

	p, err1 := s.GetProductByID(productID)
	if err1 != nil {
		return err1
	}

	// Check each field and set the value from the existing product if the new value is empty or nil
	if product.Name == nil {
		product.Name = &p.Name
	}
	if product.Price == nil {
		product.Price = &p.Price
	}
	if product.Image == nil {
		product.Image = &p.Image
	}
	if product.Description == nil {
		product.Description = &p.Description
	}
	if product.Quantity == nil {
		product.Quantity = &p.Quantity
	}

	_, err2 := s.db.Exec("UPDATE products SET name = ?, price = ?, image = ?, description = ?, quantity = ? WHERE id = ?", product.Name, product.Price, product.Image, product.Description, product.Quantity, productID)
	if err2 != nil {
		return err2
	}

	return nil
}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, price = ?, image = ?, description = ?, quantity = ? WHERE id = ?", product.Name, product.Price, product.Image, product.Description, product.Quantity, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")

	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)

	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) GetProductByID(id int) (*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	p := new(types.Product)
	for rows.Next() {
		p, err = scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.ID == 0 {
		return nil, fmt.Errorf("product not found")
	}

	return p, nil
}

func (s *Store) GetProductsByID(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	// Convert productIDs to []interface{}
	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil

}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	p := new(types.Product)
	err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}
