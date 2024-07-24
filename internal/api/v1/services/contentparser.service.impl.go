package services

type ContentParserServiceImpl struct {
}

func NewContentParserServiceImpl() *ContentParserServiceImpl {
	return &ContentParserServiceImpl{}
}

func (cps *ContentParserServiceImpl) GetMediaLinks(content string) []string {
	return []string{} // TODO: implement
}
