# K8s Housekeeper Operator

A Kubernetes operator for monitoring and managing ExternalName services in your cluster.

## Overview

The K8s Housekeeper Operator is designed to monitor and track ExternalName services in your Kubernetes cluster. It provides a custom resource `ServiceMonitor` that allows you to:

- Monitor services with specific ExternalName
- Track the count of matching services
- Get real-time status updates
- Configure check intervals

## Features

- Monitor ExternalName services in real-time
- Customizable check intervals
- Status tracking and reporting
- Kubernetes-native CRD integration
- Prometheus metrics support

## Installation

### Prerequisites

- Kubernetes cluster (v1.16+)
- kubectl configured to access your cluster
- cert-manager (for webhook support)

### Deploy the Operator

1. Install the CRDs:
```bash
make install
```

2. Deploy the operator:
```bash
make deploy
```

## Usage

### Create a ServiceMonitor

Create a ServiceMonitor to monitor services with a specific ExternalName:

```yaml
apiVersion: monitoring.cluster.local/v1alpha1
kind: ServiceMonitor
metadata:
  name: example-monitor
  namespace: default
spec:
  targetDomain: example.com
  checkInterval: 60  # Check every 60 seconds
```

### Check Status

View the status of your ServiceMonitor:

```bash
kubectl get servicemonitors.monitoring.cluster.local -A
```

## Development

### Prerequisites

- Go 1.19+
- Docker
- kubebuilder
- kustomize

### Build

```bash
make build
```

### Test

```bash
make test
```

### Run Locally

```bash
make run
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

