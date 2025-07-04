package managers

import (
	"fmt"
	"log"
	"time"

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

func (b BookManager) GetLatest(
	db *gorm.DB,
	genreSlug string,
	sectionSlug string,
	subSectionSlug string,
	tagSlug string,
	title string,
	byRating bool,
	username string,
	nameContains string,
	featured bool,
	weeklyFeatured bool,
	trending bool,
	orderBySubSection bool,
) ([]models.Book, *utils.ErrorResponse) {
	books := b.ModelList
	joinedSubSections := false

	query := db.Model(&b.Model)

	// Genre filter
	if genreSlug != "" {
		genre := models.Genre{Slug: genreSlug}
		db.Take(&genre, genre)
		if genre.ID == uuid.Nil {
			err := utils.NotFoundErr("Invalid book genre")
			return books, &err
		}
		query = query.Where(models.Book{GenreID: genre.ID})
	}

	// Section filter
	if sectionSlug != "" {
		section := models.Section{Slug: sectionSlug}
		db.Take(&section, section)
		if section.ID == uuid.Nil {
			err := utils.NotFoundErr("Invalid book section")
			return books, &err
		}
		query = query.
			Joins("JOIN book_sub_sections ON book_sub_sections.book_id = books.id").
			Joins("JOIN sub_sections ON sub_sections.id = book_sub_sections.sub_section_id").
			Where("sub_sections.section_id = ?", section.ID)
		joinedSubSections = true
	}

	// SubSection filter
	if subSectionSlug != "" {
		subSection := models.SubSection{Slug: subSectionSlug}
		db.Take(&subSection, subSection)
		if subSection.ID == uuid.Nil {
			err := utils.NotFoundErr("Invalid book subsection")
			return books, &err
		}
		if !joinedSubSections {
			query = query.
				Joins("JOIN book_sub_sections ON book_sub_sections.book_id = books.id").
				Joins("JOIN sub_sections ON sub_sections.id = book_sub_sections.sub_section_id")
			joinedSubSections = true
		}
		query = query.Where("book_sub_sections.sub_section_id = ?", subSection.ID)
	}

	// Tag filter
	if tagSlug != "" {
		tag := models.Tag{Slug: tagSlug}
		db.Take(&tag, tag)
		if tag.ID == uuid.Nil {
			return books, nil
		}
		query = query.Where("books.id IN (?)", db.Table("book_tags").Select("book_id").Where("tag_id = ?", tag.ID))
	}

	// Title filter
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	// Author username filter
	if username != "" {
		author := models.User{Username: username, AccountType: choices.ACCTYPE_AUTHOR}
		db.Take(&author, author)
		if author.ID == uuid.Nil {
			err := utils.NotFoundErr("Invalid author username")
			return books, &err
		}
		query = query.Where(models.Book{AuthorID: author.ID})
	}

	// Search by nameContains
	if nameContains != "" {
		query = query.
			Joins("LEFT JOIN users ON users.id = books.author_id").
			Where("users.username ILIKE ? OR users.name ILIKE ?", "%"+nameContains+"%", "%"+nameContains+"%")
	}

	// Featured filters
	if featured {
		query = query.Where("featured = ?", true)
	}
	if weeklyFeatured {
		query = query.Where("weekly_featured >= ? AND weekly_featured <= ?",
			time.Now().Truncate(7*24*time.Hour), time.Now().Add(7*24*time.Hour))
	}

	// Always join comments for rating computation
	query = query.Joins("LEFT JOIN comments ON comments.book_id = books.id")

	// Dynamically build SELECT clause
	selectFields := "books.*, COALESCE(AVG(comments.rating), 0) AS avg_rating"

	// If ordering by subsection name
	if orderBySubSection {
		if !joinedSubSections {
			query = query.
				Joins("LEFT JOIN book_sub_sections ON book_sub_sections.book_id = books.id").
				Joins("LEFT JOIN sub_sections ON sub_sections.id = book_sub_sections.sub_section_id")
			joinedSubSections = true
		}
		selectFields += ", MIN(sub_sections.name) AS sort_name"
	}

	// Final SELECT and GROUP
	query = query.Select(selectFields).Group("books.id")

	// Ordering
	if orderBySubSection {
		query = query.Order("sort_name ASC")
	} else if trending {
		query = query.
			Joins("LEFT JOIN book_reads ON book_reads.book_id = books.id").
			Group("books.id").
			Order("COUNT(book_reads.id) DESC")
	} else if byRating {
		query = query.Order("avg_rating DESC")
	} else {
		query = query.Order("books.created_at DESC")
	}

	// Apply preload scope
	query.Scopes(scopes.AuthorGenreTagBookPreloadScope).Find(&books)
	return books, nil
}

func (b BookManager) GetUserBookmarkedBooks(db *gorm.DB, user models.User) []models.Book {
	books := b.ModelList
	db.Joins("JOIN bookmarks ON bookmarks.book_id = books.id").
		Where("bookmarks.user_id = ?", user.ID).
		Scopes(scopes.AuthorGenreTagBookPreloadScope).
		Find(&books)
	return books
}

func (b BookManager) GetBySlug(db *gorm.DB, slug string, preload bool) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{Slug: slug}
	query := db
	if preload {
		query = query.Scopes(scopes.AuthorGenreTagBookScope)
	}
	query.Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.NotFoundErr("No book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) GetBooksOrderedByRatingAndVotes(db *gorm.DB) []schemas.BookWithStats {
	var books []schemas.BookWithStats

	db.Model(&b.Model).
		Select(`
			books.slug, 
			books.title, 
			books.cover_image, 
			users.username AS author_name, 
			COALESCE(AVG(comments.rating), 0) AS avg_rating, 
			COUNT(votes.id) AS votes_count, 
			genres.name AS genre_name, 
			genres.slug AS genre_slug
		`).
		Joins("LEFT JOIN users ON users.id = books.author_id"). // Adjust `author_id` if necessary
		Joins("LEFT JOIN comments ON comments.book_id = books.id").
		Joins("LEFT JOIN votes ON votes.book_id = books.id").
		Joins("LEFT JOIN genres ON genres.id = books.genre_id"). // Adjust `genre_id` if necessary
		Group("books.slug, books.title, books.cover_image, users.username, genres.name, genres.slug").
		Order("avg_rating DESC, votes_count DESC").
		Limit(10).
		Scan(&books)
	return books
}

func (b BookManager) GetBookContracts(db *gorm.DB, name *string, contractStatus *choices.ContractStatusChoice) []models.Book {
	books := []models.Book{}
	q := db.Not("full_name = ?", "")
	if contractStatus != nil {
		q.Where(models.Book{ContractStatus: *contractStatus})
	}
	if name != nil {
		q.Where(models.Book{FullName: *name})
	}
	q.Find(&books)
	return books
}

func (b BookManager) GetContractedBookBySlug(db *gorm.DB, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{Slug: slug, ContractStatus: choices.CTS_APPROVED}
	db.Scopes(scopes.AuthorGenreTagBookScope).Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.NotFoundErr("No contract approved book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) GetBySlugWithReviews(db *gorm.DB, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{Slug: slug}
	db.Scopes(scopes.AuthorGenreTagReviewsBookScope).
		Select("books.*, AVG(comments.rating) as avg_rating").
		Joins("LEFT JOIN comments ON comments.book_id = books.id").
		Group("books.id").
		Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.NotFoundErr("No book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) GetYearlyReadingProgress(db *gorm.DB, bookID uuid.UUID) []schemas.BookReadingProgressSchema {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -11, 0)

	var reads []models.BookRead
	db.Where("book_id = ? AND created_at >= ?", bookID, start).Find(&reads)

	monthlyData := make([]schemas.BookReadingProgressSchema, 12)
	for i := 0; i < 12; i++ {
		monthlyData[i] = schemas.BookReadingProgressSchema{
			Label:          start.AddDate(0, i, 0).Format("Jan"),
			TotalReaders:   0,
			CompletionRate: 0.0,
			NewReaders:     0,
		}
	}

	for _, r := range reads {
		monthIndex := int((r.CreatedAt.Year()-start.Year())*12 + int(r.CreatedAt.Month()) - int(start.Month()))
		if monthIndex < 0 || monthIndex >= 12 {
			continue
		}

		bucket := &monthlyData[monthIndex]
		bucket.TotalReaders++

		if r.Completed {
			prevTotal := float64(bucket.TotalReaders - 1)
			bucket.CompletionRate = (bucket.CompletionRate*prevTotal + 1) / float64(bucket.TotalReaders)
		}

		if r.FirstRead {
			bucket.NewReaders++
		}
	}

	return monthlyData
}

func (b BookManager) Get30DayReadingProgress(db *gorm.DB, bookID uuid.UUID) []schemas.BookReadingProgressSchema {
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)

	var reads []models.BookRead
	db.Where("book_id = ? AND created_at BETWEEN ? AND ?", bookID, startDate, now).Find(&reads)

	weekStats := make([]schemas.BookReadingProgressSchema, 4)
	for i := 0; i < 4; i++ {
		weekStats[i] = schemas.BookReadingProgressSchema{
			Label:          fmt.Sprintf("Week %d", i+1),
			TotalReaders:   0,
			CompletionRate: 0.0,
			NewReaders:     0,
		}
	}

	for _, r := range reads {
		weekIndex := int(now.Sub(r.CreatedAt).Hours() / 24 / 7)
		if weekIndex > 3 {
			weekIndex = 3
		}
		bucket := &weekStats[3-weekIndex]

		bucket.TotalReaders++

		if r.Completed {
			prevTotal := float64(bucket.TotalReaders - 1)
			bucket.CompletionRate = (bucket.CompletionRate*prevTotal + 1) / float64(bucket.TotalReaders)
		}

		if r.FirstRead {
			bucket.NewReaders++
		}
	}

	return weekStats
}

func (b BookManager) Get7DayReadingProgress(db *gorm.DB, bookID uuid.UUID) []schemas.BookReadingProgressSchema {
	now := time.Now()
	startDate := now.AddDate(0, 0, -6)

	var reads []models.BookRead
	db.Where("book_id = ? AND created_at BETWEEN ? AND ?", bookID, startDate, now).Find(&reads)

	dailyStats := make([]schemas.BookReadingProgressSchema, 7)
	for i := 0; i < 7; i++ {
		day := now.AddDate(0, 0, -6+i).Weekday().String()
		dailyStats[i] = schemas.BookReadingProgressSchema{
			Label:          day,
			TotalReaders:   0,
			CompletionRate: 0.0,
			NewReaders:     0,
		}
	}

	for _, r := range reads {
		dayIndex := int(now.Sub(r.CreatedAt).Hours() / 24)
		if dayIndex > 6 {
			dayIndex = 6
		}
		bucket := &dailyStats[6-dayIndex]

		bucket.TotalReaders++

		if r.Completed {
			prevTotal := float64(bucket.TotalReaders - 1)
			bucket.CompletionRate = (bucket.CompletionRate*prevTotal + 1) / float64(bucket.TotalReaders)
		}

		if r.FirstRead {
			bucket.NewReaders++
		}
	}
	return dailyStats
}

func (b BookManager) GetReaderRetentionPieData(db *gorm.DB, bookID uuid.UUID) schemas.RetentionStatsSchema {
	var reads []models.BookRead
	db.Where("book_id = ?", bookID).Find(&reads)

	stats := schemas.RetentionStatsSchema{}
	for _, r := range reads {
		switch {
		case r.Completed:
			stats.Completed++
		case r.InLibrary:
			stats.InProgress++
		default:
			stats.Dropped++
		}
	}
	stats.Total = stats.Completed + stats.InProgress + stats.Dropped


	if stats.Total != 0 {
		stats.Completed = float64(stats.Completed) / float64(stats.Total) * 100
		stats.InProgress = float64(stats.InProgress) / float64(stats.Total) * 100
		stats.Dropped = float64(stats.Dropped) / float64(stats.Total) * 100
	} else {
		stats.Completed = 0
		stats.InProgress = 0
		stats.Dropped = 0
	}
	return stats
}

func (b BookManager) GetByAuthorAndSlug(db *gorm.DB, author *models.User, slug string) (*models.Book, *utils.ErrorResponse) {
	book := models.Book{AuthorID: author.ID, Slug: slug}
	db.Scopes(scopes.AuthorGenreTagBookScope).Preload("Chapters").Take(&book, book)
	if book.ID == uuid.Nil {
		errD := utils.NotFoundErr("Author has no book with that slug")
		return nil, &errD
	}
	return &book, nil
}

func (b BookManager) Create(db *gorm.DB, author models.User, data schemas.BookCreateSchema, genre models.Genre, coverImage string, tags []models.Tag) models.Book {
	book := models.Book{
		AuthorID: author.ID, Author: author, Title: data.Title,
		Blurb: data.Blurb, AgeDiscretion: data.AgeDiscretion,
		GenreID: genre.ID,
		Tags:       tags,
		CoverImage: coverImage,
	}
	db.Omit("Tags.*").Create(&book)
	book.Tags = tags
	book.Genre = genre
	return book
}

func (b BookManager) Update(db *gorm.DB, book models.Book, data schemas.BookCreateSchema, genre models.Genre, coverImage string, tags []models.Tag) models.Book {
	book.Title = data.Title
	book.Blurb = data.Blurb
	book.AgeDiscretion = data.AgeDiscretion
	book.GenreID = genre.ID
	book.Genre = genre
	book.Tags = tags

	if coverImage != "" {
		book.CoverImage = coverImage
	}
	db.Omit("Tags.*").Save(&book)
	book.Tags = tags
	book.Genre = genre
	return book
}

func (b BookManager) DeleteBookWithAllRelations(db *gorm.DB, bookID uuid.UUID) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. Handle Book Reviews (Comments) and their nested relations
        var reviewIDs []uuid.UUID
        if err := tx.Model(&models.Comment{}).Where("book_id = ?", bookID).
            Pluck("id", &reviewIDs).Error; err != nil {
            return fmt.Errorf("failed to get review IDs: %w", err)
        }
        
        if len(reviewIDs) > 0 {
            // Delete likes for book reviews
            if err := tx.Where("comment_id IN ?", reviewIDs).Delete(&models.Like{}).Error; err != nil {
                return fmt.Errorf("failed to delete review likes: %w", err)
            }
            
            // Delete replies to book reviews (assuming replies reference parent comment)
            if err := tx.Where("parent_id IN ?", reviewIDs).Delete(&models.Comment{}).Error; err != nil {
                return fmt.Errorf("failed to delete review replies: %w", err)
            }
            
            // Delete book reviews
            if err := tx.Where("book_id = ?", bookID).Delete(&models.Comment{}).Error; err != nil {
                return fmt.Errorf("failed to delete reviews: %w", err)
            }
        }
        
        // 2. Delete Book Votes
        if err := tx.Where("book_id = ?", bookID).Delete(&models.Vote{}).Error; err != nil {
            return fmt.Errorf("failed to delete votes: %w", err)
        }
        
        // 3. Delete Book Reads
        if err := tx.Where("book_id = ?", bookID).Delete(&models.BookRead{}).Error; err != nil {
            return fmt.Errorf("failed to delete book reads: %w", err)
        }
        
        // 4. Delete Book Bookmarks
        if err := tx.Where("book_id = ?", bookID).Delete(&models.Bookmark{}).Error; err != nil {
            return fmt.Errorf("failed to delete bookmarks: %w", err)
        }
        
        // 5. Delete Many-to-Many Tags association
        if err := tx.Table("book_tags").Where("book_id = ?", bookID).Delete(&struct{}{}).Error; err != nil {
            return fmt.Errorf("failed to delete book tags: %w", err)
        }
        
        // 6. Handle Chapters and their nested relations
        var chapterIDs []uuid.UUID
        if err := tx.Model(&models.Chapter{}).Where("book_id = ?", bookID).
            Pluck("id", &chapterIDs).Error; err != nil {
            return fmt.Errorf("failed to get chapter IDs: %w", err)
        }
        
        if len(chapterIDs) > 0 {
            // Get all paragraph IDs for these chapters
            var paragraphIDs []uuid.UUID
            if err := tx.Model(&models.Paragraph{}).Where("chapter_id IN ?", chapterIDs).
                Pluck("id", &paragraphIDs).Error; err != nil {
                return fmt.Errorf("failed to get paragraph IDs: %w", err)
            }
            
            if len(paragraphIDs) > 0 {
                // Get all comment IDs for these paragraphs
                var commentIDs []uuid.UUID
                if err := tx.Model(&models.Comment{}).Where("paragraph_id IN ?", paragraphIDs).
                    Pluck("id", &commentIDs).Error; err != nil {
                    return fmt.Errorf("failed to get paragraph comment IDs: %w", err)
                }
                
                // Delete likes for paragraph comments
                if len(commentIDs) > 0 {
                    if err := tx.Where("comment_id IN ?", commentIDs).Delete(&models.Like{}).Error; err != nil {
                        return fmt.Errorf("failed to delete paragraph comment likes: %w", err)
                    }
                }
                
                // Delete paragraph comments
                if err := tx.Where("paragraph_id IN ?", paragraphIDs).Delete(&models.Comment{}).Error; err != nil {
                    return fmt.Errorf("failed to delete paragraph comments: %w", err)
                }
            }
            
            // Delete paragraphs
            if err := tx.Where("chapter_id IN ?", chapterIDs).Delete(&models.Paragraph{}).Error; err != nil {
                return fmt.Errorf("failed to delete paragraphs: %w", err)
            }
        }
        
        // Delete chapters
        if err := tx.Where("book_id = ?", bookID).Delete(&models.Chapter{}).Error; err != nil {
            return fmt.Errorf("failed to delete chapters: %w", err)
        }
        
        // 7. Finally, delete the book itself
        if err := tx.Delete(&models.Book{}, bookID).Error; err != nil {
            return fmt.Errorf("failed to delete book: %w", err)
        }
        
        return nil
    })
}

