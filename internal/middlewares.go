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
const AuthTokenKey = "authToken"

// BearerTokenMiddleware extracts the bearer token from the Authorization header
// and stores it in the context.
// The error is ignored because the token is optional and the resolver will handle the error.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := getAuthToken(r)

		if authToken == "" && err == nil {
			next.ServeHTTP(w, r)
		} else {
			if err != nil {
				slog.Warn(err.Error())
			}

			ctx := context.WithValue(r.Context(), auth.TokenKey, authToken)
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

// getAuthToken extracts the token from the source.
// The source can be a Header, if in dev mode, or a Cookie
// An error is returned if Authorization header is missing or the format is invalid.
func getAuthToken(r *http.Request) (string, error) {
	return fromCookie(r, AuthTokenKey)
}

func fromCookie(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			slog.Warn("cookie not found")
			return "", nil
		default:
			slog.Error(err.Error())
			return "", errors.New("internal server error")
		}
	}

	return cookie.Value, nil
}
