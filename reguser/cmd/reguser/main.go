package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/audetv/hex-ecample/reguser/internal/app/starter"
	"github.com/audetv/hex-ecample/reguser/internal/db/mem/usermemstore"
)

func main() {
	// Создадим глобальный стартовый контекст, относительно бэкграунд контекста,
	// он будет прерываем по ctrl+c
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	ust := usermemstore.NewUsers()
	a := starter.NewApp(ust)

	// Канцелим контекст потом дожидаемся всех горутин
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
