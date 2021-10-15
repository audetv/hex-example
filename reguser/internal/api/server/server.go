package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
)

// Server принимает запросы, и вызывает бизнес логику
// Входящий адаптер обращается в бизнес логику us user.Users
// Должен открыть листенер для http протокола, мы должны в него встроить http сервер
type Server struct {
	srv http.Server
	us  *user.Users
}

// NewServer Адрес и порт передавать через параметр
func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

// Делаем два метода, надо стартовать и остановить сервер,
// Есть два варианта: ListenAndServe() и ListenAndServeTLS() - для https надо указать сертификационные файлы,
// имена, которые в операционной системе резолвятся на файлы с сертификатами.
// В примере используем http. Но если просто вызвать, то листен остановится, значит делаем это в горутине.
// Так же можно обработать ошибку и вывести например ее в логи, в общемто больше с ней нечего делать.

func (s *Server) Start(us *user.Users) {
	s.us = us
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil {
			log.Printf("serve error %v:", err)
		}
	}()
}

// Stop метод для остановки сервера, для этого у http сервера есть Shutdown(), который принимает контекст.
// Этот контекст сделаем с таймаутом
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_ = s.srv.Shutdown(ctx)
	cancel()
}
