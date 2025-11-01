BINARY_NAME=auction
SRC_DIR=cmd/auction
DOCKER_IMAGE=auction-system:latest

build:
	@echo "🔨 Compilando binário local..."
	go build -o $(BINARY_NAME) $(SRC_DIR)/main.go

run:
	@echo "🚀 Executando aplicação localmente..."
	go run $(SRC_DIR)/main.go

test:
	@echo "🧪 Executando testes..."
	go test ./... -v

docker-build:
	@echo "🐳 Build da imagem Docker..."
	docker build -t $(DOCKER_IMAGE) .

up:
	@echo "🚀 Subindo containers (app + MongoDB)..."
	docker-compose up --build

down:
	@echo "🧹 Parando e removendo containers..."
	docker-compose down -v

logs:
	docker-compose logs -f app

lint:
	@echo "🧹 Verificando código com go vet e fmt..."
	go vet ./... && go fmt ./...
