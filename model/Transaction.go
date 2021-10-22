package model

import (
	"errors"
	"fmt"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

type TransactionType string
const(
	Debit TransactionType = "DEBIT"
	Credit	 = "CREDIT"
)

type Transaction struct{
	ID        	string    	`json:"id";gorm:"primary_key"`
	Value		float64		`gorm:"not null";sql:"type:decimal(8,2)";json:"value"`
	TransactionType	TransactionType	`json:"type"`
	OcurredIn		string		`json:"occurredIn"`
	Customer		Customer		`gorm:"embedded;embeddedPrefix:customer_";json:"customer"`
}

func (t *Transaction) Validate() error{
	if t.ID == ""{
		return errors.New("required id")
	}
	if t.TransactionType == ""{
		return errors.New("required transaction type")
	}
	if t.Customer.Name == ""{
		return errors.New("required customer name")
	}
	if t.Customer.CustomerId == ""{
		return errors.New("Customer required")
	}
	if t.Customer.Account.AccountNumber == ""{
		return errors.New("Account name required")
	}

	if err := checkmail.ValidateFormat(t.Customer.Email); err != nil {
		return errors.New("invalid email")
	}
	return nil
}

func (t *Transaction) SaveTransaction(db *gorm.DB) (*Transaction, error){
	fmt.Println("Id da transacao: "+t.ID)
	var err error
	result := db.Create(&t)
	if result.Error != nil {
		fmt.Printf(result.Error.Error())
		return &Transaction{},err
	}
	fmt.Printf(t.ID)
	return t,err
}

func (t *Transaction) FindAllTransactions( db *gorm.DB) (*[]Transaction, error){
	var err  error
	transactions :=[]Transaction{}
	err = db.Debug().Model(&Transaction{}).Limit(100).Find(&transactions).Error
	if err != nil {
		return &[]Transaction{},err
	}
	return &transactions,err
}

func (t *Transaction) FindTransactionById(db *gorm.DB,tid string) (*Transaction,error){
	var err  error
	err = db.Debug().Model(Transaction{}).Where("id = ?",tid).Take(&t).Error
	if err != nil {
		return &Transaction{},err
	}
	if gorm.IsRecordNotFoundError(err){
		return &Transaction{}, errors.New("Transaction not found")
	}
	return t, err
}

func (t *Transaction) FindTransactionByCustomerId(db *gorm.DB,cid string) (*[]Transaction,error){
	var err  error
	transactions :=[]Transaction{}
	err = db.Debug().Model(&Transaction{}).Limit(100).Find(&transactions).Where("customer.customer_id = ?",cid).Error
	if err != nil {
		return &[]Transaction{},err
	}
	if gorm.IsRecordNotFoundError(err){
		return &[]Transaction{},errors.New("Any transaction was found to this customer")
	}
	return &transactions,err
}