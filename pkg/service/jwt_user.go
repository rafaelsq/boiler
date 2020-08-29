package service

import (
	"context"
	"net/http"
	"strconv"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

// AuthUserMiddleware parse JWT Token and inject it back as a *entity.AuthUser from request if available
func (s *Service) AuthUserMiddleware(next http.Handler) http.Handler {
	prefixLen := len("Bearer ")
	fn := func(w http.ResponseWriter, r *http.Request) {
		if raw := r.Header.Get("Authorization"); len(raw) > prefixLen {
			token, err := jwt.ParseString(raw[prefixLen:], jwt.WithVerify(jwa.RS256, &s.config.JWT.PrivateKey.PublicKey))
			if err == nil && jwt.Verify(token) == nil {
				id, err := strconv.ParseInt(token.Subject(), 10, 64)
				if err == nil {
					next.ServeHTTP(w, r.WithContext(
						context.WithValue(
							r.Context(), iface.ContextKeyAuthenticationUser{}, &entity.JWTUser{
								ID: id,
							},
						),
					))
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
