## Overview
Chirpy is a web service that provides a RESTful API for managing users, chirps (messages), and authentication. It runs on port 8080 by default and uses PostgreSQL for data storage.

## Configuration
The service requires the following environment variables in a .env file:
- `DB_URL`: PostgreSQL database connection URL
- `PLATFORM`: Environment setting ("dev" or production)
- `JWT_SECRET`: Secret key for JWT token generation
- `POLKA_KEY`: API key for Polka webhook authentication

## API Endpoints

### Health Check
- `GET /api/healthz` - Check API health status

### Authentication
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh authentication token
- `POST /api/revoke` - Revoke authentication token

### Users
- `POST /api/users` - Create new user
- `PUT /api/users` - Update user information

### Chirps
- `POST /api/chirps` - Create new chirp
- `GET /api/chirps` - List all chirps
- `GET /api/chirps/{chirpID}` - Get specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete specific chirp

### Admin
- `GET /admin/metrics` - View site visit metrics
- `POST /admin/reset` - Reset application state
- `POST /api/polka/webhooks` - Handle Polka webhook events (requires Polka API key)

### Static Files
- `/app/*` - Serves static files from root directory with visit tracking

## Security
The application includes:
- JWT-based authentication
- API key validation for webhooks
- Production/development environment separation
- Database connection security

The server implements standard HTTP status codes and content-type headers for responses.