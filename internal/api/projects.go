package api

import (
	"net/url"
	"strconv"
)

func (c *Client) ListAssignedProjects(params ListParams) ([]*ProjectAssignment, error) {
	queryParams := url.Values{}
	if params.Limit > 0 {
		queryParams.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		queryParams.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Filter != "" {
		queryParams.Set("filter", params.Filter)
	}

	endpoint := "/users/me/project_assignments"
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	var response ListAssignedProjectsResponse
	err := c.makeRequest("GET", endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.ProjectAssignments, nil
}
