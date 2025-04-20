# Golang Auth JWT
---
---

## Installation

#### Clone project

```
$ git clone https://github.com/jagoankode/golang-auth-jwt.git
$ cd golang-auth-jwt
```
#### Install dependency
```
$ go mod tidy
```
#### Run api
```
$ go run main.go
```
```bash
Server API will run in http://localhost:8080
```
### Call API

#### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username": "yourusername", "password": "yourpassword"}'
```
#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "yourusername", "password": "yourpassword"}'
```
#### Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "your_refresh_token_here"}'
```

