package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
	"github.com/google/uuid"
)

// Отдельно выносим пакет, который относится к самому роутеру, чтобы нам была возможность разные роутеры подключать
// на одну и туже реализацию. Реализация у нас будет универсальная, немного модифицированный предыдущий пакет.
// Полностью убрали из нее вещи связанные со связанностью, с конкретной реализацией роутера, дефолтного http.
// Оставили 4 функции от нашего crud'а, который был. Каждая из этих функций оперирует карточкой user'а, которая
// маршалится в json и размаршаливается из json'а и в каждую функцию передается какой-то контекст.

type Handlers struct {
	us *user.Users
}

func NewRouter(us *user.Users) *Handlers {
	r := &Handlers{
		us: us,
	}
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

// AuthMiddleware принимает next http.Handler и возвращает http.Handler
// это стандартного вида middleware стандартной библиотеке к горилле к чи роутеру и т.д
func (rt *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Проверяем авторизацию, если нет то 401 и выходим, а если все хорошо, то пробрасываем
			// writer и reader дальше в next обработчик. Такими замыканиями можно выстроить целую цепочку из middlware,
			// которые что-то делаю, до того как основные хэндлеры получат writer и reader
			if u, p, ok := r.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func (rt *Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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

// ReadUser надо повторить проверку авторизации, сделаем middleware
// read?uid=...
func (rt *Router) ReadUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	suid := r.URL.Query().Get("uid")
	if suid == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(suid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (uid == uuid.UUID{}) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nbu, err := rt.us.Read(r.Context(), uid)
	// проверяем если db вернул ошибку sql.ErrNoRows, мы пробросили ее через %w, значит она будет распознана
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error when reading user", http.StatusInternalServerError)
		}
		return
	}
	_ = json.NewEncoder(w).Encode(User{
		ID:          nbu.ID,
		Name:        nbu.Name,
		Data:        nbu.Data,
		Permissions: nbu.Permissions,
	})
}

func (rt *Router) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	suid := r.URL.Query().Get("uid")
	if suid == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(suid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (uid == uuid.UUID{}) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nbu, err := rt.us.Delete(r.Context(), uid)
	// проверяем если db вернул ошибку sql.ErrNoRows, мы пробросили ее через %w, значит она будет распознана
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error when reading user", http.StatusInternalServerError)
		}
		return
	}
	_ = json.NewEncoder(w).Encode(User{
		ID:          nbu.ID,
		Name:        nbu.Name,
		Data:        nbu.Data,
		Permissions: nbu.Permissions,
	})
}

// SearchUser /search?q=...
func (rt *Router) SearchUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// передаем контекст и строку запроса q и возвращается канал ch, в котором мы будем стримить юзеров
	ch, err := rt.us.SearchUsers(r.Context(), q)
	// Ошибка может произойти если она в самом сторе произошла.
	// Там она возникает, только если мы в закрытом контексте находимся, по большому счету ее можно и проскипать.
	if err != nil {
		http.Error(w, "error when searching", http.StatusInternalServerError)
		return
	}
	// Все выполняется в горутинах, соответственно здесь у нас тоже отдельная горутина.
	// Каждый handler вызывается в отдельной горутине, которую принял http server на входе,
	// когда к нему подконнектился клиент. Под каждый запрос клиента создается отдельная горутина в http сервере.
	// И эта горутина в итоге приходит сюда в этот метод через роутер,
	// Т.е на нужные обработчики приходит ровно одна горутина, связанная с одним запросом.
	// Например, если клиент выполнил 10 запросов, то у нас запустится 10 горутин с SearchUsers хэндлером
	// и они могут параллельно исполняться. Соответственно мы в отдельной горутине и у нас ничего не зависнет,
	// если мы пытаться будем читать из этого канала. Но у нас есть два момента: есть cancel контексты,
	// у нас есть какие-то ошибки возникающие при отправке. Делаем в бесконечном цикле через select

	enc := json.NewEncoder(w)
	first := true
	fmt.Fprintf(w, "[")
	defer fmt.Fprintln(w, "]")

	for {
		select {
		case <-r.Context().Done():
			return
		case u, ok := <-ch:
			// Если у нас закрылся канал, а нам его закрыла бизнес логика,
			// а бизнес логика его закрыла если его закрыла база данных, то тогда return
			if !ok {
				return
			}
			if first {
				first = false
			} else {
				fmt.Fprintf(w, ",")
			}
			_ = enc.Encode(
				User{
					ID:          u.ID,
					Name:        u.Name,
					Data:        u.Data,
					Permissions: u.Permissions,
				},
			)
			w.(http.Flusher).Flush()
		}
	}
}
