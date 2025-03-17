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
 * Upload de arquivo processado
 * @param {string} filePath - Caminho do arquivo
 * @param {string} clientId - UUID do cliente
 * @param {string} token - Token de autenticação
 * @param {string} originalFileId - ID do arquivo original (opcional)
 * @returns {Promise<string|null>} ID do arquivo ou null
 */
async function processAndUpload(filePath, clientId, token, originalFileId = null) {
  const client = createClient(token);
  const form = new FormData();
  form.append('file', fs.createReadStream(filePath));
  
  try {
    console.log(`Upload para cliente ${clientId}. Arquivo: ${filePath}`);
    
    // Upload do arquivo
    const uploadRes = await client.post(
      `${API_URL}/${clientId}/upload`,
      form,
      { headers: { ...form.getHeaders() } }
    );
    
    const { id: fileId, originalName, size } = uploadRes.data;
    console.log(`✓ Upload OK - ID: ${fileId}, Nome: ${originalName}, Tamanho: ${(size/1024/1024).toFixed(2)}MB`);
    
    // Exclusão do arquivo original
    if (originalFileId) {
      await client.delete(`${API_URL}/${clientId}/delete-unprocessed/${originalFileId}`);
      console.log(`✓ Arquivo original ${originalFileId} excluído`);
    }
    
    return fileId;
  } catch (err) {
    console.error('Erro:', err.response?.data?.error || err.message);
    return null;
  }
}

/**
 * Lista arquivos não processados
 * @param {string} clientId - UUID do cliente
 * @param {string} token - Token de autenticação
 * @returns {Promise<Array>} Lista de arquivos
 */
async function listUnprocessedFiles(clientId, token) {
  try {
    const res = await createClient(token).get(`${API_URL}/${clientId}/unprocessed-files`);
    const files = res.data.files;
    
    console.log('\nArquivos não processados:');
    files.length === 0
      ? console.log('  Nenhum arquivo encontrado')
      : files.forEach(f => console.log(`  • ID: ${f.id}, Nome: ${f.originalName}, Tamanho: ${(f.fileSize/1048576).toFixed(2)}MB`));
    
    return files;
  } catch (err) {
    console.error('Falha ao listar arquivos:', err.message);
    return [];
  }
}

// Execução direta ou exportação como módulo
if (require.main === module) {
  const args = {
    filePath: './documento_processado.pdf',
    clientId: '<client_uuid>',
    originalFileId: '<original_file_uuid>',
    token: '<token_uuid>'
  };
  
  processAndUpload(args.filePath, args.clientId, args.token, args.originalFileId)
    .then(fileId => console.log(fileId ? `\n✅ Sucesso! Novo ID: ${fileId}` : '\n❌ Falha'))
    .catch(err => console.error('\n❌ Erro:', err.message));
} else {
  module.exports = { processAndUpload, listUnprocessedFiles };
}