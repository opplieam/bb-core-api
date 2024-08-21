package store

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/opplieam/bb-core-api/.gen/buy-better-core/public/model"
	. "github.com/opplieam/bb-core-api/.gen/buy-better-core/public/table"
)

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

func (s *UserStore) InsertOrUpdateUser(email, firstName, lastName, role string) error {
	stmt := Users.
		INSERT(Users.Email, Users.FirstName, Users.LastName, Users.Role).
		MODEL(model.Users{Email: email, FirstName: &firstName, LastName: &lastName, Role: role}).
		ON_CONFLICT(Users.Email).
		DO_UPDATE(SET(
			Users.LoginAt.SET(CURRENT_TIMESTAMP()),
		))
	_, err := stmt.Exec(s.DB)
	if err != nil {
		return DBTransformError(err)
	}
	return nil
}

func (s *UserStore) FindUserByEmail(email string) (model.Users, error) {
	stmt := SELECT(Users.AllColumns).
		FROM(Users).
		WHERE(
			Users.Email.EQ(String(email)).AND(Users.Active.IS_TRUE()),
		)
	var dest model.Users
	err := stmt.Query(s.DB, &dest)
	if err != nil {
		return model.Users{}, DBTransformError(err)
	}
	return dest, nil
}
