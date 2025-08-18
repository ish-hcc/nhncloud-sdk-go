// Proof of Concepts of NHN Cloud SDK Go
// NHN Cloud SDK Go is an SDK for developing NHN Cloud connection drivers that connect NHN Cloud to CB-Spider, a sub-framework of the Cloud-Barista multi-cloud project.
//
// * Cloud-Barista: https://github.com/cloud-barista
//
// Created by ETRI, 2025.08

package internetgateways

import (
	"encoding/json"
	"time"

	"github.com/cloud-barista/nhncloud-sdk-go"
	"github.com/cloud-barista/nhncloud-sdk-go/pagination"
)

// InternetGatewayState represents possible states of an Internet Gateway
type InternetGatewayState string

const (
	// StateAvailable indicates the gateway is in normal operational state
	StateAvailable InternetGatewayState = "available"
	
	// StateUnavailable indicates the gateway is not connected to any routing table
	StateUnavailable InternetGatewayState = "unavailable"
	
	// StateMigrating indicates the gateway is being moved to another server for maintenance
	StateMigrating InternetGatewayState = "migrating"
	
	// StateError indicates the gateway is connected to a routing table but not functioning properly
	StateError InternetGatewayState = "error"
)

// MigrateStatus represents possible migration statuses during maintenance
type MigrateStatus string

const (
	// MigrateStatusNone indicates no migration in progress or migration completed
	MigrateStatusNone MigrateStatus = "none"
	
	// MigrateStatusUnbindingProgress indicates removal from old server in progress
	MigrateStatusUnbindingProgress MigrateStatus = "unbinding_progress"
	
	// MigrateStatusUnbindingError indicates error occurred during removal from old server
	MigrateStatusUnbindingError MigrateStatus = "unbinding_error"
	
	// MigrateStatusBindingProgress indicates setup on new server in progress
	MigrateStatusBindingProgress MigrateStatus = "binding_progress"
	
	// MigrateStatusBindingError indicates error occurred during setup on new server
	MigrateStatusBindingError MigrateStatus = "binding_error"
)

// InternetGateway represents an Internet Gateway
type InternetGateway struct {
	// ID is the unique identifier for the Internet Gateway
	ID string `json:"id"`

	// Name is the name of the Internet Gateway
	Name string `json:"name"`

	// ExternalNetworkID is the ID of the external network connected to this gateway
	ExternalNetworkID string `json:"external_network_id"`

	// RoutingTableID is the ID of the routing table connected to this gateway (may be null)
	RoutingTableID *string `json:"routingtable_id"`

	// State represents the current state of the Internet Gateway
	// Possible values: available, unavailable, migrating, error
	State string `json:"state"`

	// CreateTime is the creation time of the Internet Gateway in UTC
	CreateTime time.Time `json:"-"`

	// TenantID is the tenant ID that owns this Internet Gateway
	TenantID string `json:"tenant_id"`

	// MigrateStatus represents the migration status during maintenance
	// Possible values: none, unbinding_progress, unbinding_error, binding_progress, binding_error
	MigrateStatus string `json:"migrate_status"`

	// MigrateError contains error message if migration fails
	MigrateError *string `json:"migrate_error"`
}

// UnmarshalJSON implements custom JSON unmarshaling for InternetGateway
func (r *InternetGateway) UnmarshalJSON(b []byte) error {
	type tmp InternetGateway
	var s struct {
		tmp
		CreateTime string `json:"create_time"`
	}
	
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	
	*r = InternetGateway(s.tmp)
	
	if s.CreateTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", s.CreateTime)
		if err != nil {
			return err
		}
		r.CreateTime = t
	}
	
	return nil
}

// InternetGatewayPage represents a single page of Internet Gateway results
type InternetGatewayPage struct {
	pagination.LinkedPageBase
}

// NextPageURL extracts the next page URL from the response
func (r InternetGatewayPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"internetgateways_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty determines if an InternetGatewayPage contains any results
func (r InternetGatewayPage) IsEmpty() (bool, error) {
	internetgateways, err := ExtractInternetGateways(r)
	return len(internetgateways) == 0, err
}

// ExtractInternetGateways extracts Internet Gateways from a List result
func ExtractInternetGateways(r pagination.Page) ([]InternetGateway, error) {
	var s struct {
		InternetGateways []InternetGateway `json:"internetgateways"`
	}
	err := (r.(InternetGatewayPage)).ExtractInto(&s)
	return s.InternetGateways, err
}

// GetResult represents the result of a get operation
type GetResult struct {
	gophercloud.Result
}

// Extract extracts an InternetGateway from a GetResult
func (r GetResult) Extract() (*InternetGateway, error) {
	var s struct {
		InternetGateway *InternetGateway `json:"internetgateway"`
	}
	err := r.ExtractInto(&s)
	return s.InternetGateway, err
}

// CreateResult represents the result of a create operation
type CreateResult struct {
	gophercloud.Result
}

// Extract extracts an InternetGateway from a CreateResult
func (r CreateResult) Extract() (*InternetGateway, error) {
	var s struct {
		InternetGateway *InternetGateway `json:"internetgateway"`
	}
	err := r.ExtractInto(&s)
	return s.InternetGateway, err
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gophercloud.ErrResult
}
