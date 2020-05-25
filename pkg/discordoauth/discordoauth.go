package discordoauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	endpointDiscordOauth = "https://discordapp.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s"
	endpointDiscordToken = "https://discordapp.com/api/oauth2/token"
	endpointDiscordMe    = "https://discordapp.com/api/users/@me"
)

var (
	errUnauthorized = errors.New("unauthorized")
)

type DiscordOauth struct {
	clientId     string
	clientSecret string
	redirectUri  string
	scopes       []string

	authUri string
}

func New(clientId, clientSecret, redirectUri string, scopes ...string) *DiscordOauth {
	d := &DiscordOauth{
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectUri:  redirectUri,
		scopes:       scopes,
	}

	if len(d.scopes) == 0 {
		d.scopes = append(d.scopes, "identify")
	}

	d.authUri = fmt.Sprintf(
		endpointDiscordOauth, clientId,
		url.QueryEscape(redirectUri),
		strings.Join(d.scopes, "%20"))

	return d
}

func (d *DiscordOauth) HandleInitialize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", d.authUri)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (d *DiscordOauth) HandleCallback(w http.ResponseWriter, r *http.Request) (*UserModel, error) {
	code := r.URL.Query().Get("code")

	token, err := d.getAuthToken(code)
	if err != nil {
		return nil, err
	}

	user, err := d.validateToken(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DiscordOauth) IsErrUnauthorized(err error) bool {
	return err == errUnauthorized
}

func (d *DiscordOauth) getAuthToken(code string) (string, error) {
	values := url.Values{
		"client_id":     []string{d.clientId},
		"client_secret": []string{d.clientSecret},
		"grant_type":    []string{"authorization_code"},
		"code":          []string{code},
		"redirect_uri":  []string{d.redirectUri},
		"scope":         d.scopes,
	}

	res, err := http.PostForm(endpointDiscordToken, values)
	if err != nil {
		return "", err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return "", errUnauthorized
	}

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("response code %d", res.StatusCode)
	}

	tokenResp := new(authTokenModel)
	if err = json.NewDecoder(res.Body).Decode(tokenResp); err != nil {
		return "", err
	}

	if tokenResp.Error != "" {
		return "", errors.New(tokenResp.Error)
	}

	return tokenResp.AccessToken, nil
}

func (d *DiscordOauth) validateToken(token string) (*UserModel, error) {
	req, _ := http.NewRequest("GET", endpointDiscordMe, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errUnauthorized
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("response code %d", res.StatusCode)
	}

	userRes := new(UserModel)
	if err = json.NewDecoder(res.Body).Decode(userRes); err != nil {
		return nil, err
	}

	if userRes.Error != "" {
		return nil, errors.New(userRes.Error)
	}

	return userRes, nil
}
