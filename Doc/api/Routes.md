# Testando a API SWAP-AI

Este documento fornece instruções para testar a API SWAP-AI usando curl e Postman.

## Pré-requisitos

1. O servidor da API está rodando localmente na porta 8080
2. O banco de dados PostgreSQL está configurado com os esquemas necessários

## Testando com curl

### Verificação de Saúde

```bash
curl -X GET http://192.168.10.3:8080/health
```

Resposta esperada:
```json
{"status":"ok"}
```

### Registro de Usuário

```bash
curl -X POST http://192.168.10.3:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{"email":"new_user@example.com","password":"password123"}'
```

Resposta esperada:
```json
{
    "token": "Bearer <uuid>",
    "userId": "<uuid>",
    "email": "new_user@example.com",
    "expiresAt": 1678909876
}
```

### Login de Usuário

```bash
curl -X POST http://192.168.10.3:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"new_user@example.com","password":"password123"}'
```

Resposta esperada:
```json
{
    "token": "Bearer <uuid>",
    "userId": "<uuid>",
    "email": "new_user@example.com",
    "expiresAt": 1678909876
}
```

### Criar Cliente

```bash
curl -X POST http://192.168.10.3:8080/api/v1/clients \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <token_uuid>" \
    -d '{"name":"Test Client","email":"client@example.com","whatsapp":"+551199999999"}'
```

Resposta esperada:
```json
{
    "ID": "<client_uuid>",
    "FolderID": "<folder_uuid>",
    "Name": "Test Client",
    "Status": "ativo",
    "Email": "client@example.com",
    "WhatsApp": "+551199999999",
    "CreatedAt": "2023-01-01T12:00:00Z",
    "UpdatedAt": "2023-01-01T12:00:00Z"
}
```

### Obter Clientes

```bash
curl -X GET http://192.168.10.3:8080/api/v1/clients \
    -H "Authorization: Bearer <token_uuid>"
```

Resposta esperada:
```json
{
    "clients": [
        {
            "ID": "<client_uuid>",
            "FolderID": "<folder_uuid>",
            "Name": "Test Client",
            "Status": "ativo",
            "Email": "client@example.com",
            "WhatsApp": "+551199999999",
            "CreatedAt": "2023-01-01T12:00:00Z",
            "UpdatedAt": "2023-01-01T12:00:00Z"
        }
    ]
}
```

### Obter Cliente por ID

```bash
curl -X GET http://192.168.10.3:8080/api/v1/clients/<client_uuid> \
    -H "Authorization: Bearer <token_uuid>"
```

Resposta esperada:
```json
{
    "ID": "<client_uuid>",
    "FolderID": "<folder_uuid>",
    "Name": "Test Client",
    "Status": "ativo",
    "Email": "client@example.com",
    "WhatsApp": "+551199999999",
    "CreatedAt": "2023-01-01T12:00:00Z",
    "UpdatedAt": "2023-01-01T12:00:00Z"
}
```

### Atualizar Cliente

```bash
curl -X PUT http://192.168.10.3:8080/api/v1/clients/<client_uuid> \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <token_uuid>" \
    -d '{"name":"Updated Client","email":"updated@example.com","whatsapp":"+551188888888"}'
```

Resposta esperada:
```json
{
    "ID": "<client_uuid>",
    "FolderID": "<folder_uuid>",
    "Name": "Updated Client",
    "Status": "ativo",
    "Email": "updated@example.com",
    "WhatsApp": "+551188888888",
    "CreatedAt": "2023-01-01T12:00:00Z",
    "UpdatedAt": "2023-01-01T12:30:00Z"
}
```

### Deletar Cliente

```bash
curl -X DELETE http://192.168.10.3:8080/api/v1/clients/<client_uuid> \
    -H "Authorization: Bearer <token_uuid>"
```

Resposta esperada:
```json
{
    "message": "Client deleted successfully"
}
```

## Testando com Postman

1. Importe a coleção do Postman do arquivo `testing/swap_ai_api.postman_collection.json`
2. Configure as variáveis de ambiente:
     - `base_url` = `http://192.168.10.3:8080`
     - `auth_token` (será preenchido após registro/login - guarde apenas o UUID, sem o prefixo "Bearer")

3. Execute as requisições na seguinte ordem:
     - Verificação de Saúde
     - Registrar Usuário (ou Login se já registrado)
     - Criar Cliente
     - Obter Todos os Clientes
     - Obter Cliente por ID
     - Atualizar Cliente
     - Deletar Cliente

## Solução de Problemas

Se você receber um erro 409 ao registrar um usuário, isso indica que o email já está em uso. Tente fazer login com esse email ou use um email diferente para registrar um novo usuário.

