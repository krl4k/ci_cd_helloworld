# K3s Deployment & CI/CD Learning Project

A comprehensive guide to setting up a lightweight Kubernetes environment with k3s, implementing CI/CD with GitLab, and deploying applications using Helm.

## Project Overview

This project will walk you through:
- Setting up a k3s Kubernetes cluster
- Creating and deploying applications with Helm
- Implementing CI/CD pipelines with GitLab
- Managing your infrastructure as code

## Repository Structure

```
.
├── .gitlab-ci.yml                  # Main CI/CD pipeline configuration
├── infrastructure/                 # Infrastructure setup scripts
│   ├── setup-k3s.sh                # K3s installation script
│   └── setup-gitlab-runner.sh      # GitLab runner setup script
│
├── services/                       # Application services
│   ├── hello-service/              # Example service
│   │   ├── src/                    # Source code
│   │   │   ├── app.go              # Go application
│   │   │   └── Dockerfile          # Container image definition
│   │   └── .gitlab-ci.yml          # Service-specific CI/CD (optional)
│   │
│   └── other-service/              # Additional services follow same pattern
│
├── helm/                           # Helm charts for deployment
│   ├── hello-service/              # Helm chart for hello-service
│   │   ├── Chart.yaml              # Chart metadata
│   │   ├── values.yaml             # Default configuration
│   │   ├── values-prod.yaml        # Production-specific values
│   │   ├── values-staging.yaml     # Staging-specific values
│   │   └── templates/              # Kubernetes manifest templates
│   │       ├── deployment.yaml     # Deployment specification
│   │       ├── service.yaml        # Service specification
│   │       └── ingress.yaml        # Ingress configuration
│   │
│   └── other-service/              # Charts for other services
│
└── docs/                           # Documentation
    ├── setup-guide.md              # Detailed setup instructions
    └── troubleshooting.md          # Common issues and solutions
```

## Step-by-Step Implementation Guide

### Phase 1: Server Setup & K3s Installation

