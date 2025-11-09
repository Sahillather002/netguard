# 🛡️ NetGuard - Enterprise Security Platform

Enterprise-grade network security monitoring and threat detection platform built with microservices architecture.

## 🚀 Quick Start

```bash
# Install dependencies
./setup-dependencies.ps1

# Start all services
./start-local.ps1
```

**Access:** http://localhost:3000  
**Login:** admin@netguard.com / password

---

## 🏗️ Architecture

### Backend Services
- **API Gateway** (Go) - Port 8080
- **Auth Service** (Go) - Port 8081  
- **Threat Detector** (Python) - Port 8082
- **Network Monitor** (Python) - Port 8083
- **Firewall Service** (Python) - Port 8084

### Frontend
- **Next.js 15** - Port 3000
- 10 pages with dark mode and responsive design

## ✨ Features

- JWT Authentication
- Real-time threat detection
- Network device monitoring (25+ devices)
- Firewall rule management (15+ rules)
- Attack flow visualization
- User management with RBAC

---

## 📋 Prerequisites

- Node.js 20+
- Go 1.21+
- Python 3.11+

## 🛠️ Installation

```bash
# Automated setup
./setup-dependencies.ps1
```

This installs all dependencies for frontend, backend services, and Python packages.

## 📡 API Endpoints

- `POST /api/v1/auth/login` - Authentication
- `GET /api/v1/threats` - List threats
- `GET /api/v1/network/devices` - Network devices
- `GET /api/v1/firewall/rules` - Firewall rules
- `GET /api/v1/users` - User management

## 🧪 Development

```bash
# Test API
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@netguard.com","password":"password"}'
```

## 🔧 Tech Stack

- **Frontend:** Next.js 15, TypeScript, Tailwind CSS
- **Backend:** Go (Fiber), Python (FastAPI)
- **Auth:** JWT with bcrypt
- **Features:** Dark mode, RBAC, Real-time monitoring

## 📝 License

MIT License

---

**Built with Next.js, Go, and Python**
