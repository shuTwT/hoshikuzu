package common

import (
	"context"

	comment_service "github.com/shuTwT/hoshikuzu/internal/services/content/comment"
	post_service "github.com/shuTwT/hoshikuzu/internal/services/content/post"
	user_service "github.com/shuTwT/hoshikuzu/internal/services/system/user"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type CommonService interface {
	GetHomeStatistic(c context.Context) model.HomeStatistic
}

type CommonServiceImpl struct {
	client         *ent.Client
	userService    user_service.UserService
	postService    post_service.PostService
	commentService comment_service.CommentService
}

func NewCommonServiceImpl(client *ent.Client, userService user_service.UserService, postService post_service.PostService, commentService comment_service.CommentService) *CommonServiceImpl {
	return &CommonServiceImpl{client: client, userService: userService, postService: postService, commentService: commentService}
}

func (s *CommonServiceImpl) GetHomeStatistic(c context.Context) model.HomeStatistic {
	userCount, _ := s.userService.GetUserCount(c)
	postCount, _ := s.postService.GetPostCount(c)
	commentCount, _ := s.commentService.GetCommentCount(c)
	homeStatistic := model.HomeStatistic{
		PostCount:    postCount,
		UserCount:    userCount,
		CommentCount: commentCount,
		VisitCount:   0,
	}
	return homeStatistic
}