func (b BookManager) DeleteBookWithSQL(db *gorm.DB, bookID uuid.UUID) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Step 1: Delete all likes for comments related to this book (both reviews and paragraph comments)
        if err := tx.Exec(`
            DELETE FROM likes 
            WHERE comment_id IN (
                WITH RECURSIVE all_book_comments AS (
                    -- Base case: Direct book reviews
                    SELECT id FROM comments 
                    WHERE book_id = $1
                    
                    UNION ALL
                    
                    -- Base case: Comments on paragraphs in this book's chapters
                    SELECT id FROM comments 
                    WHERE paragraph_id IN (
                        SELECT p.id FROM paragraphs p
                        INNER JOIN chapters c ON p.chapter_id = c.id
                        WHERE c.book_id = $1
                    )
                    
                    UNION ALL
                    
                    -- Recursive case: Replies to any of the above comments
                    SELECT c.id FROM comments c
                    INNER JOIN all_book_comments abc ON c.parent_id = abc.id
                )
                SELECT id FROM all_book_comments
            )
        `, bookID).Error; err != nil {
            return fmt.Errorf("failed to delete likes: %w", err)
        }

        // Step 2: Delete all comments (reviews + paragraph comments + nested replies) for this book
        if err := tx.Exec(`
            DELETE FROM comments 
            WHERE id IN (
                WITH RECURSIVE all_book_comments AS (
                    -- Base case: Direct book reviews
                    SELECT id FROM comments 
                    WHERE book_id = $1
                    
                    UNION ALL
                    
                    -- Base case: Comments on paragraphs in this book's chapters
                    SELECT id FROM comments 
                    WHERE paragraph_id IN (
                        SELECT p.id FROM paragraphs p
                        INNER JOIN chapters c ON p.chapter_id = c.id
                        WHERE c.book_id = $1
                    )
                    
                    UNION ALL
                    
                    -- Recursive case: Replies to any of the above comments
                    SELECT c.id FROM comments c
                    INNER JOIN all_book_comments abc ON c.parent_id = abc.id
                )
                SELECT id FROM all_book_comments
            )
        `, bookID).Error; err != nil {
            return fmt.Errorf("failed to delete comments: %w", err)
        }

        // Step 3: Delete book votes
        if err := tx.Exec("DELETE FROM votes WHERE book_id = $1", bookID).Error; err != nil {
            return fmt.Errorf("failed to delete votes: %w", err)
        }

        // Step 4: Delete book reads
        if err := tx.Exec("DELETE FROM book_reads WHERE book_id = $1", bookID).Error; err != nil {
            return fmt.Errorf("failed to delete book reads: %w", err)
        }

        // Step 5: Delete book bookmarks
        if err := tx.Exec("DELETE FROM bookmarks WHERE book_id = $1", bookID).Error; err != nil {
            return fmt.Errorf("failed to delete bookmarks: %w", err)
        }

        // Step 6: Delete book tags (many-to-many relationship)
        if err := tx.Exec("DELETE FROM book_tags WHERE book_id = $1", bookID).Error; err != nil {
            return fmt.Errorf("failed to delete book tags: %w", err)
        }

        // Step 7: Delete paragraphs in all chapters of this book
        if err := tx.Exec(`
            DELETE FROM paragraphs 
            WHERE chapter_id IN (
                SELECT id FROM chapters WHERE book_id = $1
            )
        `, bookID).Error; err != nil {
            return fmt.Errorf("failed to delete paragraphs: %w", err)
        }

        // Step 8: Delete chapters
        if err := tx.Exec("DELETE FROM chapters WHERE book_id = $1", bookID).Error; err != nil {
            return fmt.Errorf("failed to delete chapters: %w", err)
        }

        // Step 9: Finally, delete the book itself
        if err := tx.Exec("DELETE FROM books WHERE id = $1", bookID).Error; err != nil {
            return fmt.Errorf("failed to delete book: %w", err)
        }

        return nil
    })
}

