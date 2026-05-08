# ISP Looking Glass Platform - Architecture Documentation

## Overview
A production-ready ISP Looking Glass platform with multi-vendor router support, providing network diagnostics tools including ping, traceroute, BGP lookups, and route analysis.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Frontend (Next.js)                              │
│  Dashboard │ Looking Glass │ Router Management │ Zones │ Users │ Logs       │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Backend API (Go)                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐      │
│  │ REST API │  │ WebSocket│  │ Auth API │  │ Admin API│  │ Query API│      │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘      │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                     Middleware Layer                                  │  │
│  │  Rate Limiting │ JWT Auth │ RBAC │ Logging │ CORS │ Recovery          │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
        ┌─────────────────────────────┼─────────────────────────────┐
        ▼                             ▼                             ▼
┌───────────────┐           ┌───────────────┐           ┌───────────────┐
│   PostgreSQL  │           │    Redis      │           │   SSH Pool    │
│   (Persistent)│           │   (Cache)     │           │   (Router)    │
└───────────────┘           └───────────────┘           └───────────────┘
                                                              │
                    ┌─────────────────────────────────────────┼─────────────┐
                    ▼                                         ▼             ▼
              ┌───────────┐                           ┌───────────┐   ┌───────────┐
              │ MikroTik  │                           │ Juniper   │   │  Cisco    │
              │ RouterOS  │                           │   Junos   │   │ IOS/XE/XR │
              └───────────┘                           └───────────┘   └───────────┘
                                                                      │
                                                                  ┌───────────┐
                                                                  │  Huawei   │
                                                                  │   VRP     │
                                                                  └───────────┘
```

## Supported Vendors & Commands

### MikroTik RouterOS
| Operation | Command |
|-----------|---------|
| Ping | `/ping <host> count=<count> size=<size>` |
| Traceroute | `/tool traceroute <host>` |
| BGP Route | `/routing bgp peer print` |
| Route Table | `/routing route print where active` |
| BGP Routes | `/routing bgp route print` |
| BGP Advertised | `/routing bgp advertisement print` |

### Juniper Junos
| Operation | Command |
|-----------|---------|
| Ping | `ping <host> count <count> size <size>` |
| Traceroute | `traceroute <host>` |
| BGP Route | `show route protocol bgp` |
| Route Table | `show route` |
| BGP Summary | `show bgp summary` |
| BGP Neighbor | `show bgp neighbor` |
| BGP Advertised | `show bgp neighbor <peer> advertised-routes` |
| BGP Received | `show bgp neighbor <peer> received-routes` |

### Cisco IOS/XE/XR/NXOS
| Operation | Command |
|-----------|---------|
| Ping | `ping <host> repeat <count> size <size>` |
| Traceroute | `traceroute <host>` |
| BGP Route | `show ip bgp` |
| Route Table | `show ip route` |
| BGP Summary | `show ip bgp summary` |
| BGP Neighbor | `show ip bgp neighbors` |
| BGP Advertised | `show ip bgp neighbors <peer> advertised-routes` |
| BGP Received | `show ip bgp neighbors <peer> received-routes` |

### Huawei VRP
| Operation | Command |
|-----------|---------|
| Ping | `ping -c <count> -s <size> <host>` |
| Traceroute | `tracert <host>` |
| BGP Route | `display bgp routing-table` |
| Route Table | `display ip routing-table` |
| BGP Summary | `display bgp peer` |
| BGP Neighbor | `display bgp peer <peer>` |
| BGP Advertised | `display bgp routing-table advertised-peer <peer>` |
| BGP Received | `display bgp routing-table received-peer <peer>` |

## Database Schema

### Users Table
```sql
CREATE TABLE users (
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
```

### Zones Table
```sql
CREATE TABLE zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    location VARCHAR(255),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Routers Table
```sql
CREATE TABLE routers (
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
    zone_id UUID REFERENCES zones(id),
    is_active BOOLEAN DEFAULT true,
    is_online BOOLEAN DEFAULT false,
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Audit Logs Table
```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    ip_address VARCHAR(45),
    user_agent TEXT,
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Query History Table
```sql
CREATE TABLE query_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    router_id UUID REFERENCES routers(id),
    command_type VARCHAR(100) NOT NULL,
    target VARCHAR(255) NOT NULL,
    parameters JSONB,
    result_summary TEXT,
    execution_time_ms INTEGER,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Current user info

### Routers
- `GET /api/v1/routers` - List routers
- `GET /api/v1/routers/:id` - Get router details
- `POST /api/v1/routers` - Create router
- `PUT /api/v1/routers/:id` - Update router
- `DELETE /api/v1/routers/:id` - Delete router
- `POST /api/v1/routers/:id/test` - Test connection
- `GET /api/v1/routers/:id/status` - Get router status

### Zones
- `GET /api/v1/zones` - List zones
- `GET /api/v1/zones/:id` - Get zone details
- `POST /api/v1/zones` - Create zone
- `PUT /api/v1/zones/:id` - Update zone
- `DELETE /api/v1/zones/:id` - Delete zone

### Looking Glass Queries
- `POST /api/v1/query/ping` - Execute ping
- `POST /api/v1/query/traceroute` - Execute traceroute
- `POST /api/v1/query/bgp-route` - BGP route lookup
- `POST /api/v1/query/bgp-summary` - BGP peer summary
- `POST /api/v1/query/asn` - ASN lookup
- `POST /api/v1/query/prefix` - Prefix lookup

### WebSocket
- `WS /api/v1/ws` - Realtime query streaming

### Admin
- `GET /api/v1/admin/users` - List users
- `POST /api/v1/admin/users` - Create user
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user
- `GET /api/v1/admin/logs` - Audit logs

### Dashboard
- `GET /api/v1/dashboard/stats` - Dashboard statistics
- `GET /api/v1/dashboard/health` - System health

## Security Features

1. **Command Whitelisting** - Only allowed commands can be executed
2. **Read-only Mode** - All operations are read-only
3. **Rate Limiting** - Per-IP and per-user rate limits
4. **JWT Authentication** - Secure token-based auth
5. **RBAC** - Role-based access control
6. **Audit Logging** - All actions are logged
7. **SSH Key Support** - Secure authentication
8. **Encrypted Credentials** - Passwords are encrypted

## Deployment Options

1. **Docker** - Full Docker Compose setup
2. **Kubernetes** - K8s manifests for production
3. **Single Binary** - Standalone Go binary deployment

## Environment Variables

See `.env.example` for all configuration options.
