# TouNetCore API Test Examples

This directory contains example API requests for testing the TouNetCore system.

## Setup

Start the server:
```bash
go run cmd/server/main.go
```

The server runs on port 8080 by default. If that port is busy, use:
```bash
PORT=8081 go run cmd/server/main.go
```

## Authentication

### 1. Admin Login
```bash
curl -X POST http://localhost:8081/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 2. User Registration
```bash
curl -X POST http://localhost:8081/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "phone": "13800138000",
    "pushdeer_token": "PUSHDEER_TOKEN_HERE",
    "invite_code": "METzuzFY8KSudQOnvpS-Qw"
  }'
```

## User Operations

### 3. Get User Info
```bash
curl -X GET http://localhost:8081/api/v1/user/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Get User Apps
```bash
curl -X GET http://localhost:8081/api/v1/user/apps \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Generate NKey
```bash
curl -X POST http://localhost:8081/api/v1/nkey/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "username": ["testuser"],
    "app_ids": ["searchall", "CardPreview"]
  }'
```

### 6. Validate NKey
```bash
curl -X POST http://localhost:8081/api/v1/nkey/validate \
  -H "Content-Type: application/json" \
  -d '{
    "nkey": "TOUNET_1_XXXXXX",
    "app_id": "searchall"
  }'
```

## Admin Operations

### 7. List Users (Admin)
```bash
curl -X GET http://localhost:8081/api/v1/admin/users?page=1&size=10 \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 8. List Apps (Admin)
```bash
curl -X GET http://localhost:8081/api/v1/admin/apps \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 9. Create App (Admin)
```bash
curl -X POST http://localhost:8081/api/v1/admin/apps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "app_id": "NewApp",
    "name": "New Application",
    "description": "Test application",
    "required_permission_level": "user",
    "is_active": true
  }'
```

### 10. Generate Invite Code (Admin)
```bash
curl -X POST http://localhost:8081/api/v1/admin/invite-codes \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

## Available Invite Codes

Here are the pre-generated invite codes for testing:
- METzuzFY8KSudQOnvpS-Qw
- Wh9iukC14-VXTnDHwaeiQw
- _41LcJuBUGv4JIR0WLieVA
- HEyOOM1MA3r_wuQ0IH-D4w
- DY29kFfqgg1hWg9MLJ5jEA
- E2RNMI7j_BrIZqwuNqq2mA
- tPPczvk_DT8wpgcl3mFtGg
- 6Zbzh41vrdR6yXf2JPdh6A
- 0z55Cf8fFussBIoYeFkGNw
- jtqpvyvqlLPPQzwhVrJR8Q

## Pre-configured Applications

- **searchall**: Global search functionality
- **segaasstes**: Sega assets management  
- **dxprender**: DXP rendering service
- **CardPreview**: Card preview functionality
- **livecontent_basic**: Basic live content access
- **livecontent_admin**: Administrative live content access (admin only)

## Admin Credentials

- **Username**: admin
- **Password**: admin123

⚠️ **Important**: Change the admin password after first login in production!
