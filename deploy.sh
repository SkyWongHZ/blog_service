
docker stop blog-service
docker rm   blog-service

docker run -d --network my-network  --name blog-service -p 8080:8080  --restart always  blog-service:2.0