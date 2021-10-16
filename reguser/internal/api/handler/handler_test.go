package handler

import (
	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
	"github.com/audetv/hex-ecample/reguser/internal/db/mem/usermemstore"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_CreateUser(t *testing.T) {
	ust := usermemstore.NewUsers()
	us := user.NewUsers(ust)
	rt := NewRouter(us)
	h := rt.AuthMiddleware(http.HandlerFunc(rt.CreateUser)).ServeHTTP

	w := &httptest.ResponseRecorder{}
	r := httptest.NewRequest("POST", "/create", strings.NewReader(`{"name":"user"}`))
	r.SetBasicAuth("admin", "admin")
	h(w, r)
	if w.Code == http.StatusUnauthorized {
		t.Errorf("status unauthorized")
	}
	if w.Code != http.StatusCreated {
		t.Errorf("status created")
	}
}