func (b BookManager) SetContract(db *gorm.DB, book models.Book, idFrontImage string, idBackImage string, data schemas.ContractCreateSchema) models.Book {
	book.FullName = data.FullName
	book.Email = data.Email
	book.PenName = data.PenName
	book.Age = data.Age
	book.Country = data.Country
	book.Address = data.Address
	book.City = data.City
	book.State = data.State
	book.PostalCode = data.PostalCode
	book.TelephoneNumber = data.TelephoneNumber
	book.IDType = data.IDType
	book.BookAvailabilityLink = data.BookAvailabilityLink
	book.PlannedLength = data.PlannedLength
	book.AverageChapter = data.AverageChapter
	book.UpdateRate = data.UpdateRate
	book.IntendedContract = data.IntendedContract
	book.FullPurchaseMode = data.FullPurchaseMode
	if idFrontImage != "" {
		book.IDFrontImage = idFrontImage
	}
	if idBackImage != "" {
		book.IDBackImage = idBackImage
	}

	if book.ContractStatus == choices.CTS_DECLINED {
		book.ContractStatus = choices.CTS_UPDATED
	}
	db.Omit("Tags.*").Save(&book)
	return book
}

type ChapterManager struct {
	Model     models.Chapter
	ModelList []models.Chapter
}

