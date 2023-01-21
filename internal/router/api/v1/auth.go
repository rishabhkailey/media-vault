package v1

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session/v3"
	"github.com/sirupsen/logrus"
)

// check for token in session if there then proceed and else save request and redirect to auth server
func (server Server) AuthMiddleware(c *gin.Context) {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"function": "server.AuthMiddleware",
		}).Errorf("session start failed")
		// todo error response
		c.Status(http.StatusInternalServerError)
		return
	}
	var accessToken string
	if token, ok := store.Get("access_token"); ok {
		accessToken, ok = token.(string)
		if !ok {
			logrus.WithFields(logrus.Fields{
				"function": "server.AuthMiddleware",
			}).Errorf("access token type cast failed")
			// todo error response
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	var idToken string
	if token, ok := store.Get("id_token"); ok {
		idToken, ok = token.(string)
		if !ok {
			logrus.WithFields(logrus.Fields{
				"function": "server.AuthMiddleware",
			}).Errorf("id token type cast failed")
			// todo error response
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	// also check expire time of access token + jwt signature key verification
	// we will only verify the jwt signature after getting the access token not here, here we will only check the id token expiret time
	if len(accessToken) != 0 && len(idToken) != 0 {
		// ?
		// c.AddParam("access_token", accessToken)
		// c.AddParam("id_token", idToken)
		c.Next()
		return
	}

	// state, err := randString(16)
	// if err != nil {
	// 	http.Error(w, "Internal error", http.StatusInternalServerError)
	// 	return
	// }

	// nonce, err := randString(16)
	// if err != nil {
	// 	http.Error(w, "Internal error", http.StatusInternalServerError)
	// 	return
	// }

	// codeChallenge, err := randString(16)
	// if err != nil {
	// 	http.Error(w, "Internal error", http.StatusInternalServerError)
	// 	return
	// }

	// setCallbackCookie(w, r, "state", state)
	// log.Printf("state = %v", state)
	// setCallbackCookie(w, r, "nonce", nonce)
	// setCallbackCookie(w, r, "code_challenge", codeChallenge)

	// // code challenge vs nonce
	// // https://danielfett.de/2020/05/16/pkce-vs-nonce-equivalent-or-not/
	// u := oauth2Config.AuthCodeURL(state,
	// 	oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256(codeChallenge)),
	// 	oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	// 	oidc.Nonce(nonce),
	// )
	// http.Redirect(w, r, u, http.StatusFound)
	// c.Next()
}

// handle request from auth server and read the saved request and redirect user to the that endpoint
func (server Server) Oauth(c *gin.Context, redirectURI string, requiredScopes []string) {
	// r.ParseForm()
	// state, err := r.Cookie("state")
	// if err != nil {
	// 	http.Error(w, "state not found", http.StatusBadRequest)
	// 	return
	// }

	// if r.URL.Query().Get("state") != state.Value {
	// 	http.Error(w, "state did not match", http.StatusBadRequest)
	// 	return
	// }

	// // we get token from the code
	// code := r.Form.Get("code")
	// if code == "" {
	// 	http.Error(w, "Code not found", http.StatusBadRequest)
	// 	return
	// }

	// codeChallenge, err := r.Cookie("code_challenge")
	// if err != nil {
	// 	http.Error(w, "code_challenge not found", http.StatusBadRequest)
	// 	return
	// }

	// // convert code to token
	// token, err := oauth2Config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", codeChallenge.Value))
	// // id_token, ok := token.Extra("id_token").(string)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// globalToken = token
	// rawIDToken, ok := token.Extra("id_token").(string)
	// if !ok {
	// 	http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
	// 	return
	// }
	// idToken, err := oidcVerifier.Verify(r.Context(), rawIDToken)
	// if err != nil {
	// 	http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// nonce, err := r.Cookie("nonce")
	// if err != nil {
	// 	http.Error(w, "nonce not found", http.StatusBadRequest)
	// 	return
	// }
	// if idToken.Nonce != nonce.Value {
	// 	http.Error(w, "nonce did not match", http.StatusBadRequest)
	// 	return
	// }
	// resp := struct {
	// 	OAuth2Token   *oauth2.Token
	// 	IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	// }{token, new(json.RawMessage)}

	// if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// data, err := json.MarshalIndent(resp, "", "    ")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(data)

}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
