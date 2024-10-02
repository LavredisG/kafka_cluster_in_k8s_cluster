# Kafka-Kubernetes cluster

This project sets up a Kafka broker on a Kubernetes cluster using Terraform and Kubernetes manifests.

## Prerequisites

- Terraform
- kubectl
- kind (Kubernetes in Docker)

## Setup

1. **Create the Kind Cluster**:
   \`\`\`sh
   cd terraform
   terraform init
   terraform apply
   \`\`\`

2. **Deploy Kafka**:
   \`\`\`sh
   cd ../k8s
   kubectl apply -f kafka-deployment.yaml
   \`\`\`

## Structure

- \`terraform/\`: Contains Terraform configuration files for setting up the Kind cluster.
- \`k8s/\`: Contains Kubernetes manifests for deploying Kafka.
