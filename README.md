# backend-test

## Requirements
*  Go 1.21 or higher. If you don't have Golang installed, you can download it from https://go.dev/doc/install)
*  Docker Compose (https://docs.docker.com/compose/install/)

## Setup

1. Start  app:
    1. Dockerized:
      
         Requires docker and docker-compose
      
          ```
          make local-up
          ```

    2. Run as a go app:

       Requires a mysql instance running
   
       ```
       make install
       make configure
       make start
       ```

The project should now be running at http://localhost:8088.

2. Run unit tests:

   ```
   make test
   ```
   
3. Run test coverage:

   ```
   make cover
   ```

## Usage

### Endpoints

The API has the following endpoints:

#### User

Create a user.
Required fields:
* name
* age
* email
* password
* address

##### Request

```
POST /v1/user
Content-Type: application/json

{
    "name": "John Doe",
    "age": 18,
    "email": "john.doe@gmail.com",
    "password": "password",
    "address": "200 Random St"
}
```

##### Response

```
HTTP/1.1 201 Created
Content-Type: application/json

{
    "id": 1,
    "name": "John Doe",
    "age": 18,
    "email": "john.doe@gmail.com",
    "address": "200 Random St"
}

```

##### Request

```
GET /v1/user/{id}
```

##### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "id": 1,
    "name": "John Doe",
    "age": 18,
    "email": "john.doe@gmail.com",
    "address": "200 Random St"
}

```

##### Request

```
GET /v1/user
```

##### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
        "id": 1,
        "name": "John Doe",
        "age": 18,
        "email": "john.doe@gmail.com",
        "address": "200 Random St"
    }
]

```

##### Request

```
PUT /v1/user/1
Content-Type: application/json

{
    "name": "John Stuart",
    "age": 50,
    "email": "john.doe@gmail.com",
    "password": "password",
    "address": "200 Random St"
}
```

##### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
        "id": 1,
        "name": "John Stuart",
        "age": 50,
        "email": "john.doe@gmail.com",
        "address": "200 Random St"
    }
]

```

##### Request

```
DELETE /v1/user/1
```

##### Response

```
HTTP/1.1 204 No Content
Content-Type: application/json

```

------------------

### Error Handling

If an error occurs, the API will return a JSON object with an error message:

```
{
    "developer_message": string,
    "user_message": string,
    "status_code": int
}
```

Possible HTTP status codes for errors include:

- `400 Bad Request` for invalid request data
- `401 Unauthorized` missing a valid authorization token
- `404 Not Found` resource not found
- `409 Conflict` username already present in the database
- `500 Internal Server Error` for server-side errors

## Swagger API Documentation

This project uses Swagger for API documentation. Swagger provides a user-friendly interface for exploring and testing the API.

To access the Swagger page:

1. Start the application if it's not already running.
2. Open a web browser and navigate to `http://localhost:8088/v1/swagger/index.html#/`.
3. The Swagger page should load, displaying a list of available endpoints.

From here, you can explore the available endpoints, see what parameters they require, and test them out.

If you have any questions or issues with the Swagger page, please refer to the API documentation or contact the project maintainers.