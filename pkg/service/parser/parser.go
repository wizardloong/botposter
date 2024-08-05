package parser

type Parser interface {
	FetchArticleData(string) (string, string, string, error)
	DownloadImage(string) ([]byte, error)
}
