// Proof of Concepts of NHN Cloud SDK Go
// NHN Cloud SDK Go is an SDK for developing NHN Cloud connection drivers that connect NHN Cloud to CB-Spider, a sub-framework of the Cloud-Barista multi-cloud project.
//
// * Cloud-Barista: https://github.com/cloud-barista
//
// Created by ETRI, 2025.08

package routingtables

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-barista/nhncloud-sdk-go"
	"github.com/cloud-barista/nhncloud-sdk-go/pagination"
)

// FlexibleSubnetInfo handles both string IDs and full subnet objects from the API
type FlexibleSubnetInfo struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle both string and object formats
func (fsi *FlexibleSubnetInfo) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string first (just the ID)
	var idStr string
	if err := json.Unmarshal(data, &idStr); err == nil {
		fsi.ID = idStr
		fsi.Name = "" // Name not available when only ID is provided
		return nil
	}
	
	// If string unmarshaling fails, try as full object
	type subnetAlias FlexibleSubnetInfo
	aux := (*subnetAlias)(fsi)
	return json.Unmarshal(data, aux)
}

// MarshalJSON implements custom JSON marshaling
func (fsi FlexibleSubnetInfo) MarshalJSON() ([]byte, error) {
	// If only ID is set, marshal as string
	if fsi.Name == "" && fsi.ID != "" {
		return json.Marshal(fsi.ID)
	}
	
	// Otherwise marshal as object
	type subnetAlias FlexibleSubnetInfo
	return json.Marshal((subnetAlias)(fsi))
}

// FlexibleVPCInfo handles both string IDs and full VPC objects from the API
type FlexibleVPCInfo struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle both string and object formats
func (fvi *FlexibleVPCInfo) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string first (just the ID)
	var idStr string
	if err := json.Unmarshal(data, &idStr); err == nil {
		fvi.ID = idStr
		fvi.Name = "" // Name not available when only ID is provided
		return nil
	}
	
	// If string unmarshaling fails, try as full object
	type vpcAlias FlexibleVPCInfo
	aux := (*vpcAlias)(fvi)
	return json.Unmarshal(data, aux)
}

// MarshalJSON implements custom JSON marshaling
func (fvi FlexibleVPCInfo) MarshalJSON() ([]byte, error) {
	// If only ID is set, marshal as string
	if fvi.Name == "" && fvi.ID != "" {
		return json.Marshal(fvi.ID)
	}
	
	// Otherwise marshal as object
	type vpcAlias FlexibleVPCInfo
	return json.Marshal((vpcAlias)(fvi))
}

// NHNCloudTime handles the custom timestamp format used by NHN Cloud API
// Format: "2024-02-13 10:45:57" instead of standard RFC3339
type NHNCloudTime struct {
	time.Time
}

// UnmarshalJSON implements custom JSON unmarshaling for NHN Cloud timestamp format
func (ct *NHNCloudTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	s := strings.Trim(string(data), `"`)
	
	// Handle empty/null values
	if s == "null" || s == "" {
		return nil
	}
	
	// List of possible time formats used by NHN Cloud API
	formats := []string{
		"2006-01-02 15:04:05",           // Most common format
		"2006-01-02T15:04:05",           // Alternative format without timezone
		"2006-01-02T15:04:05Z",          // UTC format
		"2006-01-02T15:04:05Z07:00",     // Full RFC3339
		"2006-01-02 15:04:05.000000",    // With microseconds
		"2006-01-02T15:04:05.000000Z",   // RFC3339 with microseconds
	}
	
	var parseErr error
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			ct.Time = t
			return nil
		} else {
			parseErr = err
		}
	}
	
	return fmt.Errorf("unable to parse time %q with any known format: %v", s, parseErr)
}

// MarshalJSON implements custom JSON marshaling for NHN Cloud timestamp format
func (ct NHNCloudTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	// Format to match NHN Cloud API format
	return json.Marshal(ct.Time.Format("2006-01-02 15:04:05"))
}

// String returns string representation of the time
func (ct NHNCloudTime) String() string {
	return ct.Time.Format("2006-01-02 15:04:05")
}

// Helper methods for FlexibleSubnetInfo

// GetSubnetIDs returns a slice of subnet IDs from the flexible subnet info list
func (rt *RoutingTable) GetSubnetIDs() []string {
	ids := make([]string, len(rt.Subnets))
	for i, subnet := range rt.Subnets {
		ids[i] = subnet.ID
	}
	return ids
}

// GetSubnetNames returns a slice of subnet names (may be empty if API only returned IDs)
func (rt *RoutingTable) GetSubnetNames() []string {
	names := make([]string, 0)
	for _, subnet := range rt.Subnets {
		if subnet.Name != "" {
			names = append(names, subnet.Name)
		}
	}
	return names
}

