# GO-CRUD-DATABASE

A simple project to learn basic Go integration with a PostgreSQL database

## Functional requirment:

- _`user`_

  - as a user, I can login and register
  - as a member, I can update and delete data user
  - as a member, I can see all data users
  - as a member, I can see detail data users by id

## Third Libs I used

- Postgres : `github.com/lib/pq`
- crypto : `golang.org/x/crypto` -> to encrypt password
- jwt : `github.com/dgrijalva/jwt-go` -> token for authorization

## Endpoint

### Register

- URL : `http://localhost:8080/api/v1/register`
- Method: `POST`
- Curl :
  ```
  curl --location 'http://localhost:8080/api/v1/register' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "username": "member8",
  "email": "member8@gmail.com",
  "password": "password"
  }'
  ```
- Response :
  ```json
  {
    "message": "New user created successfully",
    "status": "success",
    "code": 201
  }
  ```

### Login

- URL : `http://localhost:8080/api/v1/login`
- Method: `POST`
- Curl :
  ```json
  curl --location 'http://localhost:8080/api/v1/login' \
  --header 'Content-Type: application/json' \
  --data '{
  "username": "member5",
  "password": "password"
  }'
  ```
- Response :
  ```json
  {
    "message": "Authentication successful",
    "status": "success",
    "code": 200,
    "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyMjYzMjUsImlhdCI6MTc0MjIyNjAyNSwic3ViIjoiMTA3In0.gn4bClELi3FRwk5mTya-BQl6_AEBcUAW7m1aqrP8xak"
  }
  ```

### Update Data User

- URL : `http://localhost:8080/api/v1/users`
- Method: `PUT`
- Curl :
  ```
  curl --location --request PUT 'http://localhost:8080/api/v1/users' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyODUyMzgsImlhdCI6MTc0MjI4NDkzOCwic3ViIjoiMTA3In0.ZSSuhhF7BsGdfPenvXHmS6ZaWwFisN7xWzBaKAOkx3k' \
  --data-raw '{
  "userId": 105,
  "username": "member56",
  "email": "member56@gmail.com",
  "password": "password"
  }'
  ```
- Response :

  ```json
  {
    "message": "User updated successfully",
    "status": "success",
    "code": 200
  }
  ```

### Delete User By ID

- URL : `http://localhost:8080/api/v1/users?id=104`
- Method: `DELETE`
- Curl :
  ```
  curl --location --request DELETE 'http://localhost:8080/api/v1/users?id=105' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyODU0MDUsImlhdCI6MTc0MjI4NTEwNSwic3ViIjoiMTA3In0.zKWh1TiGCtIdC7fXAi_Q0w0Bq8059A68rdXsOqLN1Hc'
  ```
- Response
  ```json
  {
    "message": "User deleted successfully",
    "status": "success",
    "code": 200
  }
  ```

### Get All Data User

- URL : `http://localhost:8080/api/v1/users`
- Method: `GET`
- Curl :
  ```
  curl --location 'http://localhost:8080/api/v1/users' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyODQ0MjEsImlhdCI6MTc0MjI4NDEyMSwic3ViIjoiMTA3In0.P_RE_SlczIeM75eplTtjuqp3m6JPWDBS3rZ3QHOkQWg'
  ```
- Response
  ```json
  {
    "message": "Successfully retrieved all users",
    "status": "success",
    "code": 200,
    "data": [
      {
        "userId": 106,
        "username": "member4",
        "email": "member4@gmail.com",
        "password": "$2a$10$AtKc3o3YLLNjnp53pQyscO5ApFTujGOAvVdbj.NSAPItEPZGQCgzm",
        "createdAt": "2025-03-15T15:45:04.808686Z",
        "updateAt": "2025-03-15T15:45:04.808686Z"
      },
      {
        "userId": 107,
        "username": "member5",
        "email": "member5@gmail.com",
        "password": "$2a$10$uLKW0v7MQL/eMgpKY/F6jO8facIZ3lZILkWsuVrUaL.YCAfPmFZLa",
        "createdAt": "2025-03-16T00:09:41.84453Z",
        "updateAt": "2025-03-16T00:09:41.84453Z"
      },
      {
        "userId": 108,
        "username": "member7",
        "email": "member7@gmail.com",
        "password": "$2a$10$5.X6F0tvzDwIsmq8Z3IdLeSkuLnPDNKAHyCXDl.Ub6Q03G9lG2v6C",
        "createdAt": "2025-03-17T22:36:53.310911Z",
        "updateAt": "2025-03-17T22:36:53.310911Z"
      },
      {
        "userId": 105,
        "username": "member56",
        "email": "member55@gmail.com",
        "password": "$2a$10$XICAYsWjCrZkEsGF5atsmeqjT4tX0VZXNe2KkbB/xGK2r7sEiqj2O",
        "createdAt": "2025-03-15T15:43:49.217183Z",
        "updateAt": "2025-03-15T15:43:49.217183Z"
      },
      {
        "userId": 109,
        "username": "member8",
        "email": "member8@gmail.com",
        "password": "$2a$10$Kgxg.4fWCoTawTvq85FTr..0fuP7hjvUQoympeoZ5Pd3azhLl33MO",
        "createdAt": "2025-03-18T14:48:02.500906Z",
        "updateAt": "2025-03-18T14:48:02.500906Z"
      }
    ]
  }
  ```

