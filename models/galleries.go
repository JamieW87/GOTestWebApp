package models

import "github.com/jinzhu/gorm"

//Gallery info for the DB
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

const (
	//ErrUserIDRequired displayed when a user hasnt entered their ID
	ErrUserIDRequired modelError = "models: user ID is required"
	//ErrTitleRequired shown when a person hasnt entered a gallery title
	ErrTitleRequired modelError = "models: title is required"
)

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

type GalleryService interface {
	GalleryDB
}

//Implements the gallery service interface
type galleryService struct {
	GalleryDB
}

//GalleryDB used to interact with the DB
type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
}

//Handles interacting with our database
type galleryGorm struct {
	db *gorm.DB
}

//Handles the validations
type galleryValidator struct {
	GalleryDB
}

type galleryValFn func(*Gallery) error

//GalleryDB ensure that gallery gorm always implements the GalleryDB interface
var _ GalleryDB = &galleryGorm{}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

//ByID accepts an integer and returns the requested gallery from the DB
func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

//
//
//
//
//VALIDATION FUNCTIONS
//
//
//
//

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}
