# Application Products and Cutomer Reviews

A full-stack application that manages products and their associated customer reviews.

## Overview

This project implements a CRUD application with:
- React TypeScript frontend for user interaction
- Go backend REST API for data operations
- PostgreSQL database running in Docker

The system allows users to create, read, update, and delete products as well as add, view customer reviews for those products.

## System Architecture

### Frontend
- React with TypeScript

### Backend
- Go REST API using standard libraries
- PostgreSQL database via GORM (Go Library for mangaing postgresql with plenty of documentation)
- Docker containerisation (allows for building from a single line)

### Data Model
- **Products**: Main resource with name, description, price
- **Reviews**: Related objects with comment, and association to a product

## Features

- Full CRUD operations for products and reviews
- Filtering products by name
- Sorting products by price, name and date
- Pagination for product listings
- Simple and intuitive UI

## Running the Application

The entire application can be started with a single command:

Windows
```bash
docker-compose up
```
Linux
```bash
docker compose up
```

This will:
1. Start the PostgreSQL database
2. Build and start the Go backend API
3. Build and serve the React frontend

### Accessing the Application

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## API Endpoints

### Products

- `GET /products` - List all products (with pagination, filtering and sorting)
- Example (`/products?page=1&limit=12&sort=price&order=asc&search=`)
- `GET /products/:id` - Get a specific product
- `POST /products` - Create a new product
- `PATCH /products/:id` - Update a product
- `DELETE /products/:id` - Delete a product

### Reviews

- `POST /reviews` - Create a new review
- `GET /products/:id/reviews` - Get all reviews for a product
- `PATCH /reviews/:id` - Update a review
- `DELETE /reviews/:id` - Delete a review

## Testing

Manually Tested

- Got GitHub actions working made some quick autotests which check for Product and Review Functionality further tests would need to be made for thorough testing

### Testing Pipeline
- `docker compose up -d`
- `cd backend`
- `go test -v`
- `docker compose down`

## Development Decisions/Assumptions

- Backend, frontend, and database are all containerised in docker to allow for single line build and deployment
- Go Backend uses GROM library, a ORM library for communicating with the postgresql database.
- Frontend uses Axios to make API calls

## Future Improvements

- Add authentication, users and privileges
- Implement more advanced filtering and search capabilities
- Add image upload for products
- Add rating system for reviews
- Enhance test coverage/Fix autotests
- Clean up CSS and separate the frontend into smaller components for readability

## Dev Comments
- Working on a fullstack project like this from scratch is quite challenging trying to figure out how to setup all the environments and just getting a basic structure setup to work on. If I could start for scratch I would make sure testing was integrated and working from the start
