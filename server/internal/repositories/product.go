package repositories

import (
	"database/sql"
	"github.com/btalv/Ecommerce/server/internal/models"
)


type PostgresProductRepository struct {
	db	*sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{
		db: db,
	}
}

func (pr *PostgresProductRepository) Create(product models.Product) (int, error) {
	row := pr.db.QueryRow("INSERT INTO products(name, quantity) VALUES ($1, $2) RETURNING id",
		product.Name, product.Quantity)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (pr *PostgresProductRepository) Read(id int) (models.Product, error){
	row := pr.db.QueryRow("SELECT id, name, quantity FROM products WHERE id = ($1)", id)
	selectedProduct := models.Product{}
	err := row.Scan(
		&selectedProduct.ID,
		&selectedProduct.Name,
		&selectedProduct.Quantity,
	)
	return selectedProduct, err
}

func (pr *PostgresProductRepository) ReadAll() ([]models.Product, error) {
	rows, err := pr.db.Query("SELECT id, name, quantity FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allProducts := make([]models.Product, 0)
	for rows.Next() {
		fetchedProduct := models.Product{}
		if err := rows.Scan(
			&fetchedProduct.ID,
			&fetchedProduct.Name,
			&fetchedProduct.Quantity,
		); err != nil {
			return nil, err
		}
		allProducts = append(allProducts, fetchedProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allProducts, nil
}

func (pr *PostgresProductRepository) Update(product models.Product) (int64, error){
	result, err := pr.db.Exec("UPDATE products SET name = ($1), quantity = ($2) WHERE id = ($3)", 
		product.Name, product.Quantity, product.ID)
	if err != nil {
		return 0, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (pr *PostgresProductRepository) Delete(id int) (int64, error){
	result, err := pr.db.Exec("DELETE FROM products WHERE id = ($1)", id)
	if err != nil {
		return 0, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affectedRows, nil
}