package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/BFamzz/comments-api/internal/comment"
	"github.com/google/uuid"
)

type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

func convertCommentRowToComment(c CommentRow) comment.Comment {
	return comment.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Body:   c.Body.String,
		Author: c.Author.String,
	}
}

func (d *Database) GetComment(ctx context.Context, uuid string) (comment.Comment, error) {
	var commentRow CommentRow
	row := d.Client.QueryRowContext(ctx, `SELECT id, slug, body, author FROM comments WHERE id = $1`, uuid)

	err := row.Scan(&commentRow.ID, &commentRow.Slug, &commentRow.Body, &commentRow.Author)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error fetching the comment by uuid: %w", err)
	}

	return convertCommentRowToComment(commentRow), nil
}

func (d *Database) PostComment(ctx context.Context, comment comment.Comment) (comment.Comment, error) {
	comment.ID = uuid.New().String()
	fmt.Println(comment.ID)

	postCommentRow := CommentRow{
		ID:     comment.ID,
		Slug:   sql.NullString{String: comment.Slug, Valid: true},
		Author: sql.NullString{String: comment.Author, Valid: true},
		Body:   sql.NullString{String: comment.Body, Valid: true},
	}
	rows, err := d.Client.NamedQueryContext(ctx, `INSERT INTO comments (id, slug, author, body)
		VALUES (:id, :slug, :author, :body)`, postCommentRow)

	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close db rows: %w", err)
	}

	return comment, nil
}
