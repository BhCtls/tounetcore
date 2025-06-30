# TouNetCore API Test Examples

This directory contains example API requests for testing the TouNetCore system.

## User Permission Levels

The system uses a hierarchical permission system with the following levels (from highest to lowest):

- **admin**: Full system access, can manage users, apps, and system settings
- **trusted**: Enhanced user with access to trusted applications and features
- **user**: Standard user access to basic applications
- **disableduser**: Minimal access, can only use applications specifically marked for disabled users

## Setup

Start the server:
```bash
go run cmd/server/main.go
```

The server runs on port 44544 by default. If that port is busy, use:
```bash
PORT=8081 go run cmd/server/main.go
```

## Authentication

### 1. Admin Login
```bash
curl -X POST http://localhost:44544/api/v1/login \
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
curl -X POST http://localhost:44544/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "phone": "13800138000",
    "pushdeer_token": "PUSHDEER_TOKEN_HERE",
    "invite_code": "METzuzFY8KSudQOnvpS-Qw"
  }'
```

Response:
```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "user_id": 2,
    "username": "testuser",
    "status": "user"
  }
}
```

## User Operations

### 3. Get User Info
```bash
curl -X GET http://localhost:44544/api/v1/user/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 2,
    "username": "testuser",
    "phone": "13800138000",
    "status": "user",
    "created_at": "2025-06-29T10:30:00Z",
    "updated_at": "2025-06-29T10:30:00Z"
  }
}
```

### 4. Get User Apps
```bash
curl -X GET http://localhost:44544/api/v1/user/apps \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "app_id": "searchall",
      "name": "Global Search",
      "description": "Global search functionality",
      "required_permission_level": "user",
      "is_active": true
    },
    {
      "app_id": "CardPreview",
      "name": "Card Preview",
      "description": "Card preview functionality",
      "required_permission_level": "user",
      "is_active": true
    },
    {
      "app_id": "livecontent_basic",
      "name": "Live Content Basic",
      "description": "Basic live content access",
      "required_permission_level": "user",
      "is_active": true
    }
  ]
}
```

### 5. Generate NKey
```bash
curl -X POST http://localhost:44544/api/v1/nkey/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "username": ["testuser"],
    "app_ids": ["searchall", "CardPreview"]
  }'
```

Response:
```json
{
  "code": 200,
  "message": "NKey generated successfully",
  "data": {
    "nkey": "TOUNET_1_AbCdEf123456",
    "expires_at": "2025-06-29T11:00:00Z",
    "valid_for_apps": ["searchall", "CardPreview"],
    "pushdeer_sent": true
  }
}
```

### 6. Validate NKey
```bash
curl -X POST http://localhost:44544/api/v1/nkey/validate \
  -H "Content-Type: application/json" \
  -d '{
    "nkey": "TOUNET_1_XXXXXX",
    "app_id": "searchall"
  }'
```

Response (Valid):
```json
{
  "code": 200,
  "message": "NKey is valid",
  "data": {
    "valid": true,
    "user_id": 2,
    "username": "testuser",
    "app_id": "searchall",
    "expires_at": "2025-06-29T11:00:00Z"
  }
}
```

Response (Invalid/Expired):
```json
{
  "code": 400,
  "message": "NKey is invalid or expired",
  "data": {
    "valid": false
  }
}
```

## Admin Operations

### 7. Create User (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "username": "newuser",
    "password": "password123",
    "status": "user",
    "phone": "13800138000",
    "pushdeer_token": "PUSHDEER_TOKEN_HERE"
  }'
```

**Creating a trusted user:**
```bash
curl -X POST http://localhost:44544/api/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "username": "trusteduser",
    "password": "password123",
    "status": "trusted",
    "phone": "13800138001",
    "pushdeer_token": "PUSHDEER_TOKEN_HERE"
  }'
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_id": 3
  }
}
```

### 8. List Users (Admin)
```bash
curl -X GET http://localhost:44544/api/v1/admin/users?page=1&size=10 \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "admin",
        "phone": "",
        "status": "admin",
        "created_at": "2025-06-29T09:00:00Z",
        "updated_at": "2025-06-29T09:00:00Z"
      },
      {
        "id": 2,
        "username": "testuser",
        "phone": "13800138000",
        "status": "user",
        "created_at": "2025-06-29T10:30:00Z",
        "updated_at": "2025-06-29T10:30:00Z"
      }
    ],
    "total": 2,
    "page": 1,
    "size": 10
  }
}
```

### 9. List Invite Codes (Admin)
```bash
curl -X GET http://localhost:44544/api/v1/admin/invite-codes?page=1&size=20 \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "invite_codes": [
      {
        "id": 1,
        "code": "METzuzFY8KSudQOnvpS-Qw",
        "used": true,
        "used_by": "testuser",
        "created_at": "2025-06-28T00:00:00Z",
        "used_at": "2025-06-29T10:30:00Z"
      },
      {
        "id": 2,
        "code": "Wh9iukC14-VXTnDHwaeiQw",
        "used": false,
        "used_by": "",
        "created_at": "2025-06-28T00:00:00Z",
        "used_at": null
      }
    ],
    "total": 10,
    "page": 1,
    "size": 20
  }
}
```

### 10. Update User (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/users/1/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "username": "updated_user",
    "status": "user",
    "phone": "13900139000"
  }'
```

Response:
```json
{
  "code": 200,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "username": "updated_user",
    "phone": "13900139000",
    "status": "user",
    "updated_at": "2025-06-29T11:15:00Z"
  }
}
```

