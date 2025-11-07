package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"example.com/greetings"
	"example.com/repository"
	"github.com/gofrs/uuid/v5"
	uuidext "github.com/jackc/pgx-gofrs-uuid"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	names := []string{"a", "b", "c"}

	message, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ordersRepository, err := repository.NewOrdersRepository(ctx, "postgresql://postgres:postgres@localhost:5432/example_data?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer ordersRepository.Close(ctx)

	orders := getAllOrders(err, ordersRepository, ctx)
	if len(orders) == 0 {
		return
	}

	getOrderById(err, ordersRepository, ctx, orders[0].ID)

	createOrder(err, ordersRepository, ctx, "apple", 1)

	deleteOrder(err, ordersRepository, ctx, orders[0].ID)
}

func getAllOrders(
	err error,
	ordersRepository *repository.OrdersRepository,
	ctx context.Context,
) []repository.Order {
	orders, err := ordersRepository.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(orders) == 0 {
		fmt.Println("no orders")
		return nil
	}

	fmt.Println("orders:")
	for _, order := range orders {
		fmt.Printf(
			"ID: %s, Product: %s, Amount: %d\n",
			uuid.UUID(order.ID),
			order.Product,
			order.Amount,
		)
	}
	return orders
}

func getOrderById(
	err error,
	ordersRepository *repository.OrdersRepository,
	ctx context.Context,
	orderUuid uuidext.UUID,
) {
	order, err := ordersRepository.GetByUuid(ctx, orderUuid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"Found order - ID: %s, Product: %s, Amount: %d\n",
		uuid.UUID(order.ID),
		order.Product,
		order.Amount,
	)
}

func createOrder(
	err error,
	ordersRepository *repository.OrdersRepository,
	ctx context.Context,
	product string,
	amount int,
) {
	order, err := ordersRepository.Create(ctx, product, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"Created order - ID: %s, Product: %s, Amount: %d\n",
		uuid.UUID(order.ID),
		order.Product,
		order.Amount,
	)
}

func deleteOrder(
	err error,
	ordersRepository *repository.OrdersRepository,
	ctx context.Context,
	orderUuid uuidext.UUID,
) {
	err = ordersRepository.DeleteByUuid(ctx, orderUuid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted order - ID: %s\n", uuid.UUID(orderUuid))
}
