<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# TouNetCore - User Management & Authorization System

This is a Go-based user management and authorization key distribution system built with:

## Technology Stack
- **Language**: Go
- **Web Framework**: Gin
- **Database**: GORM (supports SQLite and PostgreSQL)
- **Authentication**: JWT tokens
- **Authorization**: Custom NKey system

## Project Structure
- `cmd/server/` - Application entry point
- `internal/` - Private application code
  - `api/` - HTTP route definitions
  - `auth/` - Authentication and cryptographic utilities
  - `config/` - Configuration management
  - `database/` - Database initialization and migrations
  - `handlers/` - HTTP request handlers
  - `middleware/` - HTTP middleware
  - `models/` - Database models

## Key Features
1. **User Management**: Registration with invite codes, login, profile management
2. **Role-based Access**: Admin, User, DisabledUser roles
3. **NKey System**: 15-minute temporary authorization keys for app access
4. **Application Management**: Configurable apps with permission levels
5. **Audit Logging**: Track system operations
6. **Push Notifications**: PushDeer integration for key delivery

## Database Models
- `User`: User accounts with roles and permissions
- `InviteCode`: Registration invitation system
- `NKey`: Temporary authorization keys
- `App`: Protected applications
- `UserAllowedApp`: User-specific app permissions
- `AuditLog`: System audit trail

## API Endpoints
- Public: `/register`, `/login`, `/nkey/validate`
- User: `/user/me`, `/user/apps`, `/nkey/generate`
- Admin: `/admin/users`, `/admin/apps`, `/admin/logs`

## Security Considerations
- Passwords are bcrypt hashed
- JWT tokens for session management
- NKeys expire after 15 minutes
- Role-based authorization
- Secret keys for app-specific encryption

When working with this codebase:
- Follow Go naming conventions
- Use GORM for database operations
- Implement proper error handling
- Add audit logging for sensitive operations
- Validate user permissions before granting access
- Use environment variables for configuration
