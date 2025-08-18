// Proof of Concepts of NHN Cloud SDK Go
// NHN Cloud SDK Go is an SDK for developing NHN Cloud connection drivers that connect NHN Cloud to CB-Spider, a sub-framework of the Cloud-Barista multi-cloud project.
//
// * Cloud-Barista: https://github.com/cloud-barista
//
// Modified by ETRI, 2025.08

package routingtables

import "github.com/cloud-barista/nhncloud-sdk-go"

const resourcePath = "routingtables"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func attachGatewayURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "attach_gateway")
}

func detachGatewayURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "detach_gateway")
}

func setAsDefaultURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "set_as_default")
}

func relatedGatewaysURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "related_gateways")
}
