package routes

import (
	"github.com/LitPad/backend/managers"
	"github.com/LitPad/backend/models"
)

var (
	truthy                 = true
	userManager            = managers.UserManager{Model: models.User{}}
	bookManager            = managers.BookManager{Model: models.Book{}}
	chapterManager         = managers.ChapterManager{}
	tagManager             = managers.TagManager{}
	genreManager           = managers.GenreManager{}
	reviewManager          = managers.ReviewManager{}
	voteManager            = managers.VoteManager{}
	commentManager         = managers.CommentManager{}
	notificationManager    = managers.NotificationManager{}
	bookmarkManager        = managers.BookmarkManager{}
	bookReportManager      = managers.BookReportManager{}
	likeManager            = managers.LikeManager{}
	featuredContentManager = managers.FeaturedContentManager{}
)
