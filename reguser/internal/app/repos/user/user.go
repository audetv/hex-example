package user

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Name        string
	Data        string
	Permissions int
}

// UserStore интерфейс системы хранения.
// Create возвращает указатель на uuid, чтобы не передавать пустой uuid в случае ошибки.
// Read возвращает указатель на User, чтобы не передавать пустого User в случае ошибки.
// Delete из системы хранения нам не надо возвращать самого юзера, т.к мы его прочитали в бизнес логике.
type UserStore interface {
	Create(u User) (*uuid.UUID, error)
	Read(uid uuid.UUID) (*User, error)
	Delete(uid uuid.UUID) error
	SearchUsers(s string) (chan User, error)
}

// Users коллекция объектов User, для того чтобы реализовать паттерн репозиторий,
// который работает с системой хранения, у него будут некоторые методы
type Users struct {
	ustore UserStore
}

// Create чтобы не передавать пустого юзера, вернем указатель на юзера.
// Получать будем полноценную карточку в виде структуры
func (us *Users) Create(u User) (*User, error) {
	u.ID = uuid.New()
	id, err := us.ustore.Create(u)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}
	u.ID = *id
	return &u, nil
}

func (us *Users) Read(uid uuid.UUID) (*User, error) {
	u, err := us.ustore.Read(uid)
	if err != nil {
		return nil, fmt.Errorf("read user error: %w", err)
	}
	return u, nil
}

func (us *Users) Delete(uid uuid.UUID) (*User, error) {
	u, err := us.ustore.Read(uid)
	if err != nil {
		return nil, fmt.Errorf("search user err: %w", err)
	}

	// Чтобы вызвать delete, мы можем просто вызвать ошибку полученную из UserStore
	return u, us.ustore.Delete(uid)
}

// SearchUsers устанавливаем для примера permissions для юзера, на уровне бизнес логики,
// система хранения ничего об этом не знает. Берем юзера из входящего канала,
// устанавливаем permissions и передаем в исходящий канал
// вычитываем пользователей в бесконечном цикле
func (us *Users) SearchUsers(s string) (chan User, error) {
	chin, err := us.ustore.SearchUsers(s)
	if err != nil {
		return nil, err
	}
	chout := make(chan User, 100)
	// на выходе из функции закрываем канал chout, поскольку мы пишем в этой горутине,
	// то тут и закрываем. Закрытие канал chout будет зависеть от закрытия канала chin,
	// соответственно должны проверить if !ok, то должны выйти
	go func() {
		defer close(chout)
		for {
			u, ok := <-chin
			if !ok {
				return
			}
			u.Permissions = 0755
			chout <- u
		}
	}()
	return chout, nil
}
