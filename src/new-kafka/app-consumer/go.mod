module github.com/pseudoelement/new-kafka/app-consumer

go 1.24.3

// github.com/pseudoelement/new-kafka/common v0.0.0
require github.com/segmentio/kafka-go v0.4.51

require github.com/pseudoelement/new-kafka/common v0.0.0-00010101000000-000000000000

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

replace github.com/pseudoelement/new-kafka/common => ../common
