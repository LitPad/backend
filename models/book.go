package models

import (
	"time"

	"github.com/LitPad/backend/models/choices"
	"github.com/LitPad/backend/utils"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Tag struct {
	BaseModel
	Name   string  `gorm:"unique"`
	Slug   string  `gorm:"unique"`
	Genres []Genre `gorm:"many2many:genre_tags;"`
}

func (tag *Tag) BeforeSave(tx *gorm.DB) (err error) {
	tag.Slug = slug.Make(tag.Name)
	return
}

type Genre struct {
	BaseModel
	Name string `gorm:"unique"`
	Slug string `gorm:"unique"`
	Tags []Tag  `gorm:"many2many:genre_tags;"`
}

func (genre *Genre) BeforeSave(tx *gorm.DB) (err error) {
	genre.Slug = slug.Make(genre.Name)
	return
}

type Section struct {
	BaseModel
	Name        string `gorm:"unique"`
	Slug        string `gorm:"unique"`
	SubSections []SubSection
}

func (section *Section) BeforeSave(tx *gorm.DB) (err error) {
	section.Slug = slug.Make(section.Name)
	return
}

type SubSection struct {
	BaseModel
	Name      string `gorm:"unique"`
	Slug      string `gorm:"unique"`
	Books     []Book
	SectionID uuid.UUID
	Section   Section `gorm:"foreignKey:SectionID;constraint:OnDelete:SET NULL;<-:false"`
}

func (subSection *SubSection) BeforeSave(tx *gorm.DB) (err error) {
	subSection.Slug = slug.Make(subSection.Name)
	return
}

type Book struct {
	BaseModel
	AuthorID      uuid.UUID
	Author        User   `gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL;<-:false"`
	Title         string `gorm:"type: varchar(255)"`
	Slug          string `gorm:"unique"`
	Blurb         string `gorm:"type: varchar(1000)"`
	AgeDiscretion choices.AgeType

	GenreID uuid.UUID
	Genre   Genre `gorm:"foreignKey:GenreID;constraint:OnDelete:SET NULL;<-:false"`

	SubSectionID   uuid.UUID
	SubSection     SubSection `gorm:"foreignKey:SubSectionID;constraint:OnDelete:SET NULL;<-:false"`
	OrderInSection uint

	Tags       []Tag     `gorm:"many2many:book_tags"`
	Chapters   []Chapter `gorm:"<-:false"`
	CoverImage string    `gorm:"type:varchar(10000)"`

	Completed bool      `gorm:"default:false"`
	Reviews   []Comment `gorm:"<-:false"`
	Votes     []Vote    `gorm:"<-:false"`

	Featured       bool `gorm:"default:false"` //controlled by admin
	WeeklyFeatured time.Time
	Reads          []BookRead
	AvgRating      float64 // meant for query purposes. do not intentionally populate field
	Bookmark       []Bookmark

	// BOOK CONTRACT
	FullName             string `gorm:"type: varchar(1000)"`
	Email                string
	PenName              string `gorm:"type: varchar(1000)"`
	Age                  uint
	Country              string `gorm:"type: varchar(1000)"`
	Address              string `gorm:"type: varchar(1000)"`
	City                 string `gorm:"type: varchar(1000)"`
	State                string `gorm:"type: varchar(1000)"`
	PostalCode           string
	TelephoneNumber      string                       `gorm:"type: varchar(20)"`
	IDType               choices.ContractIDTypeChoice `gorm:"type: varchar(100)"`
	IDFrontImage         string
	IDBackImage          string
	BookAvailabilityLink *string
	PlannedLength        uint
	AverageChapter       uint
	UpdateRate           uint
	IntendedContract     choices.ContractTypeChoice
	FullPrice            *int
	ChapterPrice         int
	FullPurchaseMode     bool                         `gorm:"default:false"`
	ContractStatus       choices.ContractStatusChoice `gorm:"default:PENDING"`
}

func (b Book) VotesCount() int {
	return len(b.Votes)
}

func (b Book) ReadsCount() int {
	return len(b.Reads)
}

func (b Book) ChaptersCount() int {
	return len(b.Chapters)
}

func (b *Book) GenerateUniqueSlug(tx *gorm.DB) string {
	uniqueSlug := slug.Make(b.Title)
	slug := b.Slug
	if slug != "" {
		uniqueSlug = slug
	}

	existingBook := Book{Slug: uniqueSlug}
	tx.Take(&existingBook, existingBook)
	if existingBook.ID != uuid.Nil && existingBook.ID != b.ID { // slug is already taken
		// Make it unique by attaching a random string
		// to it and repeat the function
		randomStr := utils.GetRandomString(6)
		b.Slug = uniqueSlug + "-" + randomStr
		return b.GenerateUniqueSlug(tx)
	}
	return uniqueSlug
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	slug := b.GenerateUniqueSlug(tx)
	b.Slug = slug
	return
}

type BookRead struct {
	BaseModel
	UserID    uuid.UUID
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`
	BookID    uuid.UUID `json:"book_id"`
	Book      Book      `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
	Completed bool      `gorm:"default:false"`
}

type BookReport struct {
	BaseModel
	UserID                uuid.UUID
	User                  User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`
	BookID                uuid.UUID `json:"book_id"`
	Book                  Book      `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
	Reason                string    `gorm:"type: varchar(1000)"`
	AdditionalExplanation *string   `gorm:"type: varchar(1000)"`
}

type Chapter struct {
	BaseModel
	BookID     uuid.UUID
	Book       Book        `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
	Title      string      `gorm:"type: varchar(255)"`
	Slug       string      `gorm:"unique"`
	Paragraphs []Paragraph `gorm:"<-:false"`
}

func (c *Chapter) GenerateUniqueSlug(tx *gorm.DB) string {
	uniqueSlug := slug.Make(c.Title)
	slug := c.Slug
	if slug != "" {
		uniqueSlug = slug
	}

	existingChapter := Chapter{Slug: uniqueSlug}
	tx.Take(&existingChapter, existingChapter)
	if existingChapter.ID != uuid.Nil && existingChapter.ID != c.ID { // slug is already taken
		// Make it unique by attaching a random string
		// to it and repeat the function
		randomStr := utils.GetRandomString(6)
		c.Slug = uniqueSlug + "-" + randomStr
		return c.GenerateUniqueSlug(tx)
	}
	return uniqueSlug
}

func (c *Chapter) BeforeSave(tx *gorm.DB) (err error) {
	c.Slug = c.GenerateUniqueSlug(tx)
	return
}

type Paragraph struct {
	BaseModel
	ChapterID uuid.UUID
	Chapter   Chapter `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE;<-:false"`
	Index     uint
	Text      string    `gorm:"type:text"`
	Comments  []Comment `gorm:"foreignKey:ParagraphID;<-:false"`
}

func (p Paragraph) CommentsCount() int {
	return len(p.Comments)
}

type Comment struct {
	BaseModel
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	BookID *uuid.UUID // For reviews
	Book   *Book      `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
	Rating choices.RatingChoice

	ParagraphID *uuid.UUID // For praragrapj
	Paragraph   *Paragraph `gorm:"foreignKey:ParagraphID;constraint:OnDelete:CASCADE;<-:false"`
	Likes       []Like
	Text        string `gorm:"type:varchar(10000)"`

	ParentID *uuid.UUID

	Replies []Comment `gorm:"foreignKey:ParentID"`
}

func (c Comment) LikesCount() int {
	return len(c.Likes)
}

func (c Comment) RepliesCount() int {
	return len(c.Replies)
}

type Like struct {
	BaseModel
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	CommentID *uuid.UUID
	Comment   *Comment `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE;<-:false"`
}

type Vote struct {
	BaseModel
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	BookID uuid.UUID
	Book   Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
}

type Bookmark struct {
	BaseModel
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;<-:false"`

	BookID uuid.UUID
	Book   Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
}

type FeaturedContent struct {
	BaseModel
	Location choices.FeaturedContentLocationChoice
	Desc     string
	BookID   uuid.UUID
	Book     Book    `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;<-:false"`
	SeenBy   []*User `gorm:"many2many:featured_content_seen_by;"`
	IsActive bool	`gorm:"default:true;"`
}
