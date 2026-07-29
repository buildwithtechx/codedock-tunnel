package auth

import "context"

type OAuthState struct {
	Provider    string
	RedirectURI string
	Verifier    string
}

type OAuthStateStore interface {
	Save(context.Context, string, OAuthState) error
	Take(context.Context, string) (OAuthState, error)
}
