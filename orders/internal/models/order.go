package models

type Order struct {
	ID              int         `json:"id"`
	UserID          int         `json:"user_id"`
	Status          string      `json:"status"`
	Dishes          []OrderDish `json:"dishes"`
	SpecialRequests string      `json:"special_requests"`
}

type OrderDish struct {
	DishID   int     `json:"dish_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
