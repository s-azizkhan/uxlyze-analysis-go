docker build -t uxlyze-analyzer-go . 
docker run --env-file .env -p 8080:8080 uxlyze-analyzer-go