# Deployment Guide

This guide explains how to set up and use the CI/CD pipeline for deploying the hello-service to the development environment.

## Overview

The deployment process is automated using GitHub Actions and triggered by creating a new version tag. The pipeline:
1. Builds a Docker image
2. Pushes it to GitHub Container Registry (GHCR)
3. Deploys the application to your k3s cluster using Helm

## Prerequisites

1. A running k3s cluster
2. GitHub repository with the code
3. Access to configure GitHub repository settings

## Configuration Steps

### 1. Set Up GitHub Secrets

Go to your repository's Settings > Secrets and variables > Actions and add the following secret:

- `KUBE_CONFIG`: Your kubeconfig file content (base64 encoded)
  ```bash
  # Get your kubeconfig and encode it
  cat ~/.kube/config | base64
  ```

### 2. Configure k3s Access

1. On your k3s server, get the kubeconfig:
   ```bash
   sudo cat /etc/rancher/k3s/k3s.yaml
   ```

2. Update the server URL in the kubeconfig to use your server's IP address:
   ```yaml
   server: https://<your-server-ip>:6443
   ```

3. Make sure your server's firewall allows access to port 6443 (k3s API)

### 3. Update Configuration Files

1. Update `helm/hello-service/values.dev.yaml` with your domain:
   ```yaml
   ingress:
     hosts:
       - host: hello-dev.yourdomain.com  # Change this to your domain
   ```

## How to Deploy

### Creating a New Release

1. Create a new version tag:
   ```bash
   git tag v1.2.3  # Use semantic versioning
   git push origin v1.2.3
   ```

2. The deployment will automatically start when you push the tag

### Manual Approval

The deployment requires manual approval:
1. Go to your repository's Actions tab
2. Find the running workflow for your tag
3. Click "Review deployments"
4. Approve the deployment to the dev environment

## What Happens During Deployment

1. **Build Stage**:
   - Checks out the code
   - Sets up Go environment
   - Builds Docker image
   - Pushes image to GHCR with tag (e.g., v1.2.3)

2. **Deploy Stage**:
   - Installs kubectl and Helm
   - Configures kubeconfig
   - Creates dev namespace if needed
   - Deploys using Helm with values.dev.yaml

## Configuration Files

### values.dev.yaml
```yaml
replicaCount: 1

image:
  repository: krl4kk/hello-service
  tag: latest
  pullPolicy: Always

service:
  type: ClusterIP
  port: 80
  targetPort: 3000

ingress:
  enabled: true
  className: "traefik"
  hosts:
    - host: hello-dev.yourdomain.com
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

## Troubleshooting

### Common Issues

1. **Deployment Fails**
   - Check the GitHub Actions logs
   - Verify kubeconfig is correct
   - Ensure server IP is accessible

2. **Image Pull Errors**
   - Verify GHCR access
   - Check image tag exists in GHCR

3. **Ingress Not Working**
   - Verify domain DNS is pointing to your server
   - Check traefik logs: `kubectl logs -n kube-system -l app.kubernetes.io/name=traefik`

### Useful Commands

```bash
# Check deployment status
kubectl get pods -n dev
kubectl get ingress -n dev

# View logs
kubectl logs -n dev -l app=hello-service

# Describe resources
kubectl describe pod -n dev -l app=hello-service
kubectl describe ingress -n dev hello-service
```

## Security Considerations

1. The kubeconfig is stored as a GitHub secret and is base64 encoded
2. GitHub Actions uses temporary credentials for GHCR
3. Manual approval is required for deployments
4. All sensitive data is stored in GitHub secrets

## Maintenance

### Updating Configuration

1. Modify `values.dev.yaml` for configuration changes
2. Changes will be applied on the next deployment

### Cleaning Up

To remove old images from GHCR:
1. Go to your repository's Packages
2. Select the hello-service package
3. Delete old versions as needed

## Support

If you encounter any issues:
1. Check the GitHub Actions logs
2. Verify k3s cluster status
3. Review the troubleshooting section above 
