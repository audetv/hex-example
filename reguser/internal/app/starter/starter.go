package starter

import (
	"context"
	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
	"sync"
)

// App Здесь мы должны стартануть приложение
type App struct {
	us *user.Users
}

// NewApp функция инициализации приложения, котора возвращает уже заполненный апп
// не хватает стора, получим его снаружи, пробросим в параметр user.UserStore
func NewApp(ust user.UserStore) *App {
	a := &App{
		us: user.NewUsers(ust),
	}
	return a
}

// Serve пока ничего не делает, ему на вход не хватает api, который будет присылать запросы
// Пробрасываем контекст, чтобы отловить сигналы от операционной системы
// и вейт группу, и в дефере завершаем ее
func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// <-ctx.Done()
}
