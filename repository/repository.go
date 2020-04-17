package repository

import (
	"errors"
	"fmt"
	"os"
	"user-service/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DatabaseRepository struct {
	connection *gorm.DB
}

// NewDatabaseHandler creates new instance of a Database connection and returns the Handler
func NewDatabaseHandler() (DatabaseRepository, error) {

	username := os.Getenv("DB_USER")
	if username == "" {
		username = "tokyuser"
	}
	password := os.Getenv("DB_PASS")
	if password == "" {
		password = "pwd"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "usermgmt"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	var databaseRepository DatabaseRepository
	databaseRepository.connection = conn
	databaseRepository.connection.AutoMigrate(&model.User{}, &model.Application{})
	databaseRepository.connection.Model(&model.Application{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	return databaseRepository, nil
}

func (repository *DatabaseRepository) FindByUserName(userName string) (model.User, error) {
	var person model.User
	err := repository.connection.Where("user_name = ?", userName).Preload("Applications").First(&person).Error
	return person, err
}

func (repository *DatabaseRepository) FindByEmail(email string) (model.User, error) {
	var person model.User
	err := repository.connection.Where("email = ?", email).Preload("Applications").First(&person).Error
	return person, err
}

func (repository *DatabaseRepository) FindByID(id uint) (model.User, error) {
	var person model.User
	var applications []model.Application
	err := repository.connection.Where("id = ?", id).First(&person).Error
	err = repository.connection.Model(&person).Related(&applications).Error
	person.Applications = applications
	return person, err
}

// FindAllUsers returns all Users as an array
func (repository *DatabaseRepository) FindAllUsers() ([]model.User, error) {
	var persons []model.User

	repository.connection.Find(&persons).Preload("Applications")
	return persons, nil

}

// CreateUser requires an user with userName or eMail and password
func (repository *DatabaseRepository) CreateUser(user model.User) (err error) {
	err = repository.connection.Create(&user).Error
	return
}

// UpdateUser persists user
func (repository *DatabaseRepository) UpdateUser(user *model.User) (err error) {

	err = repository.connection.Save(user).Error
	return
}

// FindByEmailOrUserName returns user or error
func (repository *DatabaseRepository) FindByEmailOrUserName(userName string) (model.User, error) {

	var user model.User
	var err error
	if user, err = repository.FindByUserName(userName); gorm.IsRecordNotFoundError(err) {
		if user, err = repository.FindByEmail(userName); gorm.IsRecordNotFoundError(err) {
			return user, errors.New("no User with username or email found")
		}
	} else {
		err = nil
	}

	return user, nil

}

func (repository *DatabaseRepository) CloseConnection() {
	repository.connection.Close()
}

type DatabaseHandler interface {
	FindUserName(string) (model.User, error)
	FindById(string) (model.User, error)
	FindByEmail(string) (model.User, error)
	FindAllUsers() ([]model.User, error)
	CreateUser(model.UserDTO) (err error)
	CheckPassword(string, string) (bool, error)
	CloseConnection()
}
