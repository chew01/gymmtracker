package data

type FormattedResponse struct {
	Success bool     `json:"success"`
	OrgCode string   `json:"org_code"`
	Outlets []Outlet `json:"outlets"`
}

type Outlet struct {
	Name      string `json:"name"`
	Occupancy int    `json:"occupancy"`
	Limit     int    `json:"occupancy_limit"`
	Queue     int    `json:"queue_length"`
}
