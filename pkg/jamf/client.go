package jamf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

func (c *Client) GetTokenDetails(ctx context.Context) (TokenDetails, error) {
	url := fmt.Sprintf("%s/api/v1/auth", c.instanceUrl)

	var res TokenDetails
	if err := c.doRequest(ctx, url, &res); err != nil {
		return TokenDetails{}, err
	}

	return res, nil
}

func (c *Client) getBaseUsers(ctx context.Context) ([]BaseType, error) {
	usersUrl, err := url.JoinPath(c.baseUrl, "/users")
	if err != nil {
		return nil, err
	}

	var res struct {
		Users []BaseType `json:"users"`
	}

	if err := c.doRequest(ctx, usersUrl, &res); err != nil {
		return nil, err
	}
	return res.Users, nil
}

func (c *Client) getUserDetails(ctx context.Context, userId int) (User, error) {
	userIdString := strconv.Itoa(userId)
	usersUrl, err := url.JoinPath(c.baseUrl, "/users/id/", userIdString)

	if err != nil {
		return User{}, err
	}

	var res struct {
		User User `json:"user"`
	}

	if err := c.doRequest(ctx, usersUrl, &res); err != nil {
		return User{}, err
	}

	return res.User, nil
}

func (c *Client) getBaseUserGroups(ctx context.Context) ([]UserGroup, error) {
	accountUrl, err := url.JoinPath(c.baseUrl, "/usergroups")
	if err != nil {
		return nil, err
	}

	var res struct {
		UserGroup []UserGroup `json:"user_groups"`
	}

	if err := c.doRequest(ctx, accountUrl, &res); err != nil {
		return nil, err
	}
	return res.UserGroup, nil
}


func (c *Client) getUserGroupDetails(ctx context.Context, userGroupId int) (UserGroup, error) {
	groupIdString := strconv.Itoa(userGroupId)
	usersUrl, err := url.JoinPath(c.baseUrl, "/usergroups/id/", groupIdString)

	if err =! nil {
		return UserGroup{}, err
	}

	var res struct {
		UserGroup UserGroup `json:"user_groups"`
	}

	if err := c.doRequest(ctx, usersUrl, &res); err != nil {
		return UserGroup{}, err
	}

	return res.UserGroup, nil
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	baseUsers, err := c.getBaseUsers(ctx)

	if err != nil {
		return nil, err
	}
	
	for _, baseUser := range baseUsers {
		user, err := c.getUserDetails(ctx, baseUser.ID)
		
		if err != nil {
			return nil, err
		}
		
		users = append(users, user)
	}

	return users, nil
}

func (c *Client) GetUserGroups(ctx context.Context) ([]UserGroup, error) {
	var usarGroups []UserGroup
	baseUserGroup, err := c.getBaseUserGroups(ctx)

	if err != nil {
		return nil, err
	}

	for _,userGroup := range baseUserGroup {
		if err != nil {
			return nil, err
		}
		userGroups = append(userGroups, userGroupInfo)
	}

	return userGroups, nil
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
func (c *Client) GetUserAccountDetails(ctx context.Context, userId int) (UserAccount, error) {
	userIdString := strconv.Itoa(userId)
	usersUrl, err := url.JoinPath(c.baseUrl, "/accounts/userid/", userIdString)

	if err != nil {
		return UserAccount{}, err
	}

	var res struct {
		UserAccount UserAccount `json:"account"`
	}

	if err := c.doRequest(ctx, usersUrl, &res); err != nil {
		return UserAccount{}, err
	}

	return res.UserAccount, nil

}

// Obtener detalles de un grupo
func (c *Client) GetGroupDetails(ctx context.Context, groupId int) (Group, error) {
	groupIdString := strconv.Itoa(groupId)
	usersUrl, err := url.JoinPath(c.baseUrl, "/accounts/groupid/", groupIdString)

	if err != nil {
		return Group{}, err
	}

	var res struct {
		Group Group `json:"group"`
	}

	if err := c.doRequest(ctx, usersUrl, &res); err != nil {
		return Group{}, err
	}

	return res.Group, nil
}

func (c *Client) GetSites(ctx context.Context) ([]Site, error) {
	sitesUrl, err := url.JoinPath(c.baseUrl, "/sites")
	if err != nil {
		return nil, err
	}

	var res struct {
		Sites []Site `json:"sites"`
	}

	if err := c.doRequest(ctx, sitesUrl, &res); err != nil {
		return nil, err
	}

	return res.Sites, nil
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
