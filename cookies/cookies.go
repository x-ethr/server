package cookies

import (
	"net/http"
	"os"
	"strings"
	"time"
)

type Options struct {
	Domain string
}

type Variadic func(o *Options)

func settings() *Options {
	return &Options{}
}

func Secure(w http.ResponseWriter, name, value string, options ...Variadic) {
	var o = settings()
	for _, option := range options {
		option(o)
	}

	domain := o.Domain
	if domain == "" {
		if v := os.Getenv("CI"); strings.ToLower(v) != "true" {
			domain = ""
		} else if v := os.Getenv("NAMESPACE"); v == "development" {
			domain = ""
		}
	}

	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(3 * time.Hour),
		MaxAge:   86400,
		Secure:   true,                    // Ensure the cookie is sent only over HTTPS
		HttpOnly: true,                    // Prevent JavaScript from accessing the cookie
		SameSite: http.SameSiteStrictMode, // Enforce SameSite policy
	}

	http.SetCookie(w, &cookie)
}

// func Invalid() *http.Cookie {
// 	var secure bool = false
// 	domain := secret.Cookie().Domain()
// 	if domain != "localhost" {
// 		secure = true
// 	}
//
// 	return &http.Cookie{
// 		Name:     "token",
// 		Value:    "",
// 		Path:     "/",
// 		Domain:   domain,
// 		Expires:  time.Unix(0, 0),
// 		Secure:   secure,
// 		HttpOnly: false,
// 		SameSite: http.SameSiteLaxMode,
// 	}
//
// }
//
// func New(ctx context.Context, claims *jwts.Claims) (*http.Cookie, error) {
// 	var secure bool = false
// 	domain := secret.Cookie().Domain()
// 	if domain != "localhost" {
// 		secure = true
// 	}
//
// 	signature, err := claims.Sign()
// 	if err != nil {
// 		slog.ErrorContext(ctx, "Error Signing Claims Data", slog.String("error", err.Error()))
// 		return nil, err
// 	}
//
// 	return &http.Cookie{
// 		Name:     "token",
// 		Value:    signature,
// 		Path:     "/",
// 		Domain:   domain,
// 		Expires:  claims.ExpiresAt.Time,
// 		Secure:   secure,
// 		HttpOnly: false,
// 		SameSite: http.SameSiteLaxMode,
// 	}, nil
// }