// Helper methods for FlexibleVPCInfo

// GetVPCIDs returns a slice of VPC IDs from the flexible VPC info list
func (rt *RoutingTable) GetVPCIDs() []string {
	ids := make([]string, len(rt.VPCs))
	for i, vpc := range rt.VPCs {
		ids[i] = vpc.ID
	}
	return ids
}

// GetVPCNames returns a slice of VPC names (may be empty if API only returned IDs)
func (rt *RoutingTable) GetVPCNames() []string {
	names := make([]string, 0)
	for _, vpc := range rt.VPCs {
		if vpc.Name != "" {
			names = append(names, vpc.Name)
		}
	}
	return names
}

// RoutingTable represents a routing table resource.
type RoutingTable struct {
	// ID is the unique identifier of the routing table
	ID string `json:"id"`
	
	// Name is the name of the routing table
	Name string `json:"name"`
	
	// DefaultTable indicates if this is the default routing table
	DefaultTable bool `json:"default_table"`
	
	// Distributed indicates the routing type (true: distributed, false: centralized)
	Distributed bool `json:"distributed"`
	
	// GatewayID is the ID of the connected internet gateway (if any)
	GatewayID string `json:"gateway_id"`
	
	// GatewayName is the name of the connected internet gateway (if any)
	GatewayName string `json:"gateway_name"`
	
	// TenantID is the ID of the tenant that owns the routing table
	TenantID string `json:"tenant_id"`
	
	// State is the current state of the routing table
	State string `json:"state"`
	
	// CreateTime is when the routing table was created
	CreateTime NHNCloudTime `json:"create_time"`
	
	// VPCs is a list of VPCs this routing table belongs to (detailed view only)
	VPCs []FlexibleVPCInfo `json:"vpcs,omitempty"`
	
	// Subnets is a list of subnets connected to this routing table (detailed view only)
	Subnets []FlexibleSubnetInfo `json:"subnets,omitempty"`
	
	// Routes is a list of routes in this routing table (get operation only)
	Routes []Route `json:"routes,omitempty"`
}

// VPCInfo represents VPC information within a routing table (legacy - kept for compatibility).
type VPCInfo struct {
	// ID is the VPC ID
	ID string `json:"id"`
	
	// Name is the VPC name
	Name string `json:"name"`
}

// SubnetInfo represents subnet information within a routing table (legacy - kept for compatibility).
type SubnetInfo struct {
	// ID is the subnet ID
	ID string `json:"id"`
	
	// Name is the subnet name
	Name string `json:"name"`
}

// Route represents a route in a routing table.
type Route struct {
	// ID is the unique identifier of the route
	ID string `json:"id"`
	
	// CIDR is the destination CIDR
	CIDR string `json:"cidr"`
	
	// Mask is the netmask of the destination CIDR
	Mask int `json:"mask"`
	
	// Gateway is the gateway IP address
	Gateway string `json:"gateway"`
	
	// GatewayID is the ID of the internet gateway (for internet gateway routes)
	GatewayID string `json:"gateway_id,omitempty"`
	
	// Description is the route description
	Description *string `json:"description"`
	
	// RoutingTableID is the ID of the routing table this route belongs to
	RoutingTableID string `json:"routingtable_id"`
	
	// TenantID is the ID of the tenant that owns the route
	TenantID string `json:"tenant_id"`
	
	// Hidden indicates if the route is hidden (internal use)
	Hidden bool `json:"hidden,omitempty"`
}

// Gateway represents a gateway that can be reached through routing policies.
type Gateway struct {
	// ID is the gateway ID
	ID string `json:"id"`
	
	// Type is the gateway type
	Type string `json:"type"`
	
	// Name is the gateway name
	Name string `json:"name"`
}