### Get User By ID

- URL : `http://localhost:8080/api/v1/users?id=16`
- Method: `GET`
- Curl :
  ```
  curl --location 'http://localhost:8080/api/v1/users?id=105' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDIyODQ4NjQsImlhdCI6MTc0MjI4NDU2NCwic3ViIjoiMTA3In0.dboxvajzvXwYk6BIXXmfhWz9rqY_ekMOYu1n_6M2myc'
  ```
- Response
  ```json
  {
    "message": "Successfully retrieved user details",
    "status": "success",
    "code": 200,
    "data": {
      "userId": 105,
      "username": "member56",
      "email": "member55@gmail.com",
      "createdAt": "2025-03-15T15:43:49.217183Z",
      "updateAt": "2025-03-15T15:43:49.217183Z"
    }
  }
  ```

## Command + SQL Queries

### Run project in local

```
  go run ./cmd
```

### Enter postgre command

```
  psql --username=postgres
```

- next enter the password of postgres

### Select Database

```
  format:
    \c [name_database];
```

```
  example:
    \c pokemon;
```

### Create Table

```
  format:
    CREATE TABLE [name_table] (
        [name_field] [type_field],
        [name_field] [type_field],
    )
```

```
  example:
    CREATE TABLE users (
        user_id serial primary key,
        username varchar(50) unique not null,
        password varchar(255) not null,
        email varchar(100) unique not null,
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp
    )
```

### Delete Table

```
  format:
    DROP TABLE [name_table];
```

```
  example:
    DROP TABLE users;
```

### Show List Table

```
  command:
    \d
```

### Insert Data

```
  format:
    INSERT INTO [name_table] ([name_field_1], [name_field_1]) VALUES
    ([value_field_1], [value_field_2]),
    ([value_field_1], [value_field_2]),;
```

```
  example:
    INSERT INTO users (username, email) VALUES
    ('john_doe', 'john.doe@example.com'),
    ('jane_smith', 'jane.smith@example.com'),
    ('alice_johnson', 'alice.johnson@example.com');
```

### Update Data

```
  format:
    UPDATE [name_table] SET [name_field_1] = [value_field_1], [name_field_2] = [value_field_2] WHERE [name_field] = [value_field];
```

```
  example:
    UPDATE users SET username = 'bunbun', email = 'bunbun@gmail.com' WHERE user_id = 11;
```

### Delete 1 Data

```
  format:
    DELETE from [name_table] where [name_field] = [value_field];
```

```
  example:
    DELETE from users WHERE user_id = 2;
```

### Delete All Data from table

```
  format:
    DELETE from [name_table]
```

```
  example:
    DELETE from users;
```

## Notes

- docs http respond code : https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
- docs golang about http respond code : https://go.dev/src/net/http/status.go
- **`Database pooling`**: To improves application performance by reusing existing connections to a database, rather than creating new ones for every request. This helps manage resources efficiently and ensures faster response times
- **`db.QueryContext`**: To get multiple rows of results
- **`db.QueryRowContext`**: To get one row of results
- **`db.ExecContext`**: To run SQL commands that do not produce result rows (such as `INSERT`, `UPDATE`, `DELETE`)
- **`db.BeginTx`**: To start a transaction. This is useful when you want to run multiple SQL commands as one atomic unit (all or nothing)
- **`db.PrepareContext`**: To prepare SQL statements that can be executed multiple times with different parameters, which can improve performance
- **`Repository Pattern`** is a design approach that separates the data access logic from the business logic. It provides an abstraction layer, allowing our application to interact with data sources (like databases) without needing to know the details of how that data is stored or retrieved
- **`bycrpt`**: is a library that simplifies the process of hashing and verifying passwords, allowing us to implement secure password handling without needing to understand the underlying cryptographic principles in depth
- **`jwt-go`**: is a library that simplifies the implementation of secure, stateless authentication in our Go applications, providing a flexible and efficient way to manage user sessions.
- **`rate limit`**: Rate limiting is a best practice for maintaining the integrity, performance, and reliability of web applications and APIs. It protects both the service provider and the users by managing how resources are consumed

```
structure folders based on repository pattern

our_project
│
├── cmd/
│   └── main.go                  # Entry point of the application
│
├── config/
│   └── config.go                # Configuration loading (e.g., loading .env)
│
├── handler/
│   └── user_handler.go          # HTTP handlers for user-related operations
│
├── middleware/
│   └── jwt.go                   # Check header authorization
│
├── models/
│   └── user.go                  # User model definition
│
├── repository/
│   ├── user_repository.go       # User repository interface
│   └── user_repository_impl.go  # Implementation of the user repository
│
├── utils/
│   └── response.go              # Utility functions for writing JSON responses
│
├── migrations/
│   └── migration_script.sql      # Database migration scripts (if any)
│
├── tests/
│   └── user_handler_test.go      # Unit tests for user handler
│
└── go.mod                       # Go module file

```
