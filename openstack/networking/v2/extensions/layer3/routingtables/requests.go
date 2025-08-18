// Proof of Concepts of NHN Cloud SDK Go
// NHN Cloud SDK Go is an SDK for developing NHN Cloud connection drivers that connect NHN Cloud to CB-Spider, a sub-framework of the Cloud-Barista multi-cloud project.
//
// * Cloud-Barista: https://github.com/cloud-barista
//
// Created by ETRI, 2025.08

package routingtables

import (
	"github.com/cloud-barista/nhncloud-sdk-go"
	"github.com/cloud-barista/nhncloud-sdk-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToRoutingTableListQuery() (string, error)
}

// ListOpts represents options for listing routing tables.
type ListOpts struct {
	// TenantID filters routing tables by tenant ID
	TenantID string `q:"tenant_id"`
	
	// ID filters routing tables by ID
	ID string `q:"id"`
	
	// Name filters routing tables by name
	Name string `q:"name"`
	
	// DefaultTable filters routing tables by default table status
	DefaultTable *bool `q:"default_table"`
	
	// GatewayID filters routing tables by connected internet gateway ID
	GatewayID string `q:"gateway_id"`
	
	// Distributed filters routing tables by routing type (true: distributed, false: centralized)
	Distributed *bool `q:"distributed"`
	
	// Detail includes detailed information in the response
	Detail *bool `q:"detail"`
	
	// SortDir specifies the sort direction (asc, desc)
	SortDir string `q:"sort_dir"`
	
	// SortKey specifies the field to sort by
	SortKey string `q:"sort_key"`
}

// ToRoutingTableListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRoutingTableListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of routing tables.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToRoutingTableListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RoutingTablePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific routing table based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToRoutingTableCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a routing table.
type CreateOpts struct {
	// Name is the name of the routing table
	Name string `json:"name" required:"true"`
	
	// VPCID is the ID of the VPC this routing table belongs to
	VPCID string `json:"vpc_id" required:"true"`
	
	// Distributed specifies the routing type (true: distributed, false: centralized)
	// Defaults to true if not specified
	Distributed *bool `json:"distributed,omitempty"`
}

// ToRoutingTableCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToRoutingTableCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "routingtable")
}

// Create accepts a CreateOpts struct and creates a new routing table using the values provided.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRoutingTableCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(createURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToRoutingTableUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a routing table.
type UpdateOpts struct {
	// Name is the new name for the routing table
	Name string `json:"name,omitempty"`
	
	// Distributed specifies the routing type (true: distributed, false: centralized)
	Distributed *bool `json:"distributed,omitempty"`
}

// ToRoutingTableUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToRoutingTableUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "routingtable")
}

