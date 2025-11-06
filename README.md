# Registration & Payment Service

Microservice untuk mengelola pendaftaran dan pembayaran event seminar/daurah.

## 🚀 Features

- CRUD pendaftaran peserta event
- Upload dan verifikasi bukti pembayaran
- Integrasi Kafka untuk event-driven architecture
- PostgreSQL database dengan migrations
- RESTful API dengan Fiber framework
- Docker containerization

## 🛠️ Quick Start

```bash
# Clone dan masuk ke directory
git clone <repository-url>
cd registration-payment-service

# Start services
docker compose up -d postgres kafka
docker compose up -d --build app

# Test
curl http://localhost:3003/healthz
```

## 📚 API Endpoints

```bash
# Registration
POST   /api/v1/registrations           # Create
GET    /api/v1/registrations           # List
GET    /api/v1/registrations/:id       # Detail
PUT    /api/v1/registrations/:id       # Update
POST   /api/v1/registrations/:id/cancel # Cancel

# Payment
POST   /api/v1/registrations/:id/payment        # Upload proof
GET    /api/v1/registrations/:id/payment        # Get info
PATCH  /api/v1/registrations/:id/payment/verify # Verify (admin)
```

## 🔧 Configuration

| Port | Service | Description |
|------|---------|-------------|
| 3003 | App | HTTP API |
| 5435 | PostgreSQL | Database |
| 19093 | Kafka | Message broker |

## 🧪 Example Usage

```bash
# Create registration
curl -X POST http://localhost:3003/api/v1/registrations \
  -H "Content-Type: application/json" \
  -d '{
    "event_id":"00000000-0000-0000-0000-000000000001",
    "full_name":"John Doe",
    "gender":"male",
    "phone":"08123456789",
    "email":"john@example.com"
  }'
```

## 📨 Kafka Events

**Published**: `registration.created`, `registration.confirmed`, `payment.uploaded`, `payment.verified`

**Consumed**: `event.status.changed`

## 🔄 Development

```bash
# Local development
go mod download
go run cmd/server/main.go

# Database migration
cat migrations/0001_init_registrations.up.sql | \
docker exec -i regpay-postgres psql -U regpay -d regpay_db
```

---

**Tech Stack**: Go, Fiber, PostgreSQL, Kafka, Docker
