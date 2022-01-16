package types

import "time"

// Type RegInfo is structure for registration info
type RegInfo struct {
	Name     string
	Phone    string
	Password string
	Address  string
}

// Type Customer is structure with customer data
type Customer struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Address  string    `json:"address"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

// Type TokenInfo is structure of token info
type TokenInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Type Token is structure for token
type Token struct {
	Token      string    `json:"token"`
	CustomerID int64     `json:"customer_id"`
	Expires    time.Time `json:"expires"`
	Created    time.Time `json:"created"`
}

type Response struct {
	Status     string `json:"status"`
	Reason     string `json:"reason,omitempty"`
	CustomerId int64  `json:"customerId,omitempty"`
}

type Medicine struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Manafacturer  string    `json:"manafacturer"`
	Description   string    `json:"description"`
	Components    []string  `json:"components"`
	Recipe_needed bool      `json:"recipe_needed"`
	Price         int       `json:"price"`
	Qty           int       `json:"qty"`
	PharmacyName  string    `json:"pharmacy_name"`
	Active        bool      `json:"active"`
	Created       time.Time `json:"created"`
}
