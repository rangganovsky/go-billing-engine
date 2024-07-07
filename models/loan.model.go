package models

import "time"

type Loan struct {
	ID                uint      `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	BorrowerID        uint      `json:"borrower_id"`
	PrincipalAmount   uint      `json:"principal_amount"`
	InterestRate      float64   `json:"interest_rate"`
	TotalAmount       uint      `json:"total_amount"`
	OutstandingAmount uint      `json:"outstanding_amount"`
	LoanTermWeeks     uint      `json:"loan_term_weeks"`
	Status            string    `json:"status"`
	CurrentWeek       uint      `json:"current_week"`
	IsDelinquent      bool      `json:"deliquency_status"`
	StartLoanDate     time.Time `json:"start_loan_date,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

type ReqLoan struct {
	BorrowerID      uint    `json:"borrower_id"`
	PrincipalAmount uint    `json:"principal_amount"`
	InterestRate    float64 `json:"interest_rate"`
	LoanTermWeeks   uint    `json:"loan_term_weeks"`
}
