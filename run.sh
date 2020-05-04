./worker --controller tcp://localhost:40899 --worker-name Ciry --tags gpu,nvidia, assets, static
./worker --controller tcp://localhost:40899 --worker-name Miranda --tags gpu,nvidia, assets, static

curl -u admin:password http://localhost:8080/login
curl -F "data=boom.png" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg3MTQ2MzIsInVzZXIiOiJhZG1pbiJ9.90olqfqerEKJiZ8j-urJHR4w7RoyauIRttrn3UNlfG4" http://localhost:8080/upload
curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status
curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/logout