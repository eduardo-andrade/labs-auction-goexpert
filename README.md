# Projeto de Sistema de Leilões em Go (GoExpert)

Este projeto é um sistema de leilões completo desenvolvido em Go como parte do curso **GoExpert**. Ele inclui funcionalidades para criar leilões, fazer lances, listar leilões ativos, buscar informações de leilões e usuários, e determinar automaticamente o vencedor de um leilão após seu término.

## 🛠 Melhorias Implementadas

- **Separação de Use Cases**: Divisão clara entre operações de criação e consulta para melhor organização do código.
- **Validação de Entrada**: Implementação robusta de validação de dados nos controllers.
- **Tratamento de Erros**: Padronização de erros com mensagens claras e códigos HTTP apropriados.
- **Formatação de Timestamps**: Consistência no formato de datas em toda a API.
- **Operações de Fechamento Automático**: Sistema que fecha automaticamente leilões expirados.
- **Documentação Aprimorada**: README completo com instruções de execução e solução de problemas.
- **Testes Automatizados**: Script de teste abrangente (`test_api.sh`) para validar todos os endpoints.
- **Segurança de Dados**: Proteção contra injeção de SQL e validação de UUIDs.
- **Configuração Flexível**: Personalização da duração do leilão via variável de ambiente.
- **Arquitetura Limpa**: Separação clara entre camadas (`entity`, `usecase`, `controller`, `repository`).

## 📁 Estrutura de Diretórios

```
.
├── cmd
│   └── auction
│       └── main.go
├── configuration
│   ├── database
│   │   └── mongodb
│   ├── logger
│   └── rest_err
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── entity
│   │   ├── auction_entity
│   │   ├── bid_entity
│   │   └── user_entity
│   ├── infra
│   │   ├── api
│   │   │   └── web
│   │   │       └── controller
│   │   ├── database
│   │   │   ├── auction
│   │   │   ├── bid
│   │   │   └── user
│   ├── internal_error
│   └── usecase
│       ├── auction_usecase
│       ├── bid_usecase
│       └── user_usecase
├── LICENSE
├── mongo-init.js
└── test_api.sh
```

## ✅ Pré-requisitos

- Docker (versão 20.10+)
- Docker Compose (versão 1.29+)
- Git (para clonar o repositório)

## ▶️ Como Rodar o Projeto

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/labs-auction-goexpert.git
cd labs-auction-goexpert
```

### 2. Construa e inicie os containers

```bash
docker compose up -d --build
docker compose up
```

Este comando irá:

- Construir a imagem da aplicação Go.
- Iniciar um container MongoDB.
- Inicializar o banco de dados com um script de configuração.
- Iniciar a aplicação na porta `8080`.

### 3. Verifique se os serviços estão rodando

```bash
docker compose ps
```

Você deve ver dois serviços: `app` e `mongodb`.

## 🧪 Como Testar o Projeto

Um script de teste (`test_api.sh`) está disponível para testar todos os endpoints da API:

### 1. Dê permissão de execução ao script

```bash
chmod +x test_api.sh
```

### 2. Execute o script de teste

```bash
./test_api.sh
```

Esse script irá testar:

- Endpoint de saúde (`/health`)
- Listagem e criação de leilões
- Criação de usuário (manual)
- Lance em leilões
- Buscar vencedor após fechamento automático

## 🌐 Endpoints da API

| Método | Endpoint                     | Descrição                            |
|--------|------------------------------|--------------------------------------|
| GET    | `/health`                    | Verifica a saúde da aplicação        |
| GET    | `/auction`                   | Lista leilões                        |
| GET    | `/auction/:auctionId`        | Busca leilão por ID                  |
| POST   | `/auction`                   | Cria um novo leilão                  |
| GET    | `/auction/winner/:auctionId`| Busca vencedor de um leilão          |
| POST   | `/bid`                       | Cria um novo lance                   |
| GET    | `/bid/:auctionId`           | Busca lances por leilão              |
| GET    | `/user/:userId`             | Busca usuário por ID                 |

## ⏱️ Configuração do Tempo do Leilão

O tempo de duração de um leilão pode ser definido pela variável de ambiente `AUCTION_DURATION`.

1. **Edite o arquivo `.env`:**

Formatos válidos para a variável`AUCTION_DURATION`:

- `30s`: 30 segundos (ideal para testes)
- `5m`: 5 minutos (configuração intermediária)
- `24h`: 24 horas (modo produção)

2. **Reconstrua os containers:**

```bash
docker compose down --volumes
docker compose up -d --build
```

---

## 📡 Base URL

Todos os exemplos abaixo utilizam a base URL:

```
http://localhost:8080
```

---

## 📚 Endpoints da API

### 1. Health Check

Verifica se a API está ativa.

- **GET** `/health`

#### Exemplo curl:

```bash
curl -X GET http://localhost:8080/health
```

#### Resposta:

```json
{ "status": "ok" }
```

---

### 2. Criar Leilão

- **POST** `/auction`

#### Payload:

```json
{
  "product_name": "iPhone 13 Pro",
  "category": "Eletrônicos",
  "description": "Novo na caixa, selado",
  "condition": "new"
}
```

Condições válidas: `new`, `used`, `refurbished`

#### Exemplo curl:

```bash
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "iPhone 13 Pro",
    "category": "Eletrônicos",
    "description": "Novo na caixa, selado",
    "condition": "new"
  }'
