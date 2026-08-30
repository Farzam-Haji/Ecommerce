package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresCartRepository struct {
	db 	*sql.DB
}

func NewPostgresCartRepository(db *sql.DB) *PostgresCartRepository {
	return &PostgresCartRepository{
		db: db,
	}
}

func (cr *PostgresCartRepository) CreateCart(ctx context.Context, userID int64) (int64, error) {
	row := cr.db.QueryRowContext(ctx, "INSERT INTO carts(user_id) VALUES ($1) RETURNING id",
		userID,
	)

	var newRowID int64
	err := row.Scan(
		&newRowID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return 0, err
		}
		if pgErr.Code == "23505" {
			return 0, ErrDuplicateKey
		}
		return 0, err
	}
	return newRowID, nil
}

func (cr *PostgresCartRepository) GetByUserID(ctx context.Context, userID int64) (*models.Cart, error) {
	row := cr.db.QueryRowContext(ctx, "SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1",
		userID,
	)
		
	var cart models.Cart
	err := row.Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRowNotFound
	}
	if err != nil {
		return nil, err
	}
	return &cart, err
}

func (cr *PostgresCartRepository) DeleteCart(ctx context.Context, userID int64) (int64, error) {
	result, err := cr.db.ExecContext(ctx, "DELETE FROM carts WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (cr *PostgresCartRepository) AddItem(ctx context.Context, item *models.CartItem) (error) {
	row := cr.db.QueryRowContext(ctx, "INSERT INTO cart_items(cart_id, product_id, name, quantity) VALUES ($1, $2, $3, $4) RETURNING id",
		item.CartID, item.ProductID, item.Name, item.Quantity,
	)

	err := row.Scan(
		&item.ID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return err
		}
		if pgErr.Code == "23505" {
			return ErrDuplicateKey
		}
		return err
	}
	return nil
}

func (cr *PostgresCartRepository) GetItem(ctx context.Context, cartID int64, productID int64) (*models.CartItem, error) {
	row := cr.db.QueryRowContext(ctx, "SELECT id, cart_id, product_id, name, quantity FROM cart_items WHERE cart_id = $1 AND product_id = $2",
		cartID, productID,
	)
	var item models.CartItem
	err := row.Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Name,
		&item.Quantity,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRowNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (cr *PostgresCartRepository) GetItemByID(ctx context.Context, itemID int64) (*models.CartItem, error) {
	row := cr.db.QueryRowContext(ctx, "SELECT id, cart_id, product_id, name, quantity FROM cart_items WHERE id = $1",
		itemID,
	)
	var item models.CartItem
	err := row.Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Name,
		&item.Quantity,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRowNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (cr *PostgresCartRepository) GetCartItems(ctx context.Context, cartID int64) ([]*models.CartItem, error) {
	query := `
		SELECT ci.id, ci.cart_id, ci.product_id, ci.name, ci.quantity, p.price
		FROM cart_items ci
		LEFT JOIN products p 
		ON p.id = ci.product_id
		WHERE ci.cart_id = $1
	`
	rows, err := cr.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*models.CartItem, 0)
	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Name,
			&item.Quantity,
			&item.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (cr *PostgresCartRepository) UpdateItemQuantity(ctx context.Context, quantity int32, itemID int64) (int64, error) {
	result, err := cr.db.ExecContext(ctx, "UPDATE cart_items SET quantity = $1 WHERE id = $2",
		quantity, itemID,
	)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (cr *PostgresCartRepository) DeleteItem(ctx context.Context, itemID int64) (int64, error) {
	result, err := cr.db.ExecContext(ctx, "DELETE FROM cart_items WHERE id = $1",
		itemID,
	)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}