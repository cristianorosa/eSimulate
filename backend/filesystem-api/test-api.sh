#!/bin/bash

# Script de teste para a FileSystem API
# Execute: chmod +x test-api.sh && ./test-api.sh

API_URL="http://localhost:8081"

echo "🧪 Testando FileSystem API..."
echo "================================"

# Função para fazer requisições e mostrar resultado
test_endpoint() {
    echo -e "\n📡 $1"
    echo "Comando: $2"
    echo "Resposta:"
    eval $2
    echo -e "\n---"
}

# 1. Verificar status da API
test_endpoint "Status da API" "curl -s $API_URL/status | jq ."

# 2. Listar diretório atual
test_endpoint "Listar diretório atual" "curl -s $API_URL/read-dir/ | jq ."

# 3. Criar um arquivo de teste
test_endpoint "Criar arquivo de teste" "curl -s -X PUT $API_URL/create-file/teste.txt -H 'Content-Type: application/json' -d '{\"content\": \"Olá, mundo! Este é um arquivo de teste.\"}' | jq ."

# 4. Ler o arquivo criado
test_endpoint "Ler arquivo criado" "curl -s $API_URL/read-file/teste.txt | jq ."

# 5. Editar o arquivo
test_endpoint "Editar arquivo" "curl -s -X PATCH $API_URL/edit-file/teste.txt -H 'Content-Type: application/json' -d '{\"content\": \"Conteúdo atualizado do arquivo de teste!\"}' | jq ."

# 6. Ler arquivo editado
test_endpoint "Ler arquivo editado" "curl -s $API_URL/read-file/teste.txt | jq ."

# 7. Criar diretório
test_endpoint "Criar diretório" "curl -s -X PUT $API_URL/create-dir/pasta-teste | jq ."

# 8. Listar diretório novamente
test_endpoint "Listar diretório após criações" "curl -s $API_URL/read-dir/ | jq ."

# 9. Criar arquivo dentro do diretório
test_endpoint "Criar arquivo no diretório" "curl -s -X PUT $API_URL/create-file/pasta-teste/arquivo-interno.txt -H 'Content-Type: application/json' -d '{\"content\": \"Arquivo dentro da pasta\"}' | jq ."

# 10. Listar conteúdo do diretório criado
test_endpoint "Listar conteúdo do diretório" "curl -s $API_URL/read-dir/pasta-teste | jq ."

# 11. Executar comando ls
test_endpoint "Executar comando ls" "curl -s -X POST '$API_URL/execute-command/ls%20-la' | jq ."

# 12. Executar comando pwd
test_endpoint "Executar comando pwd" "curl -s -X POST '$API_URL/execute-command/pwd' | jq ."

# 13. Tentar comando não permitido
test_endpoint "Tentar comando não permitido (rm)" "curl -s -X POST '$API_URL/execute-command/rm%20teste.txt' | jq ."

# 14. Renomear diretório
test_endpoint "Renomear diretório" "curl -s -X PATCH $API_URL/edit-dir/pasta-teste -H 'Content-Type: application/json' -d '{\"new_path\": \"pasta-renomeada\"}' | jq ."

# 15. Listar diretório após renomeação
test_endpoint "Listar após renomeação" "curl -s $API_URL/read-dir/ | jq ."

# 16. Deletar arquivo
test_endpoint "Deletar arquivo" "curl -s -X DELETE $API_URL/delete-file/teste.txt | jq ."

# 17. Deletar diretório
test_endpoint "Deletar diretório" "curl -s -X DELETE $API_URL/delete-dir/pasta-renomeada | jq ."

# 18. Listar diretório final
test_endpoint "Listar diretório final" "curl -s $API_URL/read-dir/ | jq ."

echo -e "\n✅ Testes concluídos!"
echo "Verifique os resultados acima para validar o funcionamento da API."