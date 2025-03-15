GRANT ALL PRIVILEGES ON DATABASE nexogestor TO nexogestor_infra;
-- Criar os schemas
CREATE SCHEMA clients;
DROP SCHEMA public CASCADE;
CREATE SCHEMA products;
CREATE SCHEMA system;
CREATE SCHEMA files;
CREATE SCHEMA transactions;
GRANT USAGE ON SCHEMA products, clients, system, files, transactions TO nexogestor_infra;

ALTER DATABASE nexogestor SET default_transaction_isolation = 'READ COMMITTED';
ALTER DATABASE nexogestor SET log_statement = 'all';

CREATE TABLE clients.clients (  
  id UUID PRIMARY KEY,  
  name VARCHAR(255) NOT NULL,  
  email VARCHAR(255) UNIQUE NOT NULL,  
  password_hash TEXT NOT NULL,  
  phone VARCHAR(20),  
  cpf_cnpj VARCHAR(14),  
  status VARCHAR(10) DEFAULT 'ativo',  
  created_at TIMESTAMP DEFAULT NOW()  
);  