---
layout: "mcs"
page_title: "mcs: kubernetes_cluster"
sidebar_current: "docs-kubernetes-cluster"
description: |-
  Get information on cluster.
---

# mcs\_kubernetes\_cluster

Use this data source to get the ID of an available MCS kubernetes cluster.

## Example Usage
```hcl
data "mcs_kubernetes_cluster" "mycluster" {
  name = "myclustername"
}
```
```hcl
data "mcs_kubernetes_cluster" "mycluster" {
  cluster_id = "myclusteruuid"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the cluster.

* `cluster_id` - (Optional) The UUID of the Kubernetes cluster
    template.

* `region` - (Optional) The region in which to obtain the Container Infra
    client.
    If omitted, the `region` argument of the provider is used.
        
**Note**: Only one of `name` or `cluster_id` must be specified

    
## Attributes
`id` is set to the ID of the found cluster template. In addition, the following
attributes are exported:

* `name` - The name of the cluster.
* `project_id` - The project of the cluster.
* `created_at` - The time at which cluster was created.
* `updated_at` - The time at which cluster was created.
* `api_address` - COE API address.
* `cluster_template_id` - The UUID of the V1 Container Infra cluster template.
* `create_timeout` - The timeout (in minutes) for creating the cluster.
* `discovery_url` - The URL used for cluster node discovery.
* `flavor` - The ID of flavor for the nodes of the cluster.
* `master_flavor` - The ID of the flavor for the master nodes.
* `keypair` - The name of the Compute service SSH keypair.
* `labels` - The list of key value pairs representing additional properties of
                 the cluster.
* `master_count` - The number of master nodes for the cluster.
* `node_count` -  The number of nodes for the cluster.
* `master_addresses` - IP addresses of the master node of the cluster.
* `node_addresses` - IP addresses of the node of the cluster.
* `stack_id` - UUID of the Orchestration service stack.
* `network_id` - UUID of the cluster's network.
* `subnet_id` - UUID of the cluster's subnet.
* `status` - Current state of a cluster.
* `pods_network_cidr` - Network cidr of k8s virtual network
* `k8s_config` - Kubeconfig for cluster
