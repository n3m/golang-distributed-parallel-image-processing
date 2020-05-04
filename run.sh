./worker --controller tcp://localhost:40899 --worker-name Ciry --tags gpu,nvidia, assets, static
./worker --controller tcp://localhost:40899 --worker-name Miranda --tags gpu,nvidia, assets, static

curl -u admin:password http://localhost:8080/login
curl -F "data=boom.png" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg4MTIzMDgsInVzZXIiOiJhZG1pbiJ9.VC0gDP4yk6a8y0_i1ISxBr4t_8omMGr65h7I0Z0qTRI" http://localhost:8080/upload
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg4MTIzMDgsInVzZXIiOiJhZG1pbiJ9.VC0gDP4yk6a8y0_i1ISxBr4t_8omMGr65h7I0Z0qTRI" http://localhost:8080/status
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg4MTIzMDgsInVzZXIiOiJhZG1pbiJ9.VC0gDP4yk6a8y0_i1ISxBr4t_8omMGr65h7I0Z0qTRI" http://localhost:8080/workloads/test
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg4MTIzMDgsInVzZXIiOiJhZG1pbiJ9.VC0gDP4yk6a8y0_i1ISxBr4t_8omMGr65h7I0Z0qTRI" http://localhost:8080/logout