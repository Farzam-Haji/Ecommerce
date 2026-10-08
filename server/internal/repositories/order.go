package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/btalv/Ecommerce/server/internal/models"
)

type PostgresOrderRepository struct {
	db 		*sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}


func (or *PostgresOrderRepository) GetOrderByID(ctx context.Context, orderID int64) (*models.Order, error) {
	query := `
		SELECT id, user_id, payment, payment_id, amount, created_at, expires_at, updated_at 
		FROM orders 
		WHERE id = $1
	`
	row := or.db.QueryRowContext(ctx, query, orderID)

	var order models.Order
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Payment,
		&order.PaymentID,
		&order.Amount,
		&order.CreatedAt,
		&order.ExpiresAt,
		&order.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRowNotFound
	}
	if err != nil {
		return nil, err
	}
	return &order, nil

}

func (or *PostgresOrderRepository) GetOrderItems(ctx context.Context, orderID int64) ([]*models.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, name, price, quantity
		FROM order_items 
		WHERE order_id = $1
		ORDER BY product_id
	`
	rows, err := or.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*models.OrderItem, 0)
	for rows.Next() {
		var item models.OrderItem
		if err = rows.Scan(
			&item.ID,
			&item.OrderID, 
			&item.ProductID, 
			&item.Name,
			&item.Price,
			&item.Quantity,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (or *PostgresOrderRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	query := `
		SELECT id, user_id, payment, payment_id, amount, created_at, expires_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at
	`
	rows, err := or.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*models.Order, 0)
	for rows.Next() {
		var order models.Order
		if err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Payment,
			&order.PaymentID,
			&order.Amount,
			&order.CreatedAt,
			&order.ExpiresAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (or *PostgresOrderRepository) CreateOrderByUserID(ctx context.Context, userID int64, timePeriod time.Duration) (*models.Order, []string, error) {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	//cart
	query := `SELECT id FROM carts WHERE user_id = $1`
	row := tx.QueryRowContext(ctx, query, userID)

	var cartID int64
	err = row.Scan(
		&cartID,
	) 
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrCartNotFound
	}
	if err != nil {
		return nil, nil, err
	}
	
	//cart items
	query = `
		SELECT product_id, name, quantity 
		FROM cart_items 
		WHERE cart_id = $1 
		ORDER BY product_id
	`
	rows, err := tx.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, nil, err
	}

	cartItems := make([]*models.CartItem, 0)
	for rows.Next() {
		var item models.CartItem
		if err = rows.Scan(
			&item.ProductID,
			&item.Name,
			&item.Quantity,
		); err != nil {
			rows.Close()
			return nil, nil, err
		}
		if item.ProductID == nil {
			rows.Close()
			return nil, nil, ErrProductUnavailable
		}
		cartItems = append(cartItems, &item)
	}

	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, nil, err
	}
	if err = rows.Close(); err != nil {
		return nil, nil, err
	}
	if len(cartItems) == 0 {
		return nil, nil, ErrEmptyCart
	}

	//order
	query = `
		INSERT INTO orders(user_id, expires_at) 
		VALUES ($1, $2) 
		RETURNING id, user_id, payment, created_at, expires_at, updated_at
	`
	expiresAt := time.Now().Add(timePeriod)
	row = tx.QueryRowContext(ctx, query, userID, expiresAt)

	var order models.Order
	if err = row.Scan(
		&order.ID,
		&order.UserID,
		&order.Payment,
		&order.CreatedAt,
		&order.ExpiresAt,
		&order.UpdatedAt,	
	); err != nil {
		return nil, nil, err
	}

	//inventory reservation / order items
	query = `
		UPDATE products 
		SET quantity = quantity - $1 
		WHERE id = $2 
			AND quantity >= $1
		RETURNING id, name, price
	`
	query2 := `
		INSERT INTO order_items(order_id, product_id, name, price, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`
	insufficientItems := make([]string, 0)
	var orderAmount int64

	for _, cartItem := range cartItems {
		var orderItem models.OrderItem
		row = tx.QueryRowContext(ctx, query, cartItem.Quantity, *cartItem.ProductID)
		err = row.Scan(
			&orderItem.ProductID,
			&orderItem.Name,
			&orderItem.Price,
		)
		if errors.Is(err, sql.ErrNoRows) {
			insufficientItems = append(insufficientItems, cartItem.Name)
			continue
		}
		if err != nil {
			return nil, nil, err
		}

		_, err := tx.ExecContext(ctx, query2, order.ID, *orderItem.ProductID, orderItem.Name, orderItem.Price, cartItem.Quantity)
		if err != nil {
			return nil, nil, err
		}

		orderAmount += orderItem.Price * int64(cartItem.Quantity)
	}

	if len(insufficientItems) > 0 {
		return nil, insufficientItems, ErrInsufficientInventory
	}

	//order amount
	query = `
		UPDATE orders
		SET amount = $1
		WHERE id = $2
	`ay
	return &order, nil, nil
}

func (or *PostgresOrderRepository) ResolvePaidOrder(ctx context.Context, orderID int64) (error) {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//update payment and updated_at
	query := `
		UPDATE orders
		SET payment = 'paid',
			updated_at = NOW()
		WHERE id = $1
			AND payment = 'pending'
		RETURNING user_id
	`
	row := tx.QueryRowContext(ctx, query, orderID)

	var userID int64
	err = row.Scan(
		&userID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNonPending
	}
	if err != nil {
		return err
	}

	//delete cart and cart items
	query = `
		DELETE FROM carts
		WHERE user_id = $1
	`
	_, err = tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (or *PostgresOrderRepository) ResolveFailedOrder(ctx context.Context, orderID int64) (error) {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//update payment and updated_at
	query := `
		UPDATE orders
		SET payment = 'failed',
			updated_at = NOW()
		WHERE id = $1
			AND payment = 'pending'
	`
	result, err := tx.ExecContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrNonPending
	}

	//get order items
	query = `
		SELECT product_id, quantity
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := tx.QueryContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	items := make([]*models.OrderItem, 0)
	for rows.Next() {
		var item models.OrderItem
		if err = rows.Scan(
			&item.ProductID, 
			&item.Quantity,
		); err != nil {
			rows.Close()
			return err
		}

		if item.ProductID == nil {
			continue
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err = rows.Close(); err != nil {
		return err
	}

	//restore inventory
	query = `
		UPDATE products
		SET quantity = quantity + $1
		WHERE id = $2
	`
	for _, item := range items {
		result, err = tx.ExecContext(ctx, query, item.Quantity, *item.ProductID)
		if err != nil {
			return err
		}

		affected, err = result.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return ErrProductUnavailable
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (or *PostgresOrderRepository) ResolveExpiredOrder(ctx context.Context, orderID int64) (error) {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//update payment and updated_at
	query := `
		UPDATE orders
		SET payment = 'expired',
			updated_at = NOW()
		WHERE id = $1
			AND payment = 'pending'
			AND expires_at <= NOW()
	`
	result, err := tx.ExecContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrNonPending
	}

	//get order items
	query = `
		SELECT product_id, quantity
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := tx.QueryContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	items := make([]*models.OrderItem, 0)
	for rows.Next() {
		var item models.OrderItem
		if err = rows.Scan(
			&item.ProductID, 
			&item.Quantity,
		); err != nil {
			rows.Close()
			return err
		}

		if item.ProductID == nil {
			continue
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err = rows.Close(); err != nil {
		return err
	}

	//restore inventory
	query = `
		UPDATE products
		SET quantity = quantity + $1
		WHERE id = $2
	`
	for _, item := range items {
		result, err = tx.ExecContext(ctx, query, item.Quantity, *item.ProductID)
		if err != nil {
			return err
		}

		affected, err = result.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return ErrProductUnavailable
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (or *PostgresOrderRepository) GetExpiredOrderIDs(ctx context.Context) ([]int64, error) {
	query := `
		SELECT id
		FROM orders
		WHERE payment = 'pending'
			AND expires_at <= NOW()
		ORDER BY expires_at
	`
	rows, err := or.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	IDs := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err = rows.Scan(
			&id,
		); err != nil {
			return nil, err
		}

		IDs = append(IDs, id)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return IDs, nil
}

func (or *PostgresOrderRepository) UpdateOrdersPaymentID(ctx context.Context, orderID int64, paymentID int64) (error) {
	query := `
		UPDATE orders
		SET payment_id = $1
			WHERE id = $2
	`
	result, err := or.db.ExecContext(ctx, query, paymentID, orderID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return ErrFailUpdate
	}

	return nil
}

func (or *PostgresOrderRepository) GetPendingOrderByUserID(ctx context.Context, userID int64, timePeriod time.Duration) (*models.Order, error) {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	//get order
	query := `
		SELECT id, user_id, payment, payment_id, amount, created_at, expires_at, updated_at
		FROM orders
		WHERE user_id = $1
			AND payment = 'pending'

	`
	row := tx.QueryRowContext(ctx, query, userID)

	var order models.Order
	err = row.Scan(
		&order.ID,
		&order.UserID,
		&order.Payment,
		&order.PaymentID,
		&order.Amount,
		&order.CreatedAt,
		&order.ExpiresAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRowNotFound
	}
	if err != nil {
		return nil, err
	}

	//update expires_at
	query = `
		UPDATE orders
		SET expires_at = $1,
			updated_at = NOW()
		WHERE id = $2
			AND payment = 'pending'
			AND expires_at >= NOW()
	`
	expiresAt := time.Now().Add(timePeriod)
	result, err := tx.ExecContext(ctx, query, expiresAt, order.ID)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, ErrFailUpdate
	}

	order.ExpiresAt = expiresAt

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}

func (or *PostgresOrderRepository) UpdateOrderCancel(ctx context.Context, orderID int64) (error) {
	tx, err := or.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//update payment and updated_at
	query := `
		UPDATE orders
		SET payment = 'canceled',
			updated_at = NOW()
		WHERE id = $1
			AND payment = 'pending'
	`
	result, err := tx.ExecContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrNonPending
	}

	//get order items
	query = `
		SELECT product_id, quantity
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := tx.QueryContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	items := make([]*models.OrderItem, 0)
	for rows.Next() {
		var item models.OrderItem
		if err = rows.Scan(
			&item.ProductID, 
			&item.Quantity,
		); err != nil {
			rows.Close()
			return err
		}

		if item.ProductID == nil {
			continue
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err = rows.Close(); err != nil {
		return err
	}

	//restore inventory
	query = `
		UPDATE products
		SET quantity = quantity + $1
		WHERE id = $2
	`
	for _, item := range items {
		result, err = tx.ExecContext(ctx, query, item.Quantity, *item.ProductID)
		if err != nil {
			return err
		}

		affected, err = result.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return ErrProductUnavailable
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}