// Update accepts a UpdateOpts struct and updates an existing routing table using the values provided.
func Update(c *gophercloud.ServiceClient, routingtableID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRoutingTableUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(resourceURL(c, routingtableID), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the routing table associated with it.
func Delete(c *gophercloud.ServiceClient, routingtableID string) (r DeleteResult) {
	resp, err := c.Delete(resourceURL(c, routingtableID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AttachGatewayOptsBuilder allows extensions to add additional parameters to the AttachGateway request.
type AttachGatewayOptsBuilder interface {
	ToAttachGatewayMap() (map[string]interface{}, error)
}

// AttachGatewayOpts represents options used to attach an internet gateway to a routing table.
type AttachGatewayOpts struct {
	// GatewayID is the ID of the internet gateway to attach
	GatewayID string `json:"gateway_id" required:"true"`
}

// ToAttachGatewayMap builds a request body from AttachGatewayOpts.
func (opts AttachGatewayOpts) ToAttachGatewayMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AttachGateway attaches an internet gateway to a routing table.
func AttachGateway(c *gophercloud.ServiceClient, routingtableID string, opts AttachGatewayOptsBuilder) (r AttachGatewayResult) {
	b, err := opts.ToAttachGatewayMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(attachGatewayURL(c, routingtableID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DetachGateway detaches an internet gateway from a routing table.
func DetachGateway(c *gophercloud.ServiceClient, routingtableID string) (r DetachGatewayResult) {
	resp, err := c.Put(detachGatewayURL(c, routingtableID), nil, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// SetAsDefault sets a routing table as the default routing table for its VPC.
func SetAsDefault(c *gophercloud.ServiceClient, routingtableID string) (r SetAsDefaultResult) {
	resp, err := c.Put(setAsDefaultURL(c, routingtableID), nil, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetRelatedGateways retrieves gateways that can be reached through the routing policies set in the routing table.
func GetRelatedGateways(c *gophercloud.ServiceClient, routingtableID string) (r GetRelatedGatewaysResult) {
	resp, err := c.Get(relatedGatewaysURL(c, routingtableID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Route management functions

// RouteListOptsBuilder allows extensions to add additional parameters to the route List request.
type RouteListOptsBuilder interface {
	ToRouteListQuery() (string, error)
}

// RouteListOpts represents options for listing routes.
type RouteListOpts struct {
	// ID filters routes by ID
	ID string `q:"id"`
	
	// CIDR filters routes by destination CIDR
	CIDR string `q:"cidr"`
	
	// Mask filters routes by destination CIDR netmask (0-32)
	Mask *int `q:"mask"`
	
	// Gateway filters routes by gateway IP
	Gateway string `q:"gateway"`
	
	// RoutingTableID filters routes by routing table ID
	RoutingTableID string `q:"routingtable_id"`
	
	// GatewayID filters routes by internet gateway ID
	GatewayID string `q:"gateway_id"`
}

// ToRouteListQuery formats a RouteListOpts into a query string.
func (opts RouteListOpts) ToRouteListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListRoutes returns a Pager which allows you to iterate over a collection of routes.
func ListRoutes(c *gophercloud.ServiceClient, opts RouteListOptsBuilder) pagination.Pager {
	url := routesURL(c)
	if opts != nil {
		query, err := opts.ToRouteListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RoutePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetRoute retrieves a specific route based on its unique ID.
func GetRoute(c *gophercloud.ServiceClient, routeID string) (r GetRouteResult) {
	resp, err := c.Get(routeURL(c, routeID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateRouteOptsBuilder allows extensions to add additional parameters to the CreateRoute request.
type CreateRouteOptsBuilder interface {
	ToRouteCreateMap() (map[string]interface{}, error)
}

// CreateRouteOpts represents options used to create a route.
type CreateRouteOpts struct {
	// RoutingTableID is the ID of the routing table to add the route to
	RoutingTableID string `json:"routingtable_id" required:"true"`
	
	// CIDR is the destination CIDR for the route
	CIDR string `json:"cidr" required:"true"`
	
	// Gateway is the gateway IP for the route
	Gateway string `json:"gateway" required:"true"`
	
	// Description is the description of the route (max 256 bytes)
	Description string `json:"description" required:"true"`
}

// ToRouteCreateMap builds a request body from CreateRouteOpts.
func (opts CreateRouteOpts) ToRouteCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "route")
}

// CreateRoute accepts a CreateRouteOpts struct and creates a new route using the values provided.
func CreateRoute(c *gophercloud.ServiceClient, opts CreateRouteOptsBuilder) (r CreateRouteResult) {
	b, err := opts.ToRouteCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(routesURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateRouteOptsBuilder allows extensions to add additional parameters to the UpdateRoute request.
type UpdateRouteOptsBuilder interface {
	ToRouteUpdateMap() (map[string]interface{}, error)
}

// UpdateRouteOpts represents options used to update a route.
type UpdateRouteOpts struct {
	// CIDR is the destination CIDR for the route
	CIDR string `json:"cidr,omitempty"`
	
	// Gateway is the gateway IP for the route
	Gateway string `json:"gateway,omitempty"`
	
	// Description is the description of the route (max 256 bytes)
	Description string `json:"description,omitempty"`
}

// ToRouteUpdateMap builds a request body from UpdateRouteOpts.
func (opts UpdateRouteOpts) ToRouteUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "route")
}

// UpdateRoute accepts an UpdateRouteOpts struct and updates an existing route using the values provided.
func UpdateRoute(c *gophercloud.ServiceClient, routeID string, opts UpdateRouteOptsBuilder) (r UpdateRouteResult) {
	b, err := opts.ToRouteUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(routeURL(c, routeID), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteRoute accepts a unique ID and deletes the route associated with it.
func DeleteRoute(c *gophercloud.ServiceClient, routeID string) (r DeleteRouteResult) {
	resp, err := c.Delete(routeURL(c, routeID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// URLs for route operations
func routesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("routes")
}

func routeURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("routes", id)
}
