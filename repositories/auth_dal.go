package repositories

import (
	"context"
	"fmt"
	"killrvideo/go-backend-astra-cql/models"
	"time"

	astradb "github.com/datastax/astra-db-go"
	"github.com/datastax/astra-db-go/filter"
	"github.com/datastax/astra-db-go/update"

)

type AuthDAL struct {
	DB *astradb.Database
	Ctx context.Context
}

func NewAuthDAL(db *astradb.Database, ctx context.Context) *AuthDAL {
	return &AuthDAL{
		DB: db,
		Ctx: ctx,
	}
}

func (r *AuthDAL) GetUserById(id astradb.timeuuid) (*models.User, error) {
	collection := r.DB.GetCollection(r.Ctx, "users")
	
	var dbuser models.DBUser

	err1 := collection.FindOne(r.Ctx, filter.Eq("userid",id)).Decode(&dbuser)
	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	user := &models.User{
		Userid: dbuser.Userid,
		Email: dbuser.email,
		FirstName: dbuser.FirstName,
		LastName: dbuser.LastName,
		AccountStatus: dbuser.AccountStatus,
		CreatedDate: dbuser.CreatedDate,
		LastLoginDate: dbuser.LastLoginDate,
	}

	return user, nil
}

func (r *AuthDAL) GetUserByEmail(email string) (*models.User, error) {
	collection := r.DB.GetCollection(r.Ctx, "users")

	var dbuser models.DBUser

	err1 := collection.FindOne(r.Ctx, filter.Eq("email",email)).Decode(&dbuser)
	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	user := &models.User{
		Userid: dbuser.Userid,
		Email: dbuser.email,
		FirstName: dbuser.FirstName,
		LastName: dbuser.LastName,
		AccountStatus: dbuser.AccountStatus,
		CreatedDate: dbuser.CreatedDate,
		LastLoginDate: dbuser.LastLoginDate,
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
	collection := r.DB.GetCollection(r.Ctx, "users")

	var dbuser models.DBUser

	err1 := collection.FindOne(r.Ctx, filter.Eq("email",email)).Decode(&dbuser)
	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	user := &models.UserCredentials{
		Userid: dbuser.Userid,
		Email: dbuser.email,
		Password: dbuser.HashedPassword,
	}

	if dbuser.AccountStatus == "suspended" {
		user.AccountLocked = true
	} else {
		user.AccountLocked = false
	}

	return user, nil
}

func (r *AuthDAL) SaveUser(user models.User) {
	collection := r.DB.GetCollection(r.Ctx, "users")

	collection.InsertOne(r.Ctx, user)
}

func (r *AuthDAL) UpdateUser(user models.User) {
	collection := r.DB.GetCollection(r.Ctx, "users")

	//	Userid         astradb.timeuuid `json:"userid"`
	//	Email          string           `json:"email"`
	//	FirstName      string           `json:"firstname"`
	//	LastName       string           `json:"lastname"`
	//	AccountStatus  string           `json:"account_status"`
	//	HashedPassword string 			`json:"hashed_password"`
	//	CreatedDate    time.Time        `json:"created_date"`
	//	LastLoginDate  time.Time        `json:"last_login_date"`

	collection.ReplaceOne(filter.Eq("userid",user.Userid), user)
}

func (r *AuthDAL) SaveUserCreds(user models.UserCredentials) {
	collection := r.DB.GetCollection(r.Ctx, "users")

}
