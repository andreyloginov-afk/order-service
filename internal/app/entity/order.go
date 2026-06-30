package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

const (
	tableNameOrder     = `orders`
	tableNameOrderItem = `order_items`
)

type Order struct {
	ID         int64      `gorm:"autoIncrement"`
	GUID       uuid.UUID  `gorm:"primaryKey;not null"`
	UserGUID   *uuid.UUID `gorm:"column:user_guid"`
	TotalPrice int64      `gorm:"not null"`
	Currency   string     `gorm:"not null"`
	Status     string     `gorm:"not null"`
	CreatedAt  time.Time  `gorm:"not null;default:now()"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()"`

	Items []OrderItem `gorm:"foreignKey:OrderGUID;references:GUID"`
}

func (Order) TableName() string { return tableNameOrder }

type OrderItem struct {
	ID          int64     `gorm:"autoIncrement;unique;not null"`
	GUID        uuid.UUID `gorm:"primaryKey;not null"`
	OrderGUID   uuid.UUID `gorm:"not null;column:order_guid"`
	ProductGUID uuid.UUID `gorm:"not null;column:product_guid"`
	Quantity    int       `gorm:"not null"`
	UnitPrice   int64     `gorm:"not null;column:unit_price"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}

func (OrderItem) TableName() string { return tableNameOrderItem }

type RequestOrderCreate struct {
	UserGUID *uuid.UUID               `json:"user_guid"`
	Currency string                   `json:"currency" binding:"required,len=3"`
	Items    []RequestOrderItemCreate `json:"items"    binding:"required,min=1,dive"`
}

type RequestOrderItemCreate struct {
	ProductGUID uuid.UUID `json:"product_guid" binding:"required"`
	Quantity    int       `json:"quantity"     binding:"required,gt=0"`
	UnitPrice   int64     `json:"unit_price"   binding:"required,gt=0"`
}

type RequestOrderUpdate struct {
	Status string `json:"status" binding:"required"`
}

type RequestOrderList struct {
	Status   *string    `json:"status" binding:"omitempty"`
	UserGUID *uuid.UUID `json:"user_guid" binding:"omitempty"`
}

type ResponseOrderItem struct {
	GUID        uuid.UUID `json:"guid"`
	ProductGUID uuid.UUID `json:"product_guid"`
	Quantity    int       `json:"quantity"`
	UnitPrice   int64     `json:"unit_price"`
}

type ResponseOrderCreate struct {
	GUID       uuid.UUID  `json:"guid"`
	UserGUID   *uuid.UUID `json:"user_guid"`
	Status     string     `json:"status"`
	Currency   string     `json:"currency"`
	TotalPrice int64      `json:"total_price"`
	CreatedAt  time.Time  `json:"created_at"`

	Items []ResponseOrderItem `json:"items"`
}

type ResponseOrderGet struct {
	GUID       uuid.UUID  `json:"guid"`
	Status     string     `json:"status"`
	Currency   string     `json:"currency"`
	TotalPrice int64      `json:"total_price"`
	UserGUID   *uuid.UUID `json:"user_guid"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	Items []ResponseOrderItem `json:"items"`
}

type ResponseOrderUpdate struct {
	GUID      uuid.UUID `json:"guid"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseOrderListItem struct {
	GUID       uuid.UUID  `json:"guid"`
	UserGUID   *uuid.UUID `json:"user_guid"`
	Status     string     `json:"status"`
	TotalPrice int64      `json:"total_price"`
	Currency   string     `json:"currency"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ResponseOrderList struct {
	Data []ResponseOrderListItem `json:"data"`
}
