package handler

import (
	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
	"net/http"
)

// Здесь мы делаем наш Mux, который буде заниматься обработкой, всего чего нам надо.
// Попробуем сделать на стандартном простом ServeMux

type Router struct {
	*http.ServeMux
	us *user.Users
}

func NewRouter(us *user.Users) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		us:       us,
	}
	r.Handle("/create", http.HandlerFunc(r.CreateUser))
	r.Handle("/read", http.HandlerFunc(r.ReadUser))
	r.Handle("/delete", http.HandlerFunc(r.DeleteUser))
	r.Handle("/search", http.HandlerFunc(r.SearchUser))
	return r
}

func (*Router) CreateUser(w http.ResponseWriter, r *http.Request) {
}
func (*Router) ReadUser(w http.ResponseWriter, r *http.Request) {
}
func (*Router) DeleteUser(w http.ResponseWriter, r *http.Request) {
}
func (*Router) SearchUser(w http.ResponseWriter, r *http.Request) {
}
