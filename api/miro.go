package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MiroApiHandler struct {
	Token string
}

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	State   string `json:"state"`
	Picture struct {
		Size44 string `json:"size44"`
	} `json:"picture"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TeamRequest struct {
	Token string `json:"token" validation:"required"`
}

type TeamResponse struct {
	Title   string `json:"title"`
	ID      string `json:"id"`
	Picture struct {
		Size44 string `json:"size44"`
	} `json:"picture"`
}

type TeamMembersRequest struct {
	Token  string `json:"token" validation:"required"`
	TeamID string `json:"team_id" validation:"required"`
}

type TeamMembersResponse struct {
	Data []struct {
		User                   `json:"user"`
		Role                   string `json:"role"`
		RoleID                 string `json:"id"`
		UserAccessBoardsNumber int    `json:"userAccessBoardsNumber"`
		OrganizationConnection struct {
			Role           string `json:"role"`
			License        string `json:"license"`
			AccountsNumber int    `json:"accountsNumber"`
		} `json:"organizationConnection"`
	} `json:"data"`
}

func GetTeamID(token string) ([]byte, error) {
	var teamResponse []TeamResponse

	req, err := http.NewRequest("GET", "https://miro.com/api/v1/accounts/?fields=id,title,creatorId,picture{size44}", nil)
	if err != nil {
		return nil, err
	}
	c := &http.Cookie{Name: "token", Value: token}

	req.AddCookie(c)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &teamResponse)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(teamResponse)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetTeamMembers(token, teamID string) ([]byte, error) {
	var teamMembersResponse TeamMembersResponse

	req, err := http.NewRequest("GET", fmt.Sprintf("https://miro.com/api/v1/accounts/%s/user-connections?fields=id,user{id,email,state,name,picture{size44}},lastActivityDate,role,dayPassesActivatedInLast30Days,organizationConnection{license,role,accountsNumber},userAccessBoardsNumber&roles=ADMIN,USER&sort=name&limit=100&offset=0", teamID), nil)
	if err != nil {
		return nil, err
	}
	c := &http.Cookie{Name: "token", Value: token}

	req.AddCookie(c)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &teamMembersResponse)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(teamMembersResponse)
	if err != nil {
		return nil, err
	}

	return data, nil
}