func (c ChapterManager) GetBySlug(db *gorm.DB, slug string) (*models.Chapter, *utils.ErrorResponse) {
	chapter := models.Chapter{Slug: slug}
	db.Joins("Book").Preload("Paragraphs.Comments").Take(&chapter, chapter)
	if chapter.ID == uuid.Nil {
		errD := utils.NotFoundErr("No chapter with that slug")
		return nil, &errD
	}
	for _, p := range chapter.Paragraphs {
		log.Println(p.Comments)
	}
	return &chapter, nil
}

func (c ChapterManager) GetBySlugWithComments(db *gorm.DB, slug string, index uint) (*models.Chapter, []models.Comment, *utils.ErrorResponse) {
	chapter := models.Chapter{Slug: slug}
	db.Take(&chapter, chapter)
	if chapter.ID == uuid.Nil {
		errD := utils.NotFoundErr("No chapter with that slug")
		return nil, nil, &errD
	}
	paragraph := models.Paragraph{ChapterID: chapter.ID, Index: index}
	comments := []models.Comment{}
	db.Take(&paragraph, paragraph)
	if paragraph.ID == uuid.Nil {
		errD := utils.NotFoundErr("That chapter has no paragraph with that index")
		return nil, nil, &errD
	}
	db.Joins("User").Preload("Replies").Preload("Likes").Where("comments.paragraph_id = ?", paragraph.ID).Find(&comments)
	return &chapter, comments, nil
}

