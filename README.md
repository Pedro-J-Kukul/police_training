# Police Training API

A comprehensive REST API for managing police training programs, workshops, officers, and training sessions. Built with Go and PostgreSQL, featuring JWT authentication, role-based permissions, and comprehensive Swagger documentation.

## Features

- **JWT Authentication & Authorization** - Role-based access control with permissions
- **Officer Management** - Complete officer lifecycle with ranks, postings, and formations
- **Training Management** - Workshops, sessions, categories, and enrollment tracking
- **Progress Tracking** - Training completion, attendance, and progress monitoring
- **Multi-region Support** - Regional formations and postings management
- **Email Integration** - Password reset and notification system
- **Swagger Documentation** - Interactive API documentation
- **Rate Limiting** - Built-in request rate limiting
- **Database Migrations** - Automated schema management

## Tech Stack

- **Backend**: Go 1.21+
- **Database**: PostgreSQL
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Database Tools**: golang-migrate, Drizzle Kit
- **Email**: SMTP integration

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 13+
- Make (for using Makefile commands)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Pedro-J-Kukul/police_training.git
   cd police_training
   ```

2. **Set up environment variables**
   ```bash
   cp .envrc.example .envrc
   # Edit .envrc with your configuration
   ```

3. **Configure your database**
   Edit `.envrc` with your PostgreSQL connection details:
   ```bash
   DB_DSN="postgres://username:password@localhost/policetraining?sslmode=disable"
   ```

4. **Install dependencies and run migrations**
   ```bash
   go mod download
   make migrate/up
   ```

5. **Start the API server**
   ```bash
   make run/api
   ```

The API will be available at `http://localhost:4000`

## Makefile Commands

The project includes a comprehensive Makefile for common operations:

### Development
- `make run/api` - Start the API server
- `make run/api/win` - Start the API server on Windows (PowerShell)
- `make run/tests` - Run all tests

### Database Management
- `make migrate/up` - Run all pending migrations
- `make migrate/down` - Rollback the last migration
- `make migrate/up1` - Run one migration up
- `make migrate/down1` - Run one migration down
- `make migrate/create name=migration_name` - Create a new migration
- `make psql/login` - Connect to the database via psql
- `make psql/sudo` - Connect as postgres superuser

### Testing
- `make migrate/up/test` - Run migrations on test database
- `make migrate/down/test` - Rollback test database migrations
- `make migrate/reset/test` - Reset test database

## API Documentation

### Swagger UI

Once the server is running, access the interactive Swagger documentation at:

**http://localhost:4000/swagger/index.html**

The Swagger UI provides:
- Complete API endpoint documentation
- Interactive request testing
- Request/response examples
- Authentication testing

### Core Endpoints

#### Authentication & User Management
- `POST /v1/users` - Register a new user
- `POST /v1/tokens/authentication` - Login and get JWT token
- `POST /v1/tokens/password-reset` - Request password reset
- `PUT /v1/users/password-reset` - Reset password with token
- `GET /v1/me` - Get current user profile
- `GET /v1/users` - List all users (admin)
- `GET /v1/users/{id}` - Get user by ID
- `PATCH /v1/users/{id}` - Update user
- `DELETE /v1/users/{id}` - Soft delete user

#### Officer Management
- `POST /v1/officers` - Create new officer
- `GET /v1/officers` - List officers with filtering
- `GET /v1/officers/{id}` - Get officer details
- `GET /v1/officers/{id}/details` - Get officer with full details
- `PATCH /v1/officers/{id}` - Update officer
- `DELETE /v1/officers/{id}` - Delete officer

#### Organizational Structure
- `GET /v1/regions` - List regions
- `POST /v1/regions` - Create region
- `GET /v1/regions/{id}` - Get region details
- `PATCH /v1/regions/{id}` - Update region

- `GET /v1/formations` - List formations
- `POST /v1/formations` - Create formation
- `GET /v1/formations/{id}` - Get formation details
- `PATCH /v1/formations/{id}` - Update formation

- `GET /v1/postings` - List postings
- `POST /v1/postings` - Create posting
- `GET /v1/postings/{id}` - Get posting details
- `PATCH /v1/postings/{id}` - Update posting

