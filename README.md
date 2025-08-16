# Furniro Backend

A simple RESTful API for managing products in the Furniro project, built with Go and Gin.

## Features

- Get all products: `GET /api/products`
- Get product by ID: `GET /api/products/:id`
- Add new product: 
- Update product: 
- Delete product: 
- CORS enabled for frontend: `https://furniro-project-bryce.netlify.app`

## Product Structure 

```json
{
  "id": "1",
  "title": "Syltherine",
  "description": "Stylish cafe chair",
  "price": 2500000.00,
  "discount": "-30%",
  "isNew": false
}
```

## Requirments

- Go 1.24
- Gin Web Framework

## Test Endpoints (with curl)

- Get all products:
```bash
curl http://localhost:8080/api/products
```

- Get product by ID:
```bash
curl http://localhost:8080/api/products/1
```
