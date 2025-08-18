// Proof of Concepts of NHN Cloud SDK Go
// NHN Cloud SDK Go is an SDK for developing NHN Cloud connection drivers that connect NHN Cloud to CB-Spider, a sub-framework of the Cloud-Barista multi-cloud project.
//
// * Cloud-Barista: https://github.com/cloud-barista
//
// Created by ETRI, 2025.08

package internetgateways

import (
	"github.com/cloud-barista/nhncloud-sdk-go"
	"github.com/cloud-barista/nhncloud-sdk-go/pagination"
)

// ListOpts allows filtering and sorting of Internet Gateway collections
type ListOpts struct {
	// TenantID filters by the tenant ID
	TenantID string `q:"tenant_id"`
	
	// ID filters by the Internet Gateway ID
	ID string `q:"id"`
	
	// Name filters by the Internet Gateway name
	Name string `q:"name"`
	
	// ExternalNetworkID filters by the external network ID
	ExternalNetworkID string `q:"external_network_id"`
	
	// RoutingTableID filters by the routing table ID
	RoutingTableID string `q:"routingtable_id"`
}

// List returns a Pager which allows you to iterate over Internet Gateways
func List(client *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	
	return pagination.NewPager(client, listURL(client)+q.String(), func(r pagination.PageResult) pagination.Page {
		return InternetGatewayPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get returns details about a specific Internet Gateway
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOpts represents options for creating an Internet Gateway
type CreateOpts struct {
	// Name is the name of the Internet Gateway
	Name string `json:"name" required:"true"`
	
	// ExternalNetworkID is the ID of the external network to connect to
	ExternalNetworkID string `json:"external_network_id" required:"true"`
}

// ToInternetGatewayCreateMap builds a request body from CreateOpts
func (opts CreateOpts) ToInternetGatewayCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "internetgateway")
}

// CreateOptsBuilder allows extensions to add additional attributes to the Create request
type CreateOptsBuilder interface {
	ToInternetGatewayCreateMap() (map[string]interface{}, error)
}

// Create creates a new Internet Gateway
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToInternetGatewayCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	
	resp, err := client.Post(createURL(client), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes an Internet Gateway
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{200, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
