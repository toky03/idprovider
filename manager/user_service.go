package manager

import (
	"errors"
	"log"
	"user-service/model"
	"user-service/repository"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// UserService Business Logic for managing Users
type UserService struct {
	databaseHandler repository.DatabaseRepository
}

func NewUserService() UserService {
	databaseHandler, err := repository.NewDatabaseHandler()
	if err != nil {
		log.Fatal(err)
	}

	return UserService{
		databaseHandler: databaseHandler,
	}
}

func (s *UserService) FindAllUsers() ([]model.UserDTO, error) {

	users, err := s.databaseHandler.FindAllUsers()
	if err != nil {
		log.Println(err)
		return []model.UserDTO{}, err
	}

	userDTOs := make([]model.UserDTO, 0, len(users))

	for _, user := range users {

		userDTOs = append(userDTOs, mapUserToDTO(user))

	}

	return userDTOs, err

}

func (s *UserService) FindUser(userID uint) (model.UserDTO, error) {
	user, err := s.databaseHandler.FindByID(userID)
	if err != nil {
		log.Println(err)
		return model.UserDTO{}, err
	}
	return mapUserToDTO(user), nil
}

func (s *UserService) FindUsersFromApplication(applicationName string) ([]model.UserDTO, error) {
	return []model.UserDTO{}, nil
}

func mapUserToDTO(user model.User) model.UserDTO {

	applicationDTOs := make([]model.ApplicationRoleDTO, 0, len(user.Applications))
	for _, application := range user.Applications {
		applicationDTOs = append(applicationDTOs, model.ApplicationRoleDTO{ApplicationName: application.ApplicationName, Roles: application.Roles})
	}
	return model.UserDTO{
		UserName:     user.UserName,
		Name:         user.Name,
		LastName:     user.LastName,
		Email:        user.Email,
		ID:           user.ID,
		Applications: applicationDTOs,
	}

}

// CreateUser requires an user with userName or eMail and password
func (s *UserService) CreateUser(userDTO model.UserDTO) (err error) {

	if userDTO.ID != 0 {
		return errors.New("Cannot create a user with existing Id")
	}

	if (userDTO.UserName == "" && userDTO.Email == "") || userDTO.Password == "" {
		return errors.New("A user must have minimum username or email and a Password")
	}

	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return errors.New("could not create password Hash")
	}

	if _, err := s.databaseHandler.FindByUserName(userDTO.UserName); !gorm.IsRecordNotFoundError(err) {
		return errors.New("Username already exists")
	}

	if _, err := s.databaseHandler.FindByEmail(userDTO.Email); !gorm.IsRecordNotFoundError(err) {
		return errors.New("Email already exists")
	}

	applications := mapApplicationDTOToEntity(userDTO.Applications)

	user := model.User{
		UserName:     userDTO.UserName,
		Email:        userDTO.Email,
		LastName:     userDTO.LastName,
		Name:         userDTO.Name,
		Password:     bcryptedPassword,
		Applications: applications,
	}

	err = s.databaseHandler.CreateUser(user)
	return
}

func (s *UserService) UpdateUser(userID uint, userDTO model.UserDTO) error {
	user, err := s.databaseHandler.FindByID(userID)
	if err != nil {
		return err
	}
	if userDTO.UserName == "-" {
		return errors.New("username cannot be deleted")
	}
	if userDTO.UserName != "" && userDTO.UserName != user.UserName {
		if _, err := s.databaseHandler.FindByUserName(userDTO.UserName); !gorm.IsRecordNotFoundError(err) {
			return errors.New("Username cannot be changed as there is already a user with this username")
		} else {
			user.UserName = userDTO.UserName
		}
	}
	if userDTO.Email != "" {
		if userDTO.Email == "-" {
			user.Email = ""
		} else {
			user.Email = userDTO.Email
		}
	}
	if userDTO.LastName != "" {
		if userDTO.LastName == "-" {
			user.LastName = ""
		} else {
			user.LastName = userDTO.LastName
		}
	}
	if userDTO.Name != "" {
		if userDTO.Name == "-" {
			user.Name = ""
		} else {
			user.Name = userDTO.Name
		}
	}
	if userDTO.ClearApplications {
		user.Applications = make([]model.Application, 0)
	} else if len(userDTO.Applications) != 0 {
		user.Applications = mapApplicationDTOToEntity(userDTO.Applications)
	}

	return s.databaseHandler.UpdateUser(&user)

}

func mapApplicationDTOToEntity(applicationDTO []model.ApplicationRoleDTO) (applications []model.Application) {
	applications = make([]model.Application, 0, len(applicationDTO))

	for _, applicationDTO := range applicationDTO {
		application := model.Application{
			ApplicationName: applicationDTO.ApplicationName,
			Roles:           applicationDTO.Roles,
		}
		applications = append(applications, application)
	}
	return
}

// CheckPassword returns true if password matches
func (s *UserService) CheckPassword(userName, password string) (bool, error) {

	var user model.User
	var err error
	user, err = s.databaseHandler.FindByEmailOrUserName(userName)
	if err != nil {
		return false, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, err
	}
	return true, nil

}
