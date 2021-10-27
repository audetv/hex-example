package auth

import (
	"net/http"
)

// AuthMiddleware вынесли в отдельный пакет auth,
// пакет находится на том же самом инфраструктурном уровне api, так как он относится к этому слою.
// Принимает next http.Handler и возвращает http.Handler.
// Это стандартного вида middleware стандартной библиотеке к gorilla/mux к go-chi/chi роутеру и т.д.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Проверяем авторизацию, если нет то 401 и выходим, а если все хорошо, то пробрасываем
			// writer и reader дальше в next обработчик. Такими замыканиями можно выстроить целую цепочку из middleware,
			// которые что-то делаю, до того как основные handlers получат writer и reader
			if u, p, ok := r.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			// Можно user'а, которого мы распознали положить в контекст и сделать request с контекстом, в котором сидит user
			// В закомментированной строке пустые структуры указаны, а так надо указывать реальную структуру user.
			// r = r.WithContext(context.WithValue(r.Context(), struct {}{}, struct {}{} ))
			// r = r.WithContext(context.WithValue(r.Context(), ctxUser{}, User{ID:"", Name:""} ))
			// таким образом в контексте у нас появится user с заданными параметрами.
			next.ServeHTTP(w, r)
		},
	)
}
