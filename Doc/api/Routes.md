# API SWAP-AI

Este documento fornece instruções para testar a API SWAP-AI usando curl

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


