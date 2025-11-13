package main

type WarehouseServiceConfig struct {
	KAFKA struct {
		Host    string
		Port    string
		Topic   string
		GroupId string
	}
}
