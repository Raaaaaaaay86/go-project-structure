DB_SCHEMA=video
DB_ROOT_PASSWORD=123456

# Application
set_up:
	bash ./scripts/setup.sh
run:
	docker-compose up -d
	migrate -path migration/postgres -database "postgres://root:123456@localhost:5432/$(DB_SCHEMA)?sslmode=disable" -verbose up
	go run main.go
swag_init:
	swag init --generalInfo ./adapter/port_in/http/gin.go

# Unit test
test:
	go test -v ./...
generate_mocks:
	docker run -v $(PWD):/src -w /src vektra/mockery --all
clean_mock_containers:
	docker ps -a | grep mockery | awk '{print $1}' | xargs docker rm

# Database Version Control
migrate_new:
	migrate create -ext sql -dir migration/postgres -seq $(DB_SCHEMA)
migrate_force:
	migrate -path migration/postgres -database "postgres://root:123456@localhost:5432/$(DB_SCHEMA)?sslmode=disable" -verbose force 3
migrate_up:
	migrate -path migration/postgres -database "postgres://root:123456@localhost:5432/$(DB_SCHEMA)?sslmode=disable" -verbose up
migrate_down:
	migrate -path migration/postgres -database "postgres://root:123456@localhost:5432/$(DB_SCHEMA)?sslmode=disable" -verbose down

init_mongo:
	docker exec -it video-mongodb mongosh -u root -p 123456 --quiet --eval "rs.initiate();"