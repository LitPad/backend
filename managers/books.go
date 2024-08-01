package managers

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/models/scopes"
	"github.com/LitPad/backend/schemas"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookManager struct {
	Model     models.Book
	ModelList []models.Book
}

func (b BookManager) GetLatest(db *gorm.DB, genreSlug string, tagSlug string, title string, usernameOpts ...string) ([]models.Book, *utils.ErrorResponse) {
	books := b.ModelList

	query := db.Model(&b.Model)
	if genreSlug != "" {
		genre := models.Genre{Slug: genreSlug}
		db.Take(&genre, genre)
		if genre.ID == uuid.Nil {
			errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Invalid book genre")
			return books, &errData
		}
		query = query.Where(models.Book{GenreID: genre.ID})
	}
	if tagSlug != "" {
		tag := models.Tag{Slug: tagSlug}
		db.Take(&tag, tag)
		if tag.ID == uuid.Nil {
			errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Invalid book tag")
			return books, &errData
		}
		query = query.Where("books.id IN (?)", db.Table("book_tags").Select("book_id").Where("tag_id = ?", tag.ID))
	}

	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if len(usernameOpts) > 0 {
		username := usernameOpts[0]
		author := models.User{Username: username, AccountType: choices.ACCTYPE_WRITER}
		db.Take(&author, author)
		if author.ID == uuid.Nil {
			errData := utils.RequestErr(utils.ERR_NON_EXISTENT, "Invalid author username")
			return books, &errData
		}
		query = query.Where(models.Book{AuthorID: author.ID})
	}
	query.Scopes(scopes.AuthorGenreTagBookScope).Order("created_at DESC").Find(&books)
	return books, nil
}

func (b BookManager) GetBySlug(db *gorm.DB, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{Slug: slug}
	db.Scopes(scopes.AuthorGenreTagBookScope).Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "No book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) GetBySlugWithReviews(db *gorm.DB, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{Slug: slug}
	db.Scopes(scopes.AuthorGenreTagReviewsBookScope).Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "No book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) GetByAuthorAndSlug(db *gorm.DB, author *models.User, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{AuthorID: author.ID, Slug: slug}
	db.Scopes(scopes.AuthorGenreTagBookScope).Preload("Chapters").Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "Writer has no book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) Create(db *gorm.DB, author models.User, data schemas.BookCreateSchema, genre models.Genre, coverImage string, Tags []models.Tag) models.Book {
	book := models.Book{
		AuthorID: author.ID, Author: author, Title: data.Title,
		Blurb: data.Blurb, AgeDiscretion: data.AgeDiscretion,
		GenreID: genre.ID, Genre: genre,
		Tags:       Tags,
		CoverImage: coverImage,
	}
	db.Omit("Tags.*").Create(&book)
	return book
}

func (b BookManager) Update(db *gorm.DB, book models.Book, data schemas.BookCreateSchema, genre models.Genre, Tags []models.Tag) models.Book {
	book.Title = data.Title
	book.Blurb = data.Blurb
	book.AgeDiscretion = data.AgeDiscretion
	book.GenreID = genre.ID
	book.Genre = genre
	book.Tags = Tags
	db.Omit("Tags.*").Save(&book)
	return book
}

type ChapterManager struct {
	Model     models.Chapter
	ModelList []models.Chapter
}

func (c ChapterManager) GetBySlug(db *gorm.DB, slug string) (*models.Chapter, *utils.ErrorResponse) {
	chapter := models.Chapter{Slug: slug}
	db.Joins("Book").Take(&chapter, chapter)
	if chapter.ID == uuid.Nil {
		errD := utils.RequestErr(utils.ERR_NON_EXISTENT, "No chapter with that slug")
		return nil, &errD
	}
	return &chapter, nil
}

func (c ChapterManager) Create(db *gorm.DB, book models.Book, data schemas.ChapterCreateSchema) models.Chapter {
	chapter := models.Chapter{
		BookID:        book.ID,
		Title:         data.Title,
		Text:          data.Text,
		ChapterStatus: data.ChapterStatus,
	}
	db.Create(&chapter)
	return chapter
}

func (c ChapterManager) Update(db *gorm.DB, chapter models.Chapter, data schemas.ChapterCreateSchema) models.Chapter {
	chapter.Title = data.Title
	chapter.Text = data.Text
	chapter.ChapterStatus = data.ChapterStatus
	db.Save(&chapter)
	return chapter
}

type TagManager struct {
	Model     models.Tag
	ModelList []models.Tag
}

