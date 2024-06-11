package adaptors

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	*gorm.DB
}

type PostgresAdaptor struct {
	db *GormDB
}

func NewPostgresAdaptor(DBConnString string) *PostgresAdaptor {
	db, err := gorm.Open(postgres.Open(DBConnString), &gorm.Config{})
	if err != nil {

	}
	return &PostgresAdaptor{
		db: &GormDB{db},
	}
}

db, err := gorm.Open(postgres.Open(config.DBConnString), &gorm.Config{})
if err != nil {
panic(err)
}