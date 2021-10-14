package starter

import "github.com/audetv/hex-ecample/reguser/internal/app/repos/user"

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
func (a *App) Serve() {
}
