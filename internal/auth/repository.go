package auth

import (
	sql "github.com/tanmaygupta069/auth-service-go/pkg/mysql"
)

type RepositoryInterface interface {
	IsUserRegistered(email string) bool
	SaveUser(email string, password string) error
	GetHashedPassword(email string) (string, error)
}

type RepositoryImplementation struct {
	db sql.SqlInterface
}

func NewRepository() RepositoryInterface {
	return &RepositoryImplementation{
		db: sql.NewSqlClient(),
	}
}

func (db *RepositoryImplementation) IsUserRegistered(email string) bool {
	res, _ := db.db.GetOne(email, "email")
	if res == nil {
		return false
	}
	return true
}

func (db *RepositoryImplementation) SaveUser(email string, password string) error {

	return db.db.Insert(&sql.User{
		Email:    email,
		Password: password,
	})
}

func (db *RepositoryImplementation) GetHashedPassword(email string) (string, error) {
	userData, err := db.db.GetOne(email, "email")
	if err != nil {
		return "", err
	}

	return userData.Password, nil
}
