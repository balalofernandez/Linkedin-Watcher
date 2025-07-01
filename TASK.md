# LinkedIn Connection Tracker - Project Task Definition

## Project Overview

Build a LinkedIn connection tracking system that monitors your connections' connections, identifies new connections based on criteria, and automates networking actions.

## Tech Stack

- **Go** - Primary language
- **sqlc** - Database query generation
- **Gin** - REST API framework
- **asynq** - Background job processing
- **Swagger** - API documentation (generated with `swag init`)
- **Testify** - Testing framework with assertions and mocks
- **Revive** - Go linter
- **Secure** - Go security checker
- **Air** - Live reload for development
- **Delve** - Debugger

## Core Features

1. **OAuth2 Authentication** - Google OAuth2 with JWT tokens
2. **Connection Tracking** - Monitor LinkedIn connections and their connections
3. **Profile Management** - Store and update LinkedIn profile data
4. **Company Tracking** - Track employment history and current companies
5. **Automated Actions** - Execute tasks when new connections match criteria
6. **Background Processing** - Schedule connection checks with asynq

## Database Design

### Tables Required:

#### `users`

- `id` (UUID, PRIMARY KEY)
- `google_id` (VARCHAR, UNIQUE)
- `email` (VARCHAR, UNIQUE)
- `name` (VARCHAR)
- `access_token` (TEXT, encrypted)
- `refresh_token` (TEXT, encrypted)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

#### `companies`

- `id` (UUID, PRIMARY KEY)
- `name` (VARCHAR, UNIQUE)
- `linkedin_url` (VARCHAR, UNIQUE)
- `industry` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

#### `linkedin_profiles`

- `id` (UUID, PRIMARY KEY)
- `linkedin_url` (VARCHAR, UNIQUE) - Format: `linkedin.com/in/*`
- `name` (VARCHAR)
- `location` (VARCHAR)
- `current_company_id` (UUID, FK to companies)
- `headline` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

#### `profile_companies` (Employment History)

- `id` (UUID, PRIMARY KEY)
- `profile_id` (UUID, FK to linkedin_profiles)
- `company_id` (UUID, FK to companies)
- `position` (VARCHAR)
- `start_date` (DATE)
- `end_date` (DATE, nullable)
- `is_current` (BOOLEAN)
- `created_at` (TIMESTAMP)

#### `tracked_connections`

- `id` (UUID, PRIMARY KEY)
- `user_id` (UUID, FK to users)
- `profile_id` (UUID, FK to linkedin_profiles)
- `created_at` (TIMESTAMP)
- `last_checked_at` (TIMESTAMP)

#### `connection_relationships`

- `id` (UUID, PRIMARY KEY)
- `profile_a_id` (UUID, FK to linkedin_profiles)
- `profile_b_id` (UUID, FK to linkedin_profiles)
- `discovered_at` (TIMESTAMP)
- `discovered_by_user_id` (UUID, FK to users)

#### `automation_rules`

- `id` (UUID, PRIMARY KEY)
- `user_id` (UUID, FK to users)
- `name` (VARCHAR)
- `company_filter` (VARCHAR, nullable)
- `location_filter` (VARCHAR, nullable)
- `action_type` (ENUM: 'send_connection_request', 'save_profile', 'notify')
- `message_template` (TEXT, nullable)
- `is_active` (BOOLEAN)
- `created_at` (TIMESTAMP)

## API Endpoints

### Authentication

- `POST /auth/google` - Initiate Google OAuth2
- `POST /auth/callback` - Handle OAuth2 callback
- `POST /auth/refresh` - Refresh JWT token
- `POST /auth/logout` - Logout user

### Profile Management

- `GET /api/v1/profile` - Get user profile
- `PUT /api/v1/profile` - Update user profile

### Connection Tracking

- `GET /api/v1/connections` - List tracked connections
- `POST /api/v1/connections` - Add connection to track
- `DELETE /api/v1/connections/{id}` - Stop tracking connection
- `POST /api/v1/connections/{id}/check` - Manually trigger connection check

### Companies

- `GET /api/v1/companies` - List companies
- `POST /api/v1/companies` - Add new company
- `GET /api/v1/companies/{id}` - Get company details

