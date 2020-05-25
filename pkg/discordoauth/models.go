package discordoauth

type authTokenModel struct {
	AccessToken string `json:"access_token"`

	Error string `json:"error"`
}

type UserModel struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`

	Error string `json:"error"`
}
