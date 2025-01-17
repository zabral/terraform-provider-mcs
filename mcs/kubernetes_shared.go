package mcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/mitchellh/mapstructure"
)

func extractKubernetesGroupMap(nodeGroups []interface{}) ([]NodeGroup, error) {
	filledNodeGroups := make([]NodeGroup, len(nodeGroups))
	for i, nodeGroup := range nodeGroups {
		g := nodeGroup.(map[string]interface{})
		for k, v := range g {
			if v == 0 {
				delete(g, k)
			}
		}
		var ng NodeGroup
		err := MapStructureDecoder(&ng, &g, config)
		if err != nil {
			return nil, err
		}
		filledNodeGroups[i] = ng
	}
	return filledNodeGroups, nil
}

func extractKubernetesLabelsMap(v map[string]interface{}) (map[string]string, error) {
	m := make(map[string]string)
	for key, val := range v {
		labelValue, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("label %s value should be string", key)
		}
		m[key] = labelValue
	}
	return m, nil
}

func extractNodeGroupLabelsList(v []interface{}) ([]NodeGroupLabel, error) {
	labels := make([]NodeGroupLabel, len(v))
	for i, label := range v {
		var L NodeGroupLabel
		err := mapstructure.Decode(label.(map[string]interface{}), &L)
		if err != nil {
			return nil, err
		}
		labels[i] = L
	}
	return labels, nil
}

func extractNodeGroupTaintsList(v []interface{}) ([]NodeGroupTaint, error) {
	taints := make([]NodeGroupTaint, len(v))
	for i, taint := range v {
		var T NodeGroupTaint
		err := mapstructure.Decode(taint.(map[string]interface{}), &T)
		if err != nil {
			return nil, err
		}
		taints[i] = T
	}
	return taints, nil
}

func flattenNodeGroupLabelsList(v []NodeGroupLabel) []map[string]interface{} {
	labels := make([]map[string]interface{}, len(v))
	for i, label := range v {
		m := map[string]interface{}{"key": label.Key, "value": label.Value}
		labels[i] = m
	}
	return labels
}

func flattenNodeGroupTaintsList(v []NodeGroupTaint) []map[string]interface{} {
	taints := make([]map[string]interface{}, len(v))
	for i, taint := range v {
		m := map[string]interface{}{"key": taint.Key, "value": taint.Value, "effect": taint.Effect}
		taints[i] = m
	}
	return taints
}

func kubernetesStateRefreshFunc(client ContainerClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := ClusterGet(client, clusterID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return c, status.DELETED, nil
			}
			return nil, "", err
		}
		errorStatus := status.ERROR
		if c.NewStatus == errorStatus {
			err = fmt.Errorf("mcs_kubernetes_cluster is in an error state: %s", c.StatusReason)
			return c, c.NewStatus, err
		}
		return c, c.NewStatus, nil
	}
}