func (c ChapterManager) IsFirstChapter(db *gorm.DB, chapter models.Chapter) bool {
	firstChapter := c.Model
	db.Order("created_at ASC").First(&firstChapter)
	return firstChapter.ID == chapter.ID
}

func (c ChapterManager) Create(db *gorm.DB, book models.Book, data schemas.ChapterCreateSchema) models.Chapter {
	chapter := models.Chapter{
		BookID: book.ID,
		Title:  data.Title,
		IsLast: data.IsLast,
	}
	db.Create(&chapter)
	// Generate paragraphs
	paragraphsToCreate := []models.Paragraph{}
	for idx, paragraph := range data.Paragraphs {
		paragraphsToCreate = append(paragraphsToCreate, models.Paragraph{ChapterID: chapter.ID, Text: paragraph, Index: uint(idx + 1)})
	}
	db.Create(&paragraphsToCreate)
	chapter.Paragraphs = paragraphsToCreate
	return chapter
}

func (c ChapterManager) Update(db *gorm.DB, chapter models.Chapter, data schemas.ChapterCreateSchema) models.Chapter {
	// Update title only if changed
	if chapter.Title != data.Title {
		chapter.Title = data.Title
	}
	chapter.IsLast = data.IsLast
	db.Save(&chapter)

	existingParagraphs := chapter.Paragraphs

	existingMap := make(map[uint]models.Paragraph) // Existing paragraphs indexed by index
	for _, p := range existingParagraphs {
		existingMap[p.Index] = p
	}

	toInsert := []models.Paragraph{}
	toUpdate := []models.Paragraph{}
	toDelete := []uuid.UUID{} // Store IDs for deletion

	existingIndexes := make(map[uint]bool)
	for i, text := range data.Paragraphs {
		index := uint(i)

		if existingPara, exists := existingMap[index]; exists {
			// Update only if text has changed
			if existingPara.Text != text {
				toUpdate = append(toUpdate, models.Paragraph{BaseModel: models.BaseModel{ID: existingPara.ID}, Text: text})
			}
			existingIndexes[index] = true
		} else {
			// Insert new paragraph
			toInsert = append(toInsert, models.Paragraph{ChapterID: chapter.ID, Index: index, Text: text})
		}
	}

	// Find paragraphs to delete (those not in `existingIndexes`)
	for _, p := range existingParagraphs {
		if !existingIndexes[p.Index] {
			toDelete = append(toDelete, p.ID)
		}
	}

	// Step 5: Execute Bulk Queries Efficiently
	tx := db.Begin()

	// Bulk Update (Uses GORM's Batch Update Feature)
	if len(toUpdate) > 0 {
		for _, p := range toUpdate {
			if err := tx.Model(&models.Paragraph{}).Where("id = ?", p.ID).Update("text", p.Text).Error; err != nil {
				tx.Rollback()
				fmt.Errorf("failed to update paragraphs: %w", err)
			}
		}
	}

	// Bulk Insert
	if len(toInsert) > 0 {
		if err := tx.Create(&toInsert).Error; err != nil {
			tx.Rollback()
			fmt.Errorf("failed to insert paragraphs: %w", err)
		}
	}

	// Bulk Delete
	if len(toDelete) > 0 {
		if err := tx.Where("id IN ?", toDelete).Delete(&models.Paragraph{}).Error; err != nil {
			tx.Rollback()
			fmt.Errorf("failed to delete paragraphs: %w", err)
		}
	}

	tx.Commit()
	return chapter
}

