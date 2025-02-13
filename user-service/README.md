# Go MongoDB Starter

This project is a starter template for building web applications using Go (Golang) and MongoDB. It includes basic setup for database connections, migrations, seeding, and basic user management (registration and login).

## Features

- Go (Golang) backend
- MongoDB as the database
- Basic user management (registration, login)
- User authentication using JWT
- Database migrations and seeding
- Structured project layout

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or above)
- [MongoDB](https://www.mongodb.com/try/download/community) installed and running
- [Git](https://git-scm.com/)
- [Docker](https://docs.docker.com/get-docker/) (optional, for containerization)

## Getting Started

### 1. Clone the repository

First, clone the repository:

```bash
git clone https://github.com/Oksastyaa/go-mongoDB-starter.git
cd go-mongoDB-starter
```
### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Set Environment Variables
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret_key

### 4. Run the Application
```bash
go run main.go
```

### 5. API Endpoints
**URL: /api/v1/users/register
Method: POST
Request Body:**
```json
{
  "username": "exampleuser",
  "email": "user@example.com",
  "password": "yourpassword",
  "address": "123 Example St",
  "phone": "123-456-7890",
  "age": 25
}
```

**URL: /api/v1/users/login
Method: POST
Request Body:**
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
} "
```
### 6.Project Structure
```
.
├── cmd
│   └── app                 # Main application entry point
├── config                  # Application configuration and environment variables
├── controllers             # HTTP request handlers (controllers)
├── database
│   ├── migrations          # Database migration files
│   ├── seeder              # Seeder files for populating the database
├── models                  # Database models
├── pkg                     # Helper packages (e.g., for hashing passwords, JWT)
├── routes                  # API routes setup
├── .env.example            # Example environment file
├── Dockerfile              # Dockerfile for containerization
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies lock file
└── README.md               # Project README file
```

Database Migrations and Seeding
You can manage database migrations and seed data with the following commands:

```bash
# Run database migrations
air
```
Docker Support
You can also run the application using Docker. First, build the Docker image:

```bash
docker build -t go-mongodb-starter .
```
Then, run the Docker container:

```bash
docker run -p 8080:8080 go-mongodb-starter
```
The application will be accessible at http://localhost:8080.

Contributing
If you'd like to contribute, feel free to open a pull request or submit an issue on GitHub.

License
This project is licensed under the MIT License - see the LICENSE file for details.

markdown
Copy code

### Penjelasan:
- **Getting Started**: Petunjuk untuk memulai proyek ini, termasuk cloning repository, menginstal dependensi, dan menjalankan aplikasi.
- **API Endpoints**: Menjelaskan endpoint untuk registrasi dan login pengguna.
- **Project Structure**: Gambaran umum mengenai struktur proyek, agar lebih mudah dipahami oleh kontributor atau developer baru.
- **Database Migrations and Seeding**: Instruksi untuk migrasi dan seeding data ke database MongoDB.
- **Running Tests**: Cara menjalankan unit tests untuk memvalidasi fungsi dari proyek.
- **Docker Support**: Instruksi untuk membangun dan menjalankan aplikasi menggunakan Docker.
- **Contributing & License**: Petunjuk untuk berkontribusi dan informasi lisensi proyek.
