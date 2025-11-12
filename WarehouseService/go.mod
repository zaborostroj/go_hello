module WarehouseService

go 1.25.3

replace example.com/KafkaUtils => ../KafkaUtils

require example.com/KafkaUtils v0.0.0-00010101000000-000000000000

require (
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/segmentio/kafka-go v0.4.49 // indirect
)
