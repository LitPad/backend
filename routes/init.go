package routes

import (
	"github.com/LitPad/backend/managers"
	"github.com/LitPad/backend/models"
)

var (
	truthy                  = true
	userManager             = managers.UserManager{Model: models.User{}}
	bookManager             = managers.BookManager{Model: models.Book{}}
	chapterManager          = managers.ChapterManager{}
	tagManager              = managers.TagManager{}
	genreManager            = managers.GenreManager{}
	reviewManager           = managers.ReviewManager{}
	replyManager            = managers.ReplyManager{}
	voteManager             = managers.VoteManager{}
	paragraphCommentManager = managers.ParagraphCommentManager{}
)
