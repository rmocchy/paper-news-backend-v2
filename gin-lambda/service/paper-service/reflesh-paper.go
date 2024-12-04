package paperservice

import (
	"context"
	"database/sql"
	"time"

	"github.com/rmocchy/paper-news-backend-v2/gin-lambda/gen/xo"
	repository "github.com/rmocchy/paper-news-backend-v2/gin-lambda/repository/paper-repository"
)

type RefleshPaperService interface {
	RefleshPaper(ctx context.Context, input *RefleshPaperInput) (*RefleshPaperOutput, error)
}

type refleshPaperService struct {
	paperRepository repository.PaperRepository
	arxivClient     repository.ArxivFetchClient
}

func NewRefleshPaperService(paperRepository repository.PaperRepository, arxivClient repository.ArxivFetchClient) RefleshPaperService {
	return &refleshPaperService{
		paperRepository: paperRepository,
		arxivClient:     arxivClient,
	}
}

type RefleshPaperInput struct {
	PapersMaxSize int64
}

type RefleshPaperOutput struct {
	RefleshedPaperSize int64
}

func (r *refleshPaperService) RefleshPaper(ctx context.Context, input *RefleshPaperInput) (*RefleshPaperOutput, error) {
	// TODO: 取得したいcategoriesを動的に変更できるようにする
	categories := []string{"cs.AI", "cs.CC", "cs.CE", "cs.CL", "cs.CR", "cs.DB", "cs.DC", "cs.DM", "cs.DS", "cs.ET", "cs.FL", "cs.IR", "cs.IT", "cs.LG", "cs.NE", "cs.SE", "physics.comp-ph", "quant-ph"}

	// arxivから論文を取得
	arxivClient := repository.NewArxivFetchClient()
	arxivRequest := &repository.ArxivFetchRequestV2{
		Categories: categories,
		MaxSize:    input.PapersMaxSize,
	}
	arxivResponse, err := arxivClient.FetchV2(ctx, arxivRequest)
	if err != nil {
		return nil, err
	}

	// 古い論文を削除
	// 24時間前以前のデータを削除
	targetTime := sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true}
	if err = r.paperRepository.Delete(targetTime); err != nil {
		return nil, err
	}

	// 取得した論文をDBに保存
	papers := make([]*xo.Paper, 0, len(arxivResponse.Papers))
	for _, arxivPaper := range arxivResponse.Papers {
		papers = append(papers, &xo.Paper{
			ArticleID:  arxivPaper.ArticleID,
			Title:      arxivPaper.Title,
			Abstract:   arxivPaper.Abstract,
			AbstractJp: "",
		})
	}
	err = r.paperRepository.BulkInsert(papers)
	if err != nil {
		return nil, err
	}

	return &RefleshPaperOutput{
		RefleshedPaperSize: int64(len(papers)),
	}, nil
}
