# Local Deployment Guide

## Prerequisites
- Docker installed and running
- Helm installed
- Kubernetes cluster running locally (e.g., minikube, kind, or Docker Desktop with Kubernetes enabled)
- GitHub Personal Access Token (PAT) with `read:packages` and `write:packages` permissions

## Deployment Steps

1. Build and push the Docker image to GitHub Container Registry:
```bash
# Log in to GitHub Container Registry
echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GITHUB_USERNAME" --password-stdin

# Create and use a new builder instance
docker buildx create --use

# Build and push multi-architecture image
docker buildx build --platform linux/amd64,linux/arm64 \
    -t ghcr.io/$GITHUB_USERNAME/hello-service:latest \
    --push .
```

2. Deploy the application using Helm:
```bash
# Create the dev namespace if it doesn't exist
kubectl create namespace dev --dry-run=client -o yaml | kubectl apply -f -

# Create a secret for pulling images from GitHub Container Registry
kubectl create secret docker-registry ghcr-secret \
    --docker-server=ghcr.io \
    --docker-username=$GITHUB_USERNAME \
    --docker-password=$GHCR_TOKEN \
    --namespace=dev

# Deploy the application
helm upgrade --install hello-service ./helm/hello-service \
    --namespace dev \
    --values ./helm/hello-service/values.yaml \
    --values ./helm/hello-service/values.dev.yaml \
    --set imagePullSecrets[0].name=ghcr-secret
```

3. Verify the deployment:
```bash
# Check the pods
kubectl get pods -n dev

# Check the services
kubectl get services -n dev

# Check the ingress
kubectl get ingress -n dev
```

## Accessing the Application
Once deployed, you can access the application through the ingress URL:
- Development: http://hello-dev.yourdomain.com

Note: Make sure to replace `yourdomain.com` with your actual domain or configure your local DNS/hosts file accordingly.

## Environment Variables
Set these environment variables before running the commands:
```bash
export GITHUB_USERNAME="your-github-username"
export GHCR_TOKEN="your-github-pat-token"
```
