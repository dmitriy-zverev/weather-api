GOOS=linux GOARCH=arm64 go build -o ../weather-api .

docker build --platform linux/amd64 -t weather-api:latest .