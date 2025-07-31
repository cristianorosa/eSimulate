# FileSystem API

API REST para operações de sistema de arquivos desenvolvida em Go seguindo Clean Architecture.

## 🚀 Funcionalidades

- ✅ Operações de arquivo (ler, criar, editar, deletar)
- ✅ Operações de diretório (ler, criar, renomear, deletar)
- ✅ Execução segura de comandos do sistema
- ✅ Validação e sanitização de caminhos
- ✅ Controle de acesso por diretório base
- ✅ Middleware de CORS e logging

## 🏗️ Arquitetura

```
filesystem-api/
├── domain/                 # Entidades e interfaces
│   ├── file.go
│   └── filesystem_repository.go
├── usecase/               # Lógica de negócio
│   └── filesystem_usecase.go
├── infra/                 # Implementação de infraestrutura
│   └── filesystem_repository_os.go
├── interfaces/            # Handlers HTTP e rotas
│   ├── filesystem_handler.go
│   └── routes.go
├── main.go               # Aplicação principal
└── README.md
```

## 🔧 Instalação e Execução

```bash
# Navegar para o diretório
cd filesystem-api

# Executar a aplicação
go run main.go

# Ou compilar e executar
go build -o filesystem-api
./filesystem-api
```

## 📡 Endpoints

### Operações de Arquivo

#### 1. Ler Arquivo
```http
GET /read-file/{path}
```

**Exemplo:**
```bash
curl http://localhost:8081/read-file/exemplo.txt
```

**Resposta:**
```json
{
  "path": "exemplo.txt",
  "content": "Conteúdo do arquivo...",
  "size": 1024
}
```

#### 2. Criar Arquivo
```http
PUT /create-file/{path}
Content-Type: application/json

{
  "content": "Conteúdo do novo arquivo"
}
```

**Exemplo:**
```bash
curl -X PUT http://localhost:8081/create-file/novo-arquivo.txt \
  -H "Content-Type: application/json" \
  -d '{"content": "Olá, mundo!"}'
```

#### 3. Editar Arquivo
```http
PATCH /edit-file/{path}
Content-Type: application/json

{
  "content": "Novo conteúdo do arquivo"
}
```

**Exemplo:**
```bash
curl -X PATCH http://localhost:8081/edit-file/arquivo.txt \
  -H "Content-Type: application/json" \
  -d '{"content": "Conteúdo atualizado"}'
```

#### 4. Deletar Arquivo
```http
DELETE /delete-file/{path}
```

**Exemplo:**
```bash
curl -X DELETE http://localhost:8081/delete-file/arquivo-para-deletar.txt
```

### Operações de Diretório

#### 5. Ler Diretório
```http
GET /read-dir/{path}
```

**Exemplo:**
```bash
curl http://localhost:8081/read-dir/minha-pasta
```

**Resposta:**
```json
{
  "path": "minha-pasta",
  "files": [
    {
      "name": "arquivo1.txt",
      "path": "minha-pasta/arquivo1.txt",
      "size": 1024,
      "is_dir": false,
      "mod_time": "2024-01-15T10:30:00Z",
      "mode": "-rw-r--r--"
    }
  ],
  "total": 1
}
```

#### 6. Criar Diretório
```http
PUT /create-dir/{path}
```

**Exemplo:**
```bash
curl -X PUT http://localhost:8081/create-dir/nova-pasta
```

#### 7. Renomear Diretório
```http
PATCH /edit-dir/{path}
Content-Type: application/json

{
  "new_path": "novo-nome-da-pasta"
}
```

**Exemplo:**
```bash
curl -X PATCH http://localhost:8081/edit-dir/pasta-antiga \
  -H "Content-Type: application/json" \
  -d '{"new_path": "pasta-nova"}'
```

#### 8. Deletar Diretório
```http
DELETE /delete-dir/{path}
```

**Exemplo:**
```bash
curl -X DELETE http://localhost:8081/delete-dir/pasta-para-deletar
```

### Execução de Comandos

#### 9. Executar Comando
```http
POST /execute-command/{command}
```

**Exemplo:**
```bash
# Comando simples
curl -X POST http://localhost:8081/execute-command/ls

# Comando com parâmetros (URL encoded)
curl -X POST http://localhost:8081/execute-command/ls%20-la

# Comando pwd
curl -X POST http://localhost:8081/execute-command/pwd
```

**Resposta:**
```json
{
  "command": "ls -la",
  "output": "total 8\ndrwxr-xr-x 2 user user 4096 Jan 15 10:30 .\ndrwxr-xr-x 3 user user 4096 Jan 15 10:29 ..",
  "exit_code": 0
}
```

### Status da API

#### Status/Health Check
```http
GET /status
```

**Exemplo:**
```bash
curl http://localhost:8081/status
```

## 🔒 Segurança

### Controle de Acesso
- Todas as operações são limitadas ao diretório base configurado
- Validação de caminhos para prevenir directory traversal
- Sanitização de entrada para evitar caracteres perigosos

### Comandos Permitidos
Por segurança, apenas os seguintes comandos são permitidos:
- `ls` - Listar arquivos
- `pwd` - Diretório atual
- `echo` - Imprimir texto
- `cat` - Mostrar conteúdo de arquivo
- `head` - Primeiras linhas de arquivo
- `tail` - Últimas linhas de arquivo
- `grep` - Buscar texto
- `find` - Encontrar arquivos

### Comandos Bloqueados
Comandos perigosos são automaticamente bloqueados:
- `rm`, `rmdir`, `del` - Remoção
- `format`, `fdisk`, `mkfs` - Formatação
- `sudo`, `su` - Elevação de privilégios
- `chmod`, `chown` - Alteração de permissões

## 🌍 Variáveis de Ambiente

- `FILESYSTEM_API_BASE_DIR`: Diretório base para operações (padrão: `./sandbox`)
- `PORT`: Porta do servidor (padrão: `8081`)

## 🧪 Testando a API

### Criar estrutura de teste
```bash
# Criar arquivo de teste
curl -X PUT http://localhost:8081/create-file/teste.txt \
  -H "Content-Type: application/json" \
  -d '{"content": "Arquivo de teste"}'

# Criar diretório de teste
curl -X PUT http://localhost:8081/create-dir/pasta-teste

# Listar diretório atual
curl http://localhost:8081/read-dir/

# Executar comando
curl -X POST http://localhost:8081/execute-command/ls%20-la
```

## 🚨 Tratamento de Erros

A API retorna códigos HTTP apropriados:

- `200 OK` - Operação bem-sucedida
- `201 Created` - Recurso criado
- `400 Bad Request` - Dados inválidos
- `404 Not Found` - Recurso não encontrado
- `405 Method Not Allowed` - Método HTTP não permitido
- `500 Internal Server Error` - Erro interno

Exemplo de resposta de erro:
```json
{
  "error": "Arquivo não encontrado"
}
```

## 📝 Logs

A API registra todas as requisições:
```
2024/01/15 10:30:00 GET /read-file/teste.txt 127.0.0.1:54321
2024/01/15 10:30:05 POST /execute-command 127.0.0.1:54322
```