# Kafka-Kubernetes cluster

This project sets up a single broker Kafka cluster on a 2 node Kubernetes cluster
and deploys a webapp which consumes messages from the broker's topic in order to
display them in real time on our browser using websockets on localhost:8080.

## Prerequisites

- [Terraform (v1.9.6)](https://developer.hashicorp.com/terraform/install?product_intent=terraform)
- [kubectl (v1.30.3)](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
- [kind (v0.23.0) -> (deploys k8s v1.30.0 images)](https://kind.sigs.k8s.io/)
- [Helm (v3.16.1)](https://helm.sh/docs/intro/install/)
- [Go (v1.22.6)](https://go.dev/doc/install)

## Setup

First of all we have to clone the repo on our system (project was set up on Ubuntu 22.04
instance running on a WSL2 machine on Windows 10).

1. **Create the Kind Cluster**:
   ```sh
      cd terraform
      terraform init
      terraform apply
   ```

We now have a cluster consisting of a control plane and a worker node:

![github1](https://github.com/user-attachments/assets/c961adb0-859c-4489-981c-e194ae2c9de0)

We should create a <b>kafka</b> namespace for easier management of our resources at this point:

```sh
kubectl create namespace kafka
```

2. **Deploy Kafka**:
   ```sh
   helm repo add bitnami https://charts.bitnami.com/bitnami
   helm install my-kafka oci://registry-1.docker.io/bitnamicharts/kafka -n kafka -f helm/values.yaml
   ```

We can now view our release:

![image](https://github.com/user-attachments/assets/2d4e1f67-7568-493c-b1ad-1810517e8940)

and the resources it created in the <b>kafka</b> namespace:

![image](https://github.com/user-attachments/assets/96061b74-0ed2-4b8b-9658-de17793dc9b7)

To continue, we read the notes of our chart deployed. We can view them anytime via the
<i><b>helm get notes my-kafka</i></b> command.

First, we will create a pod to use as a Kafka client:
```sh
kubectl run my-kafka-client --restart='Never' --image docker.io/bitnami/kafka:3.8.0-debian-12-r5 --namespace kafka --command -- sleep infinity
kubectl exec --tty -i my-kafka-client --namespace kafka -- bash
```
Once we are inside the <b>my-kafka-client</b> pod we can create a topic named <b>test</b>
and then verify its creation with these commands:
```sh
kafka-topics.sh --create --topic test --bootstrap-server my-kafka:9092
kafka-topics.sh --list --bootstrap-server my-kafka:9092
```
![image](https://github.com/user-attachments/assets/949defb3-e3ae-44a4-91f9-b64d8ab537fc)

3. **Deploy webapp**:

We need to know how to connect to our controller+broker node from outside of the cluster,
so we follow the notes instructions once again to get that information:
```sh
kubectl get pods --namespace kafka -l "app.kubernetes.io/name=kafka,app.kubernetes.io/instance=my-kafka"
kubectl exec -it my-kafka-controller-0 -- cat /opt/bitnami/kafka/config/server.properties | grep advertised.listeners
```

![image](https://github.com/user-attachments/assets/c8a0f509-10ae-474f-b7d1-02964a9a7c60)

We can see that the external advertised listener in on 172.18.0.2:31551,
(<node IP on which the kafka broker pod is deployed>:<exposed PORT>). However,
our app will target 172.18.0.3:31551 which is Control Plane node's IP address, while
172.18.0.2 is Worker node's IP address on which kafka pod is deployed. Kubernetes
networking will take care of properly forwarding our request to the right node

![image](https://github.com/user-attachments/assets/ce31973d-be48-4480-a3b8-e478be5973b3)

To deploy the server:
```sh
cd webapp
go run app.go
```

The server now listens on <b>localhost:8080</b>.

4. **Real time message producing/consuming**:

We open a second terminal on which we have to log back in the kafka-client pod
and execute the producer script:

```sh
kubectl exec --tty -i my-kafka-client --namespace kafka -- bash
kafka-console-producer.sh --broker-list my-kafka-controller-0.my-kafka-controller-headless.kafka.svc.cluster.local:9092 --topic test
```

This is our final result taking place in real time:

![image](https://github.com/user-attachments/assets/4f218b71-d3de-448d-bb3e-b7ffde192e59)



## Troubleshooting

- In case you already had a cluster, you may have to execute <i><b>docker system prune</b></i>
so that every resource related to the previous cluster has being deleted.





   

