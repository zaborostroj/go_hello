package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"example.com/repository"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	uuidext "github.com/jackc/pgx-gofrs-uuid"
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

	ordersHandler := &OrdersHandler{ordersRepository: ordersRepository}

	router := gin.Default()
	router.GET("/orders", ordersHandler.getAll)
	router.GET("/orders/:uuid", ordersHandler.getByUuid)
	router.Run("localhost:8080")
}

type OrdersHandler struct {
	ordersRepository *repository.OrdersRepository
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

func (handler *OrdersHandler) createOrder(context *gin.Context) {

}
