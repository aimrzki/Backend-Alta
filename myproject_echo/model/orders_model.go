package model

import "time"

type Order struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `json:"user_id"`    // ID pengguna yang melakukan pesanan
	TotalCost  int        `json:"total_cost"` // Total biaya pesanan
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  time.Time
	OrderItems []OrderItem // Detail pesanan yang berisi tiket yang dibeli
}

type OrderItem struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	OrderID  uint `json:"order_id"`  // ID pesanan yang tiketnya terkait
	TicketID uint `json:"ticket_id"` // ID tiket yang dibeli
	Quantity int  `json:"quantity"`  // Jumlah tiket yang dibeli dalam pesanan ini
	Subtotal int  `json:"subtotal"`  // Total biaya untuk tiket ini dalam pesanan ini
}
