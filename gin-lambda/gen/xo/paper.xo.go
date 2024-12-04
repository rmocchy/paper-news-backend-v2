package xo

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Paper represents a row from 'paper-news_local.papers'.
type Paper struct {
	ID         uint         `json:"id" db:"id"`                   // id
	ArticleID  string       `json:"article_id" db:"article_id"`   // article_id
	Title      string       `json:"title" db:"title"`             // title
	Abstract   string       `json:"abstract" db:"abstract"`       // abstract
	AbstractJp string       `json:"abstract_jp" db:"abstract_jp"` // abstract_jp
	CreatedAt  sql.NullTime `json:"created_at" db:"created_at"`   // created_at
	UpdatedAt  sql.NullTime `json:"updated_at" db:"updated_at"`   // updated_at
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Paper] exists in the database.
func (p *Paper) Exists() bool {
	return p._exists
}

// Deleted returns true when the [Paper] has been marked for deletion
// from the database.
func (p *Paper) Deleted() bool {
	return p._deleted
}

// Insert inserts the [Paper] to the database.
func (p *Paper) Insert(ctx context.Context, db DB) error {
	switch {
	case p._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case p._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO paper-news_local.papers (` +
		`article_id, title, abstract, abstract_jp, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, p.ArticleID, p.Title, p.Abstract, p.AbstractJp, p.CreatedAt, p.UpdatedAt)
	res, err := db.ExecContext(ctx, sqlstr, p.ArticleID, p.Title, p.Abstract, p.AbstractJp, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	p.ID = uint(id)
	// set exists
	p._exists = true
	return nil
}

// Update updates a [Paper] in the database.
func (p *Paper) Update(ctx context.Context, db DB) error {
	switch {
	case !p._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case p._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE paper-news_local.papers SET ` +
		`article_id = ?, title = ?, abstract = ?, abstract_jp = ?, created_at = ?, updated_at = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, p.ArticleID, p.Title, p.Abstract, p.AbstractJp, p.CreatedAt, p.UpdatedAt, p.ID)
	if _, err := db.ExecContext(ctx, sqlstr, p.ArticleID, p.Title, p.Abstract, p.AbstractJp, p.CreatedAt, p.UpdatedAt, p.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Paper] to the database.
func (p *Paper) Save(ctx context.Context, db DB) error {
	if p.Exists() {
		return p.Update(ctx, db)
	}
	return p.Insert(ctx, db)
}

// Upsert performs an upsert for [Paper].
func (p *Paper) Upsert(ctx context.Context, db DB) error {
	switch {
	case p._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO paper-news_local.papers (` +
		`id, article_id, title, abstract, abstract_jp, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`article_id = VALUES(article_id), title = VALUES(title), abstract = VALUES(abstract), abstract_jp = VALUES(abstract_jp), created_at = VALUES(created_at), updated_at = VALUES(updated_at)`
	// run
	logf(sqlstr, p.ID, p.ArticleID, p.Title, p.Abstract, p.AbstractJp, p.CreatedAt, p.UpdatedAt)
	if _, err := db.ExecContext(ctx, sqlstr, p.ID, p.ArticleID, p.Title, p.Abstract, p.AbstractJp, p.CreatedAt, p.UpdatedAt); err != nil {
		return logerror(err)
	}
	// set exists
	p._exists = true
	return nil
}

// Delete deletes the [Paper] from the database.
func (p *Paper) Delete(ctx context.Context, db DB) error {
	switch {
	case !p._exists: // doesn't exist
		return nil
	case p._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM paper-news_local.papers ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, p.ID)
	if _, err := db.ExecContext(ctx, sqlstr, p.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	p._deleted = true
	return nil
}

// PapersByCreatedAt retrieves a row from 'paper-news_local.papers' as a [Paper].
//
// Generated from index 'idx_created_at'.
func PapersByCreatedAt(ctx context.Context, db DB, createdAt sql.NullTime) ([]*Paper, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, article_id, title, abstract, abstract_jp, created_at, updated_at ` +
		`FROM paper-news_local.papers ` +
		`WHERE created_at = ?`
	// run
	logf(sqlstr, createdAt)
	rows, err := db.QueryContext(ctx, sqlstr, createdAt)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Paper
	for rows.Next() {
		p := Paper{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&p.ID, &p.ArticleID, &p.Title, &p.Abstract, &p.AbstractJp, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// PaperByID retrieves a row from 'paper-news_local.papers' as a [Paper].
//
// Generated from index 'papers_id_pkey'.
func PaperByID(ctx context.Context, db DB, id uint) (*Paper, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, article_id, title, abstract, abstract_jp, created_at, updated_at ` +
		`FROM paper-news_local.papers ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	p := Paper{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&p.ID, &p.ArticleID, &p.Title, &p.Abstract, &p.AbstractJp, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, logerror(err)
	}
	return &p, nil
}

// PaperByArticleID retrieves a row from 'paper-news_local.papers' as a [Paper].
//
// Generated from index 'uq_article_id'.
func PaperByArticleID(ctx context.Context, db DB, articleID string) (*Paper, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, article_id, title, abstract, abstract_jp, created_at, updated_at ` +
		`FROM paper-news_local.papers ` +
		`WHERE article_id = ?`
	// run
	logf(sqlstr, articleID)
	p := Paper{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, articleID).Scan(&p.ID, &p.ArticleID, &p.Title, &p.Abstract, &p.AbstractJp, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, logerror(err)
	}
	return &p, nil
}