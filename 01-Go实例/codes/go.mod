module code

go 1.12

require (
	calculator v0.0.1
	github.com/go-redis/redis/v8 v8.0.0-beta.7
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/websocket v1.4.2
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0 // indirect
)

replace calculator => ./calculator