func (c ChapterManager) DeleteChapterWithAllRelations(db *gorm.DB, chapterID uuid.UUID) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Get all paragraph IDs for this chapter
        var paragraphIDs []uuid.UUID
        if err := tx.Model(&models.Paragraph{}).Where("chapter_id = ?", chapterID).
            Pluck("id", &paragraphIDs).Error; err != nil {
            return fmt.Errorf("failed to get paragraph IDs: %w", err)
        }
        
        if len(paragraphIDs) > 0 {
            // Get all comment IDs for these paragraphs
            var commentIDs []uuid.UUID
            if err := tx.Model(&models.Comment{}).Where("paragraph_id IN ?", paragraphIDs).
                Pluck("id", &commentIDs).Error; err != nil {
                return fmt.Errorf("failed to get comment IDs: %w", err)
            }
            
            // Delete likes for these comments
            if len(commentIDs) > 0 {
                if err := tx.Where("comment_id IN ?", commentIDs).Delete(&models.Like{}).Error; err != nil {
                    return fmt.Errorf("failed to delete likes: %w", err)
                }
            }
            
            // Delete comments
            if err := tx.Where("paragraph_id IN ?", paragraphIDs).Select(clause.Associations).Delete(&models.Comment{}).Error; err != nil {
                return fmt.Errorf("failed to delete comments: %w", err)
            }
        }
        
        // Delete paragraphs
        if err := tx.Where("chapter_id = ?", chapterID).Delete(&models.Paragraph{}).Error; err != nil {
            return fmt.Errorf("failed to delete paragraphs: %w", err)
        }
        
        // Delete chapter
        if err := tx.Delete(&models.Chapter{}, chapterID).Error; err != nil {
            return fmt.Errorf("failed to delete chapter: %w", err)
        }
        
        return nil
    })
}

func (c ChapterManager) DeleteChapterWithSQL(db *gorm.DB, chapterID uuid.UUID) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Step 1: Delete all likes for comments related to this chapter
        if err := tx.Exec(`
            DELETE FROM likes 
            WHERE comment_id IN (
                WITH RECURSIVE comment_tree AS (
                    -- Base case: Direct comments on paragraphs in this chapter
                    SELECT id FROM comments 
                    WHERE paragraph_id IN (
                        SELECT id FROM paragraphs WHERE chapter_id = $1
                    )
                    
                    UNION ALL
                    
                    -- Recursive case: Replies to comments
                    SELECT c.id FROM comments c
                    INNER JOIN comment_tree ct ON c.parent_id = ct.id
                )
                SELECT id FROM comment_tree
            )
        `, chapterID).Error; err != nil {
            return fmt.Errorf("failed to delete likes: %w", err)
        }

        // Step 2: Delete all comments (including nested replies) for this chapter
        if err := tx.Exec(`
            DELETE FROM comments 
            WHERE id IN (
                WITH RECURSIVE comment_tree AS (
                    -- Base case: Direct comments on paragraphs in this chapter
                    SELECT id FROM comments 
                    WHERE paragraph_id IN (
                        SELECT id FROM paragraphs WHERE chapter_id = $1
                    )
                    
                    UNION ALL
                    
                    -- Recursive case: Replies to comments
                    SELECT c.id FROM comments c
                    INNER JOIN comment_tree ct ON c.parent_id = ct.id
                )
                SELECT id FROM comment_tree
            )
        `, chapterID).Error; err != nil {
            return fmt.Errorf("failed to delete comments: %w", err)
        }

        // Step 3: Delete paragraphs
        if err := tx.Exec("DELETE FROM paragraphs WHERE chapter_id = $1", chapterID).Error; err != nil {
            return fmt.Errorf("failed to delete paragraphs: %w", err)
        }

        // Step 4: Delete chapter
        if err := tx.Exec("DELETE FROM chapters WHERE id = $1", chapterID).Error; err != nil {
            return fmt.Errorf("failed to delete chapter: %w", err)
        }

        return nil
    })
}

