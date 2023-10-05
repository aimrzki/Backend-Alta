package model

import "time"

type Event struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Title            string     `json:"title"`
	Location         string     `json:"location"`
	Description      string     `json:"description"`
	Price            int        `json:"price"`
	UserID           uint       `json:"user_id"` // ID pengguna yang membuat event
	AvailableTickets int        `json:"available_tickets"`
	CreatedAt        *time.Time `json:"created_at"` // Kolom created_at yang diharapkan tipe data *time.Time
	UpdatedAt        time.Time
}

// Metode untuk mengurangi jumlah tiket yang tersedia
func (e *Event) DecrementTickets(quantity int) {
	if e.AvailableTickets >= quantity {
		e.AvailableTickets -= quantity
	}
}
