package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

// Auth holds an OIDC ID token verifier.
type Auth struct {
	verifier *oidc.IDTokenVerifier
}

// NewAuth creates the middleware with given issuer URL.
func NewAuth(ctx context.Context, issuer string, clientID string) (*Auth, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	cfg := &oidc.Config{SkipClientIDCheck: true}
	return &Auth{verifier: provider.Verifier(cfg)}, nil
}

// Middleware verifies bearer tokens and injects claims into the context.
func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		raw := strings.TrimPrefix(auth, "Bearer ")
		idToken, err := a.verifier.Verify(r.Context(), raw)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		var claims struct {
			Email string `json:"email"`
			Scope string `json:"scope"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userEmail", claims.Email)
		ctx = context.WithValue(ctx, "scope", claims.Scope)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
