Employee Payroll System
A comprehensive Go-based payroll management system that handles employee attendance, overtime, reimbursements, and automated payroll processing.

Features
✅ Attendance Management

Record daily attendance
Track working days
Handle attendance adjustments
🕒 Overtime Tracking

Submit overtime requests
Calculate overtime pay
Period-based overtime reports
💰 Reimbursement Processing

Submit expense claims
Multi-level approval workflow
Category-based tracking
📊 Payroll Generation

Automated payslip generation
Detailed salary breakdowns
Support for multiple components
Tech Stack
Language: Go 1.21+
Database: PostgreSQL 13+
ORM: XORM
Framework: Gin
Testing: Go testing + sqlmock
Containerization: Docker
Getting Started
Prerequisites
Go 1.21 or higher
PostgreSQL 13+
Docker (optional)
Installation
Clone the repository
bash


git clone https://github.com/faisalhardin/employee-payroll-system.git
cd employee-payroll-system
Install dependencies
```
go mod download
```
Set up environment variables
```
cp files/env/envconfig.yaml.example files/env/envconfig.yaml
```
# Edit .env with your configuration

Start the application

go run cmd/main.go
Docker Setup


# Using Docker Compose
docker-compose up -d

API Documentation
Base URL
http://localhost:8080/api/v1
Authentication
All endpoints require JWT authentication:

Authorization: Bearer <jwt_token>
valid jwt_token for admin role:
```eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJleHAiOjE3NTU2MTAxODYsImlhdCI6MTc1MzAxODE4NiwicGF5bG9hZCI6eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsInJvbGUiOiJhZG1pbiJ9fQ.1wEfwfVeSaKGFESYa9c8ykaKTxhJfi9TdXFoXzoIYi4```
Key Endpoints
## Login
curl --location 'localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "employee_001",
    "password": "$2a$10$hash001"
}'
```
## Attendance
POST v1/tap-in - Record attendance
```
curl --location --request POST 'localhost:8080/v1/tap-in' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJleHAiOjE3NTU2MTAxODYsImlhdCI6MTc1MzAxODE4NiwicGF5bG9hZCI6eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsInJvbGUiOiJhZG1pbiJ9fQ.1wEfwfVeSaKGFESYa9c8ykaKTxhJfi9TdXFoXzoIYi4'
```
## Overtime
POST /overtime - Submit overtime
```
curl --location 'localhost:8080/v1/overtime' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJleHAiOjE3NTU2MTAxODYsImlhdCI6MTc1MzAxODE4NiwicGF5bG9hZCI6eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsInJvbGUiOiJhZG1pbiJ9fQ.1wEfwfVeSaKGFESYa9c8ykaKTxhJfi9TdXFoXzoIYi4' \
--data '{
    "overtime_date": "2025-07-26T00:00:00+07:00",
    "hours": 5
}'
```
## Reimbursement
POST /reimbursement - Submit reimbursement
```
curl --location 'localhost:8080/v1/reimbursement' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJleHAiOjE3NTU2MTAxODYsImlhdCI6MTc1MzAxODE4NiwicGF5bG9hZCI6eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsInJvbGUiOiJhZG1pbiJ9fQ.1wEfwfVeSaKGFESYa9c8ykaKTxhJfi9TdXFoXzoIYi4' \
--data '{
    "amount": 750000,
    "description": "therapist 5"
}'
```
## Payroll
POST /payroll/generate - Generate payslip
```
curl --location 'localhost:8080/v1/payroll/generate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJleHAiOjE3NTU2MTAxODYsImlhdCI6MTc1MzAxODE4NiwicGF5bG9hZCI6eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsInJvbGUiOiJhZG1pbiJ9fQ.1wEfwfVeSaKGFESYa9c8ykaKTxhJfi9TdXFoXzoIYi4' \
--data '{
    "payroll_period_id" : 5
}'
```
GET v1/payroll - View payroll details
```
curl --location 'localhost:8080/v1/payroll?payroll_period_id=5' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJleHAiOjE3NTU2MTAxODYsImlhdCI6MTc1MzAxODE4NiwicGF5bG9hZCI6eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsInJvbGUiOiJhZG1pbiJ9fQ.1wEfwfVeSaKGFESYa9c8ykaKTxhJfi9TdXFoXzoIYi4'
```
## Payslip
GET v1/payslip - View payslip
```
curl --location 'localhost:8080/v1/payslip?payroll_period_id=5' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJleHAiOjE3NTU3MzgwODUsImlhdCI6MTc1MzE0NjA4NSwicGF5bG9hZCI6eyJpZCI6MiwidXNlcm5hbWUiOiJlbXBsb3llZV8wMDEiLCJyb2xlIjoiZW1wbG95ZWUifX0.wV7fGrxh8tz9KMJGlJXBmta2M8prcaD5L7CiF-2_th8'
```


Project Structure
employee-payroll-system/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── entity/
│   │   └── model/              # Domain models
│   ├── usecase/                # Business logic
│   ├── handler/                # HTTP handlers
│   └── repo/
│       └── db/                 # Database repositories
├── pkg/
│   └── xorm/                   # Database utilities
├── config/                     # Configuration
├── migrations/                 # Database migrations
└── docs/                       # Documentation
Testing
Run tests:

bash


# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/repo/db/attendance/...
