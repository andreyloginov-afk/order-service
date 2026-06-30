package horder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"github.com/andreyloginov-afk/order-service/internal/app/entity"
	rhandler "github.com/andreyloginov-afk/order-service/internal/app/handler/http"
	"github.com/andreyloginov-afk/order-service/internal/app/service"
	"github.com/andreyloginov-afk/order-service/internal/pkg/http/httph"
)

type handler struct {
	srv service.Order
}

func NewHandler(srv service.Order) rhandler.Order {
	return &handler{srv: srv}
}

func (h *handler) Create(c *gin.Context) {
	var req entity.RequestOrderCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.Create(c.Request.Context(), req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	httph.SendJSON(c.Writer, http.StatusCreated, entity.ResponseOrderCreate{
		GUID:       order.GUID,
		Status:     order.Status,
		Currency:   order.Currency,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		Items:      toResponseItems(order.Items),
	})
}

func (h *handler) GetByGUID(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.GetByGUID(c.Request.Context(), guid)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	httph.SendJSON(c.Writer, http.StatusOK, entity.ResponseOrderGet{
		GUID:       order.GUID,
		Status:     order.Status,
		Currency:   order.Currency,
		TotalPrice: order.TotalPrice,
		UserGUID:   order.UserGUID,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		Items:      toResponseItems(order.Items),
	})
}

func (h *handler) Update(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	var req entity.RequestOrderUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.Update(c.Request.Context(), guid, req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	httph.SendJSON(c.Writer, http.StatusOK, entity.ResponseOrderUpdate{
		GUID:      order.GUID,
		Status:    order.Status,
		UpdatedAt: order.UpdatedAt,
	})
}

func (h *handler) Delete(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	if err := h.srv.Delete(c.Request.Context(), guid); err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	httph.SendEmpty(c.Writer, http.StatusNoContent)
}

func (h *handler) List(c *gin.Context) {
	var req entity.RequestOrderList
	if c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
			return
		}
	}

	orders, err := h.srv.List(c.Request.Context(), req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	data := make([]entity.ResponseOrderListItem, len(orders))
	for i, o := range orders {
		data[i] = entity.ResponseOrderListItem{
			GUID:       o.GUID,
			UserGUID:   o.UserGUID,
			Status:     o.Status,
			TotalPrice: o.TotalPrice,
			Currency:   o.Currency,
			CreatedAt:  o.CreatedAt,
		}
	}

	httph.SendJSON(c.Writer, http.StatusOK, entity.ResponseOrderList{Data: data})
}

func toResponseItems(items []entity.OrderItem) []entity.ResponseOrderItem {
	resp := make([]entity.ResponseOrderItem, len(items))
	for i, item := range items {
		resp[i] = entity.ResponseOrderItem{
			GUID:        item.GUID,
			ProductGUID: item.ProductGUID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}
	return resp
}
