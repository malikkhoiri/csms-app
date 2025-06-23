# CSMS (Charging Station Management System)

A backend system for managing EV charging stations, built with Go, following Clean Architecture and supporting OCPP 1.6 JSON protocol.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ HTTP Handler│  │ WS Handler  │  │   API Endpoints     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ Services    │  │ Use Cases   │  │ Business Logic      │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Models    │  │ Repository  │  │   Service           │ │
│  │             │  │ Interfaces  │  │   Interfaces        │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ PostgreSQL  │  │ Repository  │  │   Configuration    │ │
│  │   Database  │  │Implementat. │  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
csms/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── application/        # Application layer (use cases, services)
│   │   └── service/        # Business logic services (transaction, chargepoint, user, etc)
│   ├── config/             # Configuration management
│   ├── domain/             # Domain layer (core business logic, entities, interfaces)
│   ├── handler/            # Presentation layer (HTTP & WebSocket handlers)
│   │   ├── http/           # HTTP handlers (REST API)
│   │   └── ws/             # WebSocket handlers (OCPP 1.6)
│   ├── infrastructure/     # Infrastructure layer
│   │   ├── database/       # Database connection (PostgreSQL)
│   │   └── repository/     # Repository implementations (GORM)
│   └── server/             # Server setup and routing
├── config.yaml             # Main configuration file
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── main.go                 # Application entry point
└── README.md               # This file
```

## 🚀 Implemented Features

- **OCPP 1.6 JSON Protocol Support**
  - BootNotification
  - Heartbeat
  - Authorize
  - StartTransaction
  - StopTransaction
  - StatusNotification
  - MeterValues

- **Charge Point Management**
  - Registration via BootNotification
  - Status monitoring & notification
  - Connector management

- **User & RFID Management**
  - User CRUD
  - RFID (IDTag) management & authorization
  - Role-based access

- **Transaction Management**
  - Start/Stop transactions
  - Meter value tracking (real-time)
  - Energy consumption & cost calculation (configurable tariff)
  - Transaction history

- **Security**
  - JWT authentication for API access
  - Role-based access control

- **Configuration**
  - Centralized config via YAML (server, DB, JWT, tariff, etc

- **Logging & Monitoring**
  - Structured logging (planned)
  - Prometheus metrics endpoint (planned)

## 🛠️ Tech Stack

- **Language**: Go 1.23+
- **Web Framework**: Gin
- **WebSocket**: Gorilla WebSocket
- **Database**: PostgreSQL (GORM ORM)
- **Configuration**: Viper
- **Logging**: Logrus (planned)
- **Authentication**: JWT
- **Testing**: Go test
- **Monitoring**: Prometheus (planned)

## 🔄 In Progress / Planned

- Remote commands (Reset, UnlockConnector, etc)
- Advanced monitoring dashboard
- Notification system (email/SMS)
- Billing integration
- Multi-tenant support
- Docker/Kubernetes deployment

## 📋 Configuration Example

```yaml
server:
  port: "8080"
  mode: "debug"

database:
  host: "localhost"
  port: 5432
  user: "csms_users"
  password: "password"
  dbname: "csms"
  sslmode: "disable"

jwt:
  secret: "your-super-secret-jwt-key-change-this-in-production"
  expiration: "24h"
  issuer: "csms"

tariff:
  price_per_kwh: 2500
```

## 🔌 API Endpoints (Core)

- `GET /health` - Server health status
- `GET /api/v1/status` - Server status
- `GET /api/v1/connections` - Active connections
- `GET /ocpp/{chargePointID}` - OCPP WebSocket endpoint

## 🧪 Virtual Charge Point Simulation

To simulate a virtual charge point (OCPP 1.6/2.0.1), you can use:

- [solidstudiosh/ocpp-virtual-charge-point](https://github.com/solidstudiosh/ocpp-virtual-charge-point) — A simple, configurable, terminal-based OCPP Charging Station simulator written in Node.js with schema validation.

This simulator is very useful for testing your CSMS backend integration without requiring physical charging station hardware.

## 📄 License

This project is licensed under the MIT License.

## 🔮 Roadmap

- [ ] Testing
- [ ] Remote commands (Reset, UnlockConnector, etc.)
- [ ] Advanced monitoring dashboard
- [ ] Email/SMS notifications
- [ ] Mobile app API
- [ ] Billing integration
- [ ] Multi-tenant support
- [ ] Docker deployment
- [ ] Kubernetes manifests 