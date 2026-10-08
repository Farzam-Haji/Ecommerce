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
	UpdateOrdersPaymentID(ctx context.Context, orderID int64, paymentID int64) (error)
	GetPendingOrderByUserID(ctx context.Context, userID int64, timePeriod time.Duration) (*models.Order, error)
	UpdateOrderCancel(ctx context.Context, orderID int64) (error)
}

type BankClient interface {
	CreatePaymentURL(orderID int64, amount int64) (string, error)
}

type OrderService struct {
	orderRepo 	OrderRepository
	bankClient 	BankClient
}

func NewOrderService(orderRepoInterface OrderRepository, bankClient BankClient) (*OrderService) {
	return &OrderService{
		orderRepo: orderRepoInterface,
		bankClient: bankClient,
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
	//order exist
	order, err := os.orderRepo.GetPendingOrderByUserID(ctx, userID, ExpiresAfter)
	if err == nil {
		return order, nil
	}
	if !errors.Is(err, repositories.ErrRowNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}


	//new order
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
		return nil, ErrForbidden
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

func (os *OrderService) GetOrderAmount(ctx context.Context, orderID int64) (int64, error) {
	items, err := os.GetOrderItems(ctx, orderID)
	if err != nil {
		return 0, err
	}

	var amount int64
	for _, item := range items {
		amount += int64(item.Quantity) * item.Price
	}
	return amount, nil
}

func (os *OrderService) CreatePayment(ctx context.Context, orderID int64) (string, error) {
	order, err := os.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return os.bankClient.CreatePaymentURL(order.ID, order.Amount)
}

func (os *OrderService) ProcessPayment(ctx context.Context, orderID int64, paymentID int64, status string) (error) {
	var err error

	if status == "paid" {
		err = os.ResolvePaidOrder(ctx, orderID)
	}
	if status == "failed" {
		err = os.ResolveFailedOrder(ctx, orderID)
	}
	errUpdate := os.orderRepo.UpdateOrdersPaymentID(ctx, orderID, paymentID)

	if err != nil {
		return err
	}

	if errUpdate != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return nil
}

func (os *OrderService) CancelUserOrder(ctx context.Context, orderID int64, userID int64) (error) {
	order, err := os.orderRepo.GetOrderByID(ctx, orderID)

	if errors.Is(err, repositories.ErrRowNotFound) {
		return ErrOrderNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	if *order.UserID != userID {
		return ErrForbidden
	}

	err = os.orderRepo.UpdateOrderCancel(ctx, orderID);
	if errors.Is(err, repositories.ErrNonPending) {
		return ErrNonPending
	}
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return nil
}
