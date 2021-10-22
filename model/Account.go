package model

type Account struct {
	Agency        string  `gorm:"size:100;not null";json:"agency"`
	AccountNumber string  `gorm:"size:100;not null";json:"accountNumber`
	Balance       float64 `gorm:"not null;type:numeric(8,2)";json:"balance"`
}