package models

import "time"

type Borrower struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;not null" json:"id,omitempty"`
	Name      string    `gorm:"not null" json:"name,omitempty"`
	Address   string    `json:"address,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type ReqBorrower struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
}
