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

        // Return the file ID for further processing
        const fileId = response.data.id;
        
        // If successful, delete the original unprocessed file
        if (fileId) {
            await deleteOriginalFile(clientId, token);
        }
        
        return fileId;
        
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
        return null;
    }
}

async function deleteOriginalFile(clientId, token) {
    try {
        const url = `http://192.168.10.3:8080/api/v1/admin/${clientId}/delete-unprocessed`;
        
        const response = await axios.delete(url, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        console.log('Arquivo original não processado excluído com sucesso:');
        console.log('Mensagem:', response.data.message);
        console.log('Arquivos excluídos:', response.data.deletedCount);
        
        return true;
    } catch (error) {
        console.error('Erro ao excluir arquivo original:');
        if (error.response) {
            console.error(error.response.data.error);
        } else {
            console.error(error.message);
        }
        return false;
    }
}

// Função principal para ser chamada de forma assíncrona
async function processFile(filePath, clientId, token) {
    const fileId = await adminUploadFile(filePath, clientId, token);
    if (fileId) {
        console.log(`Processamento concluído. ID do arquivo processado: ${fileId}`);
        return fileId;
    } else {
        console.error('Falha no processamento do arquivo.');
        return null;
    }
}

// Configurações
const filePath = './documento_processado.pdf';
const clientId = '<client_uuid>'; // UUID do cliente
const token = '<token_uuid>'; // Token de autenticação (sem o prefixo "Bearer")

// Executar upload administrativo
processFile(filePath, clientId, token).catch(console.error);

// Export para uso como módulo se necessário
module.exports = {
    adminUploadFile,
    processFile
};
