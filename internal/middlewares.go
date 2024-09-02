package internal

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/carlqt/ezsplit/internal/auth"
)

type ContextKey string

const ContextKeySetCookie ContextKey = "setCookieFn"
const JWTCookie = "JWTCookie"

// JwtMiddleware extracts the JWT from the cookie.
// After extracting, it's decoded and stores the claim in context
// The error is ignored because the token is optional and the resolver will handle the error.
func JwtMiddleware(next http.Handler, tokenSecret []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getJwtCookieValue(r)

		// If there's an error in getting the cookie, just treat it as no cookie detected
		if token == "" || err != nil {
			next.ServeHTTP(w, r)
		} else {
			// If JWT exists. Try to validate and decode
			// if validated and decoded, add to context and call next
			claims, err := auth.ValidateJWT(token, tokenSecret)

			// If there's an error validating, treat it as no cookie
			if err != nil {
				slog.Warn("failed to validate token", "error", err.Error())
				next.ServeHTTP(w, r)
			}

			ctx := context.WithValue(r.Context(), auth.UserClaimKey, claims)
			newReq := r.WithContext(ctx)

			next.ServeHTTP(w, newReq)
		}
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

// getJwtCookieValue tries to extract the token from the Cookie.
// Supresses ErrNoCookie error since it's a valid scenario because not all requests requires a cookie
func getJwtCookieValue(r *http.Request) (string, error) {
	cookie, err := r.Cookie(JWTCookie)

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			slog.Warn("auth cookie not found")
			return "", nil
		default:
			slog.Error(err.Error())
			return "", errors.New("internal server error")
		}
	}

	return cookie.Value, nil
}
