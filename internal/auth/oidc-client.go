package auth

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/coreos/go-oidc/v3/oidc"
// 	"github.com/gin-gonic/gin"
// 	"github.com/sirupsen/logrus"
// 	"golang.org/x/oauth2"
// )

// type OidcClient struct {
// 	provider oidc.Provider
// 	verfier  oidc.IDTokenVerifier
// 	// config   oauth2.Config // oauth config contains redirect uri and each request will have different reqiest so we can not use common config for all requests
// }

// func NewClient(url, clientID, clientSecret string) (*OidcClient, error) {
// 	var err error

// 	oidcProvider, err := oidc.NewProvider(context.Background(), url)
// 	if err != nil {
// 		logrus.WithFields(logrus.Fields{
// 			"url":   url,
// 			"error": err,
// 		}).Error("failed to get oidc provider config")
// 		return nil, nil
// 	}

// 	oidcVerifier := oidcProvider.Verifier(&oidc.Config{
// 		ClientID: clientID,
// 	})

// 	return &OidcClient{
// 		provider: *oidcProvider,
// 		verfier:  *oidcVerifier,
// 	}, nil
// }

// func (client OidcClient) Authorize(c *gin.Context) {
// 	state, err := randString(16)
// 	if err != nil {
// 		http.Error(w, "Internal error", http.StatusInternalServerError)
// 		return
// 	}

// 	nonce, err := randString(16)
// 	if err != nil {
// 		http.Error(w, "Internal error", http.StatusInternalServerError)
// 		return
// 	}

// 	codeChallenge, err := randString(16)
// 	if err != nil {
// 		http.Error(w, "Internal error", http.StatusInternalServerError)
// 		return
// 	}

// 	setCallbackCookie(w, r, "state", state)
// 	log.Printf("state = %v", state)
// 	setCallbackCookie(w, r, "nonce", nonce)
// 	setCallbackCookie(w, r, "code_challenge", codeChallenge)

// 	// code challenge vs nonce
// 	// https://danielfett.de/2020/05/16/pkce-vs-nonce-equivalent-or-not/
// 	u := oauth2Config.AuthCodeURL(state,
// 		oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256(codeChallenge)),
// 		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
// 		oidc.Nonce(nonce),
// 	)
// 	http.Redirect(w, r, u, http.StatusFound)
// }

// func (client OidcClient) Authorize2(c *gin.Context, redirectURI string, requiredScopes []string) {
// 	r.ParseForm()
// 	state, err := r.Cookie("state")
// 	if err != nil {
// 		http.Error(w, "state not found", http.StatusBadRequest)
// 		return
// 	}

// 	if r.URL.Query().Get("state") != state.Value {
// 		http.Error(w, "state did not match", http.StatusBadRequest)
// 		return
// 	}

// 	// we get token from the code
// 	code := r.Form.Get("code")
// 	if code == "" {
// 		http.Error(w, "Code not found", http.StatusBadRequest)
// 		return
// 	}

// 	codeChallenge, err := r.Cookie("code_challenge")
// 	if err != nil {
// 		http.Error(w, "code_challenge not found", http.StatusBadRequest)
// 		return
// 	}

// 	// convert code to token
// 	token, err := oauth2Config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", codeChallenge.Value))
// 	// id_token, ok := token.Extra("id_token").(string)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	globalToken = token
// 	rawIDToken, ok := token.Extra("id_token").(string)
// 	if !ok {
// 		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
// 		return
// 	}
// 	idToken, err := oidcVerifier.Verify(r.Context(), rawIDToken)
// 	if err != nil {
// 		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	nonce, err := r.Cookie("nonce")
// 	if err != nil {
// 		http.Error(w, "nonce not found", http.StatusBadRequest)
// 		return
// 	}
// 	if idToken.Nonce != nonce.Value {
// 		http.Error(w, "nonce did not match", http.StatusBadRequest)
// 		return
// 	}
// 	resp := struct {
// 		OAuth2Token   *oauth2.Token
// 		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
// 	}{token, new(json.RawMessage)}

// 	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data, err := json.MarshalIndent(resp, "", "    ")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(data)

// }
