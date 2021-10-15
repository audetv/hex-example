package usermemstore

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

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
	// прежде чем вернуть канал, надо запустить горутину, в которой будем опять лочиться,
	// после лока должны будем перебрать мапу
	// Если по контексту прервали обработку и выходим, то в этом случае дефер закрывает канал
	// Горутина может злочится на канале, если мы заполнили весь буфер, мы будем ждать пока, в бизнес логике
	// вычитают данные из канала, а если бизнес	логика отвалилась и уже там не читают, то мьютекс останется залоченым
	// и хранилище будет залочено навсегда. Есть два пути решения: таймауты или проброска какого то сигнала о том, что
	// у нас бизнес логика не может воспользоваться этими данными.
	// Сигнал можно организовать в виде канала Done, который мы можем вернуть отдельно из функции SearchUser и допустим тогда,
	// бизнес логика его может закрыть на своей стороне, а тут в селект мы вставим еще один кейс на этот канал.
	// Но мы не будем усложнять и сделаем таймаут.
	// Чтобы это заработало, мы отправку должны поместить внутрь селекта
	// На стороне бизнес логики нельзя закрывать канал chout, потому что мы в него здесь пишем,
	// будет паника, поэтому нужен отдельный сигнальный канал.
	go func() {
		defer close(chout)
		us.Lock()
		defer us.Unlock()
		for _, u := range us.m {
			if strings.Contains(u.Name, s) {
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):

				case chout <- u:
				}

			}
		}

	}()
	return chout, nil
}
