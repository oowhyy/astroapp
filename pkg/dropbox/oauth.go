package dropbox

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SLToken struct {
	AccessToken string `json:"access_token"`
}

func GetAccessToken(refreshToken, appAuth string) (string, error) {
	client := &http.Client{}
	bodyS := fmt.Sprintf("refresh_token=%s&grant_type=refresh_token", refreshToken)
	req, err := http.NewRequest("POST", "https://api.dropbox.com/oauth2/token", strings.NewReader(bodyS))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(appAuth))
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	var sltoken SLToken
	if err := json.NewDecoder(res.Body).Decode(&sltoken); err != nil {
		return "", err
	}
	return sltoken.AccessToken, nil
}
