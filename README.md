# Auction System

## Descrição

Este projeto é um sistema de leilão desenvolvido em Go. Ele permite a criação de leilões, lances e fechamento automático dos leilões após um tempo definido.

## Requisitos

- Go 1.20+
- Docker
- Docker Compose

## Configuração do Ambiente

1. Clone o repositório:
```sh
   git clone <URL_DO_REPOSITORIO>
   cd auction-system
```

2. Crie um arquivo `.env` na pasta `cmd/auction` com o seguinte conteudo:
```sh
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_DURATION=20s

MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

### Executar a Aplicacao

## Usando Docker Compose

1. Construa e inicie os servicos

```sh
docker-compose up --build
```

2. A aplicacao estara disponivel em `http://localhost:8080`.

3. Rodar testes.

```sh
 go test ./...
```

## Endpoints:
* POST /auction: Cria um novo leilão.
* GET /auction: Lista todos os leilões.
* GET /auction/:auctionId: Obtém detalhes de um leilão específico.
* GET /auction/winner/:auctionId: Obtém o lance vencedor de um leilão específico.
* POST /bid: Cria um novo lance.
* GET /bid/:auctionId: Lista todos os lances de um leilão específico.
* GET /user/:userId: Obtém detalhes de um usuário específico.