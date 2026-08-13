package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrDuplicateKey = errors.New("Product Already Exist")


type PostgresProductRepository struct {
	db	*sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{
		db: db,
	}
}

func (pr *PostgresProductRepository) Create(ctx	context.Context, product *models.Product) (error) {
	row := pr.db.QueryRowContext(ctx, "INSERT INTO products(name, quantity, seller_id, price) VALUES ($1, $2, $3, $4) RETURNING id",
		product.Name, product.Quantity, product.SellerID, product.Price)

	err := row.Scan(
		&product.ID,
	)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	if pgErr.Code == "23505" {
		return ErrDuplicateKey
	}
	return err
}

func (pr *PostgresProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	row := pr.db.QueryRowContext(ctx, "SELECT id, name, quantity, seller_id, price FROM products WHERE id = ($1)", id)
	selectedProduct := models.Product{}

	err := row.Scan(
		&selectedProduct.ID,
		&selectedProduct.Name,
		&selectedProduct.Quantity,
		&selectedProduct.SellerID,
		&selectedProduct.Price,
	)
	return &selectedProduct, err
}

func (pr *PostgresProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT id, name, quantity, seller_id, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allProducts := make([]*models.Product, 0)
	for rows.Next() {
		fetchedProduct := models.Product{}
		if err := rows.Scan(
			&fetchedProduct.ID,
			&fetchedProduct.Name,
			&fetchedProduct.Quantity,
			&fetchedProduct.SellerID,
			&fetchedProduct.Price,
		); err != nil {
			return nil, err
		}
		allProducts = append(allProducts, &fetchedProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allProducts, nil
}

func (pr *PostgresProductRepository) Update(ctx	context.Context, product *models.Product) (int64, error) {
	result, err := pr.db.ExecContext(ctx, "UPDATE products SET name = ($1), quantity = ($2), price = ($3) WHERE id = ($4) AND seller_id = ($5)", 
		product.Name, product.Quantity, product.Price, product.ID, product.SellerID)

	if err != nil {
	    var pgErr *pgconn.PgError
	    if errors.As(err, &pgErr) && pgErr.Code == "23505" {
	        return 0, ErrDuplicateKey
	    }

	    return 0, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (pr *PostgresProductRepository) Delete(ctx	context.Context, id int64, seller_id int64) (int64, error) {
	result, err := pr.db.ExecContext(ctx, "DELETE FROM products WHERE id = ($1) AND seller_id = ($2)", id, seller_id)
	if err != nil {
		return 0, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affectedRows, nil
}

func (pr *PostgresProductRepository) GetBySellerID(ctx context.Context, seller_id int64) ([]*models.Product, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT id, name, quantity, seller_id, price FROM products WHERE seller_id = ($1)", seller_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userProducts := make([]*models.Product, 0)
	for rows.Next() {
		fetchedProduct := models.Product{}
		if err := rows.Scan(
			&fetchedProduct.ID,
			&fetchedProduct.Name,
			&fetchedProduct.Quantity,
			&fetchedProduct.SellerID,
			&fetchedProduct.Price,
		); err != nil {
			return nil, err
		}
		userProducts = append(userProducts, &fetchedProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userProducts, nil
}