# 🏷️ Auction System (Go + MongoDB + Fechamento Automático)

Este projeto implementa um **sistema de leilões (auctions)** em Go, utilizando **MongoDB** como banco de dados e suporte a **fechamento automático de leilões** via **goroutines**.

## 🚀 Funcionalidades

- Criar leilões com produto, categoria e condição (novo, usado, recondicionado);
- Receber lances (bids) enquanto o leilão estiver ativo;
- **Fechar automaticamente** os leilões após o tempo configurado;
- Estrutura modular (`internal/infra/database`, `internal/entity`, `internal/internal_error`);
- Logs estruturados via `zap` (`configuration/logger`).

---

## 🧩 Estrutura do Projeto


├── auction
├── cmd
│   └── auction
│       └── main.go
├── configuration
│   ├── database
│   │   └── mongodb
│   │       └── connection.go
│   ├── logger
│   │   └── logger.go
│   └── rest_err
│       └── rest_err.go
├── internal
│   ├── entity
│   │   ├── auction_entity
│   │   │   └── auction_entity.go
│   │   ├── bid_entity
│   │   │   └── bid_entity.go
│   │   └── user_entity
│   │       └── user_entity.go
│   ├── infra
│   │   ├── api
│   │   │   └── web
│   │   │       ├── controller
│   │   │       │   ├── auction_controller
│   │   │       │   │   ├── create_auction_controller.go
│   │   │       │   │   └── find_auction_controller.go
│   │   │       │   ├── bid_controller
│   │   │       │   │   ├── create_bid_controller.go
│   │   │       │   │   └── find_bid_controller.go
│   │   │       │   └── user_controller
│   │   │       │       └── find_user_controller.go
│   │   │       └── validation
│   │   │           └── validation.go
│   │   └── database
│   │       ├── auction
│   │       │   ├── create_auction.go
│   │       │   ├── create_auction_test.go
│   │       │   └── find_auction.go
│   │       ├── bid
│   │       │   ├── create_bid.go
│   │       │   └── find_bid.go
│   │       └── user
│   │           └── find_user.go
│   ├── internal_error
│   │   └── internal_error.go
│   └── usecase
│       ├── auction_usecase
│       │   ├── create_auction_usecase.go
│       │   └── find_auction_usecase.go
│       ├── bid_usecase
│       │   ├── create_bid_usecase.go
│       │   └── find_bid_usecase.go
│       └── user_usecase
│           └── find_user_usecase.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md

## 🚀 Execução
```bash
make up
```

Acesse: http://localhost:8080

## 🧪 Testes
```bash
go test ./internal/infra/database/auction -v
```

## ⚙️ Variáveis (.env)
```
MONGO_URI=mongodb://mongodb:27017
AUCTION_DURATION_SECONDS=10
APP_PORT=8080
```
