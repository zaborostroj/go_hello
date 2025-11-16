package orders

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB

func setupPostgresContainer(t *testing.T) (tc.Container, string) {
	ctx := context.Background()
	req := tc.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "example_data",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/example_data?sslmode=disable", host, port.Port())

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open DB connection: %v", err)
	}

	// Wait for DB to be ready
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Migrate DB
	_, err = db.Exec(`
		CREATE EXTENSION IF NOT EXISTS pgcrypto;
		CREATE TABLE IF NOT EXISTS orders (
			uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product TEXT NOT NULL,
			amount INTEGER NOT NULL CHECK (amount >= 0),
			user_uuid UUID DEFAULT gen_random_uuid()
		);
	`)
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	return container, port.Port()
}

func TestOrdersService_CRUD(t *testing.T) {
	ctx := context.Background()
	container, port := setupPostgresContainer(t)
	defer func(container tc.Container, ctx context.Context, opts ...tc.TerminateOption) {
		err := container.Terminate(ctx, opts...)
		if err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}(container, ctx)

	err := os.Setenv("APP_ENV", "test")
	if err != nil {
		return
	}
	err = os.Setenv("DB_PORT", port)
	if err != nil {
		return
	}

	// --- Run OrdersService ---
	go func() {
		Start()
	}()
	time.Sleep(3 * time.Second)

	client := &http.Client{}
	baseURL := "http://localhost:8082"

	// --- POST /orders ---
	body := map[string]interface{}{
		"product": "Book",
		"amount":  3,
	}
	data, _ := json.Marshal(body)
	resp, err := client.Post(baseURL+"/orders", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("POST /orders failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	// --- GET /orders ---
	resp, err = client.Get(baseURL + "/orders")
	if err != nil {
		t.Fatalf("GET /orders failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	respData, _ := io.ReadAll(resp.Body)
	t.Logf("Orders: %s", string(respData))

	// --- Get order UUID ---
	var orders []map[string]interface{}
	if err := json.Unmarshal(respData, &orders); err != nil {
		t.Fatalf("invalid JSON in GET /orders: %v", err)
	}
	if len(orders) == 0 {
		t.Fatalf("no orders returned")
	}
	uuid := orders[0]["ID"].(string)

	//--- DELETE /orders/:uuid ---
	req, _ := http.NewRequest(http.MethodDelete, baseURL+"/orders/"+uuid, nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("DELETE /orders failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 200/204 after delete, got %d", resp.StatusCode)
	}
}
