package v1

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session/v3"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type OidcClient struct {
	provider     oidc.Provider
	verfier      oidc.IDTokenVerifier
	oauth2Config oauth2.Config
}

func NewOidcClient(url, clientID, clientSecret, redirectURI string) (*OidcClient, error) {
	var err error

	oidcProvider, err := oidc.NewProvider(context.Background(), url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":   url,
			"error": err,
		}).Error("failed to get oidc provider config")
		return nil, nil
	}

	oidcVerifier := oidcProvider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{oidc.ScopeOpenID, "user", "email"},
		RedirectURL:  redirectURI,
		Endpoint:     oidcProvider.Endpoint(),
	}

	return &OidcClient{
		provider:     *oidcProvider,
		verfier:      *oidcVerifier,
		oauth2Config: oauth2Config,
	}, nil
}

// todo: oauth2.TokenSource interface implementation
// func (oidcClient *OidcClient) Token() (*oauth2.Token, error) {

// }

// check for token in session if there then proceed and else save request and redirect to auth server
func (server Server) AuthMiddleware(c *gin.Context) {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.AuthMiddleware"}).Errorf("session start failed")
		// todo error response
		c.Status(http.StatusInternalServerError)
		return
	}

	accessToken, idToken, err := getTokensFromSession(store)
	if errors.Is(err, ErrTokenNotFound) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// todo refresh?
	if err != nil || accessToken.Expiry.Before(time.Now()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// also check expire time of access token + jwt signature key verification
	// we will only verify the jwt signature after getting the access token not here, here we will only check the id token expiret time
	c.Set("access_token", *accessToken)
	c.Set("id_token", *idToken)
	c.Next()
}

func (server Server) LoginHandler(c *gin.Context) {
	// todo check if user is already logged in
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.AuthMiddleware"}).Errorf("session start failed")
		// todo error response
		c.Status(http.StatusInternalServerError)
		return
	}

	var returnUri string
	if c.Request.Form == nil {
		if err := c.Request.ParseForm(); err != nil {
			logrus.WithFields(logrus.Fields{"function": "server.AuthMiddleware"}).Errorf("request form parse failed")
			c.Status(http.StatusInternalServerError)
			return
		}
	}
	returnUri = c.Request.FormValue("returnUri")
	if len(returnUri) == 0 {
		// todo get home endpoint from config or some other source?
		returnUri = "/"
	}

	state, err := randString(16)
	if err != nil {
		// todo error response
		c.Status(http.StatusInternalServerError)
	}

	nonce, err := randString(16)
	if err != nil {
		// todo error response
		c.Status(http.StatusInternalServerError)
	}

	codeChallenge, err := randString(16)
	if err != nil {
		// todo error response
		c.Status(http.StatusInternalServerError)
	}

	store.Set("state", state)
	store.Set("nonce", nonce)
	store.Set("code_challenge", codeChallenge)
	store.Set("return_uri", returnUri)
	store.Save()

	// code challenge vs nonce
	// https://danielfett.de/2020/05/16/pkce-vs-nonce-equivalent-or-not/
	authURI := server.OidcClient.oauth2Config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256(codeChallenge)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oidc.Nonce(nonce),
	)
	c.Redirect(http.StatusFound, authURI)

}

func (server Server) LogoutHandler(c *gin.Context) {
	// todo revoke token api call
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.AuthMiddleware"}).Errorf("session start failed")
		// todo error response
		c.Status(http.StatusInternalServerError)
		return
	}

	store.Delete("id_token")
	store.Delete("access_token")
	if store.Save() != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// handle request from auth server and read the saved request and redirect user to the that endpoint
