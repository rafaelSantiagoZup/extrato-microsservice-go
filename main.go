package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ot-bank/extrato/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/segmentio/kafka-go"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)
type Server struct{
	DB *gorm.DB
}

var server Server
func init(){
	var Dbdriver = "postgres"
	var DbUser = "postgres"
	var DbPassword = "changeme"
	var DbPort = "5433"
	var DbHost = "localhost"
	var DbName = "ot-bank-extrato"
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
		server.DB.AutoMigrate(&model.Transaction{})
		
}
func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func (server *Server) getTransactionsById(c *gin.Context){
	tran := model.Transaction{}
	id := c.Param("id")

	res,err := tran.FindTransactionById(server.DB,id)
	if err != nil{
		checkErr(err)
	}
	c.IndentedJSON(http.StatusOK,res)
}

func (server *Server) getTransactionsByCustomerId(c *gin.Context){
	tran := model.Transaction{}
	id := c.Param("customerId")

	res,err := tran.FindTransactionByCustomerId(server.DB,id)
	if err != nil{
		checkErr(err)
	}
	c.IndentedJSON(http.StatusOK,res)
}

func (server *Server) getTransaction(c *gin.Context){
	tran := model.Transaction{}
	res,err := tran.FindAllTransactions(server.DB)
	if err != nil{
		checkErr(err)
	}
	c.JSON(http.StatusOK,res)
}
func NewConsumer (ctx context.Context){
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "transacoes",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
		var m model.Transaction
		er := json.Unmarshal(msg.Value, &m)
		if er != nil{
			fmt.Errorf("erro ao converter mensagem")
		}
		fmt.Println("\n recebido: \n",m)
		transactionCreated,er := m.SaveTransaction(server.DB)
		if er!=nil{
			fmt.Errorf(er.Error())
			fmt.Errorf("Erro ao salvar mensagem")
		}
		fmt.Println("\n recebido: \n",transactionCreated)
	}
}
func main(){
	ctx := context.Background()
	go NewConsumer(ctx)

	router := gin.Default()
	router.GET("/api/v1/transactions", server.getTransaction)
	router.GET("/api/v1/transaction/:id",server.getTransactionsById)
	router.GET("/api/v1/transactions/:customerId",server.getTransactionsByCustomerId)

	router.Run("localhost:8819")
}