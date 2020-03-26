# Distributed Parallel Image Processing Application

## Language
- Golang

## Contributors
- Alan Enrique Maldonado Navarro
- Guillermo Gonzalez Mena

### Dependencies Commands:
- go get -u github.com/labstack/echo/...
- go get github.com/dgrijalva/jwt-go

curl -X POST -d "username=admin&password=password" http://localhost:8080/login
curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODU0NDIzOTgsInVzZXIiOiJhZG1pbiJ9.DJzp2ttD4tIALaIDUzzOBmBTqjDw9LsjnZgTK3ivcbE" localhost:8080/logout
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODU0NDIzOTgsInVzZXIiOiJhZG1pbiJ9.DJzp2ttD4tIALaIDUzzOBmBTqjDw9LsjnZgTK3ivcbE" localhost:8080/status