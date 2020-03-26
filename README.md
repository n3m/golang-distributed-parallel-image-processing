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
curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODU0NDExNzAsInVzZXIiOiJhZG1pbiJ9.9vp0BTCLNupYmY6HtOVxLGkuD3ePNpX9NT6uH1FqB3c" localhost:8080/logout