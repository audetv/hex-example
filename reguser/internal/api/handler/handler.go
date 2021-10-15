package handler

import "net/http"

// Здесь мы делаем наш Mux, который буде заниматься обработкой, всего чего нам надо.
// Поробуем сделать на стандартном простом ServeMux

type Router struct {
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		http.NewServeMux(),
	}
}
