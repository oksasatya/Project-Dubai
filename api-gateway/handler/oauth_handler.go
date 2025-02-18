package handler

//func (h *UserHandler) GoogleLogin(c echo.Context) error {
//	url := config.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
//	return c.Redirect(http.StatusTemporaryRedirect, url)
//}
//
//func (h *UserHandler) GoogleCallback(c echo.Context) error {
//	code := c.QueryParam("code")
//	if code == "" {
//		return webResponse.ResponseJson(c, http.StatusBadRequest, nil, "Code not found")
//	}
//
//	token, err := config.GoogleOauthConfig.Exchange(c.Request().Context(), code)
//	if err != nil {
//		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to exchange token")
//	}
//
//	// fetch user info
//	userInfo, err := fetchGoogleUserInfo(token.AccessToken)
//	if err != nil {
//		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to fetch user info")
//	}
//
//	// generate JWT token
//	token, err = config.GoogleOauthConfig.Exchange(context.Background(), code)
//	if err != nil {
//		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to generate JWT token")
//	}
//
//	userInfo, err = fetchGoogleUserInfo(token.AccessToken)
//	if err != nil {
//		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to fetch user info")
//	}
//	requestBody, _ := json.Marshal(map[string]string{
//		"email":     userInfo["email"].(string),
//		"username":  userInfo["name"].(string),
//		"google_id": userInfo["id"].(string),
//	})
//	corrID := time.Now().Format("20060102150405")
//	err = messaging.PublishMessage("user_register_oauth_queue", string(requestBody), "user_register_oauth_response_queue", corrID)
//	if err != nil {
//		return webResponse.ResponseJson(c, http.StatusInternalServerError, nil, "Failed to publish message")
//	}
//
//	return webResponse.HandleServiceResponse(
//		c,
//		"user_register_oauth_response_queue",
//		corrID,
//		false,
//		http.StatusOK,
//		h.Config.RequestTimeout,
//		"User logged in via Google successfully",
//	)
//}
//
//func fetchGoogleUserInfo(accessToken string) (map[string]interface{}, error) {
//	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	var userInfo map[string]interface{}
//	if err := json.Unmarshal(body, &userInfo); err != nil {
//		return nil, err
//	}
//
//	return userInfo, nil
//}
