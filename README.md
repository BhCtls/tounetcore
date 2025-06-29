# TouNetCore - User Management & Authorization System

TouNetCore is a Go-based user management and authorization key distribution system designed for controlling access to various applications through a centralized authentication service.

## Features

### Core Functionality
- **User Management**: Registration with invite codes, login, and profile management
- **Role-based Access Control**: Three user roles (Admin, User, DisabledUser)
- **NKey Authorization**: Temporary 15-minute authorization keys for app access
- **Application Management**: Configurable applications with different permission levels
- **Audit Logging**: Comprehensive system operation tracking
- **Push Notifications**: PushDeer integration for key delivery

### Security Features
- **Password Security**: bcrypt hashing for password storage
- **JWT Authentication**: Secure session management
- **Time-limited Keys**: NKeys expire automatically after 15 minutes
- **Permission Validation**: Multi-level authorization checks
- **Encrypted Storage**: Secret keys are encrypted in database

## Quick Start

### Prerequisites
- Go 1.21 or higher
- SQLite (default) or PostgreSQL

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd tounetcore
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default.

### Configuration

Edit the `.env` file to configure the application:

```env
# Application Environment
ENVIRONMENT=development

# Database Configuration
DATABASE_URL=sqlite://./tounetcore.db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=24h

# NKey Configuration
NKEY_EXPIRATION=15m

# PushDeer Configuration
PUSHDEER_API=https://api2.pushdeer.com/message/push

# Server Configuration
PORT=8080
```

## API Documentation

### Authentication Endpoints

#### Register User
```http
POST /api/v1/register
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securePass123!",
  "phone": "13800138000",
  "pushdeer_token": "PUSHDEER_XXXXXXXXX",
  "invite_code": "INVITE123456"
}
```

#### Login
```http
POST /api/v1/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securePass123!"
}
```

### User Endpoints (Require Authentication)

#### Get User Information
```http
GET /api/v1/user/me
Authorization: Bearer <jwt_token>
```

#### Update User Profile
```http
PUT /api/v1/user/me
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "phone": "13900139000",
  "pushdeer_token": "NEW_PUSHDEER_TOKEN"
}
```

#### Get Allowed Applications
```http
GET /api/v1/user/apps
Authorization: Bearer <jwt_token>
```

### NKey Endpoints

#### Generate NKey
```http
POST /api/v1/nkey/generate
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "username": ["Jack"],
  "app_ids": ["Approval", "Edit"]
}
```

#### Validate NKey
```http
POST /api/v1/nkey/validate
Content-Type: application/json

{
  "nkey": "TOUNET_NKEY_XXXXXX",
  "app_id": "Approval"
}
```

### Admin Endpoints (Require Admin Role)

#### Create User
```http
POST /api/v1/admin/users
Authorization: Bearer <admin_jwt_token>
Content-Type: application/json

{
  "username": "new_user",
  "password": "pass123",
  "status": "user",
  "phone": "13800138000",
  "pushdeer_token": "PUSHDEER_XXX"
}
```

#### List Users
```http
GET /api/v1/admin/users?page=1&size=10
Authorization: Bearer <admin_jwt_token>
```

#### Create Application
```http
POST /api/v1/admin/apps
Authorization: Bearer <admin_jwt_token>
Content-Type: application/json

{
  "app_id": "NewApp",
  "name": "New Application",
  "description": "Application description",
  "required_permission_level": "user",
  "is_active": true
}
```

#### Generate Invite Code
```http
POST /api/v1/admin/invite-codes
Authorization: Bearer <admin_jwt_token>
```

## Database Schema

### Tables

1. **users**: User accounts and profiles
2. **invite_codes**: Registration invitation codes
3. **nkeys**: Temporary authorization keys
4. **apps**: Application definitions
5. **user_allowed_apps**: User-specific app permissions
6. **audit_logs**: System operation logs

### Pre-configured Applications

