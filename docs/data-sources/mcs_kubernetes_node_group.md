---
layout: "mcs"
page_title: "mcs: kubernetes_node_group"
sidebar_current: "docs-kubernetes-node-group"
description: |-
  Get information on clusters node group.
---

# mcs\_kubernetes\_cluster

Use this data source to get the ID of an available MCS kubernetes clusters node group.

## Example Usage
```
data "mcs_kubernetes_node_group" "mynodegroup" {
  uuid = "mynguuid"
}
```

## Argument Reference

The following arguments are supported:

* `uuid` - (Required) The UUID of the cluster's node group.

    
## Attributes
`id` is set to the ID of the found cluster template. In addition, the following
attributes are exported:

* `name` - The name of the node group.
* `cluster_id` - The UUID of cluster that node group belongs.
* `node_count` - The count of nodes in node group.
* `max_nodes` - The maximum amount of nodes in node group.
* `min_nodes` - The minimum amount of nodes in node group.
* `volume_size` - The amount of memory of volume in Gb.
* `volume_type` - The type of volume.
* `flavor_id` - The id of flavor.
* `autoscaling_enabled` - Determines whether the autoscaling is enabled.
* `uuid` - The UUID of the cluster's node group.
* `state` - Determines current state of node group (RUNNING, SHUTOFF, ERROR).
* `nodes` - The list of node group's node objects.
