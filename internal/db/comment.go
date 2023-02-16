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

func (d *Database) PostComment(ctx context.Context, newComment comment.Comment) (comment.Comment, error) {
	newComment.ID = uuid.New().String()

	postCommentRow := CommentRow{
		ID:     newComment.ID,
		Slug:   sql.NullString{String: newComment.Slug, Valid: true},
		Author: sql.NullString{String: newComment.Author, Valid: true},
		Body:   sql.NullString{String: newComment.Body, Valid: true},
	}
	rows, err := d.Client.NamedQueryContext(ctx, `INSERT INTO comments (id, slug, author, body)
		VALUES (:id, :slug, :author, :body)`, postCommentRow)

	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close db rows: %w", err)
	}

	return newComment, nil
}

func (d *Database) DeleteComment(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(ctx, `DELETE FROM comments WHERE id = $1`, id)

	if err != nil {
		return fmt.Errorf("failed to delete comment from database: %w", err)
	}

	return nil
}

func (d *Database) UpdateComment(ctx context.Context, id string, updateComment comment.Comment) (comment.Comment, error) {
	commentRow := CommentRow{
		ID: id,
		Slug: sql.NullString{String: updateComment.Slug, Valid: true},
		Author: sql.NullString{String: updateComment.Author, Valid: true},
		Body: sql.NullString{String: updateComment.Body, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(ctx,
		`UPDATE comments SET slug = :slug, author = :author, body = :body WHERE id = :id`, commentRow,
	)

	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to update comment: %w", err)
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return convertCommentRowToComment(commentRow), nil
}
