package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"boiler/pkg/entity"
	"boiler/pkg/errors"
	"boiler/pkg/store/config"

	"github.com/go-chi/chi/middleware"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
)

// Recoverer recover from panic
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logEntry := middleware.GetLogEntry(r)
				caller := errors.Caller()
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
					stackTitle, stackTrace := errors.GetStack()
					fmt.Fprintf(
						os.Stderr,
						"%s\n    %s\n",
						stackTitle,
						strings.Join(stackTrace, "\n    "),
					)
				} else if e, is := rvr.(error); is {
					log.Error().Err(e).Str("file", caller).Msg("panic")
				} else {
					log.Error().Str("file", caller).Msg(rvr.(string))
				}

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// AuthUserMiddleware parse JWT Token and inject it back as a *entity.AuthUser from request if available
func AuthUserMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	prefixLen := len("Bearer ")

	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			if raw := r.Header.Get("Authorization"); len(raw) > prefixLen {
				token, err := jwt.ParseString(raw[prefixLen:], jwt.WithVerify(jwa.RS256, &cfg.JWT.PrivateKey.PublicKey))
				if err == nil && jwt.Verify(token) == nil {
					id, err := strconv.ParseInt(token.Subject(), 10, 64)
					if err == nil {
						next.ServeHTTP(w, r.WithContext(
							context.WithValue(
								r.Context(), config.ContextKeyAuthenticationUser{}, &entity.JWTUser{
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
}
