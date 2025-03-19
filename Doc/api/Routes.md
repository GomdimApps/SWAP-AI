# API SWAP-AI

Este documento fornece instruções para testar a API SWAP-AI usando curl.  
Os tokens retornados são JWT assinados com a chave definida em SECRET_KEY e devem ser enviados no cabeçalho da requisição como:

    Authorization: Bearer <token>

## Pré-requisitos

1. O servidor da API deve estar rodando localmente na porta 8080.
2. O banco de dados PostgreSQL deve estar configurado com os esquemas necessários.

## Testando com curl

### Verificação de Saúde

```bash
curl -X GET http://192.168.10.3:8080/health
```

Resposta esperada:
```json
{"status": "ok"}
```

# Rotas de Administração

### Registro de Usuário

```bash
curl -X POST http://192.168.10.3:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{"email": "new_user@example.com", "password": "password123"}'
```

Resposta esperada:
```json
{
    "token": "Bearer <jwt_token>",
    "userId": "<uuid>",
    "email": "new_user@example.com",
    "expiresAt": 1678909876
}
```

### Login de Usuário (Admin)

Utilize o mesmo endpoint de login para acessar como usuário (admin):

```bash
curl -X POST http://192.168.10.3:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email": "new_user@example.com", "password": "password123"}'
```

Resposta esperada:
```json
{
    "token": "Bearer <jwt_token>",
    "userId": "<uuid>",
    "email": "new_user@example.com",
    "expiresAt": 1678909876
}
```

### Criar Cliente

Envie o token JWT de um usuário admin no cabeçalho para criar um cliente:

```bash
curl -X POST http://192.168.10.3:8080/api/v1/clients \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <jwt_token>" \
    -d '{"name": "Test Client", "email": "client@example.com", "whatsapp": "+551199999999", "password": "clientpassword"}'
```

Resposta esperada:
```json
{
    "ID": "<uuid>",
    "FolderID": "<uuid>",
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
    -H "Authorization: Bearer <jwt_token>"
```

Resposta esperada:
```json
{
    "clients": [
        {
            "ID": "<uuid>",
            "FolderID": "<uuid>",
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
    -H "Authorization: Bearer <jwt_token>"
```

Resposta esperada:
```json
{
    "ID": "<uuid>",
    "FolderID": "<uuid>",
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
    -H "Authorization: Bearer <jwt_token>" \
    -d '{"name": "Updated Client", "email": "updated@example.com", "whatsapp": "+551188888888"}'
```

Resposta esperada:
```json
{
    "ID": "<uuid>",
    "FolderID": "<uuid>",
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
    -H "Authorization: Bearer <jwt_token>"
```

Resposta esperada:
```json
{
    "message": "Client deleted successfully"
}
```

---

# Rotas de Cliente

### Login de Cliente

Utilize o endpoint específico para login de cliente:

```bash
curl -X POST http://192.168.10.3:8080/api/v1/auth/client/login \
    -H "Content-Type: application/json" \
    -d '{"email": "client@example.com", "password": "clientpassword"}'
```

Resposta esperada:
```json
{
    "token": "Bearer <jwt_token>",
    "clientId": "<uuid>",
    "email": "client@example.com",
    "expiresAt": 1678909876
}
```
