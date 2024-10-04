resource "kind_cluster" "lw-cluster" {
  name            = "lw-cluster"
  node_image      = "kindest/node:v1.30.0"
  kubeconfig_path = pathexpand("~/.kube/config")
  wait_for_ready  = true

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role = "control-plane"
      extra_port_mappings {
        container_port = 30092
        host_port      = 30092
      }
    }

    node {
      role = "worker"
    }
  }
}
