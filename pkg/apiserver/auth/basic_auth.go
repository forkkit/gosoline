package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/applike/gosoline/pkg/cfg"
	"github.com/applike/gosoline/pkg/mon"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	ByBasicAuth = "basicAuth"

	headerBasicAuth = "Authorization"
	configBasicAuth = "api_auth_basic_users"
	AttributeUser   = "user"
)

type basicAuthAuthenticator struct {
	logger mon.Logger
	users  map[string]string
}

func NewBasicAuthHandler(config cfg.Config, logger mon.Logger) gin.HandlerFunc {
	auth := NewBasicAuthAuthenticator(config, logger)

	appName := config.GetString("app_name")

	return func(ginCtx *gin.Context) {
		valid, err := auth.IsValid(ginCtx)

		if valid {
			return
		}

		if err == nil {
			err = fmt.Errorf("the user credentials weren't valid nor was there an error")
		}

		ginCtx.Header("www-authenticate", fmt.Sprintf("Basic realm=\"%s\"", appName))
		ginCtx.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		ginCtx.Abort()
	}
}

func NewBasicAuthAuthenticator(config cfg.Config, logger mon.Logger) Authenticator {
	userEntries := config.GetStringSlice(configBasicAuth)

	users := make(map[string]string)

	for _, user := range userEntries {
		if user == "" {
			continue
		}

		split := strings.SplitN(user, ":", 2)
		if len(split) != 2 {
			logger.Panic(
				fmt.Errorf("invalid basic auth credentials: %s", user),
				"basic auth credentials have to be in the form user:password",
			)
		}

		users[split[0]] = split[1]
	}

	return NewBasicAuthAuthenticatorWithInterfaces(logger, users)
}

func NewBasicAuthAuthenticatorWithInterfaces(logger mon.Logger, users map[string]string) Authenticator {
	return &basicAuthAuthenticator{
		logger: logger,
		users:  users,
	}
}

func (a *basicAuthAuthenticator) IsValid(ginCtx *gin.Context) (bool, error) {
	basicAuth := ginCtx.GetHeader(headerBasicAuth)

	if basicAuth == "" {
		return false, fmt.Errorf("no credentials provided")
	}

	if !strings.HasPrefix(basicAuth, "Basic ") {
		return false, fmt.Errorf("invalid credentials provided")
	}

	auth, err := base64.StdEncoding.DecodeString(basicAuth[6:])

	if err != nil {
		return false, err
	}

	split := strings.SplitN(string(auth), ":", 2)

	if len(split) != 2 {
		return false, fmt.Errorf("invalid credentials provided")
	}

	if password, ok := a.users[split[0]]; ok {
		if password != split[1] {
			return false, fmt.Errorf("invalid credentials provided")
		}

		user := &Subject{
			Name:            Anonymous,
			Anonymous:       true,
			AuthenticatedBy: ByBasicAuth,
			Attributes: map[string]interface{}{
				AttributeUser: split[0],
			},
		}

		RequestWithSubject(ginCtx, user)

		return true, nil

	}

	return false, fmt.Errorf("invalid credentials provided")
}