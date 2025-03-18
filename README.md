# Chat Go

Chat Go is a lightweight, efficient chat application built using Golang, 
designed for seamless integration as a microservice. Leveraging Golang's 
concurrency features and robust performance, it provides a scalable and 
reliable solution for real-time communication within distributed systems.

## Integration Guide

To integrate Chat Go into your system, implement these two endpoints:

1. **Current User Endpoint**: Handles authorization and retrieves current 
user details (ID, email, username, first name, last name, about me, image).

Request example:

```json
{
  "httpRequest": {
    "method": "GET",
    "path": "http://localhost:3000/users/current",
    "headers": {
      "Authorization": [
        "Bearer {AUTH_TOKEN}"
      ]
    }
  }
}
```

Response Example:

```json
{
  "httpResponse": {
    "headers": {
      "Content-Type": [
        "application/json"
      ]
    },
    "body": {
      "id": 1,
      "email": "john.doe@gmail.com",
      "username": "user1",
      "firstName": "John",
      "lastName": "Doe",
      "aboutMe": "Info",
      "image": {
        "url": "",
        "base64": ""
      }
    },
    "statusCode": 200
  }
}
```

2. **Users Search Endpoint**: Enables user search by filters and enriches chat 
and message data with createdBy user information.

```json
{
  "httpRequest": {
    "method": "GET",
    "path": "http://localhost:3000/users",
    "headers": {
      "Authorization": [
        "Bearer {AUTH_TOKEN}"
      ]
    }
  }
}
```

Response Example:

```json
{
  "httpResponse": {
    "headers": {
      "Content-Type": ["application/json"]
    },
    "body": {
      "items": [
        {
          "id": 1,
          "email": "john.doe@gmail.com",
          "username": "john.doe",
          "firstName": "John",
          "lastName": "Doe",
          "aboutMe": "Info",
          "image": {
            "url": "",
            "base64": ""
          }
        },
        {
          "id": 2,
          "email": "alex.carter@gmail.com",
          "username": "alex.carter",
          "firstName": "Alex",
          "lastName": "Carter",
          "aboutMe": "Info",
          "image": {
            "url": "",
            "base64": ""
          }
        }
      ],
      "count": 2
    },
    "statusCode": 200
  }
}
```