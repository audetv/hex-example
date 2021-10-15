package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/audetv/hex-ecample/reguser/internal/api/handler"
	"github.com/audetv/hex-ecample/reguser/internal/api/server"
	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"

	"github.com/audetv/hex-ecample/reguser/internal/app/starter"
	"github.com/audetv/hex-ecample/reguser/internal/db/mem/usermemstore"
)

func main() {
	// Создадим глобальный стартовый контекст, относительно бэкграунд контекста,
	// он будет прерываем по ctrl+c
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	ust := usermemstore.NewUsers()
	a := starter.NewApp(ust)
	us := user.NewUsers(ust)

	h := handler.NewRouter(us)

	srv := server.NewServer(":8000", h)

	// Канцелим контекст потом дожидаемся всех горутин
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
