package dto

import "time"

type CurrencyRequest struct {
	Currency string `form:"currency" binding:"required"`
	DateFrom string `form:"date_from" binding:"required,datetime=2006-01-02"`
	DateTo   string `form:"date_to" binding:"required,datetime=2006-01-02"`
}

type ParsedCurrencyRequest struct {
	Currency string
	DateFrom time.Time
	DateTo   time.Time
}

type CurrencyResponse struct {
	Currency string
	Rates    []CurrencyRate
}

type CurrencyRate struct {
	Rate float32
	Date time.Time
}

type RegisterRequest struct {
	Username string
	Password string
}

type RegisterResponse struct {
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
}
