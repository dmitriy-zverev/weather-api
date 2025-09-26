docker run -d --name redis --network weather-api-network redis:7
docker run -d --name weather-api -p 8080:8080 --platform linux/amd64 --network weather-api-network weather-api:latest
