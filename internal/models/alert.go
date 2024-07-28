package models

type Alert struct {
	Currency string  `json:"currency" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
}

type DeleteAlert struct {
	AlertID string `json:"alert_id" validate:"required"`
}
