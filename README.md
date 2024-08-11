Newsfeed

https://drive.google.com/file/d/1Qnca6dITavOpRk4mYVik2BCb-G9nsaPi/view

1. Init project structure
We will have 3 services in this project: 
- HTTP server: works like a gateway to receive HTTP requests, validate and forward data to internal GRPC services.
- User GRPC server: CRUD for users, posts, follows, likes, and comments
- Newsfeed GRPC server: creates and returns newsfeed.

The project will be structured with the following folders:

- `/cmd`: The entry point for each service, containing main function. We will initializes all configs, service dependencies and the main services will be called to run here.
- `/config`: Configurations will be stored here.
- `/pkg`: Our backend implementation will be here. We can use `/internal` with the same purpose with `/pkg`. It's based on personal choice.
  - `/pkg/handler`: This folder will handle the request/response to each service. For HTTP server, it will have a HTTP server (gin). For GRPC servers, it will have a GRPC server attached with GRPC handler. The data model for these handlers will be XxxRequest and XxxResponse structs.
  - `/pkg/service`: This folder will store the business logic. The data model is business model. `Handler` and `Repo` layers will communicate with this `Service` layer by the business model here.
  - `/pkg/repo`: This is where we perform the actual operations to the external storage, such as database, cache, message queues ... The data model is defined based on each external storage library.
  - `/pkg/util`: Utility functions defined here.

2. Handler
3. Service
4. Repo
5. Setup local external dependencies
- MySQL:
  - Steps: pull image -> run mysql container with configs -> run SQL to create schemas
  - Pull: `docker pull mysql`
  - Run: `docker run --name newsfeed-db -e MYSQL_ROOT_PASSWORD=123456 -p 33334:3306 -v {{absolute_path}}/data/mysql:/var/lib/mysql mysql`
  - Exec `mysql` client in container: `docker exec -it newsfeed-db mysql -p`
- Redis: TBD
- Kafka: TBD
- Prometheus/Grafana: TBD

7. Generate GRPC servers and clients code:
- Quick start References: https://grpc.io/docs/languages/go/quickstart/, https://grpc.io/docs/languages/go/basics/, https://grpc.io/docs/languages/go/generated-code/
- Example code: https://github.com/grpc/grpc-go/tree/master/examples/helloworld
- Steps:
  - Download and install: protoc, protoc-gen-go, protoc-gen-go-grpc (follow Quick Start reference)
  - Write proto: `/proto/user.proto`
  - Generate grpc server: `protoc --go_out=proto/gen --go_opt=paths=source_relative --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative proto/user.proto`

8. Run demo:
- Setup MySQL in Section 5.
- Run User GRPC server: `go run cmd/user_grpc/main.go`
- Run HTTP server: `go run cmd/http/main.go`
- Test request: 
```
curl --location 'localhost:33333/signup' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "user_name": "testuser",
  "password": "123456",
  "first_name": "Test",
  "last_name": "User",
  "dob": "01-01-2000",
  "email": "abcxyz@gmail.com"
  }'
```