### Automation Rules

- `GET /api/v1/rules` - List automation rules
- `POST /api/v1/rules` - Create automation rule
- `PUT /api/v1/rules/{id}` - Update automation rule
- `DELETE /api/v1/rules/{id}` - Delete automation rule

### Background Jobs

- `GET /api/v1/jobs` - List recent jobs
- `POST /api/v1/jobs/schedule-check` - Schedule connection check

## Implementation Steps

### Phase 1: Project Setup & Infrastructure

1. **Initialize Go Module**

   - Create `go.mod` with required dependencies
   - Setup project structure (`cmd/`, `internal/`, `pkg/`, `migrations/`)

2. **Database Setup**

   - Create PostgreSQL migration files
   - Setup sqlc configuration (`sqlc.yaml`)
   - Generate database models and queries

3. **Development Environment**

   - Configure Air for live reload (`.air.toml`)
   - Setup Delve debugger configuration
   - Configure linting with Revive (`.revive.toml`)

4. **Security & Dependencies**
   - Integrate Secure Go security checker
   - Setup dependency management and security scanning

### Phase 2: Authentication & User Management

1. **OAuth2 Integration**

   - Implement Google OAuth2 flow
   - Create JWT token generation and validation
   - Setup middleware for authentication

2. **User Management**
   - Create user registration/login endpoints
   - Implement token refresh mechanism
   - Add user profile management

### Phase 3: Core Data Models

1. **Database Schema**

   - Implement all required tables with migrations
   - Setup foreign key constraints and indexes
   - Create sqlc queries for CRUD operations

2. **Profile Management**
   - LinkedIn profile validation (URL format)
   - Company management system
   - Employment history tracking

### Phase 4: API Development

1. **Gin Router Setup**

   - Create router configuration
   - Setup middleware (CORS, logging, auth)
   - Implement rate limiting

2. **REST Endpoints**

   - Implement all CRUD endpoints
   - Add request validation
   - Error handling and response formatting

3. **Swagger Documentation**
   - Add swagger annotations to handlers
   - Configure swagger middleware
   - Generate API documentation with `swag init`

### Phase 5: Background Processing

1. **Asynq Integration**

   - Setup asynq server and client
   - Create job types for connection checking
   - Implement retry logic and error handling

2. **Scheduled Tasks**
   - Create periodic connection check jobs
   - Implement new connection detection logic
   - Add automation rule processing

### Phase 6: Business Logic

1. **Connection Discovery**

   - LinkedIn scraping simulation (mock data for development)
   - New connection detection algorithm
   - Graph relationship management

2. **Automation Engine**
   - Rule matching system
   - Action execution (connection requests, notifications)
   - Template processing for personalized messages

### Phase 7: Testing & Quality

1. **Unit Testing**

   - Write tests with Testify
   - Create mocks for external dependencies
   - Achieve >80% code coverage

2. **Integration Testing**
   - Database integration tests
   - API endpoint testing
   - Background job testing

### Phase 8: Deployment & Monitoring

1. **Production Setup**

   - Docker containerization
   - Environment configuration
   - Database connection pooling

2. **Monitoring & Logging**
   - Structured logging
   - Health check endpoints
   - Metrics collection

## Development Notes

### LinkedIn URL Validation

```go
// Profile URL must match: linkedin.com/in/{username}
var linkedinProfileRegex = regexp.MustCompile(`^https?://(www\.)?linkedin\.com/in/[a-zA-Z0-9-]+/?$`)
```

### Security Considerations

- Encrypt stored Google OAuth tokens
- Implement rate limiting for API calls
- Use HTTPS only
- Validate all user inputs
- Implement proper CORS policies

### Background Job Priorities

1. **High Priority**: User-triggered connection checks
2. **Medium Priority**: Scheduled connection updates
3. **Low Priority**: Historical data cleanup

## Success Metrics

- Successful OAuth2 integration with Google
- Ability to track and store connection data
- Working automation rules that trigger actions
- Background jobs processing without errors
- Comprehensive API documentation
- > 80% test coverage
- Security scan passing without critical issues
