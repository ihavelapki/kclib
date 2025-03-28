package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Функция получения access_token из Keycloak
func GetMasterToken(kcURL, username, password string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	form.Set("grant_type", "password")
	form.Set("client_id", "admin-cli")

	req, err := http.NewRequest("POST", kcURL+"/realms/master/protocol/openid-connect/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Transport = tr

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s, body: %s", resp.Status, body)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}
