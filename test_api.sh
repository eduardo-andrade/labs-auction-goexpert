#!/bin/bash

# URL base da API
BASE_URL="http://localhost:8080"

# Função para fazer requests e mostrar a resposta completa
function request {
    echo "=== $3 ==="
    echo "Request: $1 $2"
    if [ -n "$4" ]; then
        echo "Payload: $4"
    fi
    echo "Response:"
    if [ -z "$4" ]; then
        curl -s -i -X $1 "$BASE_URL$2"
    else
        curl -s -i -X $1 "$BASE_URL$2" -H "Content-Type: application/json" -d "$4"
    fi
    echo ""
    echo "----------------------------------------"
}

# 1. Testar endpoint de saúde
request GET "/health" "Testando /health"

# 2. Criar um novo leilão
AUCTION_PAYLOAD='{
    "product_name": "iPhone 13 Pro",
    "category": "Eletrônicos",
    "description": "Novo na caixa, selado",
    "condition": "new"
}'
request POST "/auction" "Criando leilão" "$AUCTION_PAYLOAD"

# Extrair o ID do leilão da resposta
AUCTION_ID=$(curl -s -X POST "$BASE_URL/auction" -H "Content-Type: application/json" -d "$AUCTION_PAYLOAD" | jq -r '.id')
if [ -z "$AUCTION_ID" ] || [ "$AUCTION_ID" == "null" ]; then
    echo "!!! ERRO: Não foi possível obter o ID do leilão"
    exit 1
fi
echo "Auction ID: $AUCTION_ID"
echo ""

# 3. Buscar leilão por ID
request GET "/auction/$AUCTION_ID" "Buscando leilão por ID"

# 4. Criar usuário manualmente no banco
USER_ID=$(uuidgen)
echo "=== Criando usuário manualmente no MongoDB ==="
echo "User ID: $USER_ID"
docker compose exec mongodb mongosh -u admin -p admin --eval "db.users.insertOne({ _id: '$USER_ID', name: 'John Doe' })" auctions
echo "----------------------------------------"
echo ""

# 5. Fazer um lance
BID_PAYLOAD='{
    "user_id": "'$USER_ID'",
    "auction_id": "'$AUCTION_ID'",
    "amount": 3500.00
}'
request POST "/bid" "Fazendo lance" "$BID_PAYLOAD"

# 6. Buscar lances por leilão
request GET "/bid/$AUCTION_ID" "Buscando lances por leilão"

# 7. Buscar usuário por ID
request GET "/user/$USER_ID" "Buscando usuário por ID"

# 8. Listar todos os leilões ativos
request GET "/auction?status=0" "Listando leilões ativos"

# 9. Aguardar fechamento automático
echo "=== Aguardando fechamento do leilão ==="
echo "Aguardando 35 segundos para fechamento automático..."
sleep 35
echo ""

# 10. Buscar vencedor do leilão
request GET "/auction/winner/$AUCTION_ID" "Buscando vencedor do leilão"

# 11. Buscar leilão após fechamento
request GET "/auction/$AUCTION_ID" "Buscando leilão após fechamento"

echo "=== Testes completos ==="