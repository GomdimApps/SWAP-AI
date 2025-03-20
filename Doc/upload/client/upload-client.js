const fs = require('fs');
const FormData = require('form-data');
const axios = require('axios');

async function uploadFiles(filePath1, filePath2, clientId, token) {
  try {
    const form = new FormData();
    form.append('file1', fs.createReadStream(filePath1));
    form.append('file2', fs.createReadStream(filePath2));
    
    const resp = await axios.create({
      maxContentLength: Infinity,
      maxBodyLength: Infinity
    }).post(
      `http://192.168.10.3:8080/api/v1/clients/${clientId}/upload`,
      form,
      { headers: { ...form.getHeaders(), 'Authorization': `Bearer ${token}` } }
    );
    
    const { file1_id, file2_id, project_id, message } = resp.data;
    console.log(`✓ Upload OK: File1 ID=${file1_id}, File2 ID=${file2_id}, ProjectID=${project_id}`);
    console.log(`  Mensagem: ${message}`);
    
    return { file1_id, file2_id, project_id };
  } catch (err) {
    const errorMsg = err.response?.data?.error || err.request ? 'Sem resposta do servidor' : err.message;
    console.error(`✗ Erro: ${errorMsg}`);
    return null;
  }
}

const filePath1 = './documento1.pdf';
const filePath2 = './documento2.pdf';
const clientId = '<client_uuid>'; // UUID do cliente
const token = '<token_uuid>'; // Token de autenticação (sem o prefixo "Bearer")

uploadFiles(filePath1, filePath2, clientId, token);