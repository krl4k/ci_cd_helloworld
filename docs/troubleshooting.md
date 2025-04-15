# Troubleshooting Guide

## Common Issues and Solutions

### 1. kubectl cannot connect to the cluster

**Symptoms:**
- `kubectl get nodes` returns connection errors
- Unable to access the Kubernetes API

**Solutions:**
1. Check your kubeconfig file:
   ```bash
   cat ~/.kube/config
   ```
2. Verify server IP/hostname is correct
3. Check firewall settings:
   ```bash
   sudo ufw status
   ```
4. Ensure k3s service is running:
   ```bash
   sudo systemctl status k3s
   ```

### 2. Pods stuck in Pending state

**Symptoms:**
- Pods remain in Pending state
- `kubectl describe pod` shows resource constraints

**Solutions:**
1. Check node resources:
   ```bash
   kubectl describe node
   ```
2. Check events:
   ```bash
   kubectl get events
   ```
3. Verify resource requests in your deployment:
   ```bash
   kubectl describe deployment <deployment-name>
   ```

### 3. CI/CD pipeline failures

**Symptoms:**
- GitLab pipeline jobs failing
- Deployment errors in CI/CD logs

**Solutions:**
1. Verify GitLab variables:
   - Check if `KUBE_CONFIG` is properly set
   - Ensure registry credentials are correct
2. Check runner connectivity:
   ```bash
   sudo gitlab-runner status
   ```
3. Validate Dockerfile and Helm chart:
   ```bash
   docker build -t test-image .
   helm lint ./helm/hello-service
   ```

### 4. Ingress not working

**Symptoms:**
- Services not accessible through ingress
- 404 errors when accessing endpoints

**Solutions:**
1. Check ingress configuration:
   ```bash
   kubectl describe ingress <ingress-name>
   ```
2. Verify Traefik is running:
   ```bash
   kubectl get pods -n kube-system | grep traefik
   ```
3. Check service endpoints:
   ```bash
   kubectl get endpoints <service-name>
   ```

## Useful Commands

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
