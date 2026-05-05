# API Documentation

This directory contains the API documentation for the Subscription Management API.

## Files

- `swagger.yaml` - OpenAPI 3.0 specification for all API endpoints

## Viewing the Documentation

### Option 1: Swagger UI (Online)

Visit [Swagger Editor](https://editor.swagger.io/) and paste the contents of `swagger.yaml`.

### Option 2: Local Swagger UI with Docker

```bash
docker run -p 8080:8080 -e SWAGGER_JSON=/api/swagger.yaml -v $(pwd):/api swaggerapi/swagger-ui
```

Then open http://localhost:8080 in your browser.

### Option 3: ReDoc (Alternative UI)

Use [ReDoc](https://redoc.ly/) online editor or install locally.

## API Overview

The API provides endpoints for managing user subscriptions with the following functionality:

### Core Operations
- **Create** - Create a new subscription
- **Read** - Retrieve a subscription by ID
- **Update** - Update an existing subscription
- **Delete** - Delete a subscription
- **List** - List all subscriptions for a user

### Aggregation Operations
- **Count by User ID** - Calculate total subscription price for a user within a date range
- **Count by Service Name** - Calculate total subscription price for a service within a date range

## Date Format

All dates must be in **MM-YYYY** format:
- Valid: `01-2025`, `12-2024`, `06-2023`
- Invalid: `2025-01`, `01/2025`, `January 2025`

## Authentication

Currently, no authentication is required. This should be added in production.

## Error Handling

All errors return appropriate HTTP status codes:
- `400` - Validation errors (invalid format, missing required fields)
- `404` - Resource not found
- `409` - Conflict (e.g., duplicate subscription)
- `500` - Internal server errors

Each error response includes:
- `code` - Error classification
- `message` - Human-readable error message