// RoutingTablePage is the page returned by a pager when traversing over a collection of routing tables.
type RoutingTablePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routing tables has reached the end of a page
// and the pager seeks to traverse over a new one.
func (r RoutingTablePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"routingtables_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RoutingTablePage struct is empty.
func (r RoutingTablePage) IsEmpty() (bool, error) {
	is, err := ExtractRoutingTables(r)
	return len(is) == 0, err
}

// ExtractRoutingTables accepts a Page struct, specifically a RoutingTablePage struct,
// and extracts the elements into a slice of RoutingTable structs.
func ExtractRoutingTables(r pagination.Page) ([]RoutingTable, error) {
	var s struct {
		RoutingTables []RoutingTable `json:"routingtables"`
	}
	err := (r.(RoutingTablePage)).ExtractInto(&s)
	return s.RoutingTables, err
}

// RoutePage is the page returned by a pager when traversing over a collection of routes.
type RoutePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routes has reached the end of a page
// and the pager seeks to traverse over a new one.
func (r RoutePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"routes_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RoutePage struct is empty.
func (r RoutePage) IsEmpty() (bool, error) {
	is, err := ExtractRoutes(r)
	return len(is) == 0, err
}

// ExtractRoutes accepts a Page struct, specifically a RoutePage struct,
// and extracts the elements into a slice of Route structs.
func ExtractRoutes(r pagination.Page) ([]Route, error) {
	var s struct {
		Routes []Route `json:"routes"`
	}
	err := (r.(RoutePage)).ExtractInto(&s)
	return s.Routes, err
}

// RoutingTableResult represents the result of routing table operations.
type RoutingTableResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a routing table resource.
func (r RoutingTableResult) Extract() (*RoutingTable, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	
	// First try the standard extraction
	var s struct {
		RoutingTable *RoutingTable `json:"routingtable"`
	}
	
	err := r.ExtractInto(&s)
	if err != nil {
		// If standard extraction fails, try alternative parsing
		return r.ExtractRoutingTableWithFallback()
	}
	
	return s.RoutingTable, nil
}

// ExtractRoutingTableWithFallback provides fallback parsing for different API response formats
func (r RoutingTableResult) ExtractRoutingTableWithFallback() (*RoutingTable, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	
	// Extract raw JSON first
	var response struct {
		RoutingTable json.RawMessage `json:"routingtable"`
	}
	
	err := r.ExtractInto(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to extract response: %w", err)
	}
	
	// Parse the routing table manually to handle edge cases
	var routingTable RoutingTable
	err = json.Unmarshal(response.RoutingTable, &routingTable)
	if err != nil {
		// If that fails, try parsing with raw map to debug
		var rawRT map[string]interface{}
		if jsonErr := json.Unmarshal(response.RoutingTable, &rawRT); jsonErr == nil {
			return r.parseRoutingTableFromMap(rawRT)
		}
		return nil, fmt.Errorf("failed to unmarshal routing table: %w", err)
	}
	
	return &routingTable, nil
}

// parseRoutingTableFromMap manually parses routing table from a map for maximum compatibility
func (r RoutingTableResult) parseRoutingTableFromMap(data map[string]interface{}) (*RoutingTable, error) {
	rt := &RoutingTable{}
	
	// Parse basic fields
	if id, ok := data["id"].(string); ok {
		rt.ID = id
	}
	if name, ok := data["name"].(string); ok {
		rt.Name = name
	}
	if defaultTable, ok := data["default_table"].(bool); ok {
		rt.DefaultTable = defaultTable
	}
	if distributed, ok := data["distributed"].(bool); ok {
		rt.Distributed = distributed
	}
	if gatewayID, ok := data["gateway_id"].(string); ok {
		rt.GatewayID = gatewayID
	}
	if gatewayName, ok := data["gateway_name"].(string); ok {
		rt.GatewayName = gatewayName
	}
	if tenantID, ok := data["tenant_id"].(string); ok {
		rt.TenantID = tenantID
	}
	if state, ok := data["state"].(string); ok {
		rt.State = state
	}
	
	// Parse create_time
	if createTimeStr, ok := data["create_time"].(string); ok {
		var ct NHNCloudTime
		if err := json.Unmarshal([]byte(`"`+createTimeStr+`"`), &ct); err == nil {
			rt.CreateTime = ct
		}
	}
	
	// Parse VPCs (handle both string array and object array)
	if vpcs, ok := data["vpcs"]; ok {
		rt.VPCs = r.parseFlexibleVPCs(vpcs)
	}
	
	// Parse Subnets (handle both string array and object array)
	if subnets, ok := data["subnets"]; ok {
		rt.Subnets = r.parseFlexibleSubnets(subnets)
	}
	
	// Parse Routes
	if routes, ok := data["routes"].([]interface{}); ok {
		rt.Routes = r.parseRoutes(routes)
	}
	
	return rt, nil
}

// parseFlexibleVPCs handles both string arrays and object arrays for VPCs
func (r RoutingTableResult) parseFlexibleVPCs(vpcs interface{}) []FlexibleVPCInfo {
	var result []FlexibleVPCInfo
	
	switch v := vpcs.(type) {
	case []interface{}:
		for _, vpc := range v {
			switch vpcData := vpc.(type) {
			case string:
				// VPC is just an ID string
				result = append(result, FlexibleVPCInfo{ID: vpcData})
			case map[string]interface{}:
				// VPC is an object
				fvi := FlexibleVPCInfo{}
				if id, ok := vpcData["id"].(string); ok {
					fvi.ID = id
				}
				if name, ok := vpcData["name"].(string); ok {
					fvi.Name = name
				}
				result = append(result, fvi)
			}
		}
	case []string:
		// Direct string array
		for _, vpcID := range v {
			result = append(result, FlexibleVPCInfo{ID: vpcID})
		}
	}
	
	return result
}

// parseFlexibleSubnets handles both string arrays and object arrays for Subnets
func (r RoutingTableResult) parseFlexibleSubnets(subnets interface{}) []FlexibleSubnetInfo {
	var result []FlexibleSubnetInfo
	
	switch v := subnets.(type) {
	case []interface{}:
		for _, subnet := range v {
			switch subnetData := subnet.(type) {
			case string:
				// Subnet is just an ID string
				result = append(result, FlexibleSubnetInfo{ID: subnetData})
			case map[string]interface{}:
				// Subnet is an object
				fsi := FlexibleSubnetInfo{}
				if id, ok := subnetData["id"].(string); ok {
					fsi.ID = id
				}
				if name, ok := subnetData["name"].(string); ok {
					fsi.Name = name
				}
				result = append(result, fsi)
			}
		}
	case []string:
		// Direct string array
		for _, subnetID := range v {
			result = append(result, FlexibleSubnetInfo{ID: subnetID})
		}
	}
	
	return result
}

// parseRoutes handles route parsing from interface{} array
func (r RoutingTableResult) parseRoutes(routes []interface{}) []Route {
	var result []Route
	
	for _, route := range routes {
		if routeMap, ok := route.(map[string]interface{}); ok {
			r := Route{}
			
			if id, ok := routeMap["id"].(string); ok {
				r.ID = id
			}
			if cidr, ok := routeMap["cidr"].(string); ok {
				r.CIDR = cidr
			}
			if mask, ok := routeMap["mask"].(float64); ok {
				r.Mask = int(mask)
			}
			if gateway, ok := routeMap["gateway"].(string); ok {
				r.Gateway = gateway
			}
			if gatewayID, ok := routeMap["gateway_id"].(string); ok {
				r.GatewayID = gatewayID
			}
			if description, ok := routeMap["description"].(string); ok {
				r.Description = &description
			}
			if routingtableID, ok := routeMap["routingtable_id"].(string); ok {
				r.RoutingTableID = routingtableID
			}
			if tenantID, ok := routeMap["tenant_id"].(string); ok {
				r.TenantID = tenantID
			}
			if hidden, ok := routeMap["hidden"].(bool); ok {
				r.Hidden = hidden
			}
			
			result = append(result, r)
		}
	}
	
	return result
}

// ExtractRoutingTable is an alternative extraction method with more control
func (r RoutingTableResult) ExtractRoutingTable() (*RoutingTable, error) {
	return r.ExtractRoutingTableWithFallback()
}

// RouteResult represents the result of route operations.
type RouteResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a route resource.
func (r RouteResult) Extract() (*Route, error) {
	var s struct {
		Route *Route `json:"route"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	
	// Additional validation can be added here if needed
	if s.Route != nil {
		// Ensure description is properly handled (can be null in API)
		// This is already handled by the pointer type
	}
	
	return s.Route, nil
}

// ExtractRoute is an alternative extraction method with more control
func (r RouteResult) ExtractRoute() (*Route, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	
	var response struct {
		Route json.RawMessage `json:"route"`
	}
	
	err := r.ExtractInto(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to extract response: %w", err)
	}
	
	var route Route
	err = json.Unmarshal(response.Route, &route)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal route: %w", err)
	}
	
	return &route, nil
}

// GatewayResult represents the result of gateway operations.
type GatewayResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts gateway resources.
func (r GatewayResult) Extract() ([]Gateway, error) {
	var s struct {
		Gateways []Gateway `json:"gateways"`
	}
	err := r.ExtractInto(&s)
	return s.Gateways, err
}

// Result types for different operations

// GetResult represents the result of a get operation.
type GetResult struct {
	RoutingTableResult
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	RoutingTableResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	RoutingTableResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// AttachGatewayResult represents the result of an attach gateway operation.
type AttachGatewayResult struct {
	RoutingTableResult
}

// DetachGatewayResult represents the result of a detach gateway operation.
type DetachGatewayResult struct {
	RoutingTableResult
}

// SetAsDefaultResult represents the result of a set as default operation.
type SetAsDefaultResult struct {
	RoutingTableResult
}

// GetRelatedGatewaysResult represents the result of a get related gateways operation.
type GetRelatedGatewaysResult struct {
	GatewayResult
}

// Route operation result types

// GetRouteResult represents the result of a get route operation.
type GetRouteResult struct {
	RouteResult
}

// CreateRouteResult represents the result of a create route operation.
type CreateRouteResult struct {
	RouteResult
}

// UpdateRouteResult represents the result of an update route operation.
type UpdateRouteResult struct {
	RouteResult
}

// DeleteRouteResult represents the result of a delete route operation.
type DeleteRouteResult struct {
	gophercloud.ErrResult
}
