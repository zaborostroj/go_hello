package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"example.com/greetings"
	"example.com/repository"
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

	orders, err := ordersRepository.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(orders) == 0 {
		fmt.Println("no orders")
		return
	}

	fmt.Println("orders:")
	for _, order := range orders {
		fmt.Printf("ID: %s, Product: %s, Amount: %d\n", order.ID, order.Product, order.Amount)
	}
}
