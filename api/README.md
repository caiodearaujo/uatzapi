# UatzAPI - API for managing WhatsApp devices and sending messages

This project provides a backend service for handling WhatsApp devices, webhooks, and messages using the [WhatsMeow](https://pkg.go.dev/go.mau.fi/whatsmeow) library. It integrates with PostgreSQL using the [Bun ORM](https://bun.uptrace.dev/) for managing devices, webhooks, and message logs, and uses the [Gin framework](https://gin-gonic.com/) for the API layer.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Technologies](#technologies)
- [Project Structure](#project-structure)
- [Development Environment Setup](#development-environment-setup)
- [Environment Variables (`.env`)](#environment-variables-env)
- [API Routes](#api-routes)
- [Models and Database](#models-and-database)
- [Error Handling](#error-handling)
- [Running the Application](#running-the-application)

## Overview

The system is responsible for managing WhatsApp devices, sending messages and stickers, receiving webhook responses, and logging interactions. It uses PostgreSQL as the database for storing information about devices, handlers, and webhooks.

## Features

- Manage WhatsApp devices (connect, list, get details)
- Send text messages and stickers
- Handle webhooks and log responses
- PostgresSQL database with Bun ORM
- REST API using Gin framework
- CORS and API token authentication middleware

## Technologies

- **Go version**: `1.23.1`
- **ORM**: [Bun](https://bun.uptrace.dev/) with PostgreSQL
- **WhatsApp API**: [WhatsMeow](https://pkg.go.dev/go.mau.fi/whatsmeow)
- **Web Framework**: [Gin](https://gin-gonic.com/)
- **Database**: PostgresSQL
- **Environment Variables**: Managed using a `.env` file (loaded via [godotenv](https://github.com/joho/godotenv))

## Project Structure

```bash
.
├── conf             # Configuration files (middleware, token, CORS)
├── data             # Models and database logic
├── events           # Event handlers for WhatsApp messages
├── helpers          # Helper functions for message and sticker processing
├── routes           # API routes for sending messages, webhooks, and managing devices
├── store            # Database connection and management logic
├── .env             # Environment variables for project configuration
└── main.go          # Main application entry point
```

## Development Environment Setup

### Prerequisites

- **Go 1.23.1**: Make sure you have Go version 1.23.1 installed.
- **PostgresSQL**: Set up a PostgresSQL instance for the database.
- **WhatsMeow**: Ensure you have the proper WhatsApp setup for connecting devices and sending messages.

### Installing Dependencies

Clone the repository and install dependencies:

```bash
git clone https://github.com/caiodearaujo/uatzapi.git
cd uatzapi/api
go mod tidy
```

### Environment Variables (`.env`)

The system requires a `.env` file for configuration. This file contains all the necessary environment variables for connecting to the database and other services.

#### Example `.env` file:

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

## API Routes

The following API routes are provided for interacting with WhatsApp devices, sending messages, and handling webhooks:

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

### Example Model:

Here’s an example of how the `Device` model is structured:

```go
type Device struct {
bun.BaseModel `bun:"table:device,alias:dvc"`
ID            int       `json:"id" bun:"id,pk,autoincrement"`
JID           string    `json:"whatsapp_id" bun:"whatsapp_id,notnull,unique"`
PushName      string    `json:"push_name" bun:"push_name,notnull"`
BusinessName  string    `json:"business_name" bun:"business_name"`
Active        bool      `json:"active" bun:"active,notnull"`
CreatedAt     time.Time `json:"created_at" bun:"created_at,notnull"`
}
```

### Database Setup

The database is set up using the `Bun` ORM, and models are automatically created if they don't exist. You can create the tables by calling the `CreateTablesFromDataPkg()` function in the `store` package.

```go
store.CreateTablesFromDataPkg()
```

## Error Handling

All errors are logged, and if critical, the application terminates gracefully using `FailOnError`. For example:

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

1. Make sure your `.env` file is set up correctly.
2. Run the following command to start the application:

```bash
go run main.go
```

The application will start on port `8080` by default. You can test the API routes using tools like `curl` or Postman.