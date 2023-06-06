package models

import "time"

type Product struct {
	Id         uint64    `json:"id" binding:"required" example:"5"`
	Name       string    `json:"name" binding:"required" example:"Product name"`
	MerchantId uint64    `json:"merchant_id"`
	CategoryId *uint64   `json:"category_id,omitempty"`
	Tags       []uint64  `json:"tags" example:"1"`
	Images     []string  `json:"image_ids" example:"image_id"`
	Sizes      []Sizes   `json:"sizes"`
	SaleDate   time.Time `json:"sale_date"`
	EndDate    time.Time `json:"end_date" binding:"required"`
	Exclude    string    `json:"-"`
	Type       string    `json:"type" binding:"required" example:"Physical" enum:"Physical|Online"`
}

type Sizes struct {
	Size string `json:"size" example:"size" binding:"required"`
}

type ProductPost struct {
	Name       string  `json:"name"`
	MerchantId uint64  `json:"merchant_id"`
	CategoryId *uint64 `json:"category_id,omitempty"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type MerchantPageResponse struct {
	Brochures    []Product `json:"products"`
	MerchantName string    `json:"merchantName"`
}
