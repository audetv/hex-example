package usermemstore

import (
	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
	"github.com/google/uuid"
	"sync"
)

// Для проверки, что соответствует интерфейсу юзер бизнес логики
var _ user.UserStore = &Users{}

// Users коллекция. Защитим мьютексом, так к этой коллекции могут обращаться
// из разных запросов внешних, а они могут приходить параллельно,
type Users struct {
	sync.Mutex
	m map[uuid.UUID]user.User
}

func NewUsers() *Users {
	return &Users{
		m: make(map[uuid.UUID]user.User),
	}
}

func (us *Users) Create(u user.User) (*uuid.UUID, error)
func (us *Users) Read(uid uuid.UUID) (*user.User, error)
func (us *Users) Delete(uid uuid.UUID) error
func (us *Users) SearchUsers(s string) (chan user.User, error)
