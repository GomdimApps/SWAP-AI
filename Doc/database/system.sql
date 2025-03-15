GRANT ALL PRIVILEGES ON DATABASE swap_database TO swap_user;
-- Criar os schemas
CREATE SCHEMA clients;
DROP SCHEMA public CASCADE;
CREATE SCHEMA products;
CREATE SCHEMA system;
CREATE SCHEMA files;
CREATE SCHEMA transactions;
GRANT USAGE ON SCHEMA products, clients, system, files, transactions TO swap_user;

ALTER DATABASE swap_database SET default_transaction_isolation = 'READ COMMITTED';
ALTER DATABASE swap_database SET log_statement = 'all';

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

-- Moved from system to files schema
CREATE TABLE files.folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clients.clients_swap (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES system.users(id) ON DELETE CASCADE,
    folder_id UUID NOT NULL REFERENCES files.folders(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('ativo', 'inativo')) DEFAULT 'ativo',
    email VARCHAR(255),
    whatsapp VARCHAR(20),  -- This field name must match the model field name
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indices
CREATE INDEX idx_clients_folder_id ON clients.clients_swap(folder_id);
CREATE INDEX idx_clients_user_id ON clients.clients_swap(user_id);
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