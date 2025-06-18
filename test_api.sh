#!/bin/bash

# URL base da API
BASE_URL="http://localhost:8080"

# Função para mostrar detalhes da requisição e resposta
test_endpoint() {
    local method=$1
    local path=$2
    local name=$3
    local payload=$4
    
    echo "=== $name ==="
    echo "URL: $BASE_URL$path"
    echo "Method: $method"
    
    if [ -n "$payload" ]; then
        echo "Payload: $payload"
        response=$(curl -s -i -X $method "$BASE_URL$path" -H "Content-Type: application/json" -d "$payload")
    else
        response=$(curl -s -i -X $method "$BASE_URL$path")
    fi
    
    echo "Response:"
    echo "$response"
    echo ""
    echo "----------------------------------------"
    echo ""
}

# 1. Testar endpoint de saúde
test_endpoint "GET" "/health" "Testando /health"

# 2. Testar listagem de leilões (deve retornar vazio)
test_endpoint "GET" "/auction" "Listar leilões"

# 3. Criar um novo leilão
AUCTION_PAYLOAD='{
    "product_name": "iPhone 13 Pro",
    "category": "Eletrônicos",
    "description": "Novo na caixa, selado",
    "condition": "new"
}'
test_endpoint "POST" "/auction" "Criar leilão" "$AUCTION_PAYLOAD"

# Extrair o ID do leilão da resposta se criado com sucesso
if [[ "$response" == *"201 Created"* ]]; then
    AUCTION_ID=$(echo "$response" | grep -Eo '"Id":"[^"]+"' | cut -d'"' -f4)
    echo "Auction ID: $AUCTION_ID"
    echo ""
else
    echo "!!! ERRO: Não foi possível criar o leilão"
    exit 1
fi

# 4. Buscar leilão por ID
test_endpoint "GET" "/auction/$AUCTION_ID" "Buscar leilão por ID"

# 5. Criar usuário manualmente no banco
USER_ID=$(uuidgen)
echo "=== Criando usuário manualmente no MongoDB ==="
echo "User ID: $USER_ID"
docker compose exec mongodb mongosh -u admin -p admin --eval "db.users.insertOne({ _id: '$USER_ID', name: 'John Doe' })" auctions
echo "----------------------------------------"
echo ""

# 6. Fazer um lance
BID_PAYLOAD='{
    "user_id": "'$USER_ID'",
    "auction_id": "'$AUCTION_ID'",
    "amount": 3500.00
}'
test_endpoint "POST" "/bid" "Fazer lance" "$BID_PAYLOAD"

# 7. Buscar lances por leilão
test_endpoint "GET" "/bid/$AUCTION_ID" "Buscar lances por leilão"

# 8. Buscar usuário por ID
test_endpoint "GET" "/user/$USER_ID" "Buscar usuário por ID"

# 9. Listar leilões ativos
test_endpoint "GET" "/auction?status=0" "Listar leilões ativos"

# 10. Aguardar fechamento automático
echo "=== Aguardando fechamento do leilão ==="
echo "Aguardando 45 segundos para fechamento automático..."
sleep 45
echo ""

# 11. Buscar vencedor do leilão
test_endpoint "GET" "/auction/winner/$AUCTION_ID" "Buscar vencedor do leilão"

# 12. Buscar leilão após fechamento
test_endpoint "GET" "/auction/$AUCTION_ID" "Buscar leilão após fechamento"

echo "=== Testes completos ==="