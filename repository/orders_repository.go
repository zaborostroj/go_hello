package repository

import (
	"context"
	"fmt"

	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v5"
)

type Order struct {
	ID      uuid.UUID
	Product string
	Amount  int
}

type OrdersRepository struct {
	connection *pgx.Conn
}

func NewOrdersRepository(context context.Context, connectionString string) (*OrdersRepository, error) {
	connection, err := pgx.Connect(context, connectionString)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	return &OrdersRepository{connection: connection}, nil
}

func (ordersRepository *OrdersRepository) Close(context context.Context) error {
	return ordersRepository.connection.Close(context)
}

func (ordersRepository *OrdersRepository) GetAll(context context.Context) ([]Order, error) {
	rows, err := ordersRepository.connection.Query(context, "SELECT uuid, product, amount FROM orders")
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос к БД: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.Product, &order.Amount); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", err)
		}

		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
