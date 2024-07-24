package services

type ContentParserService interface {
	GetMediaLinks(content string) []string
}
