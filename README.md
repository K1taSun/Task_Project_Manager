# Task Project Manager

**Professional Task and Project Management System** built with Go, featuring a modern architecture, comprehensive testing, and production-ready deployment options.

## 🚀 Features

### Core Functionality
- **Project Management**: Create, edit, delete, and track projects with status management
- **Task Management**: Comprehensive task system with priorities, deadlines, tags, and time tracking
- **Real-time Updates**: WebSocket support for live updates
- **Advanced Filtering**: Filter tasks by tags, priority, date ranges, and project
- **Data Export**: Export data in JSON and CSV formats
- **Statistics Dashboard**: Real-time project and task statistics
- **Badge System**: Visual indicators for project status and task priority

### Technical Features
- **RESTful API**: Complete CRUD operations with proper HTTP status codes
- **Data Validation**: Comprehensive input validation and error handling
- **CORS Support**: Cross-origin resource sharing enabled
- **Logging**: Structured logging with middleware
- **Configuration Management**: Environment-based configuration
- **File-based Storage**: JSON-based data persistence

## 🏗️ Architecture

The project follows a clean, modular architecture:

```
Task_Project_Manager/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routing
│   ├── config/          # Configuration management
│   ├── models/          # Data models and business logic
│   ├── storage/         # Data persistence layer
│   ├── utils/           # Utility functions
│   └── validation/      # Input validation
├── pkg/                 # Public packages (future use)
├── web/
│   └── static/          # Frontend assets
├── tests/
│   ├── unit/            # Unit tests
│   └── integration/     # Integration tests
├── deployments/         # Deployment configurations
├── docs/               # Documentation
└── scripts/            # Build and deployment scripts
```

## 🛠️ Technology Stack

- **Backend**: Go 1.21+
- **Frontend**: HTML5, CSS3, JavaScript (Vanilla)
- **Data Storage**: JSON files
- **Testing**: Go testing framework
- **Containerization**: Docker & Docker Compose
- **Build System**: Make

## 📋 Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional, for containerized deployment)
- Make (optional, for build automation)

## 🚀 Quick Start

### Option 1: Direct Go Run
```bash
# Clone the repository
git clone https://github.com/yourusername/Task_Project_Manager.git
cd Task_Project_Manager

# Install dependencies
go mod download

# Run the application
go run ./cmd/server
```

### Option 2: Using Make
```bash
# Install dependencies and run
make run

# Or build and run
make build
./build/task-manager
```

### Option 3: Docker
```bash
# Build and run with Docker
docker-compose up --build

# Or run in development mode
docker-compose --profile dev up --build
```

## 🧪 Testing

The project includes comprehensive test coverage:

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race
```

## 📊 API Documentation

### Projects Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects` | List all projects |
| POST | `/projects` | Create a new project |
| GET | `/projects/{id}` | Get project by ID |
| PUT | `/projects/{id}` | Update project |
| DELETE | `/projects/{id}` | Delete project |

### Tasks Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | List all tasks (with filtering) |
| POST | `/tasks` | Create a new task |
| GET | `/tasks/{id}` | Get task by ID |
| PUT | `/tasks/{id}` | Update task |
| DELETE | `/tasks/{id}` | Delete task |

### Additional Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/export?format=json\|csv` | Export data |
| GET | `/stats` | Get statistics |
| GET | `/badges` | Get available badges |
| GET | `/ws` | WebSocket endpoint |

### Query Parameters for Task Filtering

- `tag`: Filter by tag
- `min_priority`: Minimum priority (1-5)
- `max_priority`: Maximum priority (1-5)
- `before`: Filter tasks created before date (RFC3339)
- `after`: Filter tasks created after date (RFC3339)
- `project_id`: Filter by project ID

## 🔧 Configuration

The application can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `HOST` | `localhost` | Server host |
| `PROJECTS_FILE` | `data_projects.json` | Projects data file |
| `TASKS_FILE` | `data_tasks.json` | Tasks data file |

## 📦 Deployment

### Production Deployment

```bash
# Build release binaries
make release

# Build Docker image
make docker-build

# Run with Docker Compose
docker-compose --profile production up -d
```

### Development Deployment

```bash
# Run with hot reload
make dev

# Run with Docker Compose (development)
docker-compose --profile dev up --build
```

## 🧹 Code Quality

The project includes several tools for maintaining code quality:

```bash
# Format code
make fmt

# Check formatting
make fmt-check

# Run linter
make lint

# Run security checks
make security

# Run all checks
make check
```

## 📈 Monitoring and Health Checks

The application includes built-in health checks and monitoring:

- **Health Check Endpoint**: `GET /health`
- **Docker Health Check**: Automatic container health monitoring
- **Structured Logging**: JSON-formatted logs for easy parsing

## 🔒 Security Features

- Input validation and sanitization
- CORS configuration
- Non-root Docker container
- Security scanning with gosec
- Rate limiting (planned)

## 🚀 Performance Features

- Efficient in-memory data structures
- Optimized JSON serialization
- Minimal memory footprint
- Fast startup time

## 📝 Data Models

### Project
```json
{
  "id": 1,
  "name": "Project Name",
  "description": "Project description",
  "status": "active|completed|archived",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "badge": {
    "text": "New",
    "color": "#ffffff",
    "background": "#3b82f6",
    "type": "status"
  }
}
```

### Task
```json
{
  "id": 1,
  "project_id": 1,
  "title": "Task Title",
  "description": "Task description",
  "priority": 3,
  "done": false,
  "deadline": "2024-12-31T23:59:59Z",
  "tags": ["tag1", "tag2"],
  "estimated_hours": 5.0,
  "actual_hours": 4.0,
  "assignee": "John Doe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards
- Write tests for new features
- Update documentation
- Run all checks before submitting PR

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/Task_Project_Manager/issues)
- **Documentation**: [Wiki](https://github.com/yourusername/Task_Project_Manager/wiki)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/Task_Project_Manager/discussions)

## 🗺️ Roadmap

### Version 1.1 (Planned)
- [ ] Database integration (PostgreSQL)
- [ ] User authentication and authorization
- [ ] API rate limiting
- [ ] Advanced search functionality
- [ ] Email notifications

### Version 1.2 (Planned)
- [ ] Mobile app support
- [ ] Real-time collaboration
- [ ] File attachments
- [ ] Advanced reporting
- [ ] Integration with external tools

### Version 2.0 (Planned)
- [ ] Microservices architecture
- [ ] Kubernetes deployment
- [ ] Multi-tenant support
- [ ] Advanced analytics
- [ ] AI-powered task suggestions

---

**Built with ❤️ using Go**

