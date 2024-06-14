package stores

type UserStoreInterface interface {
	*UserPostgresStoreInterface
}

type UserStore struct {
	*UserPostgresStore
}

func NewUserStore() *UserStore {
	return &UserStore{
		UserPostgresStore: NewUserPostgresStore(),
	}
}
