package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
}

type Role string

const (
	Chef     Role = "chef"
	Customer Role = "customer"
	Manager  Role = "manager"
)
