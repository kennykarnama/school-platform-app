package interceptor

import (
	"github.com/kennykarnama/school-adminstration-api/domain/service/user"
	"github.com/kennykarnama/school-adminstration-api/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Auth struct {
	userSvc            user.Service
	whitelistEndpoints []string
}

func NewAuth(userSvc user.Service, whiteListEndpoints []string) *Auth {
	return &Auth{
		userSvc:            userSvc,
		whitelistEndpoints: whiteListEndpoints,
	}
}

func (a *Auth) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		for _, whitelistEndpoint := range a.whitelistEndpoints {
			if whitelistEndpoint == request.RequestURI {
				next.ServeHTTP(writer, request)
				return
			}
		}
		logrus.Infof("action=authMiddleware requestUri=%v", request.RequestURI)
		token, err := getTokenFromCookie(request)
		if token == "" {
			token, err = getTokenFromHeader(request)
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		userSession, err := a.userSvc.Validate(request.Context(), token)
		if err == nil {
			request = request.WithContext(userSession.NewCtx(request.Context()))
			next.ServeHTTP(writer, request)
		} else {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
	})
}

func getTokenFromHeader(request *http.Request) (string, error) {
	token, err := util.ExtractBearerToken(request.Header.Get("Authorization"))
	if err != nil {
		return "", err
	}
	return token, nil
}

func getTokenFromCookie(request *http.Request) (string, error) {
	sessionCookie, err := request.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return sessionCookie.Value, nil
}
