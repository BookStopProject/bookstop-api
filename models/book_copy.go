package models

type BookCondition string

const (
	BookConditionNew        BookCondition = "new"
	BookConditionLikeNew    BookCondition = "like_new"
	BookConditionGood       BookCondition = "good"
	BookConditionAcceptable BookCondition = "acceptable"
)

type BookCopy struct {
	ID         int           `json:"id"`
	BookID     int           `json:"bookId"`
	Condition  BookCondition `json:"condition"`
	Book       *Book         `json:"book"`
	LocationID *int          `json:"locationId"`
	Location   *Location     `json:"location"`
}
