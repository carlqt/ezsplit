package internal

import (
	"context"
	"net/http"

	"github.com/carlqt/ezsplit/internal/auth"
)

var ContextKeySetCookie string = "setCookieFn"

// BearerTokenMiddleware extracts the bearer token from the Authorization header
// and stores it in the context.
// The error is ignored because the token is optional and the resolver will handle the error.
func BearerTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, _ := auth.GetBearerToken(r)

		ctx := context.WithValue(r.Context(), auth.TokenKey, bearerToken)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}

// InjectSetCookieMiddleware adds SetCookie method to the context so the resolvers can call it
func InjectSetCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCookieFn := func(cookie *http.Cookie) {
			http.SetCookie(w, cookie)
		}

		ctx := context.WithValue(r.Context(), ContextKeySetCookie, setCookieFn)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}
