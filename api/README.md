# UatzAPI - API for managing WhatsApp devices and sending messages

This project provides a backend service for handling WhatsApp devices, webhooks, and messages using the [WhatsMeow](https://pkg.go.dev/go.mau.fi/whatsmeow) library. It integrates with PostgreSQL using the [Bun ORM](https://bun.uptrace.dev/) for managing devices, webhooks, and message logs, and uses the [Gin framework](https://gin-gonic.com/) for the API layer.

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Technologies](#technologies)
- [Development Environment Setup](#development-environment-setup)
- [Environment Variables (`.env`)](#environment-variables-env)
- [Database Setup with Flyway](#database-setup-with-flyway)
- [API Routes](#api-routes)
- [Models and Database](#models-and-database)
- [Error Handling](#error-handling)
- [Running the Application](#running-the-application)

## Overview

The UatzAPI system is responsible for managing WhatsApp devices, sending messages and stickers, handling webhooks, and logging interactions with external webhook services. The backend is built using Go with the Gin framework for API management and integrates with PostgreSQL via the Bun ORM. The database schema is managed using Flyway, with migrations handled through Docker Compose.

## Project Structure

```bash
.
├── api               # Backend API in Go
├── db                # Flyway migrations for PostgreSQL
├── web               # Frontend interface built with Vue.js
├── docker-compose.yml # Docker Compose for orchestrating services
└── README.md         # Project documentation
```

## Technologies

- **Go version**: `1.23.1` (backend API)
- **ORM**: [Bun](https://bun.uptrace.dev/) (PostgreSQL ORM)
- **WhatsApp API**: [WhatsMeow](https://pkg.go.dev/go.mau.fi/whatsmeow)
- **Web Framework**: [Gin](https://gin-gonic.com/) (API routing)
- **Database**: PostgreSQL with Flyway for database migrations
- **Containerization**: Docker & Docker Compose for service orchestration

## Development Environment Setup

### Prerequisites

- **Go 1.23.1**: Ensure Go is installed.
- **Docker & Docker Compose**: Ensure Docker and Docker Compose are installed for running Redis, PostgreSQL and Flyway migrations.

### Installing Backend Dependencies

Clone the repository and install Go dependencies:

```bash
git clone https://github.com/caiodearaujo/uatzapi.git
cd uatzapi/api
go mod tidy
```

## Environment Variables (`.env`)

The backend API requires a `.env` file for configuration. This file contains the necessary environment variables for connecting to the PostgreSQL database and managing API authentication.

### Example `.env` file for Backend (`api/.env`):

```bash
# PostgreSQL configuration
PG_USERNAME=your_pg_user
PG_PASSWORD=your_pg_password
PG_HOSTNAME=localhost
PG_PORT=5432
PG_DATABASE=whatsapp_db
PG_UA_SCHEMA=public

# API token for authenticating requests
API_KEY_TOKEN=your_secure_token
```

#### Parameters Explanation:

- **`PG_USERNAME`**: PostgreSQL username.
- **`PG_PASSWORD`**: PostgreSQL password.
- **`PG_HOSTNAME`**: Hostname or IP address of the PostgreSQL instance.
- **`PG_PORT`**: Port on which PostgreSQL is running (default is `5432`).
- **`PG_DATABASE`**: Name of the database to use for this project.
- **`PG_UA_SCHEMA`**: Database schema (usually `public`).
- **`API_KEY_TOKEN`**: A secure token used for authenticating API requests.

## Database Setup with Flyway

The project uses Flyway to manage database migrations, and these migrations are handled using Docker Compose. The migration files are located in the `db/migrations` folder.

### Running Migrations

1. Make sure Docker is running.
2. Navigate to the root directory of the project (`uatziapi`).
3. Run the following command to bring up the PostgreSQL database and apply migrations:

```bash
docker-compose up
```

This command will start PostgreSQL and Flyway, automatically applying the migrations found in `db/migrations`.

> **Note**: Flyway is configured in the `docker-compose.yml` file. It automatically runs the migrations on startup.

## API Routes

Here are the available API routes for managing WhatsApp devices, sending messages, and handling webhooks:

| HTTP Method | Route                          | Description                                  |
|-------------|--------------------------------|----------------------------------------------|
| GET         | `/connect`                     | Connect a new WhatsApp device                |
| GET         | `/device`                      | Get a list of all devices                    |
| GET         | `/device/:deviceId`            | Get information about a specific device      |
| GET         | `/start_listener`              | Start a message listener for WhatsApp        |
| POST        | `/send/message`                | Send a text message via WhatsApp             |
| POST        | `/send/sticker`                | Send a sticker via WhatsApp                  |
| GET         | `/webhook`                     | List all active webhooks                     |
| POST        | `/webhook`                     | Add a new webhook                            |
| DELETE      | `/webhook/:deviceID`           | Remove a webhook by device ID                |
| GET         | `/webhook/:deviceID`           | Get active webhook for a specific device     |
| GET         | `/webhook/:deviceID/all`       | List all webhooks for a specific device      |

## Models and Database

### Key Models:

1. **Device**: Represents a WhatsApp device connected to the system.
2. **DeviceHandler**: Tracks the state of device handlers (active/inactive).
3. **DeviceWebhook**: Stores webhook URLs and statuses.
4. **WebhookMessage**: Stores webhook interactions (messages sent and responses received).

The models are located in the `api/data` directory and are managed by the Bun ORM.

### Database Setup

Flyway is responsible for applying migrations. You can find migration scripts inside the `db/migrations` directory. Ensure that you have PostgreSQL running via Docker Compose to apply migrations.

## Error Handling

All errors are logged, and if critical, the application terminates gracefully using the `FailOnError` function. This function logs errors and ensures they are handled appropriately:

```go
package handler

import (
"log"
)

func FailOnError(err error, msg string) {
if err != nil {
log.Fatalf("%s: %s", msg, err)
}
}
```

## Running the Application

### Backend

1. Ensure that your `.env` file is set up correctly inside the `api` folder.
2. Start the database and apply migrations with Docker Compose:

```bash
docker-compose up
```

3. Navigate to the `api` directory and run the Go application:

```bash
cd api
go run main.go
```

The backend API will start on port `8080` by default.