func (c ChapterManager) GetParagraph(db *gorm.DB, chapter models.Chapter, index uint) *models.Paragraph {
	paragraph := models.Paragraph{ChapterID: chapter.ID, Index: index}
	db.Take(&paragraph, paragraph)
	if paragraph.ID == uuid.Nil {
		return nil
	}
	return &paragraph
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

func (t TagManager) GetBySlug(db *gorm.DB, slug string) *models.Tag {

	tag := models.Tag{Slug: slug}
	db.Take(&tag, tag)

	if tag.ID == uuid.Nil {
		return nil
	}

	return &tag
}

type GenreManager struct {
	Model     models.Genre
	ModelList []models.Genre
}

func (g GenreManager) GetAll(db *gorm.DB) []models.Genre {
	genres := g.ModelList
	db.Preload("Tags").Preload("Tags.Books").Find(&genres)
	return genres
}

func (g GenreManager) GetAllSections(db *gorm.DB) []models.Section {
	sections := []models.Section{}
	db.Preload("SubSections").Preload("SubSections.Books").Find(&sections)
	return sections
}

func (g GenreManager) GetAllSubSections(db *gorm.DB, sectionID *uuid.UUID) []models.SubSection {
	subSections := []models.SubSection{}
	query := db.Preload("Books")
	if sectionID != nil {
		query = query.Where(models.SubSection{SectionID: *sectionID})
	}
	query.Find(&subSections)
	return subSections
}

func (g GenreManager) GetBySlug(db *gorm.DB, slug string) *models.Genre {

	genre := models.Genre{Slug: slug}
	db.Take(&genre, genre)

	if genre.ID == uuid.Nil {
		return nil
	}

	return &genre
}

func (g GenreManager) GetSectionBySlug(db *gorm.DB, slug string) *models.Section {

	section := models.Section{Slug: slug}
	db.Take(&section, section)

	if section.ID == uuid.Nil {
		return nil
	}

	return &section
}

func (g GenreManager) GetSubSectionBySlug(db *gorm.DB, slug string) (*models.SubSection, []schemas.BookWithOrder) {
	var subSection models.SubSection
	err := db.Preload("Section").Where("slug = ?", slug).First(&subSection).Error
	if err != nil || subSection.ID == uuid.Nil {
		return nil, []schemas.BookWithOrder{}
	}

	// Fetch books with pivot data
	var booksWithOrder []schemas.BookWithOrder
	err = db.Table("books").
		Select("books.*, book_sub_sections.order_in_section, users.id AS author_id, users.username AS author_username, users.name AS author_name, users.avatar AS author_avatar").
		Joins("JOIN book_sub_sections ON book_sub_sections.book_id = books.id").
		Joins("LEFT JOIN users ON users.id = books.author_id").
		Where("book_sub_sections.sub_section_id = ?", subSection.ID).
		Order("book_sub_sections.order_in_section ASC").
		Scan(&booksWithOrder).Error

	if err != nil {
		return nil, []schemas.BookWithOrder{}
	}

	return &subSection, booksWithOrder
}

type ReviewManager struct {
	Model     models.Comment
	ModelList []models.Comment
}

func (r ReviewManager) GetByID(db *gorm.DB, id uuid.UUID) *models.Comment {
	review := r.Model
	db.Where("comments.id = ? AND comments.parent_id IS NULL", id).Joins("User").Joins("Book").Preload("Replies.User").Preload("Replies.Likes").Take(&review, review)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) GetByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.Comment {
	review := models.Comment{}
	db.Where("user_id = ?", user.ID).Joins("Book").Joins("User").Preload("Replies").Preload("Likes").Take(&review, id)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) GetByUserAndBook(db *gorm.DB, user *models.User, book models.Book) *models.Comment {
	review := models.Comment{
		UserID: user.ID,
		BookID: &book.ID,
	}
	db.Take(&review, review)
	if review.ID == uuid.Nil {
		return nil
	}
	return &review
}

func (r ReviewManager) Create(db *gorm.DB, user *models.User, book models.Book, data schemas.ReviewBookSchema) models.Comment {
	review := models.Comment{
		UserID: user.ID,
		User:   *user,
		BookID: &book.ID,
		Book:   &book,
		Rating: data.Rating,
		Text:   data.Text,
	}
	db.Create(&review)
	return review
}

func (r ReviewManager) Update(db *gorm.DB, review models.Comment, data schemas.ReviewBookSchema) models.Comment {
	review.Text = data.Text
	review.Rating = data.Rating
	db.Save(&review)
	return review
}

type CommentManager struct {
	Model     models.Comment
	ModelList []models.Comment
}

func (c CommentManager) GetByID(db *gorm.DB, id uuid.UUID, preload bool) *models.Comment {
	comment := c.Model
	query := db.Where("comments.id = ?", id)
	if preload {
		query = query.Joins("User").Joins("Paragraph").Preload("Replies").Preload("Replies.User").Preload("Replies.Likes")
	}
	query.Take(&comment, comment)
	if comment.ID == uuid.Nil {
		return nil
	}
	return &comment
}

func (c CommentManager) GetByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.Comment {
	comment := c.Model
	db.Where("user_id = ?", user.ID).Joins("Chapter").Joins("User").Preload("Replies").Preload("Likes").Take(&comment, id)
	if comment.ID == uuid.Nil {
		return nil
	}
	return &comment
}

func (c CommentManager) GetByParagraphID(db *gorm.DB, paragraphId uuid.UUID) []models.Comment {
	comments := c.ModelList
	db.Where("paragraph_id = ?", paragraphId).Find(&comments)
	return comments
}

func (c CommentManager) Create(db *gorm.DB, user *models.User, paragraphID uuid.UUID, data schemas.ParagraphCommentAddSchema) models.Comment {
	comment := models.Comment{
		UserID:      user.ID,
		User:        *user,
		ParagraphID: &paragraphID,
		Text:        data.Text,
	}
	db.Create(&comment)
	return comment
}

func (c CommentManager) Update(db *gorm.DB, comment models.Comment, data schemas.ParagraphCommentAddSchema) models.Comment {
	comment.Text = data.Text
	db.Save(&comment)
	return comment
}

func (c CommentManager) GetReplyByUserAndID(db *gorm.DB, user *models.User, id uuid.UUID) *models.Comment {
	reply := c.Model
	db.Where("user_id = ?", user.ID).Joins("User").Preload("Likes").Take(&reply, id)
	if reply.ID == uuid.Nil {
		return nil
	}
	return &reply
}

func (c CommentManager) CreateReply(db *gorm.DB, user *models.User, reviewOrParagraphComment *models.Comment, data schemas.ReplyReviewOrCommentSchema) models.Comment {
	reply := models.Comment{
		UserID:   user.ID,
		User:     *user,
		Text:     data.Text,
		ParentID: &reviewOrParagraphComment.ID,
	}
	db.Create(&reply)
	return reply
}

func (c CommentManager) UpdateReply(db *gorm.DB, reply models.Comment, data schemas.ReplyEditSchema) models.Comment {
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

type BookmarkManager struct {
	Model     models.Bookmark
	ModelList []models.Bookmark
}

func (b BookmarkManager) AddOrDelete(db *gorm.DB, user models.User, book models.Book) string {
	bookmark := models.Bookmark{UserID: user.ID, BookID: book.ID}
	db.Take(&bookmark, bookmark)
	if bookmark.ID == uuid.Nil {
		db.Create(&bookmark)
		return "Bookmarked"
	}
	db.Delete(&bookmark)
	return "Unbookmarked"
}

type BookReportManager struct {
	Model     models.BookReport
	ModelList []models.BookReport
}

func (b BookReportManager) Create(db *gorm.DB, user models.User, book models.Book, reason string, additionalExplanation *string) {
	bookReport := models.BookReport{UserID: user.ID, BookID: book.ID, Reason: reason, AdditionalExplanation: additionalExplanation}
	db.Create(&bookReport)
}

type LikeManager struct {
	Model     models.Bookmark
	ModelList []models.Bookmark
}

func (l LikeManager) AddOrDelete(db *gorm.DB, user models.User, comment models.Comment) string {
	like := models.Like{UserID: user.ID, CommentID: &comment.ID}
	db.Take(&like, like)
	if like.ID == uuid.Nil {
		db.Create(&like)
		return "Liked"
	}
	db.Delete(&like)
	return "Unliked"
}

type FeaturedContentManager struct {
	Model     models.FeaturedContent
	ModelList []models.FeaturedContent
}

func (f FeaturedContentManager) GetAll(db *gorm.DB, location *choices.FeaturedContentLocationChoice, isActive *bool) []models.FeaturedContent {
	contents := f.ModelList
	query := db.Joins("Book")
	if location != nil {
		query = query.Where("location = ?", location)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", isActive)
	}
	query.Find(&contents)
	return contents
}

func (f FeaturedContentManager) Create(db *gorm.DB, location choices.FeaturedContentLocationChoice, desc string, book models.Book) models.FeaturedContent {
	content := models.FeaturedContent{Location: location, Desc: desc, BookID: book.ID}
	db.Create(&content)
	content.Book = book
	return content
}

func (f FeaturedContentManager) GetByID(db *gorm.DB, id uuid.UUID) *models.FeaturedContent {
	content := models.FeaturedContent{}
	db.Where("id = ?", id).Take(&content, content)
	if content.ID == uuid.Nil {
		return nil
	}
	return &content
}

func (f FeaturedContentManager) Update(db *gorm.DB, content models.FeaturedContent, location choices.FeaturedContentLocationChoice, desc string, book models.Book) models.FeaturedContent {
	content.Location = location
	content.Desc = desc
	content.BookID = book.ID
	db.Save(&content)
	content.Book = book
	return content
}
