package paperservice

import (
	"context"

	"github.com/rmocchy/paper-news-backend-v2/gin-lambda/gen/xo"
	repository "github.com/rmocchy/paper-news-backend-v2/gin-lambda/repository/paper-repository"
)

type ListPaperService interface {
	ListPapers(ctx context.Context, input *ListPapersInput) (*ListPapersOutput, error)
}

type listPaperService struct {
	paperRepository repository.PaperRepository
}

func NewListPaperService(paperRepository repository.PaperRepository) ListPaperService {
	return &listPaperService{
		paperRepository: paperRepository,
	}
}

type ListPapersInput struct{}

type ListPapersOutput struct {
	Papers []*xo.Paper
}

func (l *listPaperService) ListPapers(ctx context.Context, input *ListPapersInput) (*ListPapersOutput, error) {
	// todo add params: pageSize, pageNumber, filters
	pageSize := 10
	pageNumber := 1

	res, err := l.paperRepository.List(int64(pageSize), int64(pageNumber), repository.Filter{})
	if err != nil {
		return nil, err
	}
	return &ListPapersOutput{
		Papers: res,
	}, nil
}
