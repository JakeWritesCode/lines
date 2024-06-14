package stores

import "lines/lines/store"

type UserStoreInterface interface {
	*UserPostgresStoreInterface
}

type UserStore struct {
	*UserPostgresStore
}

func NewUserStore() *UserStore {
	return &UserStore{
		UserPostgresStore: NewUserPostgresStore(
			[]store.PostgresModel{
				&User{},
			},
		),
	}
}
