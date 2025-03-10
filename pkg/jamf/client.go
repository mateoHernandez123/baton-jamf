package jamf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

type Client struct {
	httpClient  *http.Client
	token       string
	baseUrl     string
	instanceUrl string
}

func NewClient(httpClient *http.Client, token, baseUrl, instanceUrl string) *Client {
	return &Client{
		httpClient:  httpClient,
		token:       token,
		baseUrl:     baseUrl,
		instanceUrl: instanceUrl,
	}
}

type AuthResponse struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

func CreateBearerToken(ctx context.Context, username, password, serverInstance string) (string, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/api/v1/auth/token", serverInstance)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Token, nil
}

// Función para obtener la base de usuarios y grupos
func (c *Client) getBaseAccounts(ctx context.Context) (*BaseAccount, error) {
	url := fmt.Sprintf("%s/api/v1/accounts", c.baseUrl)
	var baseAccounts BaseAccount
	err := c.doRequest(ctx, url, &baseAccounts)
	if err != nil {
		return nil, err
	}
	return &baseAccounts, nil
}

// Obtener detalles de un usuario
func (c *Client) GetUserAccountDetails(ctx context.Context, userID int) (*UserAccount, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", c.baseUrl, userID)
	var userAccount UserAccount
	err := c.doRequest(ctx, url, &userAccount)
	if err != nil {
		return nil, err
	}
	return &userAccount, nil
}

// Obtener detalles de un grupo
func (c *Client) GetGroupDetails(ctx context.Context, groupID int) (*Group, error) {
	url := fmt.Sprintf("%s/api/v1/groups/%d", c.baseUrl, groupID)
	var group Group
	err := c.doRequest(ctx, url, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// Obtener todas las cuentas de usuario y grupos
func (c *Client) GetAccounts(ctx context.Context) ([]UserAccount, []Group, error) {
	var userAccounts []UserAccount
	var groups []Group

	baseAccounts, err := c.getBaseAccounts(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Obtener información detallada de cada usuario
	for _, user := range baseAccounts.Users {
		userAccountInfo, err := c.GetUserAccountDetails(ctx, user.ID)
		if err != nil {
			return nil, nil, err
		}
		userAccounts = append(userAccounts, *userAccountInfo)
	}

	// Obtener información detallada de cada grupo
	for _, group := range baseAccounts.Groups {
		groupInfo, err := c.GetGroupDetails(ctx, group.ID)
		if err != nil {
			return nil, nil, err
		}
		groups = append(groups, *groupInfo)
	}

	return userAccounts, groups, nil
}

// Realizar peticiones GET genéricas
func (c *Client) doRequest(ctx context.Context, url string, res interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return err
	}
	return nil
}