func (t TagManager) GetAll(db *gorm.DB) []models.Tag {
	tags := t.ModelList
	db.Find(&tags)
	return tags
}

type GenreManager struct {
	Model     models.Genre
	ModelList []models.Genre
}

func (t GenreManager) GetAll(db *gorm.DB) []models.Genre {
	genres := t.ModelList
	db.Preload("Tags").Find(&genres)
	return genres
}

type BoughtChapterManager struct {
	Model     models.BoughtChapter
	ModelList []models.BoughtChapter
}

func (b BoughtChapterManager) GetBoughtChapters(db *gorm.DB, buyer *models.User, book *models.Book) []models.Chapter {
	boughtChapters := b.ModelList
	chapters := []models.Chapter{}
	db.Joins("JOIN chapters ON chapters.id = bought_chapters.chapter_id").Where("bought_chapters.buyer_id = ? AND chapters.book_id = ?", buyer.ID, book.ID).Scopes(scopes.BoughtChapterScope).Find(&boughtChapters)
	for i := range boughtChapters {
		chapters = append(chapters, boughtChapters[i].Chapter)
	}
	if len(chapters) == 0 {
		// If the user hasn't bought any, he should see the first chapter for free
		chapters = book.Chapters[:1]
	}
	return chapters
}

func (b BoughtChapterManager) GetBoughtBooks(db *gorm.DB, buyer *models.User) []models.Book {
	// Return books in which the user has bought at least a chapter
	books := []models.Book{}
	subQuery := db.Model(&BoughtChapterManager{}).Select("chapter_id").Where("user_id = ?", buyer.ID)

	// Main query to find all books that have at least one chapter in the subquery
	db.Model(&models.Book{}).
		Joins("JOIN chapters ON chapters.book_id = books.id").
		Where("chapters.id IN (?)", subQuery).
		Distinct("books.*"). // To ensure unique books
		Scopes(scopes.AuthorGenreTagReviewsBookScope).
		Find(&books)
	return books
}

func (b BoughtChapterManager) CheckAllChaptersBought(db *gorm.DB, buyer *models.User, book *models.Book) bool {
	var chapterIDs []uuid.UUID
	for _, chapter := range book.Chapters {
		chapterIDs = append(chapterIDs, chapter.ID)
	}

	var count int64
	db.Model(&models.BoughtChapter{}).
		Where("buyer_id = ? AND chapter_id IN ?", buyer.ID, chapterIDs).
		Group("buyer_id").
		Having("COUNT(DISTINCT chapter_id) = ?", len(chapterIDs)).
		Count(&count)

	return count > 0
}

func (b BoughtChapterManager) CheckIfAtLeastAChapterWasBought(db *gorm.DB, buyer *models.User, book models.Book) bool {
	var chapterIDs []uuid.UUID
	for _, chapter := range book.Chapters {
		chapterIDs = append(chapterIDs, chapter.ID)
	}

	var boughtChaptersCount int64
	db.Model(&models.BoughtChapter{}).Where("chapter_id IN ?", chapterIDs).Count(&boughtChaptersCount)
	return boughtChaptersCount > 0
}

func (b BoughtChapterManager) GetByBuyerAndChapter(db *gorm.DB, buyer *models.User, chapter models.Chapter) *models.BoughtChapter {
	boughtChapter := models.BoughtChapter{
		BuyerID:   buyer.ID,
		ChapterID: chapter.ID,
	}
	db.Joins("Chapter").Take(&boughtChapter, boughtChapter)
	if boughtChapter.ID == uuid.Nil {
		return nil
	}
	return &boughtChapter
}

func (b BoughtChapterManager) BuyAChapter(db *gorm.DB, buyer *models.User, book *models.Book) models.BoughtChapter {
	// Get the chapter that the user doesn't have yet
	nextChapter := models.Chapter{}
	subQuery := db.Model(&models.BoughtChapter{}).Select("chapter_id").Where("buyer_id = ?", buyer.ID)

	db.Model(&models.Chapter{}).
		Where("book_id = ? AND id NOT IN (?)", book.ID, subQuery).
		Order("created_at ASC"). // Assuming chapters are ordered by created_at descending order
		First(&nextChapter)

	secondChapter := models.Chapter{}
	bookChapters := book.Chapters
	if len(bookChapters) > 1 {
		firstChapter := bookChapters[0]
		if firstChapter.ID == nextChapter.ID {
			// Get second chapter
			secondChapter = bookChapters[1]
		}
	}

	boughtChapters := []models.BoughtChapter{{
		BuyerID:   buyer.ID,
		ChapterID: nextChapter.ID,
		Chapter:   nextChapter,
	}}
	if secondChapter.ID != uuid.Nil {
		// Add second chapter too
		secondBoughtChapter := models.BoughtChapter{
			BuyerID:   buyer.ID,
			ChapterID: secondChapter.ID,
			Chapter:   secondChapter,
		}
		boughtChapters = append(boughtChapters, secondBoughtChapter)
	}
	db.Create(&boughtChapters)
	chapterPrice := book.ChapterPrice

	// Move coins from buyer to author
	buyer.Coins = buyer.Coins - chapterPrice
	db.Save(&buyer)

	// Increase user coins
	author := book.Author
	author.Coins = author.Coins + chapterPrice
	db.Save(&author)
	return boughtChapters[len(boughtChapters)-1]
}

