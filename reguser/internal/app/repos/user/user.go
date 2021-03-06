package user

import (
	"context"
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
	Create(ctx context.Context, u User) (*uuid.UUID, error)
	Read(ctx context.Context, uid uuid.UUID) (*User, error)
	Delete(ctx context.Context, uid uuid.UUID) error
	SearchUsers(ctx context.Context, s string) (chan User, error)
}

// Users коллекция объектов User, для того чтобы реализовать паттерн репозиторий,
// который работает с системой хранения, у него будут некоторые методы
type Users struct {
	ustore UserStore
}

// NewUsers функция инициализации, пробрасываем систему хранения в виде UserStore, будем возвращать Users,
// Но не с пустым store, его надо принять на вход, возьмем в параметр: ustore UserStore и присвоим ustore: ustore
func NewUsers(ustore UserStore) *Users {
	return &Users{
		ustore: ustore,
	}
}

// Create чтобы не передавать пустого пользователя, вернем указатель на него.
// Получать будем полноценную карточку в виде структуры.
func (us *Users) Create(ctx context.Context, u User) (*User, error) {
	// FIXME: здесь нужно использовать паттерн Unit of Work
	// бизнес-транзакция, нужно создать транзакцию либо внутри бизнес логики, например заложить сюда мьютекс отдельный
	// и везде его использовать в бизнес логике, либо создать транзакцию на уровне базы данных и внутри нее выполнять операции
	u.ID = uuid.New()
	id, err := us.ustore.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}
	u.ID = *id
	return &u, nil
}

func (us *Users) Read(ctx context.Context, uid uuid.UUID) (*User, error) {
	// FIXME: здесь нужно использовать паттерн Unit of Work
	// бизнес-транзакция
	u, err := us.ustore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("read user error: %w", err)
	}
	return u, nil
}

func (us *Users) Delete(ctx context.Context, uid uuid.UUID) (*User, error) {
	// FIXME: здесь нужно использовать паттерн Unit of Work
	// бизнес-транзакция
	u, err := us.ustore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("search user err: %w", err)
	}

	// Чтобы вызвать delete, мы можем просто вызвать ошибку полученную из UserStore
	return u, us.ustore.Delete(ctx, uid)
}

// SearchUsers устанавливаем для примера permissions для юзера, на уровне бизнес логики,
// система хранения ничего об этом не знает. Берем пользователя из входящего канала,
// устанавливаем permissions и передаем в исходящий канал
// вычитываем пользователей в бесконечном цикле
func (us *Users) SearchUsers(ctx context.Context, s string) (chan User, error) {
	// FIXME: здесь нужно использовать паттерн Unit of Work
	// бизнес-транзакция
	chin, err := us.ustore.SearchUsers(ctx, s)
	if err != nil {
		return nil, err
	}
	chout := make(chan User, 100)
	// На выходе из функции закрываем канал chout, поскольку мы пишем в этой горутине,
	// то тут и закрываем. Закрытие канала chout будет зависеть от закрытия канала chin,
	// соответственно должны проверить if !ok, то должны выйти
	go func() {
		defer close(chout)
		for {
			// Select - сидим и ждем, если ни один канал ничего не выдает просто ждем,
			// нам не надо в цикле крутиться для этого.
			select {
			case <-ctx.Done():
				return
			case u, ok := <-chin:
				if !ok {
					return
				}
				u.Permissions = 0755
				chout <- u
			}
		}
	}()
	return chout, nil
}
