package models

import "time"

type BillingSchedule struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	LoanID        uint      `json:"loan_id"`
	Amount        uint      `json:"amount"`
	ScheduledWeek uint      `json:"scheduled_week"`
	Status        string    `json:"status"`
	StartWeekDate time.Time `json:"start_week_date"`
	EndWeekDate   time.Time `json:"end_week_date"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type MakePaymentRequest struct {
	Amount      uint `json:"amount"`
	PaymentWeek uint `json:"payment_week"`
}

type IsDelinquentRequest struct {
	CurrentDate time.Time `json:"current_date"`
}
type IsDelinquentResponse struct {
	IsDelinquent bool `json:"is_delinquent"`
}

type GetOutstandingResponse struct {
	OutstandingAmount uint   `json:"outstanding_amount"`
	LatestPaymentWeek uint   `json:"latest_payment_week"`
	LoanStatus        string `json:"loan_status"`
	NextPaymentWeek   uint   `json:"next_payment_week,omitempty"`
}

type MakePaymentResponse struct {
}
