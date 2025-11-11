package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"example.com/kafka-client"
	"example.com/repository"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	uuidext "github.com/jackc/pgx-gofrs-uuid"
	"github.com/segmentio/kafka-go"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ordersRepository, err := repository.NewOrdersRepository(ctx, "postgresql://postgres:postgres@localhost:5432/example_data?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ordersRepository.Close(ctx)

	kafkaClient := kafka_client.NewClient(kafka_client.Config{
		Brokers: []string{"localhost:29092"},
		Topic:   "orders",
		GroupID: "orders-group",
	})
	defer kafkaClient.Close()

	ordersHandler := &OrdersHandler{
		ordersRepository: ordersRepository,
		kafkaClient:      kafkaClient,
	}

	router := gin.Default()
	router.GET("/orders", ordersHandler.getAll)
	router.GET("/orders/:uuid", ordersHandler.getByUuid)
	router.POST("/orders", ordersHandler.createOrder)
	router.DELETE("/orders/:uuid", ordersHandler.deleteByUuid)
	router.Run("localhost:8080")
}

type OrdersHandler struct {
	ordersRepository *repository.OrdersRepository
	kafkaClient      *kafka_client.Client
}

func (handler *OrdersHandler) getAll(context *gin.Context) {
	orders, err := handler.ordersRepository.GetAll(context)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Something went wrong"})
	}
	context.IndentedJSON(http.StatusOK, orders)
}

func (handler *OrdersHandler) getByUuid(context *gin.Context) {
	orderUuid, err := uuid.FromString(context.Param("uuid"))
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error parsing uuid: " + context.Param("uuid")})
	}
	order, err := handler.ordersRepository.GetByUuid(context, uuidext.UUID(orderUuid))
	context.IndentedJSON(http.StatusOK, order)
}

type ProductRequest struct {
	Product string `json:"product"`
	Amount  int    `json:"amount"`
}

func (handler *OrdersHandler) createOrder(context *gin.Context) {
	var request ProductRequest

	if context.Request.Body == nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Request body required"})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if request.Product == "" || request.Amount == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Product or amount required"})
		return
	}

	result, err := handler.ordersRepository.Create(context, request.Product, request.Amount)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = handler.kafkaClient.Writer().WriteMessages(context, getOrderCreatedMessage(result))
	if err != nil {
		log.Print(err)
	}
	context.IndentedJSON(http.StatusCreated, result)
}

// Order serialized here just for example.
// The message should be like 'The new order created' to prevent synchronization of messages
// on the consumer side.
func getOrderCreatedMessage(order repository.OrderDTO) kafka.Message {
	serializedData, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error serializing order: %s", err)
	}

	return kafka.Message{Key: []byte("orders list"), Value: serializedData}
}

func (handler *OrdersHandler) deleteByUuid(context *gin.Context) {
	orderUuid, err := uuid.FromString(context.Param("uuid"))
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Error parsing uuid: " + context.Param("uuid")})
		return
	}
	err = handler.ordersRepository.DeleteByUuid(context, orderUuid)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
