
# Devices API

## Overview

The **Devices API** is a REST service built in **Go**, designed to manage devices.  
It supports creating, updating, listing, filtering, and deleting devices.  

The project follows best practices such as clean architecture, middleware-based HTTP handling, request IDs, structured logs, and full API documentation using **Swagger**.

---

## Technologies Used

- **Go 1.22+**
- Native `net/http` server
- **Swagger / Swaggo** for API documentation
- **UUID** for request ID generation
- PostgreSQL


---

# 🛠 How to Run the Project

## ▶️ Run locally

```bash
go mod tidy
go run cmd/api/main.go
```

The API will be available at:

```
http://localhost:8080
```

---

## 🐳 Running with Docker

### Build the image

```bash
docker build -t devices-api .
```

### Run the container

```bash
docker run -p 8080:8080 devices-api
```

---


# API Endpoints

## Create Device  
**POST /devices**

### Request Body

```json
{
  "name": "iPhone 14 Pro",
  "brand": "Apple",
  "state": "active"
}
```

### Response 201 Created

```json
{
  "id": "3a298e4b-1f12-4060-aeb8-1ec54430ea67"
}
```

---

## Update Device  
**PUT /devices/{id}**

### Request Body

```json
{
  "name": "Galaxy S22",
  "brand": "Samsung",
  "state": "active"
}
```

### Response Example

```json
{
  "updated_fields": ["name", "brand"],
  "ignored_fields": ["state"],
  "device": {
    "id": "3a298e4b-1f12-4060-aeb8-1ec54430ea67",
    "name": "Galaxy S22",
    "brand": "Samsung",
    "state": "active",
    "created_at": "2025-01-10T15:04:05Z"
  }
}
```

---

## Delete Device  
**DELETE /devices/{id}**

### Response 204 No Content  
(empty response body)

---

## Get Device by ID  
**GET /devices/{id}**

### Response Example

```json
{
  "id": "3a298e4b-1f12-4060-aeb8-1ec54430ea67",
  "name": "iPhone 14 Pro",
  "brand": "Apple",
  "state": "active",
  "created_at": "2025-01-10T15:04:05Z"
}
```

---

## List all Devices  
**GET /devices**

### Response Example

```json
[
  {
    "id": "3a298e4b-1f12-4060-aeb8-1ec54430ea67",
    "name": "iPhone 12",
    "brand": "Apple",
    "state": "active",
    "created_at": "2025-01-10T15:04:05Z"
  }
]
```

---

## Filter by brand  
**GET /devices?brand=Apple**

## Filter by state  
**GET /devices?state=active**

---

# Middlewares Included

- ✔ **Recover** — prevents server crashes on panic  
- ✔ **RequestID** — injects a unique `X-Request-ID` into each request  
- ✔ **Logger** — logs all requests with method, path, status & duration  
- ✔ **Timeout** — ensures long-running requests are aborted safely  

---


