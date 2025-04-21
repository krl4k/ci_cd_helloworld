# PostgreSQL Installation Guide

This guide explains how to install PostgreSQL as a dependency for the hello-service application using Helm.

## Prerequisites

- Helm installed and configured
- Access to a Kubernetes cluster
- Appropriate permissions to create namespaces and deploy resources

## Installation Steps

1. Add PostgreSQL as a dependency in the chart file
2. Update Helm dependencies:
   ```bash
   helm dependency update helm/hello-service
   ```
3. Deploy to the cluster:
   ```bash
   cd helm/hello-service
   helm upgrade --install hello-service . --namespace dev --create-namespace
   ```

## Verification

After installation, you can verify the PostgreSQL deployment by checking the pods in the dev namespace:

```bash
kubectl get pods -n dev
```

## Notes

- The installation creates a new namespace called `dev` if it doesn't exist
- The PostgreSQL instance will be deployed as part of the hello-service chart
