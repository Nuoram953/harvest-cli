package api

import (
	"net/url"
	"strconv"
)

func (c *Client) ListTasks(projectId int64, params ListParams) ([]*TaskAssignment, error) {
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

	for _, projectAssignment := range response.ProjectAssignments {
		if projectAssignment.Project.ID == projectId {
			return projectAssignment.TaskAssignments, nil
		}
	}

	return nil, nil
}
