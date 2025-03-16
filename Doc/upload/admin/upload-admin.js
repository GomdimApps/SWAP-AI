const fs = require('fs');
const FormData = require('form-data');
const axios = require('axios');

async function adminUploadFile(filePath, clientId, token) {
    const url = `http://192.168.10.3:8080/api/v1/admin/${clientId}/upload`;

    try {
        const form = new FormData();
        form.append('file', fs.createReadStream(filePath));

        // Configurar limites maiores para permitir arquivos grandes (110MB)
        const axiosInstance = axios.create({
            maxContentLength: Infinity,
            maxBodyLength: Infinity,
        });

        const response = await axiosInstance.post(url, form, {
            headers: {
                ...form.getHeaders(),
                'Authorization': `Bearer ${token}`
            }
        });

        console.log('Upload administrativo realizado com sucesso:');
        console.log('ID do Arquivo:', response.data.id);
        console.log('Nome do Arquivo:', response.data.originalName);
        console.log('Tamanho:', (response.data.size / 1024 / 1024).toFixed(2), 'MB');
        console.log('Cliente ID:', response.data.clientId);
        console.log('Processado por:', response.data.processedBy);
        console.log('Mensagem:', response.data.message);
        
    } catch (error) {
        if (error.response) {
            console.error('Erro na resposta do servidor:', error.response.data.error);
            console.error('Status:', error.response.status);
        } else if (error.request) {
            console.error('Erro na requisição - não houve resposta do servidor');
            console.error(error.message);
        } else {
            console.error('Erro ao configurar requisição:', error.message);
        }
    }
}

// Configurações
const filePath = './documento_processado.pdf';
const clientId = '<client_uuid>'; // UUID do cliente
const token = '<token_uuid>'; // Token de autenticação (sem o prefixo "Bearer")

// Executar upload administrativo
adminUploadFile(filePath, clientId, token);
