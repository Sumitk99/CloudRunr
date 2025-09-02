# CloudRunr

> A modern Platform-as-a-Service (PaaS) for deploying frontend applications with real-time logging and subdomain-based routing.

[![System Design](assets/system_design.png)](https://app.eraser.io/workspace/jmmdy7Cc7P13jBA1lK7Q?origin=share&elements=mKTok0gCqqoNiXvLTik8BQ)

##  Overview

CloudRunr is a cloud-native platform that automates the deployment of frontend applications (React/Angular) from GitHub repositories. It provides a complete CI/CD pipeline with real-time build logging, subdomain-based hosting, and multi-tenant architecture.

## Features

- **Secure Authentication**: JWT-based user authentication with GitHub integration
- **Multi-Framework Support**: React and Angular applications
- **Automated CI/CD**: GitHub integration with automatic builds and deployments
- **Real-time Logging**: Live build logs with persistent storage
- **Subdomain Routing**: Each project gets its own subdomain (`project-id.domain.com`)
- **Cloud-Native**: Built on AWS with ECS, S3, and CloudFront
- **Scalable Architecture**: Microservices-based design with event-driven logging

## Architecture

### System Components

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Server     │    │  Build Server   │
│   (Angular)     │───▶│   (Go/Gin)       │───▶│   (Go/Docker)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Proxy Server   │    │   PostgreSQL     │    │   AWS S3        │
│   (Go/Gin)      │    │   TimescaleDB    │    │   (Static       │
└─────────────────┘    └──────────────────┘    │    Assets)      │
                                │               └─────────────────┘
                                ▼
                       ┌──────────────────┐
                       │  Log Consumer    │
                       │  (Kafka→DB)      │
                       └──────────────────┘
```

### 🔄 Data Flow

1. **User Authentication**: Secure login/signup with JWT tokens
2. **Project Creation**: User provides GitHub repository URL and configuration
3. **Deployment Trigger**: API server launches containerized build process
4. **Build Process**: Clone → Install → Build → Upload to S3
5. **Log Streaming**: Real-time logs via Kafka to TimescaleDB
6. **Application Serving**: Subdomain-based routing via AWS CloudFront

## 🛠️ Technology Stack

### Backend Services
- **Language**: Go 1.23+
- **Web Framework**: Gin
- **Authentication**: JWT with bcrypt password hashing
- **Databases**: PostgreSQL + TimescaleDB
- **Message Queue**: Apache Kafka
- **Cloud**: AWS (ECS, S3, CloudFront)
- **Containerization**: Docker

### Frontend
- **Framework**: Angular 18+
- **Language**: TypeScript
- **Styling**: SCSS
- **Build Tool**: Angular CLI

### Infrastructure
- **Container Orchestration**: AWS ECS Fargate
- **Static Hosting**: AWS S3 + CloudFront
- **Load Balancing**: Application Load Balancer
- **Monitoring**: CloudWatch Logs

## Prerequisites

- Go 1.23+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL 15+
- Apache Kafka
- AWS Account with ECS, S3, CloudFront access

## Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/Sumitk99/CloudRunr.git
cd CloudRunr
```

### 2. Environment Setup

Create `.env` files in each service directory:

#### API Server (`.env`)
```env
# Database
PG_URL=postgresql://username:password@localhost:5432/cloudrunr
TS_URL=postgresql://username:password@localhost:5432/timescale

# AWS Configuration
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
AWS_ENDPOINT=https://ecs.us-east-1.amazonaws.com

# ECS Configuration
ECS_CLUSTER_ARN=arn:aws:ecs:region:account:cluster/cloudrunr
ECS_TASK_DEF_ARN=arn:aws:ecs:region:account:task-definition/build-server
SUBNETS=subnet-xxx,subnet-yyy
SECURITY_GROUPS=sg-xxx,sg-yyy
```

#### Build Server & Log Consumer
```env
# AWS S3
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
AWS_ENDPOINT=https://s3.us-east-1.amazonaws.com

# Database (Log Consumer)
TS_URL=postgresql://username:password@localhost:5432/timescale
```

### 3. Kafka Configuration

Create `client.properties` files:
```properties
bootstrap.servers=localhost:9092
security.protocol=PLAINTEXT
session.timeout.ms=30000
auto.offset.reset=earliest
```

### 4. Database Setup

```sql
-- PostgreSQL Database
CREATE DATABASE cloudrunr;

-- Users table
CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    github_id VARCHAR(255)
);

-- Projects table
CREATE TABLE projects (
    project_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(user_id),
    github_url TEXT NOT NULL,
    name VARCHAR(255) NOT NULL,
    framework VARCHAR(50) NOT NULL,
    dist_folder VARCHAR(255) DEFAULT 'dist',
    subdomain VARCHAR(255),
    custom_subdomain VARCHAR(255)
);

-- Deployments table
CREATE TABLE deployments (
    deployment_id VARCHAR(255) PRIMARY KEY,
    project_id VARCHAR(255) REFERENCES projects(project_id),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- TimescaleDB for logs
CREATE TABLE log_statements (
    deployment_id VARCHAR(255),
    project_id VARCHAR(255),
    log_statement TEXT,
    ts TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Convert to hypertable (TimescaleDB)
SELECT create_hypertable('log_statements', 'ts');
```

### 5. Start Services

#### Using Docker Compose (Recommended)
```bash
# Start infrastructure services
docker-compose up -d postgres kafka

# Build and start application services
docker-compose up --build
```

#### Manual Setup
```bash
# Start API Server
cd api-server
go mod download
go run cmd/api-server/main.go

# Start Proxy Server
cd proxy-server
go run cmd/proxy-server/main.go

# Start Log Consumer
cd log-consumer
go run main.go

# Start Frontend
cd frontend
npm install
ng serve
```

## 📚 API Documentation

### Authentication Endpoints

#### Sign Up
```http
POST /signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword"
}
```

### Project Management

#### Create Project
```http
POST /project
Authorization: token your_jwt_token
Content-Type: application/json

{
  "git_url": "https://github.com/username/react-app.git",
  "name": "My React App",
  "framework": "REACT",
  "dist_folder": "build",
  "project_id": "my-react-app"
}
```

#### Get User Projects
```http
GET /projects
Authorization: token your_jwt_token
```

#### Deploy Project
```http
POST /deploy/{project_id}
Authorization: token your_jwt_token
```

#### Get Deployment Logs
```http
GET /logs/{deployment_id}/{offset}
Authorization: token your_jwt_token
```

## 🔧 Configuration

### Supported Frameworks

| Framework | Build Command | Default Dist Folder |
|-----------|---------------|-------------------|
| React | `npm run build` | `build` |
| Angular | `npx ng build --configuration=production` | `dist` |

### Environment Variables

#### API Server
- `PG_URL`: PostgreSQL connection string
- `TS_URL`: TimescaleDB connection string
- `AWS_*`: AWS credentials and configuration
- `ECS_*`: ECS cluster and task definition ARNs

#### Build Server
- `GIT_REPOSITORY_URL`: Repository to clone (set by ECS)
- `PROJECT_ID`: Unique project identifier
- `FRAMEWORK`: Target framework (REACT/ANGULAR)
- `DEPLOYMENT_ID`: Unique deployment identifier
- `DEFAULT_DIST_FOLDER`: Build output directory

## Development

### Project Structure
```
CloudRunr/
├── api-server/          # Main backend API
│   ├── cmd/api-server/  # Application entry point
│   ├── internal/        # Internal packages
│   │   ├── constants/   # Application constants
│   │   ├── handler/     # HTTP handlers
│   │   ├── middleware/  # HTTP middleware
│   │   ├── models/      # Data models
│   │   ├── repository/  # Database layer
│   │   ├── routes/      # Route definitions
│   │   ├── server/      # ECS integration
│   │   └── service/     # Business logic
│   └── Dockerfile
├── build-server/        # Containerized build service
│   ├── constants/       # Build constants
│   ├── helper/          # Utility functions
│   ├── script/          # Build orchestration
│   ├── server/          # S3 & Kafka integration
│   └── Dockerfile
├── log-consumer/        # Kafka to DB consumer
├── proxy-server/        # Application proxy
├── frontend/            # Angular frontend
└── assets/              # Documentation assets
```


### Building Docker Images
```bash
# Build API Server
cd api-server
docker build -t cloudrunr/api-server .

# Build Build Server
cd build-server
docker build -t cloudrunr/build-server .
```

## Security

### Authentication
- JWT tokens with 30-day expiration
- Bcrypt password hashing (cost factor: 14)
- Secure HTTP-only cookies (production)

### Authorization
- User-project ownership validation
- Deployment access control
- Multi-tenant isolation via subdomain-based routing (CloudFront wildcard domain)
- S3 bucket path isolation per project (`s3://bucket/project_id/`)

### Infrastructure
- VPC with private subnets
- Security groups with minimal access
- IAM roles with least privilege

## Monitoring & Logging

### Application Logs
- Structured JSON logging
- Real-time log streaming via Kafka
- Persistent storage in TimescaleDB
- Log aggregation and search

## Deployment

### Production Setup

1. **Infrastructure**: Deploy AWS resources (ECS, S3, Cloudfront)
2. **Database**: Set up PostgreSQL with TimescaleDB extension
3. **Messaging**: Configure Apache Kafka cluster
4. **Container Registry**: Push images to AWS ECR
5. **Load Balancer**: Configure ALB with SSL termination
6. **CDN**: Set up CloudFront with custom domain

### Environment-Specific Configuration
- **Development**: Local services with Docker Compose
- **Staging**: Reduced resource allocation
- **Production**: Auto-scaling, monitoring, and backup strategies

## Acknowledgments

- Built with Go, Angular, and modern cloud technologies
- Inspired by platforms like Vercel and Netlify

---
