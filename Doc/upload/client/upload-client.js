const fs = require('fs');
const FormData = require('form-data');
const axios = require('axios');

async function uploadFile(filePath, clientId, token) {
    const url = `http://192.168.10.3:8080/api/v1/clients/${clientId}/upload`;

    try {
        const form = new FormData();
        form.append('file', fs.createReadStream(filePath));

        const response = await axios.post(url, form, {
            headers: {
                ...form.getHeaders(),
                'Authorization': `Bearer ${token}`
            },
            maxContentLength: 100 * 1024 * 1024, // 100MB
            maxBodyLength: 100 * 1024 * 1024,
        });

        console.log('Upload realizado com sucesso:');
        console.log('ID do Arquivo:', response.data.id);
        console.log('Nome do Arquivo:', response.data.originalName);
        console.log('Tamanho:', (response.data.size / 1024).toFixed(2), 'KB');
        console.log('Mensagem:', response.data.message);
        
    } catch (error) {
        if (error.response) {
            console.error('Erro na resposta do servidor:', error.response.data.error);
            console.error('Status:', error.response.status);
        } else {
            console.error('Erro na requisição:', error.message);
        }
    }
}

// Configurações
const filePath = './documento.pdf';
const clientId = '<client_uuid>'; // UUID do cliente
const token = '<token_uuid>'; // Token de autenticação (sem o prefixo "Bearer")

// Executar upload
uploadFile(filePath, clientId, token);