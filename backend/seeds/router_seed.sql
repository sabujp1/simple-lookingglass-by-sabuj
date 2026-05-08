-- Router Seed Data
-- Add this to your database or use the Admin UI at http://localhost:3000/admin/routers

-- Example MikroTik Router
INSERT INTO routers (id, name, hostname, vendor, port, username, password_encrypted, zone_id, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Core Router MK-01',
    '10.0.1.1',
    'mikrotik',
    22,
    'admin',
    'encrypted_password_here',
    (SELECT id FROM zones WHERE name = 'Core'),
    true,
    NOW(),
    NOW()
);

-- Example Juniper Router
INSERT INTO routers (id, name, hostname, vendor, port, username, password_encrypted, zone_id, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Edge Router JNPR-01',
    '10.0.2.1',
    'juniper',
    22,
    'admin',
    'encrypted_password_here',
    (SELECT id FROM zones WHERE name = 'Edge'),
    true,
    NOW(),
    NOW()
);

-- Example Cisco Router
INSERT INTO routers (id, name, hostname, vendor, port, username, password_encrypted, zone_id, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Distribution Cisco-01',
    '10.0.3.1',
    'cisco',
    22,
    'admin',
    'encrypted_password_here',
    (SELECT id FROM zones WHERE name = 'Distribution'),
    true,
    NOW(),
    NOW()
);

-- Example Huawei Router
INSERT INTO routers (id, name, hostname, vendor, port, username, password_encrypted, zone_id, is_active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Access Huawei-01',
    '10.0.4.1',
    'huawei',
    22,
    'admin',
    'encrypted_password_here',
    (SELECT id FROM zones WHERE name = 'Access'),
    true,
    NOW(),
    NOW()
);