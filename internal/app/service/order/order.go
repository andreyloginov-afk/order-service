package sorder

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/andreyloginov-afk/order-service/internal/app/entity"
	"github.com/andreyloginov-afk/order-service/internal/app/repository"
	"github.com/andreyloginov-afk/order-service/internal/app/service"
)

type srv struct {
	repoOrder repository.Order
}

func NewService(repoOrder repository.Order) service.Order {
	return &srv{repoOrder: repoOrder}
}

func (s *srv) Create(ctx context.Context, req entity.RequestOrderCreate) (entity.Order, error) {
	now := time.Now()

	items := make([]entity.OrderItem, 0, len(req.Items))
	var totalPrice int64

	for _, item := range req.Items {
		items = append(items, entity.OrderItem{
			GUID:        uuid.Must(uuid.NewV4()),
			ProductGUID: item.ProductGUID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		totalPrice += int64(item.Quantity) * item.UnitPrice
	}

	order := entity.Order{
		GUID:       uuid.Must(uuid.NewV4()),
		UserGUID:   req.UserGUID,
		Currency:   req.Currency,
		Status:     entity.OrderStatusPending,
		TotalPrice: totalPrice,
		CreatedAt:  now,
		UpdatedAt:  now,
		Items:      items,
	}

	err := s.repoOrder.Create(ctx, order)
	return order, err
}

func (s *srv) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Order, error) {
	return s.repoOrder.GetByGUID(ctx, guid)
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	return s.repoOrder.InsideTx(ctx, func(ctx context.Context) error {
		_, err := s.repoOrder.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}
		return s.repoOrder.Delete(ctx, guid)
	})
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestOrderUpdate) (entity.Order, error) {
	var order entity.Order
	err := s.repoOrder.InsideTx(ctx, func(ctx context.Context) error {
		var err error
		order, err = s.repoOrder.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}
		order.Status = req.Status
		order.UpdatedAt = time.Now()
		return s.repoOrder.Update(ctx, order)
	})
	return order, err
}

func (s *srv) List(ctx context.Context, req entity.RequestOrderList) ([]entity.Order, error) {
	return s.repoOrder.List(ctx, req.Status, req.UserGUID)
}
