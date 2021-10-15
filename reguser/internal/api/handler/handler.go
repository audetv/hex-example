package handler

import (
	"encoding/json"
	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
	"github.com/google/uuid"
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

// User - реализует отдельную структуру, которая не зависит от бизнес логики.
// Парсим, декодируем.
type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Data string    `json:"data"`
}

func (*Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию, если нет то 401 и выходим.
	if u, p, ok := r.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Body нужно закрывать, если начали из него читать, по умолчанию в handler в r http.Request приходят только заголовки
	// На бади можно не смотреть и работать и читать только заголовки, а оставшееся тело будет проигнорировано
	// и не будет даже загружено в память, если начинаем работать с телом, то реквест превращается в такой объект,
	// который аллоцирует память, т.е уже начинается накопление в памяти и по завершению этого хзндлера,
	// го должен явно знать, что мы закончили с ним работать, т.е его надо явно закрыть.
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
}
func (*Router) ReadUser(w http.ResponseWriter, r *http.Request) {
}
func (*Router) DeleteUser(w http.ResponseWriter, r *http.Request) {
}
func (*Router) SearchUser(w http.ResponseWriter, r *http.Request) {
}