```

#### Resposta:

```json
{
  "Id": "...",
  "ProductName": "...",
  "Category": "...",
  "Description": "...",
  "Condition": 1,
  "Status": 0,
  "Timestamp": "..."
}
```

---

### 3. Listar Leilões

- **GET** `/auction`

Parâmetros opcionais:
- `status`: 0 (ativos), 1 (finalizados)
- `category`
- `productName` (busca parcial)

#### Exemplo curl:

```bash
curl "http://localhost:8080/auction?status=0&category=Eletrônicos&productName=iPhone"
```

#### Resposta:

```json
[
  {
    "Id": "...",
    "ProductName": "...",
    "Category": "...",
    "Description": "...",
    "Condition": 1,
    "Status": 0,
    "Timestamp": "..."
  }
]
```

---

### 4. Buscar Leilão por ID

- **GET** `/auction/:auctionId`

#### Exemplo curl:

```bash
curl http://localhost:8080/auction/acde3b18-3328-4c00-966d-9571e604640b
```

---

### 5. Buscar Vencedor do Leilão

- **GET** `/auction/winner/:auctionId`

#### Exemplo curl:

```bash
curl http://localhost:8080/auction/winner/acde3b18-3328-4c00-966d-9571e604640b
```

#### Resposta:

```json
{
  "Auction": { ... },
  "Bid": { ... }
}
```

---

### 6. Criar Lance

- **POST** `/bid`

#### Payload:

```json
{
  "user_id": "uuid-do-usuario",
  "auction_id": "uuid-do-leilao",
  "amount": 3500.00
}
```

#### Exemplo curl:

```bash
curl -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "e2042afe-9664-4967-8132-1a26430c6219",
    "auction_id": "acde3b18-3328-4c00-966d-9571e604640b",
    "amount": 3500.00
  }'
```

---

### 7. Buscar Lances por Leilão

- **GET** `/bid/:auctionId`

#### Exemplo curl:

```bash
curl http://localhost:8080/bid/acde3b18-3328-4c00-966d-9571e604640b
```

---

### 8. Buscar Usuário por ID

- **GET** `/user/:userId`

#### Exemplo curl:

```bash
curl http://localhost:8080/user/e2042afe-9664-4967-8132-1a26430c6219
```

---

## 🛠️ Solução de Problemas

### Leilão não fecha automaticamente?

- Verifique logs:

```bash
docker compose logs app | grep "auction_closer"
```

- Confirme a variável:

```bash
docker compose exec app env | grep AUCTION_DURATION
```

- Verifique no MongoDB:

```bash
docker compose exec mongodb mongosh -u admin -p admin \
  --eval "db.auctions.find({}, {end_time:1, status:1, product_name:1})" auctions
```

---

### MongoDB não inicializa?

```bash
docker compose logs mongodb
docker compose exec mongodb mongosh -u admin -p admin \
  --eval "$(cat mongo-init.js)" auctions
```

---

### Erro ao construir imagem Docker

```bash
ERROR: failed to solve: process "/bin/sh -c ..." did not complete successfully
```

Solução:

```bash
docker compose down --volumes
docker compose build --no-cache
docker compose up -d
```

### Conexão recusada na porta 8080

```bash
curl: (7) Failed to connect to localhost port 8080
```

Solução:

```bash
docker compose ps
docker compose logs app
```

### Erro ao conectar ao MongoDB

```text
Failed to connect to MongoDB: ...
```

Solução:

```bash
docker compose ps
docker compose logs mongodb
```

### Erros 400/404 durante testes

Verifique:

```bash
docker compose logs -f app
```

## 🧾 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.