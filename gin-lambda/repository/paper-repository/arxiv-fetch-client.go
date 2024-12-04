package repository

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ArxivFetchClient interface {
	FetchV2(ctx context.Context, request *ArxivFetchRequestV2) (*ArxivFetchResponseV2, error)
}

type arxivFetchClient struct{}

func NewArxivFetchClient() ArxivFetchClient {
	return &arxivFetchClient{}
}

type ArxivFetchRequestV2 struct {
	Categories []string
	MaxSize    int64
}

type ArxivFetchResponseV2 struct {
	Papers []ArxivPaperV2
}

type ArxivPaperV2 struct {
	ArticleID string
	Title     string
	Abstract  string
	Author    string
}

type ArxivResponseFeed struct {
	Entry []ArxivResponseEntry `xml:"entry"`
}

type ArxivResponseEntry struct {
	ID         string                  `xml:"id"`
	Title      string                  `xml:"title"`
	Published  string                  `xml:"published"`
	Summary    string                  `xml:"summary"`
	Authors    []ArxivResponseAuthor   `xml:"author"`
	Categories []ArxivResponseCategory `xml:"category"`
}

type ArxivResponseAuthor struct {
	Name string `xml:"name"`
}

type ArxivResponseCategory struct {
	Term string `xml:"term,attr"`
}

type arxivFetchUrlBuilder struct {
	searchQuery struct {
		categories []string
	}
	sortBy     string
	maxResults int
}

func (a *arxivFetchClient) FetchV2(ctx context.Context, request *ArxivFetchRequestV2) (*ArxivFetchResponseV2, error) {
	baseUrl := "http://export.arxiv.org/api/query"
	builder := NewQueryBuilder("lastUpdatedDate", int(request.MaxSize)).SetCategories(request.Categories)
	url := builder.Build(baseUrl)

	res, err := httpGetBody(url)
	if err != nil {
		return nil, err
	}
	var ResponseFeed ArxivResponseFeed
	if err := xml.Unmarshal(res, &ResponseFeed); err != nil {
		return nil, err
	}

	papers := make([]ArxivPaperV2, 0, len(ResponseFeed.Entry))
	for _, entry := range ResponseFeed.Entry {
		papers = append(papers, ArxivPaperV2{
			ArticleID: entry.ID,
			Title:     entry.Title,
			Abstract:  entry.Summary,
			Author:    entry.Authors[0].Name,
		})
	}
	return &ArxivFetchResponseV2{
		Papers: papers,
	}, nil
}

func httpGetBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return resBody, nil
}

func NewQueryBuilder(
	sortBy string,
	maxResults int,
) *arxivFetchUrlBuilder {
	return &arxivFetchUrlBuilder{
		sortBy:     sortBy,
		maxResults: maxResults,
	}
}

func (b *arxivFetchUrlBuilder) Build(baseUrl string) string {
	var serchQueryContent = ""
	if b.searchQuery.categories != nil && len(b.searchQuery.categories) > 0 {
		serchQueryContent = "search_query=" + fmt.Sprintf("cat:%s", strings.Join(b.searchQuery.categories, "+"))
	}
	url := fmt.Sprintf("%s?%s&sortBy=%s&max_results=%d", baseUrl, serchQueryContent, b.sortBy, b.maxResults)
	return url
}

func (b *arxivFetchUrlBuilder) SetCategories(categories []string) *arxivFetchUrlBuilder {
	b.searchQuery.categories = categories
	return b
}