- `GET /v1/ranks` - List ranks
- `POST /v1/ranks` - Create rank
- `GET /v1/ranks/{id}` - Get rank details
- `PATCH /v1/ranks/{id}` - Update rank

#### Training Management
- `GET /v1/workshops` - List workshops
- `POST /v1/workshops` - Create workshop
- `GET /v1/workshops/{id}` - Get workshop details
- `PATCH /v1/workshops/{id}` - Update workshop

- `GET /v1/training/categories` - List training categories
- `POST /v1/training/categories` - Create category
- `GET /v1/training/categories/{id}` - Get category details
- `PATCH /v1/training/categories/{id}` - Update category

- `GET /v1/training/types` - List training types
- `POST /v1/training/types` - Create type
- `GET /v1/training/types/{id}` - Get type details
- `PATCH /v1/training/types/{id}` - Update type

#### Training Sessions & Enrollment
- `GET /v1/training-sessions` - List training sessions
- `POST /v1/training-sessions` - Create training session
- `GET /v1/training-sessions/{id}` - Get session details
- `PATCH /v1/training-sessions/{id}` - Update session

- `GET /v1/training-enrollments` - List enrollments
- `POST /v1/training-enrollments` - Create enrollment
- `GET /v1/training-enrollments/{id}` - Get enrollment details
- `PATCH /v1/training-enrollments/{id}` - Update enrollment

#### Status Management
- `GET /v1/attendance/status` - List attendance statuses
- `POST /v1/attendance/status` - Create attendance status
- `GET /v1/attendance/status/{id}` - Get attendance status
- `PATCH /v1/attendance/status/{id}` - Update attendance status

- `GET /v1/enrollment/status` - List enrollment statuses
- `POST /v1/enrollment/status` - Create enrollment status
- `GET /v1/enrollment/status/{id}` - Get enrollment status
- `PATCH /v1/enrollment/status/{id}` - Update enrollment status

- `GET /v1/progress/status` - List progress statuses
- `POST /v1/progress/status` - Create progress status
- `GET /v1/progress/status/{id}` - Get progress status
- `PATCH /v1/progress/status/{id}` - Update progress status

## Sample cURL Requests

Here are some example requests you can try:

### 1. Health Check
```bash
curl -X GET http://localhost:4000/v1/healthcheck \
  -H "Content-Type: application/json"
```

### 2. Register a New User
```bash
curl -X POST http://localhost:4000/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@police.gov",
    "password": "securepassword123"
  }'
```

### 3. Login and Get Authentication Token
```bash
curl -X POST http://localhost:4000/v1/tokens/authentication \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@police.gov",
    "password": "securepassword123"
  }'
```

After getting your token from the login response, use it in subsequent requests:

### 4. Get Current User Profile (Authenticated)
```bash
curl -X GET http://localhost:4000/v1/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### 5. Create a Region
```bash
curl -X POST http://localhost:4000/v1/regions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "region": "Northern Region"
  }'
```

### 6. Create a Formation
```bash
curl -X POST http://localhost:4000/v1/formations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "formation": "Metropolitan Police Division",
    "region_id": 1
  }'
```

### 7. Create a Rank
```bash
curl -X POST http://localhost:4000/v1/ranks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "rank": "Constable",
    "code": "PC",
    "annual_training_hours_required": 40
  }'
```

### 8. Create an Officer
```bash
curl -X POST http://localhost:4000/v1/officers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "user_id": 1,
    "regulation_number": "12345",
    "rank_id": 1,
    "posting_id": 1,
    "formation_id": 1,
    "date_of_appointment": "2020-01-15"
  }'
```

### 9. Create a Training Category
```bash
curl -X POST http://localhost:4000/v1/training/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "name": "Firearms Training",
    "is_active": true
  }'
```

### 10. Create a Workshop
```bash
curl -X POST http://localhost:4000/v1/workshops \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "workshop_name": "Basic Firearms Safety",
    "category_id": 1,
    "training_type_id": 1,
    "credit_hours": 8,
    "description": "Comprehensive firearms safety and handling course",
    "objectives": "Learn proper firearm handling and safety protocols",
    "is_active": true
  }'
