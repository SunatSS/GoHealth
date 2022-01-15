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
