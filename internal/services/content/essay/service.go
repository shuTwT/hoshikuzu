package essay

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/essay"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type EssayService interface {
	CreateEssay(ctx context.Context, userId int, req *model.EssayCreateReq) (*ent.Essay, error)
	UpdateEssay(ctx context.Context, id int, req *model.EssayUpdateReq) error
	GetEssay(ctx context.Context, id int) (*ent.Essay, error)
	GetEssayPage(ctx context.Context, page, pageSize int) ([]*ent.Essay, int, error)
	GetEssayList(ctx context.Context, limit int) ([]*ent.Essay, error)
	DeleteEssay(ctx context.Context, id int) error
}

type EssayServiceImpl struct {
	client *ent.Client
}

func NewEssayServiceImpl(client *ent.Client) EssayService {
	return &EssayServiceImpl{client: client}
}

func (s *EssayServiceImpl) CreateEssay(ctx context.Context, userId int, req *model.EssayCreateReq) (*ent.Essay, error) {
	essay := s.client.Essay.Create().
		SetContent(req.Content).
		SetDraft(req.Draft).
		SetImages(req.Images).
		SetUserID(userId).
		SaveX(ctx)
	return essay, nil
}

func (s *EssayServiceImpl) UpdateEssay(ctx context.Context, id int, req *model.EssayUpdateReq) error {
	return s.client.Essay.UpdateOneID(id).
		SetContent(req.Content).
		SetDraft(req.Draft).
		SetImages(req.Images).
		Exec(ctx)
}

func (s *EssayServiceImpl) GetEssay(ctx context.Context, id int) (*ent.Essay, error) {
	essay, err := s.client.Essay.Query().
		Where(essay.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return essay, nil
}

func (s *EssayServiceImpl) GetEssayPage(ctx context.Context, page, pageSize int) ([]*ent.Essay, int, error) {
	count, err := s.client.Essay.Query().
		Order(ent.Desc(essay.FieldID)).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	essays, err := s.client.Essay.Query().
		Order(ent.Desc(essay.FieldID)).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return essays, count, nil
}

func (s *EssayServiceImpl) GetEssayList(ctx context.Context, limit int) ([]*ent.Essay, error) {
	essays, err := s.client.Essay.Query().
		Order(ent.Desc(essay.FieldID)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return essays, nil
}

func (s *EssayServiceImpl) DeleteEssay(ctx context.Context, id int) error {
	return s.client.Essay.DeleteOneID(id).Exec(ctx)
}
