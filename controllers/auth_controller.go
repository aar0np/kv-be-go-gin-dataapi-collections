package controllers

import (
	"context"
	"errors"
	"killrvideo/go-backend-astra-dataapi/models"
	repo "killrvideo/go-backend-astra-dataapi/repositories"
	"net/http"
	"os"
	"time"

	"github.com/datastax/astra-db-go/astra"
	astratypes "github.com/datastax/astra-db-go/astra/datatypes"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Time
	ExpRef   time.Time
	UserID   string
	Issuer   string
	Audience string
}

type AuthController struct {
	authDAL repo.AuthDAL
}

func NewAuthController(db *astra.Db, ctx context.Context) *AuthController {
	return &AuthController{
		authDAL: *repo.NewAuthDAL(db, ctx),
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var newUserReg models.UserRegistrationRequest

	if err1 := c.BindJSON(&newUserReg); err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
	}

	//var creds models.UserCredentials
	var user models.DBUser
	var jwtResp models.JwtResponse

	// generate userid
	var newUserId = astratypes.NewUUID().String()

	// generate user
	user.Email = newUserReg.Email
	user.Userid = newUserId
	user.FirstName = newUserReg.FirstName
	user.LastName = newUserReg.LastName
	user.CreatedDate = time.Now()
	user.LastLoginDate = time.Now()
	user.HashedPassword = hashPassword(newUserReg.Password)
	user.AccountStatus = "active"

	// save to DB
	ac.authDAL.SaveUser(user)
	//ac.authDAL.SaveUserCreds(creds)

	// gen token
	token, err3 := issueToken(newUserId, newUserReg.Email)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
		return
	}

	jwtResp.Email = newUserReg.Email
	jwtResp.UserID = newUserId
	jwtResp.Token = token.Access

	c.JSON(http.StatusOK, jwtResp)
}

func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err1 := c.BindJSON(&req); err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err1.Error()})
	}

	user, err2 := ac.authDAL.GetUserCredsByEmail(req.Email)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err2.Error()})
		return
	}

	hashedPassword := user.Password

	if !validatePassword(req.Password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error4": "Invalid password."})
		return
	}

	id := user.Userid
	token, err3 := issueToken(id, req.Email)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err3.Error()})
		return
	}

	// register login with database
	ac.authDAL.RegisterLogin(id)

	var jwtResp models.JwtResponse
	jwtResp.Email = req.Email
	jwtResp.UserID = id
	jwtResp.Token = token.Access

	c.JSON(http.StatusOK, jwtResp)
}

func (ac *AuthController) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err2 := ac.authDAL.GetUserById(id)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
	}

	c.JSON(http.StatusOK, user)
}

func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	// parse UserID from request
	userid, err1 := getUserIdFromToken(c)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
	}

	// get User from DB
	user, err2 := ac.authDAL.GetUserById(userid)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
	}

	c.JSON(http.StatusOK, user)
}

func (ac *AuthController) UpdateCurrentUser(c *gin.Context) {
	// parse UserID from context
	userid, err1 := getUserIdFromToken(c)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
		return
	}

	// bind request body
	var user models.User
	var updateUserReq models.UserUpdateRequest
	var userCreds models.UserCredentials
	var password string
	var passwordChanged bool
	var origEmail string
	var emailChanged bool

	// get User from DB
	userFromDB, err2 := ac.authDAL.GetUserById(userid)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		return
	}

	user.Userid = userid

	if err3 := c.BindJSON(&updateUserReq); err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
	} else {
		// Update only the fields that are present in the request
		if updateUserReq.FirstName != uuid.Nil.String() {
			user.FirstName = updateUserReq.FirstName
		}

		if updateUserReq.LastName != uuid.Nil.String() {
			user.LastName = updateUserReq.LastName
		}

		if updateUserReq.Password != uuid.Nil.String() {
			password = hashPassword(updateUserReq.Password)
			passwordChanged = true
			userCreds.Password = password
		}

		if updateUserReq.Email != "" && userFromDB.Email != updateUserReq.Email {
			// check if email already exists
			userCreds, err4 := ac.authDAL.GetUserCredsByEmail(updateUserReq.Email)

			if err4 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err4.Error()})
			}

			if userCreds != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "email already exists"})
				return
			}
			emailChanged = true
			origEmail = userFromDB.Email
			// set new email
			user.Email = updateUserReq.Email
			userCreds.Email = updateUserReq.Email
		} else {
			userCreds.Email = user.Email
		}

		if passwordChanged || emailChanged {
			userCreds.Userid = user.Userid
			ac.authDAL.UpdatePassword(userCreds)

			if emailChanged {
				ac.authDAL.DeleteUserCreds(origEmail)
			}
		}

		ac.authDAL.UpdateUser(user)
	}
}

func getUserIdFromToken(c *gin.Context) (string, error) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return astratypes.NewUUID().String(), errors.New("authorization header is required")
	}

	// Extract token from "Bearer <token>" format
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	claims, err1 := parseWithSecret(token)
	if err1 != nil {
		return astratypes.NewUUID().String(), err1
	}

	// get UserID from Subject
	//uuid, err2 := astratypes.ParseUUID(claims.Subject)
	//return uuid, err2
	return claims.Subject, nil
}

func validatePassword(password string, hashedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == nil {
		return true
	} else {
		return false
	}
}

func hashPassword(password string) string {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashed)
}

func issueToken(userID string, email string) (*Tokens, error) {
	now := time.Now().UTC()

	key := os.Getenv("JWT_KEY")

	t := &Tokens{
		UserID:   userID,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   now.Add(15 * time.Minute),
		ExpRef:   now.Add(7 * 24 * time.Hour),
		Issuer:   "kv-be-go-gin-dataapi-collections",
		Audience: "killrvideo-react-frontend",
	}

	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIAcc,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{userID, email},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIRef,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{userID, email},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpRef),
	})

	var err error
	t.Access, err = acc.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}
	t.Refresh, err = ref.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func parseWithSecret(tokenStr string) (*jwt.RegisteredClaims, error) {

	secret := os.Getenv("JWT_KEY")

	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Extra safety: ensure HMAC family
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
