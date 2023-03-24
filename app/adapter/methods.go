package adapter

import (
	"time"

	"github.com/daugminas/kyc/app/domain"
	"github.com/daugminas/kyc/app/utils"
	"github.com/daugminas/kyc/lib/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a userAdapter) CreateUser(u *domain.User, makeActive bool) (userId string, err error) {

	// hass password
	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return
	}

	// insert to db
	u.Password = hash
	u.CreatedAt = time.Now()
	u.Active = makeActive
	result, err := a.db.InsertOne(a.userCollection, nil, u)
	if err != nil {
		return
	}
	id := result.InsertedID.(primitive.ObjectID)
	userId = id.String()

	return
}

func (a userAdapter) GetUser(userId string) (u *domain.User, err error) {
	u = &domain.User{}
	filter := make(db.Filter)
	mgoId, _ := primitive.ObjectIDFromHex(userId)
	filter["_id"] = mgoId

	err = a.db.FindOne(a.userCollection, filter, nil, u)

	return
}

func (a userAdapter) UpdateUser(userId string, userUpdate *domain.User) (updated *domain.User, err error) {
	filter := make(db.Filter)
	mgoId, _ := primitive.ObjectIDFromHex(userId)
	filter["_id"] = mgoId

	update := make(db.Change)
	if userUpdate.Email != "" {
		update["email"] = userUpdate.Email
	}
	if userUpdate.Password != "" {
		var hash string
		hash, err = utils.HashPassword(userUpdate.Password)
		if err != nil {
			return
		}
		update["password"] = hash
	}
	update["active"] = userUpdate.Active
	update["updated_at"] = time.Now()

	result, e := a.db.FindOneAndUpdate(a.userCollection, filter, db.Change{"$set": update})
	if e != nil {
		err = e
		return
	}
	updated = &domain.User{} // the updated user in DB
	result.Decode(updated)

	return
}

func (a userAdapter) ActivateUser(userId string) (err error) {
	filter := make(db.Filter)
	mgoId, _ := primitive.ObjectIDFromHex(userId)
	filter["_id"] = mgoId

	update := make(db.Change)
	update["active"] = true
	update["updated_at"] = time.Now()

	_, err = a.db.UpdateOne(a.userCollection, filter, update, nil)

	return
}

func (a userAdapter) DeActivateUser(userId string) (err error) {
	filter := make(db.Filter)
	mgoId, _ := primitive.ObjectIDFromHex(userId)
	filter["_id"] = mgoId

	update := make(db.Change)
	update["active"] = false
	update["updated_at"] = time.Now()

	_, err = a.db.UpdateOne(a.userCollection, filter, update, nil)

	return
}

func (a userAdapter) DeleteUser(userId string) (err error) {
	filter := make(db.Filter)
	mgoId, _ := primitive.ObjectIDFromHex(userId)
	filter["_id"] = mgoId

	u := &domain.User{} // need this to decode result into, can't be nil
	err = a.db.FindOneAndDelete(a.userCollection, filter, nil, u)

	return
}
