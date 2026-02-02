# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Hoshikuzu is a modern Content Management System (CMS) built with Go 1.25+ and Fiber web framework. It features a clean, modular architecture with comprehensive content management capabilities, user management, payment integration, and a Vue.js 3 frontend.

## Development Commands

### Backend Development

```bash
# Hot reload development (recommended)
air -c .air.toml

# Or use make
make dev          # Build and run in one step
make run          # Build frontend + backend, then run

# Build only
make build        # Build both frontend and backend
make build-backend # Build backend only
make build-frontend # Build frontend only

# Clean build artifacts
make clean

# Cross-platform builds
make build-all-platforms  # Build for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
make package              # Build and package all platforms
```

### Frontend Development

```bash
cd ui
pnpm install       # Install dependencies
pnpm run dev       # Start dev server on http://localhost:5173
pnpm run build-only # Build for production
```

### Database Operations

```bash
# Generate Ent ORM code (after schema changes)
go generate ./ent

# Run database migrations (automatic on first run)
go run main.go
```

## Architecture

### Project Structure

```
hoshikuzu/
├── cmd/              # Application entry point
│   └── server.go     # Fiber server initialization
├── internal/         # Private application code
│   ├── handlers/     # HTTP request handlers (by domain)
│   ├── services/     # Business logic layer (by domain)
│   ├── middleware/   # Fiber middleware
│   ├── router/       # Route registration
│   ├── database/     # Database client initialization
│   └── schedule/     # Scheduled job management
├── pkg/              # Public packages
├── ent/              # Ent ORM schema definitions and generated code
├── ui/               # Vue.js 3 + TypeScript frontend
├── views/            # HTML templates
├── assets/           # Static assets (includes frontend build)
└── main.go           # Application entry point
```

### Layered Architecture

The codebase follows a clean three-layer architecture:

1. **Handler Layer** (`internal/handlers/`): HTTP request/response handling
2. **Service Layer** (`internal/services/`): Business logic and domain operations
3. **Repository Layer** (Ent ORM): Data access and persistence

Each domain module (auth, user, post, pay_order, etc.) has its own handler and service packages.

### Key Technologies

- **Backend**: Fiber v2 web framework, Ent ORM, JWT auth
- **Frontend**: Vue 3, TypeScript, Vite, Naive UI, Pinia
- **Database**: Supports SQLite (default), MySQL, PostgreSQL
- **Storage**: Local, FTP, AWS S3
- **Build**: CGO enabled for cross-platform compilation

### Domain Modules

Major domain modules include:
- **Authentication**: JWT-based auth with email/phone verification, PAT management
- **User Management**: User CRUD, roles, permissions, profiles
- **Content Management**: Posts, essays, dynamic content models
- **Payment**: Multi-channel payment (ePay integration), order lifecycle
- **File Management**: Multi-storage backend support, image processing
- **Album/Gallery**: Photo management, categorization, optimization
- **Comment System**: Twikango integration

### API Structure

- RESTful API with version prefix: `/api/v1/`
- Resource-based routing conventions
- Swagger/OpenAPI documentation available
- JWT authentication middleware for protected routes

### Configuration

- **Config Framework**: Viper with TOML support
- **Environment Variables**: Supports `.env` files
- **Key Settings**: Database connection, server port, JWT secrets, payment gateways
- **Default Database**: SQLite with `file:ent?mode=memory&cache=shared&_fk=1`

### Frontend-Backend Integration

- Backend runs on port 13000 (configurable)
- Frontend dev server runs on port 5173
- Production builds embed frontend in `assets/` directory
- API communication via REST with JWT tokens

### Important Notes

- **CGO Required**: Must use `CGO_ENABLED=1` for builds (required by SQLite and Ent)
- **Cross-compilation**: Requires appropriate compilers for each platform (see Makefile)
- **Hot Reload**: Air excludes `ui/` directory - frontend must be developed separately
- **Ent Schema**: After modifying `ent/schema/*.go`, run `go generate ./ent` to regenerate code

### Testing

Testing infrastructure is in place but no explicit test files exist yet. When adding tests:
- Use `ent/enttest/` for database test utilities
- Services use interface-based design for easy mocking
- Fiber provides good support for integration testing