func (b BoughtChapterManager) BuyWholeBook(db *gorm.DB, buyer *models.User, book models.Book) models.Book {
	chaptersToBuy := []models.BoughtChapter{}
	for _, chapter := range book.Chapters {
		chapterToBuy := models.BoughtChapter{
			BuyerID:   buyer.ID,
			ChapterID: chapter.ID,
			Chapter:   chapter,
		}
		chaptersToBuy = append(chaptersToBuy, chapterToBuy)
	}

	db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&chaptersToBuy)

	bookPrice := book.FullPrice

	// Move coins from buyer to author
	buyer.Coins = buyer.Coins - *bookPrice
	db.Save(&buyer)

	// Increase user coins
	author := book.Author
	author.Coins = author.Coins + *bookPrice
	db.Save(&author)
	return book
}

type ReviewManager struct {
	Model     models.Review
	ModelList []models.Review
}

func (r ReviewManager) GetByID(db *gorm.DB, id uuid.UUID) *models.Review {
	review := r.Model
	db.Where("reviews.id = ?", id).Joins("User").Joins("Book").Preload("Replies").Preload("Replies.User").Preload("Replies.Likes").Take(&review, review)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) GetByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.Review {
	review := models.Review{}
	db.Where("user_id = ?", user.ID).Joins("Book").Joins("User").Preload("Replies").Preload("Likes").Take(&review, id)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) GetByUserAndBook(db *gorm.DB, user *models.User, book models.Book) *models.Review {
	review := models.Review{
		UserID: user.ID,
		BookID: book.ID,
	}
	db.Take(&review, review)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) Create(db *gorm.DB, user *models.User, book models.Book, data schemas.ReviewBookSchema) models.Review {
	review := models.Review{
		UserID: user.ID,
		User:   *user,
		BookID: book.ID,
		Book:   book,
		Rating: data.Rating,
		Text:   data.Text,
	}
	db.Create(&review)
	return review
}

func (r ReviewManager) Update(db *gorm.DB, review models.Review, data schemas.ReviewBookSchema) models.Review {
	review.Text = data.Text
	review.Rating = data.Rating
	db.Save(&review)
	return review
}

type ReplyManager struct {
	Model     models.Reply
	ModelList []models.Reply
}

func (r ReplyManager) GetByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.Reply {
	reply := models.Reply{}
	db.Where("user_id = ?", user.ID).Joins("User").Preload("Likes").Take(&reply, id)
	if reply.ID == uuid.Nil {
		return nil
	}
	return &reply
}

func (r ReplyManager) Create(db *gorm.DB, user *models.User, review *models.Review, data schemas.ReplyReviewSchema) models.Reply {
	reply := models.Reply{
		UserID:   user.ID,
		User:     *user,
		ReviewID: review.ID,
		Text:     data.Text,
	}
	db.Create(&reply)
	return reply
}

func (r ReplyManager) Update(db *gorm.DB, reply models.Reply, data schemas.ReplyReviewSchema) models.Reply {
	reply.Text = data.Text
	db.Save(&reply)
	return reply
}

type VoteManager struct {
	Model     models.Vote
	ModelList []models.Vote
}

func (v VoteManager) GetByUserAndBook(db *gorm.DB, user *models.User, book *models.Book) *models.Vote {
	vote := models.Vote{UserID: user.ID, BookID: book.ID}
	db.Joins("User").Joins("Book").Take(&vote, vote)
	if vote.ID == uuid.Nil {
		return nil
	}
	return &vote
}

func (v VoteManager) Create(db *gorm.DB, user *models.User, book *models.Book) models.Vote {
	vote := models.Vote{UserID: user.ID, User: *user, Book: *book, BookID: book.ID}
	db.Create(&vote)
	return vote
}
