package model

type Customer struct {
	CustomerId string  `json:"customerId"`
	Name       string  `gorm:"size:100;not null";json:"name"`
	Email      string  `gorm:"size:100;not null";json:"email"`
	Account    Account `gorm:"embedded;embeddedPrefix:account_";json:"account"`
}