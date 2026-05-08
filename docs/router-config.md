# Looking Glass - Router Configuration

This directory contains router configuration examples and documentation for different vendors.

## Adding Routers

### Method 1: Web UI (Recommended)
1. Navigate to `http://localhost:3000/admin/routers`
2. Click "Add Router"
3. Fill in the details:
   - **Name**: Display name (e.g., "Core Router MK-01")
   - **Hostname/IP**: Router's management IP
   - **Vendor**: Choose from MikroTik, Juniper, Cisco, Huawei
   - **Port**: SSH port (default: 22)
   - **Username**: SSH username
   - **Password**: SSH password
   - **Zone**: Select appropriate zone
   - **Tags**: Optional tags for filtering

### Method 2: Database Direct
Run SQL from `seeds/router_seed.sql` after modifying credentials.

### Method 3: API
```bash
curl -X POST http://localhost:8080/api/v1/routers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "name": "Core Router MK-01",
    "hostname": "10.0.1.1",
    "vendor": "mikrotik",
    "port": 22,
    "username": "admin",
    "password": "your_password",
    "zone_id": "<zone-uuid>",
    "is_active": true
  }'
```

## Supported Vendors

### MikroTik RouterOS
```bash
# Required SSH config
ssh admin@10.0.1.1
# RouterOS 6.x or 7.x supported
```

### Juniper Junos
```bash
# Required SSH config
ssh admin@10.0.2.1
# Junos OS 12.x+
```

### Cisco IOS/IOS-XE
```bash
# Required SSH config
ssh admin@10.0.3.1
# IOS 15.x+, IOS-XE 16.x+
```

### Huawei VRP
```bash
# Required SSH config
ssh admin@10.0.4.1
# VRP 5.x, 8.x
```

## Router Requirements

Each router must:
1. Have SSH enabled and accessible
2. Support the following commands:
   - `ping <target>` - ICMP ping
   - `traceroute <target>` - Traceroute
   - `show bgp <prefix>` or `display bgp routing-table <prefix>` - BGP lookup

## Zones

Group routers by zone for easier management:
- Core
- Edge
- Distribution
- Access
- Peering

## Example Router Configuration

### MikroTik
```
/user
add name=lookingglass group=full password=your_password

/ip service
set ssh address=0.0.0.0/0 disabled=no
```

### Juniper
```
set system services ssh
set system login user lookingglass class super-user authentication plain-text-password
```

### Cisco
```
username lookingglass privilege 15 secret your_password
ip ssh source-interface Loopback0
```

## Troubleshooting

### SSH Connection Issues
- Verify firewall allows port 22
- Check SSH service is enabled
- Verify credentials are correct
- Ensure router has enough resources

### Command Not Supported
- Some routers may have limited command support
- Check vendor documentation for exact command syntax
- Contact support if issues persist