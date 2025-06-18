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

## 🧯 Possíveis Erros e Soluções

### 1. Erro ao construir imagem Docker

```bash
ERROR: failed to solve: process "/bin/sh -c ..." did not complete successfully
```

Solução:

```bash
docker compose down --volumes
docker compose build --no-cache
docker compose up -d
```

### 2. Conexão recusada na porta 8080

```bash
curl: (7) Failed to connect to localhost port 8080
```

Solução:

```bash
docker compose ps
docker compose logs app
```

### 3. Erro ao conectar ao MongoDB

```text
Failed to connect to MongoDB: ...
```

Solução:

```bash
docker compose ps
docker compose logs mongodb
```

### 4. Erros 400/404 durante testes

Verifique:

```bash
docker compose logs -f app
```

Variáveis de ambiente:

- `MONGODB_URL`
- `AUCTION_DURATION`

### 5. Script falha na criação do leilão

```bash
curl -X POST http://localhost:8080/auction   -H "Content-Type: application/json"   -d '{ "product_name": "iPhone 13 Pro", "category": "Eletrônicos", "description": "Novo na caixa, selado", "condition": "new" }'
```

## ⚙️ Configuração Avançada

### Variáveis de Ambiente (`.env`)

```env
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
AUCTION_DURATION=30s
```

### Exemplos de `AUCTION_DURATION`:

- `30s` - 30 segundos (teste)
- `5m` - 5 minutos
- `24h` - 24 horas (produção)

## 📜 Licença

Este projeto está sob a licença MIT. Consulte o arquivo `LICENSE` para mais detalhes.