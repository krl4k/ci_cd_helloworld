# Log in to Docker


echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

# Create and use a new builder instance
docker buildx create --use

# Build and push multi-architecture image
docker buildx build --platform linux/amd64,linux/arm64 \
    -t "$DOCKER_USERNAME/hello-service:latest" \
    --push .
