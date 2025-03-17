const fs = require('fs');
const FormData = require('form-data');
const axios = require('axios');

async function uploadFile(filePath, clientId, token) {
  try {
    const form = new FormData();
    form.append('file', fs.createReadStream(filePath));
    
    const resp = await axios.create({
      maxContentLength: Infinity,
      maxBodyLength: Infinity
    }).post(
      `http://192.168.10.3:8080/api/v1/clients/${clientId}/upload`,
      form,
      { headers: { ...form.getHeaders(), 'Authorization': `Bearer ${token}` } }
    );
    
    const { id, originalName, size, message } = resp.data;
    console.log(`✓ Upload OK: ID=${id}, Nome=${originalName}, Tamanho=${(size/1024/1024).toFixed(2)}MB`);
    console.log(`  Mensagem: ${message}`);
    
    return id;
  } catch (err) {
    const errorMsg = err.response?.data?.error || err.request ? 'Sem resposta do servidor' : err.message;
    console.error(`✗ Erro: ${errorMsg}`);
    return null;
  }
}

const filePath = './documento.pdf';
const clientId = '<client_uuid>'; // UUID do cliente
const token = '<token_uuid>'; // Token de autenticação (sem o prefixo "Bearer")

uploadFile(filePath, clientId, token);