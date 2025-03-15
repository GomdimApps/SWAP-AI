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

-- New tables for the SWAP-AI system
CREATE TABLE system.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE system.folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES system.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clients.clients_swap (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    folder_id UUID NOT NULL REFERENCES system.folders(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('ativo', 'inativo')) DEFAULT 'ativo',
    email VARCHAR(255),
    whatsapp VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indices
CREATE INDEX idx_clients_folder_id ON clients.clients_swap(folder_id);
CREATE INDEX idx_clients_status ON clients.clients_swap(status);

-- Create trigger for updating updated_at timestamp
CREATE OR REPLACE FUNCTION system.update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON system.users
FOR EACH ROW EXECUTE FUNCTION system.update_modified_column();

CREATE TRIGGER update_clients_updated_at
BEFORE UPDATE ON clients.clients_swap
FOR EACH ROW EXECUTE FUNCTION system.update_modified_column();