package repositories

import (
	"context"
	"fmt"
	"killrvideo/go-backend-astra-dataapi/models"
	"time"

	astradb "github.com/datastax/astra-db-go"
	astratypes "github.com/datastax/astra-db-go/datatypes"

	"github.com/datastax/astra-db-go/filter"
	"github.com/datastax/astra-db-go/update"
)

type AuthDAL struct {
	DB  *astradb.Db
	Ctx context.Context
}

func NewAuthDAL(db *astradb.Db, ctx context.Context) *AuthDAL {
	return &AuthDAL{
		DB:  db,
		Ctx: ctx,
	}
}

func (r *AuthDAL) GetUserById(id astratypes.UUID) (*models.User, error) {
	collection := r.DB.Collection("users")

	var dbuser models.DBUser

	err1 := collection.FindOne(r.Ctx, filter.Eq("userid", id)).Decode(&dbuser)
	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	user := &models.User{
		Userid:        dbuser.Userid,
		Email:         dbuser.Email,
		FirstName:     dbuser.FirstName,
		LastName:      dbuser.LastName,
		AccountStatus: dbuser.AccountStatus,
		CreatedDate:   dbuser.CreatedDate,
		LastLoginDate: dbuser.LastLoginDate,
	}

	return user, nil
}

func (r *AuthDAL) GetUserByEmail(email string) (*models.User, error) {
	collection := r.DB.Collection("users")

	var dbuser models.DBUser

	err1 := collection.FindOne(r.Ctx, filter.Eq("email", email)).Decode(&dbuser)
	if err1 != nil {
		return nil, fmt.Errorf("GUbE1 query has failed: %w", err1)
	}

	user := &models.User{
		Userid:        dbuser.Userid,
		Email:         dbuser.Email,
		FirstName:     dbuser.FirstName,
		LastName:      dbuser.LastName,
		AccountStatus: dbuser.AccountStatus,
		CreatedDate:   dbuser.CreatedDate,
		//LastLoginDate: dbuser.LastLoginDate,
	}

	return user, nil
}

func (r *AuthDAL) ExistsByEmail(email string) bool {
	user, err := r.GetUserByEmail(email)
	if user != nil && err == nil {
		return true
	}
	return false
}

func (r *AuthDAL) GetUserCredsByEmail(email string) (*models.UserCredentials, error) {
	collection := r.DB.Collection("users")

	var dbuser models.DBUser
	//var raw json.RawMessage

	err1 := collection.FindOne(r.Ctx, filter.Eq("email", email)).Decode(&dbuser)
	if err1 != nil {
		return nil, fmt.Errorf("GUCbE1 query has failed: %w", err1)
	}

	//fmt.Printf("raw == %s\n", raw)

	user := &models.UserCredentials{
		Userid:   dbuser.Userid,
		Email:    dbuser.Email,
		Password: dbuser.HashedPassword,
	}

	if dbuser.AccountStatus == "suspended" {
		user.AccountLocked = true
	} else {
		user.AccountLocked = false
	}

	return user, nil
}

func (r *AuthDAL) SaveUser(user models.DBUser) {
	collection := r.DB.Collection("users")

	collection.InsertOne(r.Ctx, user)
}

func (r *AuthDAL) UpdateUser(user models.User) {
	collection := r.DB.Collection("users")

	upd := update.Coll()
	upd.Set("email", user.Email)
	upd.Set("firstname", user.FirstName)
	upd.Set("lastname", user.LastName)
	upd.Set("account_status", user.AccountStatus)
	upd.Set("last_login_date", user.LastLoginDate)
	upd.Set("created_date", user.CreatedDate)

	collection.UpdateOne(r.Ctx, filter.Eq("userid", user.Userid), upd)
}

func (r *AuthDAL) SaveUserCreds(userCreds models.UserCredentials) {
	collection := r.DB.Collection("users")

	collection.UpdateOne(r.Ctx, filter.Eq("userid", userCreds.Userid),
		update.Coll().Set("password", userCreds.Password).Set("email", userCreds.Email))
}

func (r *AuthDAL) UpdatePassword(userCreds models.UserCredentials) {
	collection := r.DB.Collection("users")

	collection.UpdateOne(r.Ctx, filter.Eq("userid", userCreds.Userid),
		update.Coll().Set("password", userCreds.Password))
}

func (r *AuthDAL) DeleteUserCreds(email string) {
	collection := r.DB.Collection("users")
	collection.DeleteOne(r.Ctx, filter.Eq("email", email))
}

func (r *AuthDAL) RegisterLogin(userid string) {
	collection := r.DB.Collection("users")

	collection.UpdateOne(r.Ctx, filter.Eq("userid", userid),
		update.Coll().Set("last_login_date", time.Now()))
}
