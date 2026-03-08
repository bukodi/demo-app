# Demo App Architecture

## Overview

Demo App is a Go-based web application demonstrating modern architectural patterns including:
- Plugin-based dependency injection
- Multiple storage backend support (GORM, DynamoDB)
- Authentication and authorization
- Embedded static web UI (Vue 3 + TypeScript)
- AWS Lambda deployment support

**Version:** 0.0.2
**Build:** Go 1.24

## Project Structure

```
demo-app/
├── cmd/
│   ├── demo/              # Main CLI application
│   └── demo-aws-lambda/   # AWS Lambda handler
├── pkg/
│   ├── authn/             # Authentication layer
│   ├── authz/             # Authorization layer
│   ├── config/            # Configuration management (koanf-based)
│   ├── data/              # Data access layer
│   │   ├── dyndb/         # DynamoDB provider
│   │   └── gormdb/        # GORM provider
│   ├── domain/            # Business domain logic
│   │   └── user/          # User domain
│   ├── init_by_tags/      # Build tag-based initialization
│   ├── rtenv/             # Runtime environment
│   ├── server/            # HTTP server with plugin system
│   ├── util/              # Utility packages
│   └── webui/             # Static web UI handler
└── _webui/                # Vue 3 + TypeScript frontend
```

## Architectural Patterns

### 1. Plugin System

The application uses a plugin-based architecture where components register themselves during initialization:

**Server Plugin Registration** (`pkg/server/handler_spi.go`):
```go
server.RegisterPlugin(name string, srvPlugin ServerPluginInit)
```

**Active Plugins:**
- `authn` - Authentication middleware and `/api/v1/authorize` endpoint
- `user` - User management endpoints (`GET /user/list`, `POST /user`)
- `webui` - Static file serving for Vue frontend
- `version` - Version information endpoint

### 2. Dependency Hierarchy

```
Util (lowest)
  ↓
rtenv (runtime environment)
  ↓
config (configuration)
  ↓
Authn (authentication)
  ↓
Authz (authorization)
  ↓
Server (HTTP listeners)
  ↓
Data (persistence)
  ↓
Domain (business logic)
```

### 3. Service Provider Interface (SPI)

**Generic SPI Pattern** (`pkg/util/generic_spi.go`):
Components use a generic SPI pattern for swappable implementations:

```go
spi = util.NewGenericSpi[Store]()
spi.RegisterProvider(providerName, initFn)
```

**User Store Providers:**
- `gorm` - Relational database via GORM (MySQL, SQLite)
- `dynamodb` - AWS DynamoDB NoSQL

### 4. Identity Provider System

The authentication layer supports pluggable identity providers (`pkg/authn/idp_spi.go`):

```go
authn.RegisterIdentityProvider(name string, idp IdentityProvider)
```

**Registered IDP:**
- `built-in` - Local user authentication via password hashing

## Domain: User

### Data Model

```go
type User struct {
    Email        string `json:"email"`
    PasswordHash string `json:"password_hash"`
    Role         string `json:"role"`
}
```

**Roles:**
- `end_user` - Standard user
- `admin` - Administrative user

### API Endpoints

#### `POST /api/v1/user`
Create a new user.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

**Response:** 201 Created

#### `GET /api/v1/user/list`
List all users (requires admin role).

**Response:**
```json
[
  {
    "email": "user@example.com",
    "role": "end_user"
  }
]
```

### Storage Implementations

**GORM Store** (`pkg/domain/user/store_gorm.go`):
- Auto-migration support
- Transaction handling
- Supports MySQL, SQLite

**DynamoDB Store** (`pkg/domain/user/store_dynamodb.go`):
- Table auto-creation
- Pay-per-request billing
- Email as hash key

## Authentication & Authorization

### Authentication Flow

1. Client sends credentials to `POST /api/v1/authorize`
2. Identity provider validates credentials
3. Cookie-based session created via middleware
4. Subsequent requests authenticated via cookie

**Middleware:** `authn.CookieMiddleware`

### Authorization

Simple role-based access control:
```go
authz.CheckHttpHandler(handler, relation, object)
```

Currently implements basic role checking (admin vs end_user).

## Commands

### `serve` (default)
Start HTTP server on `localhost:8080`.

```bash
./demo serve
# or simply
./demo
```

### `version`
Display version and build information.

```bash
./demo version
```

### `cdn-export`
Export static web UI as `webui_static.zip` for CDN deployment.

```bash
./demo cdn-export
```

## Web UI

**Technology Stack:**
- Vue 3 with Composition API
- TypeScript
- Vite build system

**Features:**
- User list with refresh
- User creation form
- Tabbed navigation
- Responsive design

**Development:**
```bash
cd _webui
npm install
npm run dev
```

**Build:**
```bash
cd _webui
npm run build
```

Static assets are embedded into the Go binary via `//go:embed`.

## Build Tags

**`no_user` build tag:** Exclude user domain and web UI:
```bash
go build -tags no_user ./cmd/demo
```

## Configuration

Uses `koanf` library for configuration management:
- JSON file support
- YAML file support
- Hierarchical configuration keys

## Testing Support

**Test Utilities:**
- Docker-based Keycloak (`pkg/util/tkeycloak/`)
- Mock OAuth2 servers (`pkg/util/toauthsrv/`)
- Selenium integration (`pkg/util/selenium_test.go`)
- Fake users for testing (`pkg/domain/user/fakeuser_test.go`)

## AWS Lambda Support

**Lambda Handler** (`cmd/demo-aws-lambda/main.go`):
Deploy as AWS Lambda function with API Gateway proxy integration.

## Key Features

1. **Version Embedding** - Git commit SHA embedded at build time
2. **Static Assets** - Web UI embedded in binary or exportable for CDN
3. **Multi-Backend** - Supports SQL (via GORM) and NoSQL (DynamoDB)
4. **Plugin Architecture** - Clean separation of concerns via plugin registration
5. **Build Tags** - Conditional compilation for different deployment scenarios
6. **SPA Support** - Frontend routing with fallback to `index.html`

## Security

- Password hashing via `bcrypt` (`pkg/domain/user/util.go`)
- JWT support (`pkg/authn/jwt_test.go`)
- Cookie-based session management
- Role-based authorization

## Dependencies

**Key Libraries:**
- `cobra` - CLI framework
- `gorm` - ORM
- `aws-sdk-go-v2` - AWS integration
- `koanf` - Configuration
- `jwt` - Token handling

**Frontend:**
- Vue 3.4+
- Vite 5.4+
- TypeScript 5.2+

## Extension Points

To add new features:

1. **New Domain:** Create package under `pkg/domain/`
2. **New API Endpoints:** Register via `server.RegisterPlugin()`
3. **New Storage Backend:** Implement Store interface and register provider
4. **New Authentication Method:** Implement `authn.IdentityProvider`
5. **New UI Components:** Add to `_webui/src/components/`

## Architecture Validation

Uses `arch-go` for architecture rule enforcement (`arch-go.yml`):
- Dependency rules (e.g., authn can only depend on server)
- Content rules (no interfaces in `*.model` packages)
- Function complexity rules
- Naming conventions
