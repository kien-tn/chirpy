# Chirpy

Chirpy is a lightweight Go-based web application for creating and managing short messages called "chirps." It provides a RESTful API for users to create, retrieve, update, and delete chirps, with features like user authentication and sorting.

## Features

- **Create Chirps**: Users can post short messages (up to 140 characters).
- **Retrieve Chirps**: Fetch all chirps or filter by author ID, with optional sorting by creation date.
- **Delete Chirps**: Users can delete their own chirps.
- **Authentication**: JWT-based authentication for secure user access.

## Project Structure

```
/chirpy
├── handler_chirps.go   # Contains HTTP handlers for chirp-related operations
├── internal
│   ├── auth            # Authentication utilities
│   ├── database        # Database interaction logic
└── README.md           # Project documentation
```

## Prerequisites

- Go 1.18 or later
- A PostgreSQL database instance
- Environment variables for database connection and JWT secret key

## Setup Instructions

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/kien-tn/chirpy.git
   cd chirpy
   ```

2. **Set Environment Variables**:
   Create a `.env` file in the project root with the following:
   ```
   DATABASE_URL=your_database_connection_string
   JWT_SECRET_KEY=your_secret_key
   ```

3. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run the Application**:
   ```bash
   go run main.go
   ```

5. **Access the API**:
   The API will be available at `http://localhost:8080`.

## API Endpoints

### Chirps

- **POST /chirps**: Create a new chirp.
- **GET /chirps**: Retrieve all chirps (supports filtering by `author_id` and sorting with `sort=asc|desc`).
- **GET /chirps/{chirp_id}**: Retrieve a chirp by its ID.
- **DELETE /chirps/{chirp_id}**: Delete a chirp (requires authentication).

## Example Usage

### Create a Chirp
```bash
curl -X POST http://localhost:8080/chirps \
-H "Authorization: Bearer <your_token>" \
-H "Content-Type: application/json" \
-d '{"body": "Hello, world!"}'
```

### Get All Chirps
```bash
curl http://localhost:8080/chirps
```

### Delete a Chirp
```bash
curl -X DELETE http://localhost:8080/chirps/<chirp_id> \
-H "Authorization: Bearer <your_token>"
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.