# Upload de Arquivo para Cliente

#### Usando curl

```bash
curl -X POST http://192.168.10.3:8080/api/v1/clients/<client_uuid>/upload \
    -H "Authorization: Bearer <token_uuid>" \
    -F "file=@/path/to/file.pdf"
```

Resposta esperada:
```json
{
    "message": "File uploaded successfully",
    "id": "<file_uuid>",
    "originalName": "file.pdf",
    "size": 12345,
    "clientId": "<client_uuid>"
}
```

> **Nota sobre armazenamento de arquivos:** Os arquivos são armazenados em `/var/upload/<folder_uuid>/` onde `<folder_uuid>` é o ID da pasta do cliente.

### Obter Arquivos do Cliente

```bash
curl -X GET http://192.168.10.3:8080/api/v1/clients/<client_uuid>/files \
    -H "Authorization: Bearer <token_uuid>"
```

Resposta esperada:
```json
{
    "files": [
        {
            "ID": "<file_uuid>",
            "ClientID": "<client_uuid>",
            "FolderID": "<folder_uuid>",
            "OriginalName": "documento.pdf",
            "FilePath": "<client_uuid>/12345678-unique-file.pdf",
            "FileSize": 12345,
            "FileType": "application/pdf",
            "CreatedAt": "2023-01-01T12:00:00Z",
            "UpdatedAt": "2023-01-01T12:00:00Z"
        }
    ]
}
```

