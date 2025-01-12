package schemas

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
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
	Synopsis             string                       `json:"synopsis"`
	Outline              string                       `json:"outline"`
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
	GenreName  string  `json:"genre_name"`
	GenreSlug  string  `json:"genre_slug"`
}
