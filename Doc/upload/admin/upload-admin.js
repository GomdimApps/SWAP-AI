const fs = require('fs');
const FormData = require('form-data');
const axios = require('axios');

// Base configuration
const API_URL = 'http://192.168.10.3:8080/api/v1/admin';
const MAX_UPLOAD_SIZE = 110; // MB

/**
 * Upload a processed file to the server
 * @param {string} filePath - Path to the file to upload
 * @param {string} clientId - UUID of the client
 * @param {string} token - Authentication token
 * @param {string} originalFileId - Optional ID of the original file to delete
 * @returns {Promise<string|null>} The ID of the uploaded file or null on failure
 */
async function processAndUpload(filePath, clientId, token, originalFileId = null) {
  try {
    console.log(`Starting admin upload for client ${clientId}...`);
    console.log(`File: ${filePath}`);
    
    // Create form with file
    const form = new FormData();
    form.append('file', fs.createReadStream(filePath));
    
    // Configure axios with unlimited size
    const client = axios.create({
      maxContentLength: Infinity,
      maxBodyLength: Infinity
    });
    
    // 1. Upload the processed file
    console.log('Uploading processed file...');
    const uploadResponse = await client.post(
      `${API_URL}/${clientId}/upload`,
      form,
      {
        headers: {
          ...form.getHeaders(),
          'Authorization': `Bearer ${token}`
        }
      }
    );
    
    // Extract file information
    const { id: fileId, originalName, size } = uploadResponse.data;
    console.log(`✓ Upload successful - File ID: ${fileId}`);
    console.log(`  Name: ${originalName}`);
    console.log(`  Size: ${(size / 1024 / 1024).toFixed(2)} MB`);
    
    // 2. Delete the original unprocessed file if ID was provided
    if (originalFileId) {
      console.log(`Deleting original file (ID: ${originalFileId})...`);
      const deleteUrl = `${API_URL}/${clientId}/delete-unprocessed/${originalFileId}`;
      
      const deleteResponse = await client.delete(deleteUrl, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      
      console.log(`✓ Original file deleted: ${deleteResponse.data.message}`);
    }
    
    return fileId;
    
  } catch (error) {
    console.error('Error during file processing:');
    
    if (error.response) {
      console.error(`Server error (${error.response.status}): ${error.response.data.error}`);
    } else if (error.request) {
      console.error('No response received from server. Check network connection.');
    } else {
      console.error(`Request error: ${error.message}`);
    }
    
    return null;
  }
}

/**
 * List all unprocessed files for a client
 * @param {string} clientId - UUID of the client
 * @param {string} token - Authentication token
 * @returns {Promise<Array>} Array of unprocessed files
 */
async function listUnprocessedFiles(clientId, token) {
  try {
    const response = await axios.get(
      `${API_URL}/${clientId}/unprocessed-files`,
      { headers: { 'Authorization': `Bearer ${token}` } }
    );
    
    console.log('\nUnprocessed files:');
    if (response.data.files.length === 0) {
      console.log('  No unprocessed files found');
    } else {
      response.data.files.forEach(file => {
        console.log(`  • ID: ${file.id}`);
        console.log(`    Name: ${file.originalName}`);
        console.log(`    Size: ${(file.fileSize / (1024 * 1024)).toFixed(2)} MB`);
        console.log('');
      });
    }
    
    return response.data.files;
  } catch (error) {
    console.error('Failed to list files:', error.message);
    return [];
  }
}

// Run as standalone script or export as module
if (require.main === module) {
  // CLI execution
  const filePath = './documento_processado.pdf';
  const clientId = '<client_uuid>'; // UUID of the client
  const originalFileId = '<original_file_uuid>'; // Optional: UUID of original file to delete
  const token = '<token_uuid>'; // Auth token (without Bearer prefix)
  
  // Run the processor
  processAndUpload(filePath, clientId, token, originalFileId)
    .then(fileId => {
      if (fileId) {
        console.log(`\n✅ Process completed successfully. New file ID: ${fileId}`);
      } else {
        console.error('\n❌ Process failed');
      }
    })
    .catch(console.error);
} else {
  // Module export
  module.exports = {
    processAndUpload,
    listUnprocessedFiles
  };
}
