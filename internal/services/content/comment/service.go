package comment

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/comment"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/medama-io/go-useragent"
)

type CommentService interface {
	GetCommentCount(c context.Context) (int, error)
	ListCommentPage(c context.Context, pageQuery model.PageQuery) (*model.PageResult[*ent.Comment], error)
	ListComment(c context.Context, url string) ([]*ent.Comment, error)
	GetComment(c context.Context, id int) (*ent.Comment, error)
	CountComment(c context.Context, includeReply bool, urls []string) (int64, error)
	CreateComment(c context.Context, comment string, href string, link string, mail string, nick string, ua string, url string, ipAddress string) (*int, error)
	GetRecentComment(c context.Context, pageSize int) ([]*ent.Comment, error)
	ParseUserAgent(ua string) (browser string, os string)
}

type CommentServiceImpl struct {
	client *ent.Client
}

func NewCommentServiceImpl(client *ent.Client) *CommentServiceImpl {
	return &CommentServiceImpl{client: client}
}

func (s *CommentServiceImpl) ListCommentPage(c context.Context, pageQuery model.PageQuery) (*model.PageResult[*ent.Comment], error) {
	client := s.client

	count, err := client.Comment.Query().Count(c)
	if err != nil {
		return nil, err
	}
	comments, err := client.Comment.Query().
		Order(ent.Desc(comment.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(c)
	if err != nil {
		return nil, err
	}
	pageResult := &model.PageResult[*ent.Comment]{
		Total:   int64(count),
		Records: comments,
	}
	return pageResult, nil
}

func (s *CommentServiceImpl) ListComment(c context.Context, url string) ([]*ent.Comment, error) {
	comments, err := s.client.Comment.Query().
		Order(ent.Desc(comment.FieldID)).
		Where(comment.URLEQ(url)).
		All(c)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentServiceImpl) GetComment(c context.Context, id int) (*ent.Comment, error) {
	comment, err := s.client.Comment.Query().
		Where(comment.IDEQ(id)).
		First(c)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentServiceImpl) CountComment(c context.Context, includeReply bool, urls []string) (int64, error) {
	if includeReply {
		count, err := s.client.Comment.Query().
			Where(comment.URLIn(urls...)).
			Where(comment.Or(comment.ParentIDIsNil(), comment.ParentIDNotIn(0))).
			Count(c)
		if err != nil {
			return 0, err
		}
		return int64(count), nil
	}
	count, err := s.client.Comment.Query().
		Where(comment.URLIn(urls...)).
		Where(comment.ParentIDIsNil()).
		Count(c)
	if err != nil {
		return 0, err
	}
	return int64(count), nil
}

func (s *CommentServiceImpl) CreateComment(c context.Context, comment string, href string, link string, mail string, nick string, ua string, url string, ipAddress string) (*int, error) {
	entity, err := s.client.Comment.Create().
		SetContent(comment).
		SetUserAgent(ua).
		SetURL(url).
		SetIPAddress(ipAddress).
		Save(c)
	if err != nil {
		return nil, err
	}
	return &entity.ID, nil
}

func (s *CommentServiceImpl) ParseUserAgent(ua string) (browser string, os string) {
	// Create a new parser. Initialize only once during application startup.
	parser := useragent.NewParser()
	agent := parser.Parse(ua)
	return agent.BrowserVersion(), agent.OS().String()
}

func (s *CommentServiceImpl) GetRecentComment(c context.Context, pageSize int) ([]*ent.Comment, error) {
	comments, err := s.client.Comment.Query().
		Order(ent.Desc(comment.FieldCreatedAt)).
		Limit(pageSize).
		All(c)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentServiceImpl) GetCommentCount(c context.Context) (int, error) {
	count, err := s.client.Comment.Query().Count(c)
	if err != nil {
		return 0, err
	}
	return count, nil
}
