package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.21 DO NOT EDIT.

import (
	"bookstop/auth"
	"bookstop/models"
	"context"
	"fmt"
	"strconv"
)

// PostCreate is the resolver for the postCreate field.
func (r *mutationResolver) PostCreate(ctx context.Context, text string, bookID string, isRecommending bool) (*models.Post, error) {
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	bookIDNum, _ := strconv.Atoi(bookID)

	post := models.Post{
		Text:           text,
		BookID:         bookIDNum,
		UserID:         usr.ID,
		IsRecommending: isRecommending,
	}

	return models.CreatePost(ctx, &post)
}

// PostUpdate is the resolver for the postUpdate field.
func (r *mutationResolver) PostUpdate(ctx context.Context, id string, text string, isRecommending bool) (*models.Post, error) {
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	postID, _ := strconv.Atoi(id)
	post, err := models.FindPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if !post.IsOwner(ctx, usr.ID) {
		return nil, auth.ErrUnauthorized
	}
	post.Text = text
	post.IsRecommending = isRecommending
	return models.UpdatePost(ctx, post)
}

// PostDelete is the resolver for the postDelete field.
func (r *mutationResolver) PostDelete(ctx context.Context, id string) (bool, error) {
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return false, err
	}
	if usr == nil {
		return false, auth.ErrUnauthorized
	}
	postID, _ := strconv.Atoi(id)
	post, err := models.FindPostByID(ctx, postID)
	if err != nil {
		return false, err
	}
	if post == nil {
		return false, fmt.Errorf("post not found")
	}
	if !post.IsOwner(ctx, usr.ID) {
		return false, auth.ErrUnauthorized
	}
	err = models.DeletePost(ctx, post.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, userID *string, bookID *string, limit *int, offset *int) ([]*models.Post, error) {
	if userID != nil {
		usedIDNum, _ := strconv.Atoi(*userID)
		return models.FindPostsByUserID(ctx, usedIDNum)
	}
	if bookID != nil {
		bookIDNum, _ := strconv.Atoi(*bookID)
		return models.FindPostsByBookID(ctx, bookIDNum)
	}
	return models.FindAllPosts(ctx)
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	postID, _ := strconv.Atoi(id)
	return models.FindPostByID(ctx, postID)
}
