package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type accessToken struct {
	AccessToken string
	ExpiresAt   time.Time
}

type KCClient struct {
	baseUrl      string
	realmName    string
	clientId     string
	clientSecret string
	token        *accessToken
}

func New(baseUrl, realmName, clientId, clientSecret string) *KCClient {
	return &KCClient{
		baseUrl:      baseUrl,
		realmName:    realmName,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (c *KCClient) refreshToken(ctx context.Context) error {
	tokenUrl := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseUrl, c.realmName)

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", c.clientId)
	form.Add("client_secret", c.clientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewBadResponseErr(resp)
	}

	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	token := tokenResponse{}

	tokenDecoder := json.NewDecoder(resp.Body)
	err = tokenDecoder.Decode(&token)
	if err != nil {
		return err
	}

	c.token = &accessToken{
		AccessToken: token.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}

	return nil
}

func (c *KCClient) sendRequestWithToken(ctx context.Context, request *http.Request) (*http.Response, error) {
	if c.token == nil || c.token.ExpiresAt.Before(time.Now()) {
		err := c.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

	body, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *KCClient) GetUserByID(ctx context.Context, userId string) (*KCUser, error) {
	userUrl := fmt.Sprintf(
		"%s/admin/realms/%s/users/%s",
		c.baseUrl,
		c.realmName,
		userId,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", userUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.sendRequestWithToken(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, UserNotFoundErr
	} else if resp.StatusCode != 200 {
		return nil, NewBadResponseErr(resp)
	}

	user := KCUser{}
	userDecoder := json.NewDecoder(resp.Body)
	err = userDecoder.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
