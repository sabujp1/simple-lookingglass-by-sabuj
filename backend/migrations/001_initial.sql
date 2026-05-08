-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    api_token VARCHAR(255) UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Zones table
CREATE TABLE IF NOT EXISTS zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    location VARCHAR(255),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Routers table
CREATE TABLE IF NOT EXISTS routers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    vendor VARCHAR(50) NOT NULL,
    model VARCHAR(255),
    asn INTEGER,
    ssh_port INTEGER DEFAULT 22,
    ssh_username VARCHAR(255) NOT NULL,
    ssh_password_encrypted TEXT,
    ssh_key_path TEXT,
    zone_id UUID REFERENCES zones(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT true,
    is_online BOOLEAN DEFAULT false,
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    ip_address VARCHAR(45),
    user_agent TEXT,
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Query history table
CREATE TABLE IF NOT EXISTS query_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    router_id UUID REFERENCES routers(id) ON DELETE SET NULL,
    command_type VARCHAR(100) NOT NULL,
    target VARCHAR(255) NOT NULL,
    parameters JSONB,
    result_summary TEXT,
    execution_time_ms INTEGER,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_routers_zone ON routers(zone_id);
CREATE INDEX IF NOT EXISTS idx_routers_vendor ON routers(vendor);
CREATE INDEX IF NOT EXISTS idx_routers_active ON routers(is_active, is_online);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_query_history_user ON query_history(user_id);
CREATE INDEX IF NOT EXISTS idx_query_history_router ON query_history(router_id);
CREATE INDEX IF NOT EXISTS idx_query_history_created ON query_history(created_at);

-- Insert default admin user (password: admin123)
-- bcrypt hash for 'admin123'
INSERT INTO users (username, email, password_hash, role)
VALUES ('admin', 'admin@lookingglass.local', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.RpITzQ.8uYvR8n4K.a', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Insert default zones
INSERT INTO zones (name, code, location, description)
VALUES 
    ('Dhaka', 'DHK', 'Dhaka, Bangladesh', 'Primary datacenter in Dhaka'),
    ('Chittagong', 'CGD', 'Chittagong, Bangladesh', 'Secondary datacenter in Chittagong'),
    ('Singapore', 'SIN', 'Singapore', 'Singapore PoP for international connectivity'),
    ('BDIX', 'BDX', 'Dhaka, Bangladesh', 'BDIX exchange point')
ON CONFLICT (code) DO NOTHING;