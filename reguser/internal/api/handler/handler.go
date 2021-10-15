package handler

import (
	"encoding/json"
	"net/http"

	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
	"github.com/google/uuid"
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
// Используем ее для получения данных юзера от клиента или отправки данных клиенту
// Парсим, декодируем.
type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Data        string    `json:"data"`
	Permissions int       `json:"permissions"`
}

func (rt *Router) CreateUser(w http.ResponseWriter, r *http.Request) {
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
	u := User{}
	if err := dec.Decode(&u); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	bu := user.User{
		Name: u.Name,
		Data: u.Data,
	}
	// У каждого запроса приходящего есть контекст внутри и его мы можем
	// использовать и пробрасывать дальше в нужные нам методы, этот контекст канцелится если мы остановим сервер.
	nbu, err := rt.us.Create(r.Context(), bu)
	if err != nil {
		http.Error(w, "error when creating user", http.StatusInternalServerError)
		return
	}
	// Если создание пользователя произошло корректно, появляется заполненный айди у юзера,
	// пермишены и мы можем вернуть это обратно клиенту. Создаем енкодер, ошибку проскипаем,
	// тут маловероятно возникновение ошибки, разве что на самом потоке, если сетевой
	// поток прервался, сетевое соединение, но тогда нам не о чем и некому сообщать возвращать эту ошибку,
	// разве что залогировать. Но при успешном создании нужно вернуть код 201 Created,
	// по умолчанию Encode возвращает код 200 OK, для этого надо указать код ответа.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(User{
		ID:          nbu.ID,
		Name:        nbu.Name,
		Data:        nbu.Data,
		Permissions: nbu.Permissions,
	})
}
func (*Router) ReadUser(w http.ResponseWriter, r *http.Request) {
}
func (*Router) DeleteUser(w http.ResponseWriter, r *http.Request) {
}
func (*Router) SearchUser(w http.ResponseWriter, r *http.Request) {
}
