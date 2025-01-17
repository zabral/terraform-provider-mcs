package mcs

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func nodeGroupFixture(name, flavorID string, count, max, min int, autoscaling bool) *NodeGroupCreateOpts {
	return &NodeGroupCreateOpts{
		Name:        name,
		FlavorID:    flavorID,
		NodeCount:   count,
		MaxNodes:    max,
		MinNodes:    min,
		Autoscaling: autoscaling,
	}
}

const nodeGroupResourceFixture = `
		%s

		resource "mcs_kubernetes_node_group" "%[2]s" {
          cluster_id          = mcs_kubernetes_cluster.%s.id
		  name                = "%[2]s"
		  flavor_id           = "%[4]s"
		  node_count          =  "%d"
		  max_nodes           =  "%d"
		  min_nodes           =  "%d"
		  autoscaling_enabled =  "%t"
		}`

func TestAccKubernetesNodeGroup_basic(t *testing.T) {
	var cluster Cluster
	var nodeGroup NodeGroup

	clusterName := "testcluster" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	createClusterFixture := clusterFixture(clusterName, ClusterTemplateID, OSFlavorID,
		OSKeypairName, OSNetworkID, OSSubnetworkID, 1)
	clusterResourceName := "mcs_kubernetes_cluster." + clusterName

	nodeGroupName := "testng" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	ngFixture := nodeGroupFixture(nodeGroupName, OSFlavorID, 1, 5, 1, false)
	nodeGroupResourceName := "mcs_kubernetes_node_group." + nodeGroupName

	ngNodeCountScaleFixture := nodeGroupFixture(nodeGroupName, OSFlavorID, 2, 5, 1, false)
	ngPatchOptsFixture := nodeGroupFixture(nodeGroupName, OSFlavorID, 2, 4, 2, true)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKubernetes(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesNodeGroupBasic(clusterName, testAccKubernetesClusterBasic(createClusterFixture), ngFixture),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists(clusterResourceName, &cluster),
					testAccCheckKubernetesNodeGroupExists(nodeGroupResourceName, clusterResourceName, &nodeGroup),
					checkNodeGroupAttrs(nodeGroupResourceName, ngFixture),
				),
			},
			{
				Config: testAccKubernetesNodeGroupBasic(clusterName, testAccKubernetesClusterBasic(createClusterFixture), ngNodeCountScaleFixture),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(nodeGroupResourceName, "node_count", strconv.Itoa(ngNodeCountScaleFixture.NodeCount)),
					testAccCheckKubernetesNodeGroupScaled(nodeGroupResourceName),
				),
			},
			{
				Config: testAccKubernetesNodeGroupBasic(clusterName, testAccKubernetesClusterBasic(createClusterFixture), ngPatchOptsFixture),
				Check: resource.ComposeTestCheckFunc(
					checkNodeGroupPatchAttrs(nodeGroupResourceName, ngPatchOptsFixture),
					testAccCheckKubernetesNodeGroupPatched(nodeGroupResourceName),
				),
			},
		},
	})
}

func checkNodeGroupAttrs(resourceName string, nodeGroup *NodeGroupCreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("state not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "name", nodeGroup.Name),
			resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(nodeGroup.NodeCount)),
			resource.TestCheckResourceAttr(resourceName, "flavor_id", nodeGroup.FlavorID),
			resource.TestCheckResourceAttr(resourceName, "max_nodes", strconv.Itoa(nodeGroup.MaxNodes)),
			resource.TestCheckResourceAttr(resourceName, "min_nodes", strconv.Itoa(nodeGroup.MinNodes)),
			resource.TestCheckResourceAttr(resourceName, "autoscaling_enabled", strconv.FormatBool(nodeGroup.Autoscaling)),
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func checkNodeGroupPatchAttrs(resourceName string, nodeGroup *NodeGroupCreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("state not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "max_nodes", strconv.Itoa(nodeGroup.MaxNodes)),
			resource.TestCheckResourceAttr(resourceName, "min_nodes", strconv.Itoa(nodeGroup.MinNodes)),
			resource.TestCheckResourceAttr(resourceName, "autoscaling_enabled", strconv.FormatBool(nodeGroup.Autoscaling)),
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func testAccCheckKubernetesNodeGroupExists(n, clusterResourceName string, nodeGroup *NodeGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, found, err := getNgAndResource(n, s)
		if err != nil {
			return err
		}
		cluster, _, err := getClusterAndResource(clusterResourceName, s)
		if err != nil {
			return err
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("node group not found")
		}

		if cluster.Primary.ID != rs.Primary.Attributes["cluster_id"] {
			return fmt.Errorf(
				"mismatched cluster id in node_group; expected %s, but got %s",
				cluster.Primary.ID, rs.Primary.Attributes["cluster_id"])
		}

		*nodeGroup = *found

		return nil
	}
}

func testAccCheckKubernetesNodeGroupScaled(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, found, err := getNgAndResource(n, s)
		if err != nil {
			return err
		}

		if strconv.Itoa(found.NodeCount) != rs.Primary.Attributes["node_count"] {
			return fmt.Errorf("mismatched node_count")
		}
		return nil
	}
}

func testAccCheckKubernetesNodeGroupPatched(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, found, err := getNgAndResource(n, s)
		if err != nil {
			return err
		}

		if strconv.Itoa(found.MaxNodes) != rs.Primary.Attributes["max_nodes"] {
			return fmt.Errorf("mismatched max_nodes")
		}
		if strconv.Itoa(found.MinNodes) != rs.Primary.Attributes["min_nodes"] {
			return fmt.Errorf("mismatched min_nodes")
		}
		if strconv.FormatBool(found.Autoscaling) != rs.Primary.Attributes["autoscaling_enabled"] {
			return fmt.Errorf("mismatched autoscaling")
		}
		return nil
	}
}

func getNgAndResource(n string, s *terraform.State) (*terraform.ResourceState, *NodeGroup, error) {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return nil, nil, fmt.Errorf("node group not found: %s", n)
	}

	config := testAccProvider.Meta().(*ConfigImpl)
	containerInfraClient, err := config.ContainerInfraV1Client(OSRegionName)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating container infra client: %s", err)
	}

	found, err := NodeGroupGet(containerInfraClient, rs.Primary.ID).Extract()
	if err != nil {
		return nil, nil, err
	}
	return rs, found, nil
}

func testAccKubernetesNodeGroupBasic(clusterName, clusterResource string, fixture *NodeGroupCreateOpts) string {

	return fmt.Sprintf(
		nodeGroupResourceFixture,
		clusterResource,
		fixture.Name,
		clusterName,
		fixture.FlavorID,
		fixture.NodeCount,
		fixture.MaxNodes,
		fixture.MinNodes,
		fixture.Autoscaling,
	)
}
