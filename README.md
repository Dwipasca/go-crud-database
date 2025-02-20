# GO-CRUD-DATABASE

A simple project to learn basic Go integration with a PostgreSQL database

## Functional requirment:

- as a user, I can create, update and delete data user
- as a user, I can see all the data users
- as a user, I can see detail data user by Id

## Third Libs I used

- Postgres : `github.com/lib/pq`

## Endpoint

### Create User

- URL : `http://localhost:8080/api/v1/users`
- Method: `POST`
- request body
  ```json
  {
    "username": "bunbun",
    "email": "bunbun@gmail.com"
  }
  ```
- response
  ```json
  {
    "message": "New user created successfully",
    "status": "success",
    "code": 201,
    "data": {
      "userId": 45,
      "username": "bunbun2",
      "email": "bunbun@gmail.com",
      "createdAt": "2025-02-19T17:57:26.729882Z"
    }
  }
  ```

### Update Data User

- URL : `http://localhost:8080/api/v1/users`
- Method: `PUT`
- request body
  ```json
  {
    "username": "bunbun",
    "email": "bunbun@gmail.com"
  }
  ```
- response

  ```json
  {
    "message": "User updated successfully",
    "status": "success",
    "code": 200,
    "data": {
      "userId": 13,
      "username": "tadadaadadsadxa",
      "email": "bunbun_builderzk@sandrock.com",
      "createdAt": "0001-01-01T00:00:00Z"
    }
  }
  ```

### Delete User By ID

- URL : `http://localhost:8080/api/v1/users?id=14`
- Method: `DELETE`
- Response
  ```json
  {
    "message": "Id User 14 deleted successfully",
    "status": "success",
    "code": 200
  }
  ```

### Get All Data User

- URL : `http://localhost:8080/api/v1/users`
- Method: `GET`
- Response
  ```json
  {
    "message": "successfully retrieved data from cache",
    "status": "success",
    "code": 200,
    "data": [
      {
        "user_id": 6,
        "username": "grace2",
        "email": "grace2@sandrock.com",
        "CreatedAt": "2025-01-22T07:53:02.327409Z"
      },
      {
        "user_id": 10,
        "username": "user1",
        "email": "bunbun_builderk@sandrocks.com",
        "CreatedAt": "2025-01-22T10:25:20.435874Z"
      },
      {
        "user_id": 11,
        "username": "'); truncate users; --",
        "email": "bunbun_builderzk@sandrock.com",
        "CreatedAt": "2025-01-22T11:56:14.120701Z"
      }
    ]
  }
  ```

### Get User By ID

- URL : `http://localhost:8080/api/v1/users?id=16`
- Method: `GET`
- Response
  ```json
  {
    "message": "Successfully get detail data",
    "status": "success",
    "code": 200,
    "data": {
      "userId": 16,
      "username": "bunbunsh",
      "email": "bunbun_buxilderh@sandrock.com",
      "createdAt": "2025-01-22T18:01:00.359879Z"
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
        email varchar(100) unique not null,
        created_at timestamp default current_timestamp
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
    UPDATE users SET username = 'new_username', email = 'new_email@example.com' WHERE user_id = 123;
```

```
  example:
    UPDATE users SET username = 'bunbun', email = 'bunbun@gmail.com' WHERE user_id = 11;
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
├── models/
│   └── user.go                  # User model definition
│
├── repository/
│   ├── user_repository.go       # User repository interface
│   └── user_repository_impl.go  # Implementation of the user repository
│
├── handler/
│   └── user_handler.go          # HTTP handlers for user-related operations
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
