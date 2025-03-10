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

// Client representa un cliente HTTP para interactuar con la API de Jamf.
type Client struct {
	httpClient  *http.Client
	token       string
	baseUrl     string
	instanceUrl string
}

// NewClient crea una nueva instancia de Client.
// Parámetros:
// - httpClient: Cliente HTTP a utilizar.
// - token: Token de autenticación.
// - baseUrl: URL base de la API.
// - instanceUrl: URL de la instancia de la API.
// Retorna: Una instancia de Client.
func NewClient(httpClient *http.Client, token, baseUrl, instanceUrl string) *Client {
	return &Client{
		httpClient:  httpClient,
		token:       token,
		baseUrl:     baseUrl,
		instanceUrl: instanceUrl,
	}
}

// AuthResponse representa la respuesta de autenticación.
type AuthResponse struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

// CreateBearerToken genera un token de autenticación.
// Parámetros:
// - ctx: Contexto de ejecución.
// - username: Nombre de usuario.
// - password: Contraseña.
// - serverInstance: URL del servidor.
// Retorna: Un token de autenticación y un error si ocurre algún problema.
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

// GetTokenDetails obtiene detalles del token de autenticación actual.
// Retorna: Un objeto TokenDetails con la información del token y un error si falla la operación.
func (c *Client) GetTokenDetails(ctx context.Context) (TokenDetails, error) {
	url := fmt.Sprintf("%s/api/v1/auth", c.instanceUrl)

	var res TokenDetails
	if err := c.doRequest(ctx, url, &res); err != nil {
		return TokenDetails{}, err
	}

	return res, nil
}

// getBaseUsers obtiene la lista de usuarios base.
// Retorna: Un slice de BaseType y un error si falla la operación.
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

// getUserDetails obtiene los detalles de un usuario dado su ID.
// Retorna: Un objeto User y un error si falla la operación.
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

// GetUsers obtiene la lista de usuarios con información detallada.
// Retorna: Un slice de User y un error si ocurre un problema.
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

// doRequest ejecuta una petición HTTP GET y decodifica la respuesta JSON.
// Parámetros:
// - ctx: Contexto de ejecución.
// - url: URL de la petición.
// - res: Interfaz donde se decodificará la respuesta.
// Retorna: Un error si la operación falla.
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
