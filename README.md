# Go Redis-like Key-Value Store

This project is a simple Redis-like in-memory key-value store implemented in Go. It supports basic operations such as GET, PUT, and DELETE, and uses sharding for concurrency. The server exposes a RESTful API and logs all events and errors to a file.

## Features
- In-memory key-value storage
- Sharding for concurrent access
- RESTful API using Gorilla Mux
- File-based event and error logging
- Simple custom hash function for sharding

## API Endpoints

### PUT (Create a new key)
- **URL:** `/v1/{key}`
- **Method:** `PUT`
- **Body:** Raw value (string)
- **Response:** `key added successfully{key:<key>}`
- **Error:** 409 Conflict if key already exists

### GET (Retrieve a value)
- **URL:** `/v1/{key}`
- **Method:** `GET`
- **Response:** Value as plain text
- **Error:** 409 Conflict if key does not exist

### DELETE (Delete a key)
- **URL:** `/v1/delete/{key}`
- **Method:** `GET`
- **Response:** `sucessfully deleted key :<key>`
- **Error:** 409 Conflict if key does not exist

## Logging
- All events and errors are logged to `tmp.log` in the project directory.

## Running the Server

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```
2. **Run the server:**
   ```bash
   go run cmd/main.go
   ```
3. The server will start on `http://localhost:8080`

## Example Usage

- **Add a key:**
  ```bash
  curl -X PUT http://localhost:8080/v1/mykey -d 'myvalue'
  ```
- **Get a key:**
  ```bash
  curl http://localhost:8080/v1/mykey
  ```
- **Delete a key:**
  ```bash
  curl http://localhost:8080/v1/delete/mykey
  ```

## Project Structure
- `cmd/main.go`: Main application code
- `go.mod`, `go.sum`: Go module files
- `tmp.log`: Log file (created at runtime)

## License
This project is for educational purposes.
