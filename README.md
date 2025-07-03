# LinkedIn Connection Tracker

A Go-based application that monitors your LinkedIn connections' connections, tracks new networking opportunities, and automates outreach based on customizable criteria.

## 🎯 Project Overview

This system helps you expand your professional network by:

- Tracking your LinkedIn connections and their connections
- Identifying new networking opportunities automatically
- Executing automated actions (connection requests, notifications) based on your criteria
- Managing your professional network data in a structured way

## 🛠️ Tech Stack

### Backend

- **[Go](https://golang.org/)** - Primary programming language
- **[Gin](https://gin-gonic.com/)** - HTTP web framework
- **[sqlc](https://sqlc.dev/)** - Generate type-safe code from SQL
- **[PostgreSQL](https://www.postgresql.org/)** - Primary database
- **[asynq](https://github.com/hibiken/asynq)** - Background job processing

### Authentication & Security

- **Dual Authentication** - Google OAuth2 and traditional username/password
- **JWT** - Token-based authentication
- **[Secure](https://github.com/unrolled/secure)** - Security middleware

### Development Tools

- **[Swagger/OpenAPI](https://swagger.io/)** - API documentation (via `swag`)
- **[Testify](https://github.com/stretchr/testify)** - Testing framework
- **[Air](https://github.com/air-verse/air)** - Live reload for development
- **[Delve](https://github.com/go-delve/delve)** - Go debugger
- **[Revive](https://github.com/mgechev/revive)** - Go linter

## ✨ Features

### Core Functionality

- **Dual Authentication** - Google OAuth2 and traditional username/password authentication
- **Connection Tracking** - Monitor your connections and their networks
- **Profile Management** - Store and update LinkedIn profile data
- **Company Tracking** - Track employment history and current companies
- **Automated Outreach** - Send connection requests based on criteria

### Automation Engine

- **Custom Rules** - Define criteria for automatic actions
- **Background Processing** - Scheduled checks for new connections
- **Smart Filtering** - Filter by company, location, and other attributes
- **Template Messages** - Personalized connection request messages

### API & Documentation

- **RESTful API** - Complete REST API for all functionality
- **Swagger Documentation** - Auto-generated API documentation
- **Rate Limiting** - Protect against API abuse
- **Comprehensive Logging** - Structured logging for monitoring

## 🏗️ Architecture

### Database Schema

```
users → tracked_connections → linkedin_profiles
                           ↓
                    connection_relationships
                           ↓
                       companies ← profile_companies
                           ↑
                    automation_rules
```

### Key Components

- **Authentication Service** - Handle OAuth2 and JWT tokens
- **Profile Service** - Manage LinkedIn profile data
- **Connection Service** - Track and discover connections
- **Automation Engine** - Process rules and execute actions
- **Background Workers** - Scheduled tasks and job processing

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Redis (for asynq)
- Google Cloud Platform Account (for OAuth2, optional)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/linkedin-watcher.git
   cd linkedin-watcher
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Setup environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start the development server**
   ```bash
   make dev
   ```

### Development Commands

```bash
# Start development server with live reload
make dev

# Run tests
make test

# Run linter
make lint

# Generate API documentation
make docs

# Run security check
make security-check

# Build for production
make build
```

## 📚 API Documentation

Once the server is running, access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

### Authentication Endpoints

The application implements a comprehensive authentication system with the following endpoints:

#### Public Endpoints

- **`POST /auth/register`** - Register a new user with email and password

  - Requires: `email`, `name`, `password` (min 8 characters)
  - Returns: JWT access token, refresh token, and user info

- **`POST /auth/login`** - Login with email and password

  - Requires: `email`, `password`
  - Returns: JWT access token, refresh token, and user info

- **`POST /auth/refresh`** - Refresh access token using refresh token

  - Requires: `refresh_token`
  - Returns: New JWT access token and refresh token

- **`POST /auth/logout`** - Logout user (client should discard tokens)
  - No authentication required
  - Returns: Success message

#### Protected Endpoints

- **`POST /auth/change-password`** - Change user password
  - Requires: Bearer token, `current_password`, `new_password`
  - Returns: Success message

#### Authentication Flow

1. **Registration**: User registers with email/password → receives JWT tokens
2. **Login**: User logs in with credentials → receives JWT tokens
3. **API Access**: Include `Authorization: Bearer <token>` header for protected endpoints
4. **Token Refresh**: Use refresh token to get new access token when expired
5. **Password Change**: Authenticated users can change their password

#### JWT Token Details

- **Access Token**: Valid for 15 minutes
- **Refresh Token**: Valid for 7 days
- **Algorithm**: HS256
- **Claims**: User ID and email

### Other Key Endpoints

- **Connections**

  - `GET /api/v1/connections` - List tracked connections
  - `POST /api/v1/connections` - Add connection to track
  - `POST /api/v1/connections/{id}/check` - Check for new connections

- **Automation**
  - `GET /api/v1/rules` - List automation rules
  - `POST /api/v1/rules` - Create automation rule

## 🔧 Configuration

### Environment Variables

```bash
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=linkedin_tracker
DB_USER=postgres
DB_PASSWORD=password

# Redis
REDIS_URL=redis://localhost:6379

# Authentication
JWT_SECRET=your_jwt_secret
JWT_EXPIRY=24h

# Google OAuth2 (optional)
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/callback

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRY=24h

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1h
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test package
go test ./internal/services/...
```

### Test Structure

- Unit tests for all services and handlers
- Integration tests for database operations
- API endpoint tests
- Mock implementations for external services

## 🔒 Security

- **Token Encryption** - OAuth tokens and password hashes stored securely
- **Rate Limiting** - API endpoint protection
- **Input Validation** - All user inputs validated
- **CORS Protection** - Proper CORS policies
- **Security Headers** - Security middleware enabled

## 📊 Monitoring & Observability

- **Structured Logging** - JSON formatted logs
- **Health Checks** - `/health` endpoint
- **Metrics Collection** - Application metrics
- **Error Tracking** - Comprehensive error handling

## 🛣️ Roadmap

- [ ] Phase 1: Project setup and infrastructure
- [ ] Phase 2: Authentication and user management
- [ ] Phase 3: Core data models and database
- [ ] Phase 4: API development and documentation
- [ ] Phase 5: Background processing with asynq
- [ ] Phase 6: Business logic and automation engine
- [ ] Phase 7: Testing and quality assurance
- [ ] Phase 8: Deployment and monitoring

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This project is for educational purposes. Please ensure compliance with LinkedIn's Terms of Service and API usage policies. The authors are not responsible for any misuse of this software.
