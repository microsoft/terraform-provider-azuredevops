package organization

type OrganizationMeta struct {
	Id                 *string     `json:"id,omitempty"`
	Name               *string     `json:"name,omitempty"`
	Status             *string     `json:"status,omitempty"`
	Owner              *string     `json:"owner,omitempty"`
	TenantId           *string     `json:"tenantId,omitempty"`
	DateCreated        *string     `json:"dateCreated,omitempty"`
	LastUpdated        *string     `json:"lastUpdated,omitempty"`
	PreferredRegion    *string     `json:"preferredRegion,omitempty"`
	PreferredGeography *string     `json:"preferredGeography,omitempty"`
	Properties         interface{} `json:"properties,omitempty"`
	Data               interface{} `json:"data,omitempty"`
}
