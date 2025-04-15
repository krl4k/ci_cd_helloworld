# Server Setup Guide for k3s and GitHub Actions Runner

## Server Requirements

### Hardware Requirements
- CPU: 2+ cores (4+ recommended for production)
- RAM: 4GB minimum (8GB recommended)
- Storage: 20GB minimum (SSD recommended)
- Network: Stable internet connection

### Software Requirements
- Operating System: Ubuntu Server 22.04 LTS (recommended)
- Docker
- kubectl
- Helm
- Go (version 1.21 or later)

## Initial Server Setup

### User Management

1. Create a dedicated user for k3s and runner:
   ```bash
   sudo useradd -m -s /bin/bash infra-user
   sudo usermod -aG sudo infra-user
   ```

2. Set up SSH access:
   ```bash
   sudo mkdir -p /home/infra-user/.ssh
   sudo chmod 700 /home/infra-user/.ssh
   sudo touch /home/infra-user/.ssh/authorized_keys
   sudo chmod 600 /home/infra-user/.ssh/authorized_keys
   sudo chown -R infra-user:infra-user /home/infra-user/.ssh
   // 
   ```

### Copy Key from Root(not good)

1. Copy the SSH public key from the root user to the new infra-user:
   ```bash
   sudo cp /root/.ssh/authorized_keys /home/infra-user/.ssh/authorized_keys
   sudo chown infra-user:infra-user /home/infra-user/.ssh/authorized_keys
   ```
   ```

### Security Configuration

#### Firewall Configuration

Required Open Ports:
- 22222 (SSH) - For remote access (non-standard port)
- 443 (HTTPS)
- 80 (HTTP)
- 6443 (k3s API server)

Recommended Firewall Rules:
```bash
# Install UFW if not present
apt install ufw

ufw disable

# Allow SSH on non-standard port
ufw allow 22222/tcp

# Allow HTTP/HTTPS
ufw allow 80/tcp
ufw allow 443/tcp

# Allow k3s ports
ufw allow 6443/tcp #apiserver
ufw allow from 10.42.0.0/16 to any #pods
ufw allow from 10.43.0.0/16 to any #services

# Enable firewall
ufw enable
```

## Software Installation

### Docker Installation
```bash
# Install Docker
sudo apt update
sudo apt install -y docker.io
sudo systemctl enable docker
sudo systemctl start docker

# Add user to docker group
sudo usermod -aG docker infra-user
```

### k3s Installation
```bash
# Install k3s
curl -sfL https://get.k3s.io | sh -

# Verify installation
sudo k3s kubectl get nodes
```

### GitHub Actions Runner Installation

1. Create runner directory:
   ```bash
   sudo mkdir -p /opt/github-runner
   sudo chown infra-user:infra-user /opt/github-runner
   ```

2. Switch to infra-user:
   ```bash
   sudo su - infra-user  # Switch to infra-user
   cd /opt/github-runner
   ```

3. Download and configure runner:
   ```bash
   # Download runner
   curl -o actions-runner-linux-x64-2.311.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz
   
   # Extract
   tar xzf actions-runner-linux-x64-2.311.0.tar.gz
   
   # Configure
   
   ./config.sh --url https://github.com/your-username/your-repo --token YOUR_TOKEN
   ```

4. Create systemd service:
   ```bash
   sudo nano /etc/systemd/system/github-runner.service
   ```

   Add the following content:
   ```ini
   [Unit]
   Description=GitHub Actions Runner
   After=network.target

   [Service]
   Type=simple
   User=infra-user
   WorkingDirectory=/opt/github-runner
   ExecStart=/opt/github-runner/run.sh
   Restart=always

   [Install]
   WantedBy=multi-user.target
   ```

5. Enable and start the service:
   ```bash
   sudo systemctl enable github-runner
   sudo systemctl start github-runner
   ```

## Monitoring and Maintenance

### Logs
```bash
# View k3s logs
sudo journalctl -u k3s -f

# View runner logs
sudo journalctl -u github-runner -f

# View Docker logs
sudo docker logs <container_id>
```

### Updates
1. k3s updates:
   ```bash
   curl -sfL https://get.k3s.io | sh -
   ```

2. Runner updates:
   ```bash
   cd /opt/github-runner
   sudo systemctl stop github-runner
   ./config.sh remove --token YOUR_TOKEN
   rm -rf *
   curl -o actions-runner-linux-x64-2.311.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz
   tar xzf actions-runner-linux-x64-2.311.0.tar.gz
   ./config.sh --url https://github.com/your-username/your-repo --token YOUR_TOKEN
   sudo systemctl start github-runner
   ```

### Backup
1. Create backup script:
   ```bash
   sudo nano /opt/backup.sh
   ```

   Add the following content:
   ```bash
   #!/bin/bash
   TIMESTAMP=$(date +%Y%m%d_%H%M%S)
   BACKUP_DIR="/opt/backups"
   mkdir -p $BACKUP_DIR
   
   # Backup k3s
   sudo k3s etcd-snapshot save --name k3s-backup-$TIMESTAMP
   
   # Backup runner
   tar -czf $BACKUP_DIR/runner-backup-$TIMESTAMP.tar.gz /opt/github-runner
   ```

2. Make script executable:
   ```bash
   sudo chmod +x /opt/backup.sh
   ```

## Troubleshooting

### Common Issues
1. k3s not starting:
   - Check logs: `sudo journalctl -u k3s -f`
   - Verify ports: `sudo netstat -tulpn | grep k3s`
   - Check service status: `sudo systemctl status k3s`

2. Runner not starting:
   - Check logs: `sudo journalctl -u github-runner -f`
   - Verify permissions: `ls -la /opt/github-runner`
   - Check service status: `sudo systemctl status github-runner`

3. Network issues:
   - Check connectivity: `ping github.com`
   - Verify firewall: `sudo ufw status`
   - Check DNS: `nslookup github.com`

### Maintenance Tasks
1. Regular updates:
   - Update system packages: `sudo apt update && sudo apt upgrade -y`
   - Update k3s and runner as needed
   - Monitor security advisories

2. Disk cleanup:
   - Clean Docker: `sudo docker system prune -a`
   - Remove old backups: `find /opt/backups -type f -mtime +30 -delete`

3. Security updates:
   - Regular security patches: `sudo unattended-upgrade`
   - Monitor security advisories
   - Regular password rotation
