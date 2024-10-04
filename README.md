# Kafka-Kubernetes cluster

This project sets up a Kafka cluster with a single broker on a Kubernetes cluster
and deploys a webapp which consumes messages from the broker's topic in order to
display them in real time on our browser using websockets.

## Prerequisites

- Terraform (v1.9.6)
- kubectl (v1.30.3)
- kind (v0.23.0) -> (deploys k8s v1.30.0 images)
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

![github1](https://github.com/user-attachments/assets/c961adb0-859c-4489-981c-e194ae2c9de0)

We should create a <b>kafka</b> namespace for easier management of our resources at this point:

```sh
kubectl create namespace kafka
```

2. **Deploy Kafka**:
   ```sh
   helm repo add bitnami https://charts.bitnami.com/bitnami
   helm install my-kafka oci://registry-1.docker.io/bitnamicharts/kafka -n kafka -f helm/values.yml
   ```
   

