package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rmocchy/paper-news-backend-v2/gin-lambda/gen/xo"
)

type PaperRepository interface {
	BulkInsert(paper []*xo.Paper) error
	List(pageSize int64, pageNumber int64, filter Filter) ([]*xo.Paper, error)
	Delete(targetCreatedAt sql.NullTime) error
}

type Filter struct {
	Title      string
	Authors    []string
	Categories []string
}

type paperRepositoryImpl struct {
	db *sqlx.DB
}

func NewPaperRepository(db *sqlx.DB) PaperRepository {
	return &paperRepositoryImpl{
		db: db,
	}
}

func (p *paperRepositoryImpl) BulkInsert(papers []*xo.Paper) error {
	if len(papers) == 0 {
		return nil
	}
	var query = "INSERT INTO papers (article_id, title, abstract, abstract_jp) VALUES "
	args := make([]any, 0, 4*len(papers))
	for i, paper := range papers {
		query += "(?, ?, ?, ?)"
		if i != len(papers)-1 {
			query += ", "
		} else {
			query += " ON DUPLICATE KEY UPDATE title = VALUES(title), abstract = VALUES(abstract), abstract_jp = VALUES(abstract_jp);"
		}
		args = append(args, paper.ArticleID, paper.Title, paper.Abstract, paper.AbstractJp)
	}
	_, err := p.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (p *paperRepositoryImpl) Delete(targetCreatedAt sql.NullTime) error {
	query := "DELETE FROM papers WHERE `created_at` < ?;"
	_, err := p.db.Exec(query, targetCreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *paperRepositoryImpl) List(pageSize int64, pageNumber int64, filter Filter) ([]*xo.Paper, error) {
	query := "SELECT * FROM papers ORDER BY `created_at` DESC LIMIT ? OFFSET ?;"
	rows, err := p.db.Queryx(query, pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	results := make([]*xo.Paper, 0, pageSize)
	for rows.Next() {
		var paper xo.Paper
		err := rows.StructScan(&paper)
		if err != nil {
			return nil, err
		}
		results = append(results, &paper)
	}

	return results, nil
}
