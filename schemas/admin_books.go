package schemas

import (
	"time"

	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
	"github.com/google/uuid"
)

type ContractSchema2 struct {
	FullName             string                       `json:"full_name"`
	Email                string                       `json:"email"`
	PenName              string                       `json:"pen_name"`
	Age                  uint                         `json:"age"`
	Country              string                       `json:"country"`
	Address              string                       `json:"address"`
	City                 string                       `json:"city"`
	State                string                       `json:"state"`
	PostalCode           uint                         `json:"postal_code"`
	TelephoneNumber      string                       `json:"telephone_number"`
	IDType               choices.ContractIDTypeChoice `json:"id_type"`
	IDFrontImage         string                       `json:"id_front_image"`
	IDBackImage          string                       `json:"id_back_image"`
	BookAvailabilityLink *string                      `json:"book_availability_link"`
	PlannedLength        uint                         `json:"planned_length"`
	AverageChapter       uint                         `json:"average_chapter"`
	UpdateRate           uint                         `json:"update_rate"`
	IntendedContract     choices.ContractTypeChoice   `json:"intended_contract"`
	FullPrice            *int                         `json:"full_price"`
	ChapterPrice         int                          `json:"chapter_price"`
	FullPurchaseMode     bool                         `json:"full_purchase_mode"`
	ContractStatus       choices.ContractStatusChoice `json:"contract_status"`
}

type ContractsResponseDataSchema struct {
	PaginatedResponseDataSchema
	Items []ContractSchema `json:"contracts"`
}

func (c ContractsResponseDataSchema) Init(contracts []models.Book) ContractsResponseDataSchema {
	// Set Initial Data
	contractItems := []ContractSchema{}
	for _, contract := range contracts {
		contractItems = append(contractItems, ContractSchema{}.Init(contract))
	}
	c.Items = contractItems
	return c
}

type ContractsResponseSchema struct {
	ResponseSchema
	Data ContractsResponseDataSchema `json:"data"`
}

type GenreAddSchema struct {
	Name     string   `json:"name" validate:"required"`
	TagSlugs []string `json:"tag_slugs"`
}

type TagsAddSchema struct {
	Name string `json:"name" validate:"required"`
}

type BookWithStats struct {
	Slug       string  `json:"slug"`
	Title      string  `json:"title"`
	CoverImage string  `json:"cover_image"`
	AuthorName string  `json:"author_name"`
	AvgRating  float64 `json:"avg_rating"`
	VotesCount int     `json:"votes_count"`
	ReadsCount int     `json:"reads_count"`
	GenreName  string  `json:"genre_name"`
	GenreSlug  string  `json:"genre_slug"`
}

type SectionWithSubsectionsSchema struct {
	Name             string             `json:"name"`
	Slug             string             `json:"slug"`
	SubSections      []SubSectionSchema `json:"sub_sections"`
	SubSectionsCount int                `json:"sub_sections_count"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

func (s SectionWithSubsectionsSchema) Init(section models.Section) SectionWithSubsectionsSchema {
	s.Name = section.Name
	s.Slug = section.Slug
	s.CreatedAt = section.CreatedAt
	s.UpdatedAt = section.UpdatedAt
	subsections := []SubSectionSchema{}
	for _, item := range section.SubSections {
		subsections = append(subsections, SubSectionSchema{}.Init(&item))
	}
	s.SubSections = subsections
	s.SubSectionsCount = len(section.SubSections)
	return s
}

type SectionsWithSubSectionsSchema struct {
	ResponseSchema
	Data []SectionWithSubsectionsSchema `json:"data"`
}

func (s SectionsWithSubSectionsSchema) Init(sections []models.Section) SectionsWithSubSectionsSchema {
	sectionsData := []SectionWithSubsectionsSchema{}
	for _, item := range sections {
		sectionsData = append(sectionsData, SectionWithSubsectionsSchema{}.Init(item))
	}
	s.Data = sectionsData
	return s
}

type SubSectionBookSchema struct {
	OrderInSection uint           `json:"order_in_section"`
	Title          string         `json:"title"`
	Author         UserDataSchema `json:"author"`
}

type SubSectionBookResponseSchema struct {
	PaginatedResponseDataSchema
	Items []SubSectionBookSchema `json:"items"`
}

type SubSectionWithBooksSchema struct {
	SubSectionSchema
	Section string                       `json:"section"`
	Books   SubSectionBookResponseSchema `json:"books"`
}

func (s SubSectionWithBooksSchema) Init(subSection models.SubSection, books []models.Book, paginatedData PaginatedResponseDataSchema) SubSectionWithBooksSchema {
	s.SubSectionSchema = s.SubSectionSchema.Init(&subSection)
	s.Section = subSection.Section.Name
	bookItems := []SubSectionBookSchema{}
	for _, item := range books {
		bookItems = append(bookItems, SubSectionBookSchema{
			OrderInSection: item.OrderInSection,
			Title:          item.Title,
			Author:         UserDataSchema{}.Init(item.Author),
		})
	}
	s.Books = SubSectionBookResponseSchema{
		PaginatedResponseDataSchema: paginatedData,
		Items:                       bookItems,
	}
	return s
}

type SubSectionWithBooksResponseSchema struct {
	ResponseSchema
	Data SubSectionWithBooksSchema `json:"data"`
}

type FeaturedContentBookSchema struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	CoverImage string `json:"cover_image"`
	Blurb      string `json:"blurb"`
}

type FeaturedContentSchema struct {
	ID       uuid.UUID                             `json:"id"`
	Location choices.FeaturedContentLocationChoice `json:"location"`
	Desc     string                                `json:"desc"`
	Book     FeaturedContentBookSchema             `json:"book"`
	IsActive bool                                  `json:"is_active"`
}

func (f FeaturedContentSchema) Init(featuredContent models.FeaturedContent) FeaturedContentSchema {
	f.ID = featuredContent.ID
	f.Location = featuredContent.Location
	f.Desc = featuredContent.Desc
	book := featuredContent.Book
	f.Book = FeaturedContentBookSchema{
		Title: book.Title, Slug: book.Slug,
		CoverImage: book.CoverImage, Blurb: book.Blurb,
	}
	f.IsActive = featuredContent.IsActive
	return f
}

type FeaturedContentsResponseSchema struct {
	ResponseSchema
	Data []FeaturedContentSchema `json:"data"`
}

func (f FeaturedContentsResponseSchema) Init(featuredContents []models.FeaturedContent) FeaturedContentsResponseSchema {
	// Set Initial Data
	contents := []FeaturedContentSchema{}
	for _, content := range featuredContents {
		contents = append(contents, FeaturedContentSchema{}.Init(content))
	}
	f.Data = contents
	return f
}

type FeaturedContentResponseSchema struct {
	ResponseSchema
	Data FeaturedContentSchema `json:"data"`
}

type FeaturedContentEntrySchema struct {
	Location choices.FeaturedContentLocationChoice `json:"location" validate:"featured_content_location_choice_validator"`
	BookSlug string                                `json:"book_slug" validate:"required"`
	Desc     string                                `json:"desc" validate:"required"`
}

type BookCompletionStatusSchema struct {
	Completed bool `json:"completed"`
}

type BookCompletionStatusResponseSchema struct {
	ResponseSchema
	Data BookCompletionStatusSchema `json:"data"`
}
