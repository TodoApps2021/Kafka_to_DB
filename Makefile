run:
	go run cmd/consumer/main.go
	
docker-up:
	docker-compose up

docker-down:
	docker-compose down

lint: 
	golangci-lint run