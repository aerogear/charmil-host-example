package login

import (
	"context"
	// embed static HTML file
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aerogear/charmil-host-example/pkg/auth/token"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	"github.com/aerogear/charmil/core/utils/logging"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

//go:embed static/mas-sso-redirect-page.html
var masSSOredirectHTMLPage string

// handler for the MAS-SSO redirect page
type masRedirectPageHandler struct {
	IO            *iostreams.IOStreams
	CfgHandler    *config.CfgHandler
	Logger        logging.Logger
	ServerAddr    string
	Port          int
	AuthURL       *url.URL
	AuthOptions   []oauth2.AuthCodeOption
	State         string
	Oauth2Config  *oauth2.Config
	Ctx           context.Context
	TokenVerifier *oidc.IDTokenVerifier
	CancelContext context.CancelFunc
	Localizer     localize.Localizer
}

// nolint:funlen
func (h *masRedirectPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger

	callbackURL := fmt.Sprintf("%v%v", h.ServerAddr, r.URL.String())
	logger.Infoln("Redirected to callback URL:", callbackURL)

	if r.URL.Query().Get("state") != h.State {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	// nolint:govet
	oauth2Token, err := h.Oauth2Config.Exchange(h.Ctx, r.URL.Query().Get("code"), h.AuthOptions...)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}
	idToken, err := h.TokenVerifier.Verify(h.Ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err = idToken.Claims(&resp.IDTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessTkn, _ := token.Parse(resp.OAuth2Token.AccessToken)
	tknClaims, _ := token.MapClaims(accessTkn)
	userName, ok := tknClaims["preferred_username"]
	rawUsername := "unknown"
	if ok {
		rawUsername = fmt.Sprintf("%v", userName)
	}

	pageTitle := h.Localizer.LocalizeByID("login.redirectPage.title")
	pageBody := h.Localizer.LocalizeByID("login.masRedirectPage.body", localize.NewEntry("Host", h.AuthURL.Host), localize.NewEntry("Username", rawUsername))

	redirectPage := fmt.Sprintf(masSSOredirectHTMLPage, pageTitle, pageTitle, pageBody)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, redirectPage)

	// save the received tokens to the user's config
	h.CfgHandler.Cfg.MasAccessToken = oauth2Token.AccessToken
	h.CfgHandler.Cfg.MasRefreshToken = oauth2Token.RefreshToken

	h.CancelContext()
}