func (server Server) AuthHandler(c *gin.Context) {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.AuthMiddleware"}).Errorf("session start failed")
		// todo error response
		c.Status(http.StatusInternalServerError)
		return
	}
	var state string
	if value, ok := store.Get("state"); ok {
		state, ok = value.(string)
		if !ok {
			logrus.WithFields(logrus.Fields{"function": "server.loginHandler"}).Errorf("state type cast failed")
			// todo error response
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Request.ParseForm()
	if c.Request.Form.Get("state") != state {
		logrus.Info("[server.loginHandler]: state mismatched")
		c.Status(http.StatusBadRequest)
		return
	}

	// we get token from the code
	code := c.Request.Form.Get("code")
	if code == "" {
		logrus.Info("[server.loginHandler]: code challenge missing")
		c.Status(http.StatusBadRequest)
		return
	}

	var codeChallenge string
	if value, ok := store.Get("code_challenge"); ok {
		codeChallenge, ok = value.(string)
		if !ok {
			logrus.WithFields(logrus.Fields{"function": "server.loginHandler"}).Errorf("code_challenge type cast failed")
			// todo error response
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	// convert code to token
	token, err := server.OidcClient.oauth2Config.Exchange(c.Request.Context(), code, oauth2.SetAuthURLParam("code_verifier", codeChallenge))
	// id_token, ok := token.Extra("id_token").(string)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	idToken, err := server.OidcClient.verfier.Verify(c.Request.Context(), rawIDToken)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var nonce string
	if value, ok := store.Get("nonce"); ok {
		nonce, ok = value.(string)
		if !ok {
			logrus.WithFields(logrus.Fields{"function": "server.loginHandler"}).Errorf("nonce type cast failed")
			// todo error response
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	if idToken.Nonce != nonce {
		c.Status(http.StatusBadRequest)
		return
	}

	var idTokenClaims struct {
		Email string `json:"email"`
	}

	// todo redirect to error page instead of returning 500 status
	if err := idToken.Claims(&idTokenClaims); err != nil {
		logrus.Errorf("[server.loginHandler]: claims unmarshell failed %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var returnUri string
	if value, ok := store.Get("return_uri"); ok {
		if returnUri, ok = value.(string); !ok {
			logrus.Errorf("[PostOAuthAuthorizeHandler]: get return_uri from session failed %v")
			c.Status(http.StatusInternalServerError)
			return
		}
	}
	store.Delete("return_uri")
	store.Save()

	store.Set("id_token", idToken)
	store.Set("access_token", token)
	store.Save()

	c.Redirect(http.StatusFound, returnUri)

}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

type UserInfoResponse struct {
	Email    string `json:"email"`
	Profile  string `json:"profile,omitempty"`
	UserName string `json:"userName"`
}

func (server Server) UserInfo(c *gin.Context) {
	var accessToken oauth2.Token
	{
		var value any
		var exists bool
		if value, exists = c.Get("access_token"); !exists {
			c.Status(http.StatusInternalServerError)
			return
		}
		var ok bool
		if accessToken, ok = value.(oauth2.Token); !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
	}
	// todo refresh?
	if accessToken.Expiry.Before(time.Now()) {
		// todo error response
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userInfo, err := server.OidcClient.provider.UserInfo(c.Request.Context(), oauth2.StaticTokenSource(&accessToken))
	if err != nil {
		// todo error response
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(userInfo.Email) == 0 {
		logrus.WithFields(logrus.Fields{"error": err, "function": "userInfo"}).Error("empty email response from server")
		// todo error response
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &UserInfoResponse{
		Email:    userInfo.Email,
		Profile:  userInfo.Profile,
		UserName: userInfo.Email,
	})
}

var ErrTokenNotFound error = errors.New("token not found in session store")

func getTokensFromSession(store session.Store) (*oauth2.Token, *oidc.IDToken, error) {
	var accessToken *oauth2.Token
	var idToken *oidc.IDToken
	if value, ok := store.Get("access_token"); ok {
		if err := utils.UnmarshalInterface(value, &accessToken); err != nil {
			logrus.WithFields(logrus.Fields{
				"function": "server.AuthMiddleware",
				"error":    err,
			}).Errorf("access token type cast failed")
			return accessToken, idToken, fmt.Errorf("[getTokensFromSession]: access token type cast failed")
		}
	}

	if value, ok := store.Get("id_token"); ok {
		if err := utils.UnmarshalInterface(value, &idToken); err != nil {
			logrus.WithFields(logrus.Fields{
				"function": "server.AuthMiddleware",
				"error":    err,
			}).Errorf("id token type cast failed")
			// todo error response
			return accessToken, idToken, fmt.Errorf("[getTokensFromSession]: access token type cast failed")
		}
	}

	// also check expire time of access token + jwt signature key verification
	// we will only verify the jwt signature after getting the access token not here, here we will only check the id token expiret time
	if accessToken == nil || idToken == nil {
		return accessToken, idToken, ErrTokenNotFound
	}
	return accessToken, idToken, nil
}

// func (server Server) SaveRequestAndRedicrectAuthMiddleware(c *gin.Context) {
// 	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
// 	if err != nil {
// 		logrus.WithFields(logrus.Fields{
// 			"error":    err,
// 			"function": "server.AuthMiddleware",
// 		}).Errorf("session start failed")
// 		// todo error response
// 		c.Status(http.StatusInternalServerError)
// 		return
// 	}
// 	var accessToken oauth2.Token
// 	if value, ok := store.Get("access_token"); ok {
// 		if err := utils.UnmarshalInterface(value, &accessToken); err != nil {
// 			logrus.WithFields(logrus.Fields{
// 				"function": "server.AuthMiddleware",
// 				"error":    err,
// 			}).Errorf("access token type cast failed")
// 			// todo error response
// 			c.Status(http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	var idToken oidc.IDToken
// 	if value, ok := store.Get("id_token"); ok {
// 		if err := utils.UnmarshalInterface(value, &idToken); err != nil {
// 			logrus.WithFields(logrus.Fields{
// 				"function": "server.AuthMiddleware",
// 				"error":    err,
// 			}).Errorf("id token type cast failed")
// 			// todo error response
// 			c.Status(http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	// also check expire time of access token + jwt signature key verification
// 	// we will only verify the jwt signature after getting the access token not here, here we will only check the id token expiret time
// 	if len(idToken.Subject) != 0 && len(accessToken.AccessToken) != 0 && accessToken.Expiry.After(time.Now()) {
// 		c.Set("access_token", accessToken)
// 		c.Set("id_token", idToken)
// 		// ?
// 		// c.AddParam("access_token", accessToken)
// 		// c.AddParam("id_token", idToken)
// 		c.Next()
// 		return
// 	}

// 	state, err := randString(16)
// 	if err != nil {
// 		// todo error response
// 		c.Status(http.StatusInternalServerError)
// 	}

// 	nonce, err := randString(16)
// 	if err != nil {
// 		// todo error response
// 		c.Status(http.StatusInternalServerError)
// 	}

// 	codeChallenge, err := randString(16)
// 	if err != nil {
// 		// todo error response
// 		c.Status(http.StatusInternalServerError)
// 	}
// 	if c.Request.Form == nil {
// 		if err := c.Request.ParseForm(); err != nil {
// 			logrus.WithFields(logrus.Fields{
// 				"function": "server.AuthMiddleware",
// 			}).Errorf("request form parse failed")
// 			c.Status(http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	store.Set("state", state)
// 	store.Set("nonce", nonce)
// 	store.Set("code_challenge", codeChallenge)
// 	store.Set("original_request_uri", c.Request.RequestURI)
// 	store.Set("original_request_form", c.Request.Form)
// 	store.Save()

// 	// code challenge vs nonce
// 	// https://danielfett.de/2020/05/16/pkce-vs-nonce-equivalent-or-not/
// 	authURI := server.OidcClient.oauth2Config.AuthCodeURL(state,
// 		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256(codeChallenge)),
// 		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
// 		oidc.Nonce(nonce),
// 	)
// 	c.Redirect(http.StatusFound, authURI)
// }
