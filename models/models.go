package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleSalesman UserRole = "salesman"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Phone     string             `bson:"phone" json:"phone"`
	Password  string             `bson:"password" json:"-"`
	Role      UserRole           `bson:"role" json:"role"`
	DeletedAt *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type Product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BasepackCode string             `bson:"basepack_code" json:"basepack_code"`
	SKU7         string             `bson:"sku7" json:"sku7"`
	ProductName  string             `bson:"product_name" json:"product_name"`
	HSNNumber    string             `bson:"hsn_number" json:"hsn_number"`
	Location     string             `bson:"location" json:"location"`
	Category     string             `bson:"category" json:"category"`
	ExpiryDate   *time.Time         `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	UPC          string             `bson:"upc" json:"upc"`
	Units        int                `bson:"units" json:"units"`
	StocksInDays int                `bson:"stocks_in_days" json:"stocks_in_days"`
	PurRate      float64            `bson:"pur_rate" json:"pur_rate"`
	MRP          float64            `bson:"mrp" json:"mrp"`
	CurStkValue  float64            `bson:"cur_stk_value" json:"cur_stk_value"`
	DeletedAt    *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type UserProductChange struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProductID     primitive.ObjectID `bson:"product_id" json:"product_id"`
	ChangeType    string             `bson:"change_type" json:"change_type"`
	ChangeDetails string             `bson:"change_details" json:"change_details"`
	Timestamp     time.Time          `bson:"timestamp" json:"timestamp"`
}
