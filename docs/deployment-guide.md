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

## How to Deploy

### Creating a New Release

1. Create a new version tag:
   ```bash
   git tag v1.2.3  # Use semantic versioning
   git push origin v1.2.3
   ```

2. The deployment will automatically start when you push the tag

