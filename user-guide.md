# Distributed Parallel Image Processing Application

## Language

- Golang

## Contributors

- Alan Enrique Maldonado Navarro
- Guillermo Gonzalez Mena

### Dependencies Commands:

- go get -u github.com/labstack/echo/...
- go get github.com/dgrijalva/jwt-go

## Commands

- go run . --controller tcp://localhost:40899 --node-name Miranda --tags gpu,nvidia,assets,static --image-store-endpoint localhost:8080 --image-store-token t0k3n-01010

- go run . --controller tcp://localhost:40899 --node-name Ciry --tags gpu,nvidia,assets,static --image-store-endpoint localhost:8080 --image-store-token t0k3n-11125

- curl -u admin:password http://localhost:8080/login

- curl -F "data=boom.png" -H "Authorization: Bearer <token>" http://localhost:8080/upload

- curl -H "Authorization: Bearer <token>" http://localhost:8080/status

- curl -H "Authorization: Bearer <token>" http://localhost:8080/status/<Workername>

- curl -H "Authorization: Bearer <token>" http://localhost:8080/workloads/test

- curl -H "Authorization: Bearer <token>" http://localhost:8080/logout
