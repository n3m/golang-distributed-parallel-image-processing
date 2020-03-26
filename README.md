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
curl -X POST -H "Authorization: Bearer <ACCESS_TOKEN>" localhost:8080/logout
curl -H "Authorization: Bearer <ACCESS_TOKEN>" localhost:8080/status
curl -F "data=@image.jpg" -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/upload

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODU0NTE2OTEsInVzZXIiOiJhZG1pbiJ9.wGtG-F58W2Gq_Zkk8XouPrN27xClGhbEa8kNwMALv3s