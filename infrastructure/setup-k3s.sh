#!/bin/bash

# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y curl openssh-server

# Configure firewall
sudo ufw allow 6443/tcp  # Kubernetes API
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# Install k3s
curl -sfL https://get.k3s.io | sh -

# Verify installation
sudo systemctl status k3s
sudo kubectl get nodes 
