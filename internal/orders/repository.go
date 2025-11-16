package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	uuidext "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
)

type OrderDTO struct {
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

	uuidext.Register(connection.TypeMap())
	return &OrdersRepository{connection: connection}, nil
}

func (ordersRepository *OrdersRepository) Close(context context.Context) error {
	return ordersRepository.connection.Close(context)
}

func (ordersRepository *OrdersRepository) GetAll(context context.Context) ([]OrderDTO, error) {
	rows, err := ordersRepository.connection.Query(context, "SELECT uuid, product, amount FROM orders")
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос к БД: %w", err)
	}
	defer rows.Close()

	var orders []OrderDTO
	for rows.Next() {
		var order OrderDTO
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

func (ordersRepository *OrdersRepository) GetByUuid(context context.Context, orderUuid uuidext.UUID) (OrderDTO, error) {
	var order OrderDTO

	row := ordersRepository.connection.QueryRow(
		context,
		"SELECT uuid, product, amount FROM orders WHERE uuid = $1",
		orderUuid,
	)
	if err := row.Scan(&order.ID, &order.Product, &order.Amount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order, fmt.Errorf("запись не найдена по uuid: %s", uuid.UUID(orderUuid))
		}
		return order, fmt.Errorf("ошибка чтения строки: %w", err)
	}

	return order, nil
}

func (ordersRepository *OrdersRepository) Create(context context.Context, product string, amount int) (OrderDTO, error) {
	var order OrderDTO

	if product == "" {
		return order, errors.New("product is empty")
	}
	if amount <= 0 {
		return order, errors.New("amount must be greater than 0")
	}

	row := ordersRepository.connection.QueryRow(
		context,
		"INSERT INTO orders (product, amount) VALUES ($1, $2) RETURNING uuid, product, amount",
		product,
		amount,
	)
	if err := row.Scan(&order.ID, &order.Product, &order.Amount); err != nil {
		return order, fmt.Errorf("не удалось создать запись: %w", err)
	}

	return order, nil
}

func (ordersRepository *OrdersRepository) DeleteByUuid(context context.Context, orderUuid uuid.UUID) error {
	tag, err := ordersRepository.connection.Exec(
		context,
		"DELETE FROM orders WHERE uuid = $1",
		orderUuid,
	)
	if err != nil {
		return fmt.Errorf("не удалось выполнить удаление: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("запись не найдена по uuid: %s", uuid.UUID(orderUuid))
	}

	return nil
}
