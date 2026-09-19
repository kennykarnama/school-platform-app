package interceptor

import (
	userEntity "github.com/kennykarnama/school-adminstration-api/domain/entity/user"
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
			teacher, profileErr := a.userSvc.Profile(request.Context(), userSession.UserId)
			if profileErr != nil || !teacher.Active || (teacher.SchoolID != nil && !teacher.SchoolActive) {
				http.Error(writer, userEntity.ErrAccountInactive.Error(), http.StatusUnauthorized)
				return
			}
			schoolID := ""
			if teacher.SchoolID != nil {
				schoolID = *teacher.SchoolID
			}
			principal := userEntity.Principal{
				UserID: teacher.Id, SchoolID: schoolID, Role: teacher.Role, MustChangePassword: teacher.MustChangePassword,
			}
			ctx := userSession.NewCtx(request.Context())
			ctx = principal.NewCtx(ctx)
			request = request.WithContext(ctx)
			if principal.MustChangePassword && request.URL.Path != "/api/v1/teacher/password" && request.URL.Path != "/api/v1/teacher/logout" && request.URL.Path != "/api/v1/teacher/session/validate" {
				http.Error(writer, userEntity.ErrPasswordChangeRequired.Error(), http.StatusForbidden)
				return
			}
			next.ServeHTTP(writer, request)
		} else {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
	})
}

func DevelopmentPrincipal(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal := userEntity.Principal{
			UserID: "30000000-0000-4000-8000-000000000001", SchoolID: "00000000-0000-4000-8000-000000000001", Role: userEntity.RoleSchoolAdmin,
		}
		r = r.WithContext(principal.NewCtx(r.Context()))
		next.ServeHTTP(w, r)
	})
}

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowed[role] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, err := userEntity.NewPrincipalFromCtx(r.Context())
			if err != nil || !allowed[principal.Role] {
				http.Error(w, userEntity.ErrForbidden.Error(), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
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
