package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//User model
//Defines the database table
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

//UserService - The abstraction layer
//Provided for methods for creating, querying and updating users
type UserService struct {
	db *gorm.DB
}

//NewUserService is our connection to the database
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

//Close the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// ErrNotFound is returned when a resource cannot be found
// in the database.
var (
	ErrNotFound = errors.New("models: resource not found")
)

//ByID is a method that will look for a user with a provided ID.
//If user is found, returns a nil error, if no user is found returns ErrNotFound (As above)
//If different error, returns an error with more info about what went wrong
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}
