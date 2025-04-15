# K3s Deployment Setup Guide

## Prerequisites

- Ubuntu Server 20.04 LTS or later
- Root or sudo access
- Minimum 2GB RAM
- Minimum 2 CPU cores
- 20GB free disk space

## Installation Steps

### 1. Server Preparation

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

### 2. K3s Installation

```bash
# Install k3s
curl -sfL https://get.k3s.io | sh -

# Verify installation
sudo systemctl status k3s
sudo kubectl get nodes

```

### 3. Local Development Setup

```bash
# Install kubectl
curl -LO "https://dl.k8s.io/release/stable.txt"
curl -LO "https://dl.k8s.io/release/$(cat stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Install Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

### 4. Configure Local Access

```bash
# On your server
sudo cat /etc/rancher/k3s/k3s.yaml

# On your local machine, create ~/.kube/config with the above content
# Replace the server IP with your actual server IP address
# Example: server: https://your-server-ip:6443
```

## Verification

```bash
# Check cluster status
kubectl get nodes
kubectl get pods -A

# Check Helm installation
helm version
```
