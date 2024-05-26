package managers

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookManager struct {
	Model	models.Book
	ModelList	[]models.Book
}

func (b BookManager) GetLatest(db *gorm.DB, tagSlug string) ([]models.Book, *utils.ErrorResponse) {
	books := b.ModelList

	query := db.Model(&b.Model)
	if tagSlug != "" {
		tag := models.Tag{Slug: tagSlug}
		db.Take(&tag, tag)
		if tag.ID == uuid.Nil {
			errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Invalid book tag")
			return nil, &errData
		}
		query = query.Scopes(scopes.AuthorGenreTagBookScope).Where("books.id IN (?)", db.Table("book_tags").Select("book_id").Where("tag_id = ?", tag.ID))
	}
	query.Order("created_at DESC").Find(&books)
	return books, nil
}

type TagManager struct {
	Model	models.Tag
	ModelList	[]models.Tag
}

func (t TagManager) GetAll(db *gorm.DB) []models.Tag {
	tags := t.ModelList
	db.Find(&tags)
	return tags
}

type GenreManager struct {
	Model	models.Genre
	ModelList	[]models.Genre
}

func (t GenreManager) GetAll(db *gorm.DB) []models.Genre {
	genres := t.ModelList
	db.Preload("Tags").Find(&genres)
	return genres
}