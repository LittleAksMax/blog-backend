package services

import "regexp"

const (
	mediaTag = "\\!\\[([a-zA-Z0-9-_])\\]"
)

type ContentParserServiceImpl struct {
	re *regexp.Regexp // compiled RegEx expression for media tags in the content
}

func NewContentParserServiceImpl() *ContentParserServiceImpl {
	return &ContentParserServiceImpl{
		re: regexp.MustCompile(mediaTag),
	}
}

func (cps *ContentParserServiceImpl) GetMediaLinks(content string) []string {
	// https://gobyexample.com/regular-expressions
	media := cps.re.FindAllString(content, -1)

	if media == nil {
		return []string{}
	} else {
		return media
	}
}
