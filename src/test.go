package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"errors"
	"fmt"
	"strings"
	"github.com/gocql/gocql"
	"github.com/gin-gonic/gin"
	"net/http"
)


func main() {
	app.Run()
}


type RestErr struct {
	ErrMessage string        `json:"message"`
	ErrStatus  int           `json:"status"`
	ErrError   string        `json:"error"`
	ErrCauses  []interface{} `json:"causes"`
}


func NewRestError(message string, status int, err string, causes []interface{}) *RestErr {
	return &RestErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  causes
	}
}


func NewRestErrorFromBytes(bytes []byte) (*RestErr, error) {
	var apiErr RestErr
	
	err := json.Unmarshal(bytes, &apiErr)
	
	if (err != nil) {
		return nil, errors.New("Invalid JSON")
	}
	
	return &apiErr, nil
}


func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request"
	}
}


func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found"
	}
}


func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized"
	}
}


func NewInternalServerError(message string, err error) *RestErr {
	result := &RestErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error"
	}
	
	if (err != nil) {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	
	return result
}


const (
	queryGetUserById string = "SELECT user_id, first_name, last_name, full_name, age, email FROM users WHERE user_id = ?"
	queryCreateUser  string = "INSERT INTO users (user_id, first_name, last_name, full_name, age, email) VALUES (?, ?, ?, ?, ?, ?)"
)


func NewDbRepository() DbRepository {
	return &dbRepository{}
}


type DbRepository interface {
	GetByID(userID string) (*User, *RestErr)
	Create(user model.User) (*User, *RestErr)
}


type dbRepository struct {}


func (repo *dbRepository) Create(user User) (*User, *RestErr) {
	if err := cassandra.GetSession().Query(queryCreateUser,
		user.ID, user.FirstName, user.LastName, user.FullName, user.Age, user.EmailId).Exec(); err != nil {
		return nil, NewInternalServerError("unable to insert user in db", errors.New(err.Error()))
	}
	return &user, nil
}


func (repo *dbRepository) GetByID(userID string) (*User, *RestErr) {
	var user User
	
	if err := cassandra.GetSession().Query(queryGetUserById, userID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.FullName, &user.Age, &user.EmailId); err != nil {
		if err.Error() == "not found" {
			fmt.Println("here")
			return nil, NewInternalServerError("no user for given user id", errors.New(err.Error()))
		}
		return nil, NewInternalServerError("unable to find user in db", errors.New(err.Error()))
	}
	return &user, nil
}


type User struct {
	ID string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	FullName string `json:"full_name"`
	Age int `json:"age"`
	EmailId string `json:"email_id"`
}


func (user *User) ValidateUser() *RestErr {
	if user.Age < 0 {
		return errors.NewInternalServerError("age can't be less than 0", nil)
	}

	if user.EmailId == "" {
		return errors.NewInternalServerError("email id can't be blank", nil)
	}
	// TODO:: Add more validation
	return nil
}


type UserRepository interface {
	GetByID(userID string) (*User, *RestErr)
	Create(user User) (*User, *RestErr)
}


type UserService interface {
	GetByID(userID string) (*User, *RestErr)
	Create(user User) (*User, *RestErr)
}


type userService struct {
	repository UserRepository
}


func (s *userService) Create(user User) (*User, *RestErr) {
	if err := user.ValidateUser(); err != nil {
		return nil, err
	}
	return s.repository.Create(user)
}


func NewService(repository UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}


func (s *userService) GetByID(userID string) (*User, *RestErr) {
	userID = strings.TrimSpace(userID)
	if len(userID) == 0 {
		return nil, NewBadRequestError("invalid user id. UserId can't be empty")
	}
	user, err := s.repository.GetByID(userID)
	if err != nil {
		userNotFoundErr := fmt.Sprintf("user not found for user id %s", userID)
		return nil, NewInternalServerError(userNotFoundErr, errors.New("here"))
	}
	return user, nil
}


var (
	session *gocql.Session
)

func init() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "user"
	cluster.Consistency = gocql.Quorum

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
}

func GetSession() *gocql.Session {
	return session
}


var (
	router = gin.Default()
)

func StartApp()  {
	userHandler := http.NewHandler(NewService(NewDbRepository()))
	router.GET("/users/:user_id", userHandler.GetById)
	router.POST("/user", userHandler.Create)

	_ = router.Run(":8888")
}


type UserHandler interface {
	GetById(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type userHandler struct {
	userService UserService
}

func (userHandler userHandler) GetById(ctx *gin.Context) {
	userID := strings.TrimSpace(ctx.Param("user_id"))
	user, err := userHandler.userService.GetByID(userID)
	if err != nil {
		ctx.JSON(err.ErrStatus, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func NewHandler(userService UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}


func (userHandler *userHandler) Create(ctx *gin.Context)  {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.ErrStatus, restErr)
	}
	if _, err := userHandler.userService.Create(user); err != nil {
		ctx.JSON(err.ErrStatus, err)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
