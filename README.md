# Maths Solution Backend

A modern Go backend for the AI-powered math solution platform with JWT authentication, PostgreSQL database, and AI service integration.

## Features

- **JWT Authentication**: Secure user registration and login
- **PostgreSQL Database**: Persistent storage with GORM ORM
- **AI Service Integration**: HTTP client for math solving service
- **CORS Support**: Cross-origin resource sharing for frontend
- **Docker Support**: Containerized deployment
- **Graceful Shutdown**: Proper signal handling
- **Error Handling**: Comprehensive error management

## API Endpoints

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login

### Math Solving (Protected)
- `POST /api/solve-math` - Solve math expression
- `GET /api/history` - Get user's solution history

### Health Check
- `GET /health` - Server health status

## Environment Variables

Copy `config.env.example` to `.env` and configure:

```bash
# Database
DATABASE_URL=postgres://username:password@localhost:5432/maths_solution_db?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY_HOURS=24

# Server
PORT=8000
GIN_MODE=debug

# AI Service
AI_SERVICE_URL=http://localhost:5000
AI_SERVICE_TIMEOUT=30

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
```

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker (optional)

### Local Development

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL:**
   ```bash
   # Using Docker
   docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15-alpine
   
   # Create database
   createdb maths_solution_db
   ```

3. **Configure environment:**
   ```bash
   cp config.env.example .env
   # Edit .env with your settings
   ```

4. **Run the server:**
   ```bash
   go run main.go
   ```

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down
```

## Project Structure

```
backend/
├── auth/           # JWT authentication
├── config/         # Configuration management
├── database/       # Database connection and migrations
├── handlers/       # HTTP request handlers
├── middleware/     # HTTP middleware (auth, CORS, error)
├── models/         # Data models and DTOs
├── routes/         # Route definitions
├── services/       # Business logic services
├── main.go         # Application entry point
├── Dockerfile      # Container configuration
└── docker-compose.yml
```

## Testing the API

### Register a user:
```bash
curl -X POST http://localhost:8000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","full_name":"John Doe","password":"password123"}'
```

### Login:
```bash
curl -X POST http://localhost:8000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Solve math (with JWT token):
```bash
curl -X POST http://localhost:8000/api/solve-math \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"expression":"integrate x^2 dx"}'
```

## Integration with Frontend

The backend is designed to work seamlessly with the Next.js frontend:

1. **CORS**: Configured to allow requests from `http://localhost:3000`
2. **API Routes**: Matches frontend expectations (`/solve-math`, `/history`)
3. **JWT Tokens**: Compatible with frontend authentication flow
4. **Error Handling**: Consistent error response format

## Production Deployment

1. **Set strong JWT secret**
2. **Use production database**
3. **Configure proper CORS origins**
4. **Set GIN_MODE=release**
5. **Use HTTPS in production**
6. **Set up proper logging and monitoring**
# AI Math Solver Backend (Golang)
