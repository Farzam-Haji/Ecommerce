package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/btalv/Ecommerce/server/internal/repositories"
)

const ExpiresAfter time.Duration = 10 * time.Minute

type OrderRepository interface {
	GetOrderByID(ctx context.Context, orderID int64) (*models.Order, error)
	GetOrderItems(ctx context.Context, orderID int64) ([]*models.OrderItem, error)
	CreateOrderByUserID(ctx context.Context, userID int64, timePeriod time.Duration) (*models.Order, []string, error)
	ResolvePaidOrder(ctx context.Context, orderID int64) (error)
	ResolveFailedOrder(ctx context.Context, orderID int64) (error)
	ResolveExpiredOrder(ctx context.Context, orderID int64) (error)
	GetExpiredOrderIDs(ctx context.Context) ([]int64, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error)
}

type OrderService struct {
	orderRepo 	OrderRepository
}

func NewOrderService(orderRepoInterface OrderRepository) (*OrderService) {
	return &OrderService{
		orderRepo: orderRepoInterface,
	}
}

func (os *OrderService) ExpireOrders(ctx context.Context) (error) {
	expiredIDs, err := os.orderRepo.GetExpiredOrderIDs(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	for _, id := range expiredIDs {
		err = os.orderRepo.ResolveExpiredOrder(ctx, id)
		if err != nil {	
			return fmt.Errorf("%w: %w", ErrInternal, err)
		}
	}

	return nil
}

func (os *OrderService) CreateOrder(ctx context.Context, userID int64) (*models.Order, error) {
	order, insufficientItems, err := os.orderRepo.CreateOrderByUserID(ctx, userID, ExpiresAfter)
	switch{
	case err == nil:
		return order, nil
	case errors.Is(err, repositories.ErrCartNotFound):
		return nil, ErrCartNotFound
	case errors.Is(err, repositories.ErrProductUnavailable):
		return nil, fmt.Errorf("Check Out Failed: %w", ErrProductNotFound)
	case errors.Is(err, repositories.ErrEmptyCart):
		return nil, fmt.Errorf("Check Out Failed: %w", ErrEmptyOrder)
	case errors.Is(err, repositories.ErrInsufficientInventory):
		return nil, fmt.Errorf("%w: Insufficient Items: %v", ErrInsufficientQuantity, insufficientItems)
	default:
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
}

func (os *OrderService) ResolvePaidOrder(ctx context.Context, orderID int64) (error) {
	if err := os.orderRepo.ResolvePaidOrder(ctx, orderID); err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return nil
}

func (os *OrderService) ResolveFailedOrder(ctx context.Context, orderID int64) (error) {
	if err := os.orderRepo.ResolveFailedOrder(ctx, orderID); err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return nil
}

func (os *OrderService) GetUserOrders(ctx context.Context, userID int64) ([]*models.Order, error) {
	orders, err := os.orderRepo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return orders, nil
}

func (os *OrderService) GetUserOrderByID(ctx context.Context, orderID int64, userID int64) (*models.Order, error) {
	order, err := os.orderRepo.GetOrderByID(ctx, orderID)

	if errors.Is(err, repositories.ErrRowNotFound) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	if *order.UserID != userID {
		return nil, ErrUnauthorized
	}
	return order, nil
}

func (os *OrderService) GetOrderItems(ctx context.Context, orderID int64) ([]*models.OrderItem, error) {
	items, err := os.orderRepo.GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return items, nil
}