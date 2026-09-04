package interceptor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kennykarnama/school-adminstration-api/domain/entity/user"
)

func TestRequireRoles(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	handler := RequireRoles(user.RoleSchoolAdmin)(next)

	tests := []struct {
		name      string
		principal *user.Principal
		want      int
	}{
		{name: "school administrator is allowed", principal: &user.Principal{Role: user.RoleSchoolAdmin}, want: http.StatusNoContent},
		{name: "teacher is forbidden", principal: &user.Principal{Role: user.RoleTeacher}, want: http.StatusForbidden},
		{name: "missing principal is forbidden", want: http.StatusForbidden},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.principal != nil {
				req = req.WithContext(test.principal.NewCtx(req.Context()))
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)
			if response.Code != test.want {
				t.Fatalf("status = %d, want %d", response.Code, test.want)
			}
		})
	}
}

func TestDevelopmentPrincipal(t *testing.T) {
	handler := DevelopmentPrincipal(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, err := user.NewPrincipalFromCtx(r.Context())
		if err != nil {
			t.Fatalf("principal: %v", err)
		}
		if principal.Role != user.RoleSchoolAdmin || principal.SchoolID == "" {
			t.Fatalf("unexpected principal: %+v", principal)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d", response.Code)
	}
}