1. **Prepare Your Server**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install dependencies
   sudo apt install -y curl openssh-server
   
   # Configure firewall
   sudo ufw allow 6443/tcp  # Kubernetes API
   sudo ufw allow 80/tcp    # HTTP
   sudo ufw allow 443/tcp   # HTTPS
   sudo ufw enable
   ```

2. **Install K3s**
   ```bash
   # Install k3s server
   curl -sfL https://get.k3s.io | sh -
   
   # Verify installation
   sudo systemctl status k3s
   sudo kubectl get nodes
   ```

3. **Configure Local Access**
   ```bash
   # On your server
   sudo cat /etc/rancher/k3s/k3s.yaml
   
   # On your local machine, create ~/.kube/config with the above content
   # Replace the server IP with your actual server IP address
   # Example: server: https://your-server-ip:6443
   ```

### Phase 2: Local Development Environment

1. **Install Required Tools**
   ```bash
   # Install kubectl
   curl -LO "https://dl.k8s.io/release/stable.txt"
   curl -LO "https://dl.k8s.io/release/$(cat stable.txt)/bin/linux/amd64/kubectl"
   chmod +x kubectl && sudo mv kubectl /usr/local/bin/
   
   # Install Helm
   curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
   ```

2. **Verify Connectivity**
   ```bash
   kubectl get nodes  # Should show your k3s server
   helm version       # Should show Helm client and server version
   ```

### Phase 3: Create Your Application

1. **Create a Simple Application**
   ```bash
   # Create application structure
   mkdir -p services/hello-service/src
   cd services/hello-service/src
   ```

2. **Write Application Code (app.go)**
   ```go
   package main

   import (
       "fmt"
       "log"
       "net/http"
   )

   func helloHandler(w http.ResponseWriter, r *http.Request) {
       fmt.Fprintf(w, "Hello World from K3s!\n")
   }

   func main() {
       http.HandleFunc("/", helloHandler)
       
       port := ":3000"
       fmt.Printf("Server running at http://0.0.0.0%s/\n", port)
       log.Fatal(http.ListenAndServe(port, nil))
   }
   ```

3. **Create Dockerfile**
   ```dockerfile
   FROM golang:1.21-alpine
   WORKDIR /app
   COPY app.go .
   RUN go build -o hello-service
   EXPOSE 3000
   CMD ["./hello-service"]
   ```

### Phase 4: Create Helm Chart

1. **Create Helm Chart Structure**
   ```bash
   mkdir -p helm/hello-service
   cd helm
   helm create hello-service
   ```

2. **Customize values.yaml**
   ```yaml
   # Modify helm/hello-service/values.yaml
   replicaCount: 1
   
   image:
     repository: registry.gitlab.com/your-username/your-project/hello-service
     tag: latest
     pullPolicy: IfNotPresent
   
   service:
     type: ClusterIP
     port: 80
     targetPort: 3000
   
   ingress:
     enabled: true
     className: "traefik"
     hosts:
       - host: hello.yourdomain.com
         paths:
           - path: /
             pathType: Prefix
   
   resources:
     limits:
       cpu: 100m
       memory: 128Mi
     requests:
       cpu: 50m
       memory: 64Mi
   ```

### Phase 5: Set Up GitLab CI/CD

1. **Create .gitlab-ci.yml**
   ```yaml
   # In the root of your repository
   stages:
     - build
     - test
     - deploy
   
   variables:
     DOCKER_REGISTRY: registry.gitlab.com
     PROJECT_PATH: your-username/your-project
   
   .build-template: &build-definition
     stage: build
     image: docker:20.10.16
     services:
       - docker:20.10.16-dind
     before_script:
       - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
     only:
       - main
       - merge_requests
   
   .deploy-template: &deploy-definition
     stage: deploy
     image: dtzar/helm-kubectl:latest
     before_script:
       - mkdir -p ~/.kube
       - echo "$KUBE_CONFIG" > ~/.kube/config
       - chmod 400 ~/.kube/config
   
   build-hello-service:
     <<: *build-definition
     script:
       - docker build -t $CI_REGISTRY_IMAGE/hello-service:$CI_COMMIT_SHORT_SHA ./services/hello-service/src
       - docker push $CI_REGISTRY_IMAGE/hello-service:$CI_COMMIT_SHORT_SHA
     only:
       changes:
         - services/hello-service/**/*
   
   test-hello-service:
     stage: test
     image: node:16-alpine
     script:
       - cd services/hello-service
       - echo "Add your tests here"
     only:
       changes:
         - services/hello-service/**/*
   
   deploy-to-staging:
     <<: *deploy-definition
     script:
       - helm upgrade --install hello-staging ./helm/hello-service
         --set image.repository=$CI_REGISTRY_IMAGE/hello-service
         --set image.tag=$CI_COMMIT_SHORT_SHA
         --values ./helm/hello-service/values-staging.yaml
         --namespace staging
     only:
       - merge_requests
     environment:
       name: staging
   
   deploy-to-production:
     <<: *deploy-definition
     script:
       - helm upgrade --install hello-production ./helm/hello-service
         --set image.repository=$CI_REGISTRY_IMAGE/hello-service
         --set image.tag=$CI_COMMIT_SHORT_SHA
         --values ./helm/hello-service/values-prod.yaml
         --namespace production
     only:
       - main
     environment:
       name: production
     when: manual
   ```

2. **Set Up GitLab Variables**
   - In GitLab, go to Settings > CI/CD > Variables
   - Add `KUBE_CONFIG` variable with your kubeconfig (base64 encoded for security)

### Phase 6: Deploy to K3s

1. **Create Namespaces**
   ```bash
   kubectl create namespace staging
   kubectl create namespace production
   ```

2. **Deploy Manually for Testing**
   ```bash
   # Build and push your image
   docker build -t registry.gitlab.com/your-username/your-project/hello-service:test ./services/hello-service/src
   docker push registry.gitlab.com/your-username/your-project/hello-service:test
   
   # Deploy with Helm
   helm upgrade --install hello-test ./helm/hello-service \
     --set image.repository=registry.gitlab.com/your-username/your-project/hello-service \
     --set image.tag=test \
     --namespace default
   ```

3. **Verify Deployment**
   ```bash
   kubectl get pods
   kubectl get services
   kubectl get ingress
   ```

### Phase 7: Set Up GitLab Runner (Optional)

1. **Install GitLab Runner**
   ```bash
   # Install GitLab Runner
   curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh" | sudo bash
   sudo apt install gitlab-runner
   
   # Register the runner with your GitLab project
   sudo gitlab-runner register
   ```

2. **Configure Runner**
   - Enter your GitLab instance URL
   - Enter the registration token from your GitLab project
   - Add tags for your runner (e.g., `k3s`, `production`)
   - Choose `docker` as executor

## Advanced Features

### Add Monitoring

1. **Install Prometheus and Grafana**
   ```bash
   # Add Helm repository
   helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
   
   # Install monitoring stack
   helm install monitoring prometheus-community/kube-prometheus-stack \
     --namespace monitoring --create-namespace
   ```

### SSL with cert-manager

1. **Install cert-manager**
   ```bash
   # Add Helm repository
   helm repo add jetstack https://charts.jetstack.io
   
   # Install cert-manager
   helm install cert-manager jetstack/cert-manager \
     --namespace cert-manager --create-namespace \
     --set installCRDs=true
   ```

## Troubleshooting

### Common Issues

1. **kubectl cannot connect to the cluster**
   - Check your kubeconfig file
   - Ensure your server's IP/hostname is correct
   - Verify firewall settings

2. **Pods stuck in Pending state**
   - Check node resources: `kubectl describe node`
   - Check events: `kubectl get events`

3. **CI/CD pipeline failures**
   - Verify GitLab variables
   - Check runner connectivity
   - Validate Dockerfile and Helm chart

### Useful Commands

```bash
# View detailed pod information
kubectl describe pod <pod-name>

# Check pod logs
kubectl logs <pod-name>

# View all resources in a namespace
kubectl get all -n <namespace>

# Execute a command in a running container
kubectl exec -it <pod-name> -- /bin/sh

# Check K3s logs
sudo journalctl -u k3s

# Restart K3s
sudo systemctl restart k3s
```

## Additional Resources

- [K3s Documentation](https://rancher.com/docs/k3s/latest/en/)
- [Helm Documentation](https://helm.sh/docs/)
- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
