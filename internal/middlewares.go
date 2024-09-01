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
const BearerTokenCookie = "bearerTokenCookie"

// BearerTokenMiddleware extracts the bearer token from the BearerTokenCookie
// and stores it in the context.
// The error is ignored because the token is optional and the resolver will handle the error.
func BearerTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := getBearerToken(r)

		if bearerToken == "" && err == nil {
			next.ServeHTTP(w, r)
		} else {
			if err != nil {
				slog.Warn(err.Error())
			}

			ctx := context.WithValue(r.Context(), auth.TokenKey, bearerToken)
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

// getBearerToken tries to extract the token from the Cookie.
// Supresses ErrNoCookie error since it's a valid scenario because not all requests requires a cookie
func getBearerToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(BearerTokenCookie)

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
