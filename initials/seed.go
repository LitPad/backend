package initials

import (
	"github.com/LitPad/backend/models"
	"github.com/LitPad/backend/models/choices"
)

var (
	TAGS []string = []string{
		"Campus", "Revenge", "Second Chance", "Age gap", "Regret", "Misunderstanding",
		"Hatred", "Kiss", "Runaway", "Hate to love", "Strong female lead", "School",
		"Forgiveness", "Royal", "Arranged marriage", "Betrayal", "Secret crush", "Pregnant",
		"Character growth", "He", "She", "Destiny", "Painful love", "Rejected", "Enemies",
		"Prophecy", "Shifter", "Omega verse", "Alternate universe", "Controlling", "World domination",
		"Karma", "Beautiful female lead", "Hardship", "Hard-working protagonist", "Future",	
	}

	GENRES []string = []string{
		"Werewolf", "Romance", "Billonaire", "Fantasy", "Mafia",
		"Historical", "YA/Teen", "Paranormal", "Urban", "Sci-fi", "Chicklit",
	}

	GIFTNAMES []string = []string{"Red rose", "Black dahlia", "Scroll", "Magic wand", "Wolf", "Baby Dragon"}

	BOOK = models.Book{
		Title:  "The Mysterious Island",
		Slug:   "the-mysterious-island",
		Blurb:  "An adventure novel about survival and discovery on an uncharted island.",
		AgeDiscretion: choices.ATYPE_EIGHTEEN,
		CoverImage: "https://example.com/cover.jpg",
		Completed:  false,
		Views:      "127.0.0.1 169.254.1.1 192.168.1.1 10.0.0.1 172.16.0.1 192.0.2.1 1.1.1.1 8.8.8.8",
	}
)


