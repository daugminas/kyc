package adapter

import (
	"github.com/daugminas/kyc/app/domain"
	"github.com/daugminas/kyc/lib/db"
)

type userAdapter struct {
	db             *db.MongoDB
	userCollection string
}

func NewUserAdapter(db *db.MongoDB, userCollection string) UserAdapter {
	return &userAdapter{
		db:             db,
		userCollection: userCollection,
	}
}

type UserAdapter interface {
	CreateUser(u *domain.User, makeActive bool) (userId string, err error)
	GetUser(userId string) (u *domain.User, err error)
	UpdateUser(userId string, userUpdate *domain.User) (updated *domain.User, err error)
	ActivateUser(userId string) (err error)
	DeActivateUser(userId string) (err error)
	DeleteUser(userId string) (err error)
}