```

### 11. List Officers with Filtering
```bash
curl -X GET "http://localhost:4000/v1/officers?rank_id=1&page=1&page_size=10" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### 12. Create a Training Session
```bash
curl -X POST http://localhost:4000/v1/training-sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "workshop_id": 1,
    "facilitator_id": 1,
    "formation_id": 1,
    "session_date": "2025-11-01",
    "start_time": "09:00:00",
    "end_time": "17:00:00",
    "location": "Training Center A",
    "max_participants": 20,
    "training_status_id": 1
  }'
```

## Environment Configuration

Key environment variables in `.envrc`:

```bash
# Database
DB_DSN="postgres://postgres:postgres@localhost/policetraining?sslmode=disable"
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_MAX_IDLE_TIME="15m"

# Server
PORT="4000"
ENVIRONMENT="development"
API_VERSION="v1"

# CORS
CORS_ALLOWED_ORIGINS="http://localhost:3000,http://localhost:5173"

# Rate Limiting
RATE_LIMITER_ENABLED=true
RATE_LIMITER_RPS=5
RATE_LIMITER_BURST=10

# Email (SMTP)
SMTP_HOST="smtp.mailtrap.io"
SMTP_PORT=587
SMTP_USERNAME="your_username"
SMTP_PASSWORD="your_password"
SMTP_SENDER="Police Training <noreply@policetraining.gov>"
```

## Architecture

### Project Structure
```
├── cmd/api/           # Application entry point and handlers
├── internal/data/     # Data models and database logic
├── internal/mailer/   # Email functionality
├── migrations/        # Database migrations
├── docs/             # Swagger documentation
├── Makefile          # Build and development commands
└── .envrc.example    # Environment configuration template
```

### Key Components

- **Authentication**: JWT-based with refresh tokens
- **Authorization**: Role-based permissions system
- **Database**: PostgreSQL with migration support
- **Validation**: Comprehensive input validation
- **Logging**: Structured logging with slog
- **Error Handling**: Consistent error responses
- **Rate Limiting**: Configurable request rate limiting

### Database Models

The API manages several core entities:
- **Users** - System users with authentication
- **Officers** - Police officers with ranks and postings
- **Regions** - Geographical regions
- **Formations** - Police formations within regions
- **Postings** - Officer posting assignments
- **Ranks** - Police ranks with training requirements
- **Workshops** - Training workshops and courses
- **Training Sessions** - Scheduled training instances
- **Training Categories** - Workshop categorization
- **Training Types** - Types of training offered
- **Enrollments** - Officer training enrollments
- **Attendance Status** - Training attendance tracking
- **Enrollment Status** - Enrollment state management
- **Progress Status** - Training progress tracking

## Authentication & Authorization

The API uses JWT tokens for authentication. All endpoints except registration, login, and health check require authentication. Many endpoints also require specific permissions based on user roles.

### Getting Started with Authentication

1. Register a user account
2. Activate the account (if email verification is enabled)
3. Login to receive a JWT token
4. Include the token in the Authorization header: `Bearer YOUR_TOKEN`

## Filtering and Pagination

Most list endpoints support filtering and pagination:

- **page**: Page number (default: 1)
- **page_size**: Items per page (default: 20, max: 100)
- **sort**: Sort field (varies by endpoint)

Example: `GET /v1/officers?rank_id=1&page=2&page_size=10&sort=regulation_number`

## Development

### Running Tests
```bash
make run/tests
```

### Database Management
```bash
# Create a new migration
make migrate/create name=add_new_table

# Run migrations
make migrate/up

# Rollback migrations
make migrate/down
```

### API Documentation
The API uses Swagger annotations in the code. Documentation is automatically generated and served at `/swagger/index.html`.

## Production Deployment

1. Set `ENVIRONMENT=production` in your environment
2. Configure production database credentials
3. Set up proper SMTP credentials for email
4. Configure CORS origins for your frontend
5. Set appropriate rate limiting values
6. Use a reverse proxy (nginx) for SSL termination

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License. See LICENSE file for details.

## Support

For questions or support, please open an issue on GitHub or contact the development team.