The system comes with pre-configured applications:
- `searchall`: Global search functionality
- `segaasstes`: Sega assets management
- `dxprender`: DXP rendering service
- `CardPreview`: Card preview functionality
- `livecontent_basic`: Basic live content access
- `livecontent_admin`: Administrative live content access (admin only)

## Development

### Project Structure
```
tounetcore/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/             # HTTP routes
│   ├── auth/            # Authentication utilities
│   ├── config/          # Configuration management
│   ├── database/        # Database operations
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   └── models/          # Database models
├── migrations/          # Database migrations
└── .github/            # GitHub configuration
```

### Building for Production

```bash
# Build binary
go build -o tounetcore cmd/server/main.go

# Run in production
ENVIRONMENT=production ./tounetcore
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Security Considerations

1. **Change Default Secrets**: Update JWT_SECRET and app secret keys in production
2. **Use HTTPS**: Always use HTTPS in production environments
3. **Database Security**: Use proper database credentials and connections
4. **Rate Limiting**: Implement rate limiting for API endpoints
5. **Input Validation**: All inputs are validated and sanitized
6. **Audit Logging**: All sensitive operations are logged

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## Support

For support and questions, please create an issue in the repository.



## 🎉 TouNetCore Project Successfully Created!

I've successfully created a comprehensive Go-based user management and authorization key distribution system called **TouNetCore**. Here's what has been implemented:

### ✅ **Core Features Completed**

1. **🏗️ Project Structure**
   - Clean architecture with cmd, internal, and pkg directories
   - Proper Go module setup with all dependencies

2. **🗄️ Database Layer**
   - GORM-based models for all entities (User, InviteCode, NKey, App, UserAllowedApp, AuditLog)
   - SQLite and PostgreSQL support
   - Auto-migration system
   - Pre-configured applications seeding

3. **🔐 Authentication & Authorization**
   - JWT-based authentication
   - bcrypt password hashing
   - Role-based access control (Admin, User, DisabledUser)
   - 15-minute temporary NKey system
   - Invite code system for registration

4. **🌐 REST API**
   - Complete API endpoints for all functionality
   - Gin web framework with middleware
   - CORS support
   - Proper error handling and JSON responses

5. **👥 User Management**
   - User registration with invite codes
   - Login/logout functionality
   - Profile management
   - Role-based permissions

6. **🔑 NKey System**
   - Generate temporary authorization keys
   - Validate keys with app-specific permissions
   - Time-based expiration (15 minutes)
   - First-use tracking

7. **⚙️ Admin Features**
   - User management (create, list, modify)
   - Application management
   - Invite code generation
   - Audit log viewing

8. **🛠️ Development Tools**
   - Makefile for common tasks
   - Seeding commands for initial data
   - VS Code tasks configuration
   - Comprehensive documentation

### 🚀 **Getting Started**

1. **Initialize the project:**
   ```bash
   make init
   ```

2. **Start the server:**
   ```bash
   make run
   # or
   go run cmd/server/main.go
   ```

3. **Access the API:**
   - Server runs on http://localhost:8081
   - Health check: `GET /health`
   - Admin login: `admin` / `admin123`

### 📚 **Key Resources**

- **README.md**: Complete project documentation
- **API_EXAMPLES.md**: API testing examples with curl commands
- **Makefile**: Development commands
- **.env.example**: Environment configuration template
- **Copilot instructions**: Workspace-specific AI coding guidelines

### 🔧 **Pre-configured Apps**

The system comes with 6 pre-configured applications:
- `searchall`, `segaasstes`, `dxprender`, `CardPreview`
- `livecontent_basic`, `livecontent_admin` (admin-only)

### 📋 **Available Commands**

```bash
make help        # Show all available commands
make build       # Build the application
make run         # Run in development mode
make test        # Run tests
make seed-apps   # Seed default applications
make seed-admin  # Create admin user
make seed-invite # Generate invite codes
```

The project is now fully functional and ready for development! You can start the server, test the APIs, and begin customizing it according to your specific needs. The architecture is designed to be scalable and maintainable, following Go best practices.