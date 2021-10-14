package usermemstore

import (
	"context"
	"database/sql"
	"sync"

	"github.com/audetv/hex-ecample/reguser/internal/app/repos/user"
	"github.com/google/uuid"
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

func (us *Users) Create(ctx context.Context, u user.User) (*uuid.UUID, error) {
	us.Lock()
	defer us.Unlock()

	// make select, если контекст прервался, вернем нил и ошибку из контекста, почему контекст прервался,
	// а если не был прерван, то ничего не делаем - default
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	us.m[u.ID] = u
	return &u.ID, nil
}
func (us *Users) Read(ctx context.Context, uid uuid.UUID) (*user.User, error) {
	us.Lock()
	defer us.Unlock()

	// контекст нужно проверять только после того как залочились
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	u, ok := us.m[uid]
	if ok {
		return &u, nil
	}
	// Если ничего не найдено, то нил и для красоты типизированная ошибка,
	// которую можно проверять на равенство sql.ErrNoRows
	return nil, sql.ErrNoRows
}

// Delete не возвращает ошибку если не нашли
func (us *Users) Delete(ctx context.Context, uid uuid.UUID) error {
	us.Lock()
	defer us.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	delete(us.m, uid)
	return nil
}
func (us *Users) SearchUsers(ctx context.Context, s string) (chan user.User, error) {
	us.Lock()
	defer us.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// FIXME: переделать на дерево остатков

	// Мы будем просто проходить мапу, чтобы пройти надо создать канал, мы же возвращаем канал.
	chout := make(chan user.User, 100)
	// прежде чем вернуть канал, надо запустить горутину, в которой будем опять лочиться
	go func() {

	}()
	return chout, nil
}
