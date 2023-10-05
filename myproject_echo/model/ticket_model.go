package model

import "time"

type Ticket struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	EventID       uint       `json:"event_id"` // ID acara yang tiketnya terkait
	UserID        uint       `json:"user_id"`  // ID pengguna yang membeli tiket
	Quantity      int        `json:"quantity"` // Jumlah tiket yang dibeli
	KodeVoucher   string     `json:"kode_voucher"`
	TotalCost     int        `json:"total_cost"`
	InvoiceNumber string     `json:"invoice_number"` // Nomor invoice untuk tiket
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     time.Time
}
