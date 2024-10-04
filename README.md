# Kafka-Kubernetes cluster

This project sets up a Kafka cluster with a single broker on a Kubernetes cluster
and deploys a webapp which consumes messages from the broker's topic in order to
display them in real time on our browser using websockets.

## Prerequisites

- Terraform (v1.9.6)
- kubectl (v1.30.3)
- kind (v0.23.0)
- Helm (v3.16.1)
- Go (v1.22.6)

## Setup

First of all we have to clone the repo on our system (project was set up on Ubuntu 22.04).

1. **Create the Kind Cluster**:
   ```sh
      cd terraform
      terraform init
      terraform apply
   ```

We now have a cluster consisting of a control plane and a worker node.

2. **Deploy Kafka**:
   ```sh
   cd ../k8s
   kubectl apply -f kafka-deployment.yaml
   ```

## Structure

- \`terraform/\`: Contains Terraform configuration files for setting up the Kind cluster.
- \`k8s/\`: Contains Kubernetes manifests for deploying Kafka.
