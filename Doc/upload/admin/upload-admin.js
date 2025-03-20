const fs = require('fs');
const FormData = require('form-data');
const axios = require('axios');

// Configuração compacta
const API_URL = 'http://192.168.10.3:8080/api/v1/admin';

// Cliente axios configurado
const createClient = (token) => axios.create({
  maxContentLength: Infinity,
  maxBodyLength: Infinity,
  headers: { 'Authorization': `Bearer ${token}` }
});

/**
 * Upload de arquivo processado pelo admin
 * @param {string} filePath - Caminho do arquivo
 * @param {string} clientId - UUID do cliente
 * @param {string} token - Token de autenticação
 * @param {string} projectId - ProjectID referente ao projeto do client (opcional)
 * @returns {Promise<Object|null>} Objeto com fileId e project_id ou null
 */
async function processAndUpload(filePath, clientId, token, projectId = '') {
  const client = createClient(token);
  const form = new FormData();
  form.append('file', fs.createReadStream(filePath));
  // Se projectId for fornecido, envia junto
  if (projectId) {
    form.append('project_id', projectId);
  }
  
  try {
    console.log(`Upload admin para cliente ${clientId}. Arquivo: ${filePath}`);
    
    // Upload do arquivo
    const uploadRes = await client.post(
      `${API_URL}/${clientId}/upload`,
      form,
      { headers: { ...form.getHeaders() } }
    );
    
    // Recebe os dados, incluindo project_id
    const { id: fileId, originalName, size, project_id } = uploadRes.data;
    console.log(`✓ Upload Admin OK - ID: ${fileId}, Nome: ${originalName}, Tamanho: ${(size/1024/1024).toFixed(2)}MB, ProjectID: ${project_id}`);
    
    return { fileId, project_id };
  } catch (err) {
    console.error('Erro:', err.response?.data?.error || err.message);
    return null;
  }
}

// Execução direta ou exportação como módulo
if (require.main === module) {
  const args = {
    filePath: './documento_processado.pdf',
    clientId: '<client_uuid>',
    token: '<token_uuid>',
    projectId: '<project_uuid>' // Passe aqui o projectId desejado ou deixe vazio
  };
  
  processAndUpload(args.filePath, args.clientId, args.token, args.projectId)
    .then(result => console.log(result ? `\n✅ Sucesso! ProjectID: ${result.project_id}` : '\n❌ Falha'))
    .catch(err => console.error('\n❌ Erro:', err.message));
} else {
  module.exports = { processAndUpload };
}