### 11. Delete User (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/users/1/delete \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "User deleted successfully",
  "data": {
    "deleted_user_id": 1
  }
}
```

### 12. List Apps (Admin)
```bash
curl -X GET http://localhost:44544/api/v1/admin/apps \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "app_id": "searchall",
      "name": "Global Search",
      "description": "Global search functionality",
      "required_permission_level": "user",
      "is_active": true,
      "secret_key": "searchall_secret_key_123",
      "created_at": "2025-06-28T00:00:00Z",
      "updated_at": "2025-06-28T00:00:00Z"
    },
    {
      "app_id": "livecontent_admin",
      "name": "Live Content Admin",
      "description": "Administrative live content access",
      "required_permission_level": "admin",
      "is_active": true,
      "secret_key": "livecontent_admin_secret_456",
      "created_at": "2025-06-28T00:00:00Z",
      "updated_at": "2025-06-28T00:00:00Z"
    }
  ]
}
```

### 13. Create App (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/apps \
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

Response:
```json
{
  "code": 200,
  "message": "App created successfully",
  "data": {
    "app_id": "NewApp",
    "name": "New Application",
    "description": "Test application",
    "required_permission_level": "user",
    "is_active": true,
    "secret_key": "NewApp_secret_key_789",
    "created_at": "2025-06-29T11:20:00Z"
  }
}
```

### 14. Update App (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/apps/NewApp/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "name": "Updated Application Name",
    "description": "Updated description"
  }'
```

Response:
```json
{
  "code": 200,
  "message": "App updated successfully",
  "data": {
    "app_id": "NewApp",
    "name": "Updated Application Name",
    "description": "Updated description",
    "required_permission_level": "user",
    "is_active": true,
    "updated_at": "2025-06-29T11:25:00Z"
  }
}
```

### 15. Toggle App Status (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/apps/NewApp/toggle \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "App status toggled successfully",
  "data": {
    "app_id": "NewApp",
    "is_active": false,
    "updated_at": "2025-06-29T11:30:00Z"
  }
}
```

### 16. Delete App (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/apps/NewApp/delete \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "App deleted successfully",
  "data": {
    "deleted_app_id": "NewApp"
  }
}
```

### 17. Generate Invite Code (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/invite-codes \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "Invite code generated successfully",
  "data": {
    "id": 11,
    "code": "xYz789AbC123DefGhi456Jkl",
    "used": false,
    "created_at": "2025-06-29T11:35:00Z"
  }
}
```

### 18. Delete Invite Code (Admin)
```bash
curl -X POST http://localhost:44544/api/v1/admin/invite-codes/METzuzFY8KSudQOnvpS-Qw/delete \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

Response:
```json
{
  "code": 200,
  "message": "Invite code deleted successfully",
  "data": {
    "deleted_code": "METzuzFY8KSudQOnvpS-Qw"
  }
}
```

## Common Error Responses

### Authentication Errors

**Missing or Invalid Token:**
```json
{
  "code": 401,
  "message": "Missing or invalid authorization token",
  "data": null
}
```

**Insufficient Permissions:**
```json
{
  "code": 403,
  "message": "Insufficient permissions to access this resource",
  "data": null
}
```

### Validation Errors

**Invalid Input Data:**
```json
{
  "code": 400,
  "message": "Invalid input data",
  "data": {
    "errors": [
      "username: cannot be empty",
      "password: must be at least 6 characters long"
    ]
  }
}
```

**User Already Exists:**
```json
{
  "code": 409,
  "message": "User already exists",
  "data": null
}
```

**Invalid Invite Code:**
```json
{
  "code": 400,
  "message": "Invalid or expired invite code",
  "data": null
}
```

### Resource Not Found

**User Not Found:**
```json
{
  "code": 404,
  "message": "User not found",
  "data": null
}
```

**App Not Found:**
```json
{
  "code": 404,
  "message": "App not found",
  "data": null
}
```

### Server Errors

**Internal Server Error:**
```json
{
  "code": 500,
  "message": "Internal server error",
  "data": null
}
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

## Testing Flow

Here's a recommended testing flow to verify all API functionality:

1. **Start the server** and ensure it's running on the correct port
2. **Login as admin** using the default credentials
3. **Create a new invite code** if needed
4. **Register a new user** with the invite code
5. **Login as the new user** to get a user token
6. **Test user operations**: Get user info, list available apps
7. **Generate an NKey** for the user
8. **Validate the NKey** to ensure it works correctly
9. **Test admin operations**: List users, manage apps, view logs

## Environment Variables

The following environment variables can be used to configure the server:

- `PORT`: Server port (default: 44544)
- `DB_TYPE`: Database type (sqlite/postgres, default: sqlite)
- `DB_PATH`: SQLite database file path (default: ./tounetcore.db)
- `JWT_SECRET`: JWT signing secret (auto-generated if not provided)
- `PUSHDEER_ENDPOINT`: PushDeer API endpoint for notifications

## Rate Limiting

Some endpoints may have rate limiting in place:
- Registration: Limited to prevent abuse
- NKey generation: Limited per user per time period
- Login attempts: Limited to prevent brute force attacks

When rate limited, you'll receive a 429 status code:
```json
{
  "code": 429,
  "message": "Rate limit exceeded",
  "data": {
    "retry_after": 60
  }
}
```

⚠️ **Important**: Change the admin password after first login in production!
