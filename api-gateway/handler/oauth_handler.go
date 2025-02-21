package handler

import (
	"api-gateway/config"
	"api-gateway/utils"
	"api-gateway/webResponse"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

func (h *UserHandler) GoogleLogin(c echo.Context) error {
	url := config.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return webResponse.ResponseJson(c, http.StatusBadRequest, nil, "Code not found")
	}

	token, err := config.GoogleOauthConfig.Exchange(c.Request().Context(), code)
	if err != nil {
		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to exchange token")
	}

	// fetch user info
	userInfo, err := fetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to fetch user info")
	}

	// generate JWT token
	token, err = config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to generate JWT token")
	}

	userInfo, err = fetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to fetch user info")
	}
	requestBody, _ := json.Marshal(map[string]string{
		"email":     userInfo["email"].(string),
		"username":  userInfo["name"].(string),
		"google_id": userInfo["id"].(string),
	})
	corrID := utils.GenerateCorrelationID()
	err = h.SendMessage.SendingToMessage("UserRegisteredGoogle", corrID, requestBody)
	if err != nil {
		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to publish message")
	}

	return h.ResponseHandler.HandleEventResponse(
		c,
		false,
		http.StatusCreated,
		h.Config.RequestTimeout,
		"User registered successfully",
		"UserRegisteredGoogleSuccess",
		"UserRegisteredGoogleFailed",
	)
}

func fetchGoogleUserInfo(accessToken string) (map[string]interface{}, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
