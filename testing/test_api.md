# Testing the SWAP-AI API

This document provides instructions for testing the SWAP-AI API using curl and Postman.

## Prerequisites

1. The API server is running locally on port 8080
2. PostgreSQL database is set up with the required schemas

## Testing with curl

### Health Check

```bash
curl -X GET http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

### User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

Expected response:
```json
{
  "token": "sample-jwt-token-<uuid>",
  "userId": "<uuid>",
  "email": "test@example.com",
  "expiresAt": 1678909876
}
```

### User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

Expected response:
```json
{
  "token": "sample-jwt-token-<uuid>",
  "userId": "<uuid>",
  "email": "test@example.com",
  "expiresAt": 1678909876
}
```

### Create Client

Save the userId from the registration/login response and use it here:

```bash
curl -X POST http://localhost:8080/api/v1/clients?userId=<user_id_here> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sample-jwt-token" \
  -d '{"name":"Test Client","email":"client@example.com","whatsapp":"+551199999999"}'
```

Expected response:
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

### Get Clients

```bash
curl -X GET http://localhost:8080/api/v1/clients?userId=<user_id_here> \
  -H "Authorization: Bearer sample-jwt-token"
```

Expected response:
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

### Get Client by ID

```bash
curl -X GET http://localhost:8080/api/v1/clients/<client_uuid> \
  -H "Authorization: Bearer sample-jwt-token"
```

Expected response:
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

### Update Client

```bash
curl -X PUT http://localhost:8080/api/v1/clients/<client_uuid> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sample-jwt-token" \
  -d '{"name":"Updated Client","email":"updated@example.com","whatsapp":"+551188888888"}'
```

Expected response:
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

### Delete Client

```bash
curl -X DELETE http://localhost:8080/api/v1/clients/<client_uuid> \
  -H "Authorization: Bearer sample-jwt-token"
```

Expected response:
```json
{
  "message": "Client deleted successfully"
}
```

## Testing with Postman

1. Import the Postman collection from the file `testing/swap_ai_api.postman_collection.json`
2. Set up environment variables:
   - `base_url` = `http://localhost:8080`
   - `user_id` (will be populated after registration/login)
   - `auth_token` (will be populated after registration/login)

3. Execute the requests in the following order:
   - Health Check
   - Register User (or Login if already registered)
   - Create Client
   - Get All Clients
   - Get Client by ID
   - Update Client
   - Delete Client
