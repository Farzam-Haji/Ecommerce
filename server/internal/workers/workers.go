package workers

import (
	"context"
	"log"
	"time"

	"github.com/Farzam-Haji/Ecommerce/server/internal/services"
)

const ExpireWorkerTicker = 30 * time.Second


func OrderExpirationWorker(ctx context.Context, orderService *services.OrderService) {
	ticker := time.NewTicker(ExpireWorkerTicker)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			if err := orderService.ExpireOrders(ctx); err != nil {
				log.Println(err)
			}

		case <- ctx.Done():
			return
		}
	}
}