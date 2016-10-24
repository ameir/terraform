package loadbalancers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToLoadBalancerListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the Loadbalancer attributes you want to see returned. SortKey allows you to
// sort by a particular attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Description        string `q:"description"`
	AdminStateUp       *bool  `q:"admin_state_up"`
	TenantID           string `q:"tenant_id"`
	ProvisioningStatus string `q:"provisioning_status"`
	VipAddress         string `q:"vip_address"`
	VipPortID          string `q:"vip_port_id"`
	VipSubnetID        string `q:"vip_subnet_id"`
	ID                 string `q:"id"`
	OperatingStatus    string `q:"operating_status"`
	Name               string `q:"name"`
	Flavor             string `q:"flavor"`
	Provider           string `q:"provider"`
	Limit              int    `q:"limit"`
	Marker             string `q:"marker"`
	SortKey            string `q:"sort_key"`
	SortDir            string `q:"sort_dir"`
}

// ToLoadbalancerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLoadBalancerListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// routers. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those routers that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToLoadBalancerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return LoadBalancerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToLoadBalancerCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Optional. Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`
	// Optional. Human-readable description for the Loadbalancer.
	Description string `json:"description,omitempty"`
	// Required. The network on which to allocate the Loadbalancer's address. A tenant can
	// only create Loadbalancers on networks authorized by policy (e.g. networks that
	// belong to them or networks that are shared).
	VipSubnetID string `json:"vip_subnet_id" required:"true"`
	// Required for admins. The UUID of the tenant who owns the Loadbalancer.
	// Only administrative users can specify a tenant UUID other than their own.
	TenantID string `json:"tenant_id,omitempty"`
	// Optional. The IP address of the Loadbalancer.
	VipAddress string `json:"vip_address,omitempty"`
	// Optional. The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
	// Optional. The UUID of a flavor.
	Flavor string `json:"flavor,omitempty"`
	// Optional. The name of the provider.
	Provider string `json:"provider,omitempty"`
}

// ToLoadBalancerCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToLoadBalancerCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "loadbalancer")
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLoadBalancerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular Loadbalancer based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToLoadBalancerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Optional. Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`
	// Optional. Human-readable description for the Loadbalancer.
	Description string `json:"description,omitempty"`
	// Optional. The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

// ToLoadBalancerUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToLoadBalancerUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "loadbalancer")
}

// Update is an operation which modifies the attributes of the specified LoadBalancer.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToLoadBalancerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete will permanently delete a particular LoadBalancer based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}

func GetStatuses(c *gophercloud.ServiceClient, id string) (r GetStatusesResult) {
	_, r.Err = c.Get(statusRootURL(c, id), &r.Body, nil)
	return
}