package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

//User model
//Defines the database table
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

//UserService - The abstraction layer
//Provided for methods for creating, querying and updating users
type UserService struct {
	db *gorm.DB
}

var (
	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")
	//value that is used as the password pepper
	userPwPepper = "secret-random-string"
)

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
//Returns the user if no errors
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//ByEmail is simlar to ByID just a different SQL statement.
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the
// users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
//Uses the bcrypt library to hash the password and store it in user.PasswordHash as the variable hashedBytes
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

// first will query using the provided gorm.DB and it will
// get the first item returned and place it into dst. If
// nothing is found returns error
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Update will update the provided user with all of the data
// in the provided user object.
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

//Delete will delete the user with the provided ID
//GORM automatically deletes all users if ID isnt provided, so need to create this var and create a method that
//returns an error if the ID isnt valid.
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}
