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

Envie uma requisição POST para registrar um novo usuário (admin), informando também o período e o valor para a expiração do usuário.  
Os campos "expirationPeriod" e "expirationValue" são obrigatórios.  
Valores permitidos para "expirationPeriod": "days", "weeks", "months" ou "years".  

Exemplo de requisição:
```bash
curl -X POST http://192.168.10.3:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{
          "email": "new_user@example.com",
          "password": "password123",
          "expirationPeriod": "days",
          "expirationValue": 30
        }'
```

Resposta esperada:
```json
{
    "token": "<jwt_token>",
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

### Listar Usuários de Administração

Utilize o token JWT de um usuário admin para listar os usuários de administração.  
A resposta inclui, para cada usuário, um array "clientes" com os UUIDs dos clientes que o admin criou.

```bash
curl -X GET http://192.168.10.3:8080/api/v1/admin/users \
    -H "Authorization: Bearer <jwt_token>"
```

Resposta esperada:
```json
{
    "users": [
        {
            "ID": "<uuid>",
            "Email": "admin@example.com",
            "PasswordHash": "<hash>",
            "CreatedAt": "<timestamp>",
            "UpdatedAt": "<timestamp>",
            "clientes": ["<client_uuid1>", "<client_uuid2>"]
        }
    ]
}
```

### Atualizar Status do Usuário Admin

Utilize o token JWT de um usuário admin para atualizar o status (ativo ou inativo) de um usuário admin.  
Envie uma requisição PUT com o novo status no corpo em formato JSON.

```bash
curl -X PUT http://192.168.10.3:8080/api/v1/admin/users/<user_uuid>/status \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <jwt_token>" \
    -d '{"status": "ativo"}'
```

Resposta esperada:
```json
{
    "message": "User status updated successfully",
    "user": {
        "ID": "<uuid>",
        "Email": "admin@example.com",
        "PasswordHash": "<hash>",
        "CreatedAt": "<timestamp>",
        "UpdatedAt": "<timestamp>",
        "ExpiresAt": "<timestamp>",
        "Status": "ativo"
    }
}
```

### Atualizar Informações do Usuário Admin

Utilize o token JWT de um usuário admin para atualizar o email e/ou senha do usuário admin.  
Envie uma requisição PUT com o novo email e/ou a nova senha no corpo em formato JSON.

```bash
curl -X PUT http://192.168.10.3:8080/api/v1/admin/users/<user_uuid>/info \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <jwt_token>" \
    -d '{"email": "new_email@example.com", "password": "newpassword123"}'
```

Resposta esperada:
```json
{
    "message": "Admin info updated successfully",
    "user": {
        "ID": "<uuid>",
        "Email": "new_email@example.com",
        "PasswordHash": "<new_hash>",
        "CreatedAt": "<timestamp>",
        "UpdatedAt": "<timestamp>",
        "ExpiresAt": "<timestamp>",
        "Status": "ativo"
    }
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

### Listar Project IDs de Arquivos Pendentes do Cliente

Envie o token JWT de um usuário autenticado para obter os project IDs e seus arquivos relacionados que não foram processados:

```bash
curl -X GET http://192.168.10.3:8080/api/v1/clients/<client_uuid>/pending-files \
    -H "Authorization: Bearer <jwt_token>"
```

Resposta esperada:
```json
{
    "pending_projects": [
        {
            "project_id": "<project_uuid>",
            "files": [
                {
                    "ID": "<uuid>",
                    "ClientID": "<uuid>",
                    "FolderID": "<uuid>",
                    "OriginalName": "doc1.png",
                    "FilePath": "<path_to_file>",
                    "FileSize": 104857600,
                    "FileType": "text/plain; charset=utf-8",
                    "ProcessedBy": "00000000-0000-0000-0000-000000000000",
                    "IsProcessed": false,
                    "CreatedAt": "<timestamp>",
                    "UpdatedAt": "<timestamp>",
                    "download_url": "http://192.168.10.3:8080/download/unprocessed/<clientid>/<projectid>/<filename>"
                }
            ]
        }
    ]
}
```

### Listar Project IDs de Arquivos Processados do Cliente

Envie o token JWT de um usuário autenticado para obter os project IDs e seus arquivos relacionados que foram processados:

```bash
curl -X GET http://192.168.10.3:8080/api/v1/clients/<client_uuid>/processed-files \
    -H "Authorization: Bearer <jwt_token>"
```

Resposta esperada:
```json
{
    "processed_projects": [
        {
            "project_id": "<project_uuid>",
            "files": [
                {
                    "ID": "<uuid>",
                    "ClientID": "<uuid>",
                    "FolderID": "<uuid>",
                    "OriginalName": "doc1.png",
                    "FilePath": "<path_to_file>",
                    "FileSize": 104857600,
                    "FileType": "text/plain; charset=utf-8",
                    "ProcessedBy": "00000000-0000-0000-0000-000000000000",
                    "IsProcessed": true,
                    "CreatedAt": "<timestamp>",
                    "UpdatedAt": "<timestamp>",
                    "download_url": "http://192.168.10.3:8080/download/processed/<clientid>/<projectid>/<filename>"
                }
            ]
        }
    ]
}
```

## Rotas de Download de Arquivos

### Arquivos Não Processados
Endpoint para baixar arquivos não processados:  
Retorna, além dos metadados, um link de download para acessar o arquivo, utilizando a rota:
```bash
curl -X GET http://192.168.10.3:8080/download/unprocessed/<clientid>/<projectid>/<filename>
```

### Arquivos Processados
Endpoint para baixar arquivos processados:  
Retorna, além dos metadados, um link de download para acessar o arquivo, utilizando a rota:
```bash
curl -X GET http://192.168.10.3:8080/download/processed/<clientid>/<projectid>/<filename>
```
