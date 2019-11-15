package redisKeys

import "fmt"

const (
	Resource = "Resource:"
	Group    = "Group:"
	//最新动态
	GroupLatestNews = "GroupLatestNews:%d"
	Chapter         = "Chapter:"
	GroupTag        = "GroupTag:"
	GroupCategory   = "GroupCategory:"
	GroupChapter    = "GroupChapter:"
	ChapterResource = "ChapterResource:"
	Categories      = "Categories"
	IndexCtrl       = "IndexCtrl:"
	NewsCtrl        = "NewsCtrl:"
	CategoryCtrl    = "CategoryCtrl:"

	User = "UserInfo:"

	MessageCheckCode = "MessageCheckCode:"
)

func Key(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args)
}
