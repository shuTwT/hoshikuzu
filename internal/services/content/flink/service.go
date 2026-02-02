package flink

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/flink"
	"github.com/shuTwT/hoshikuzu/ent/flinkgroup"
	"github.com/shuTwT/hoshikuzu/ent/predicate"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type FlinkService interface {
	CountFlinkByGroupID(c context.Context, groupID int) (int, error)
	ListFlink(c context.Context, listQuery model.FlinkListReq) ([]*ent.FLink, error)
	ListFlinkPage(c context.Context, pageQuery model.FlinkPageReq) ([]*ent.FLink, int, error)
}

type FlinkServiceImpl struct {
	client *ent.Client
}

func NewFlinkServiceImpl(client *ent.Client) *FlinkServiceImpl {
	return &FlinkServiceImpl{
		client: client,
	}
}

func (s *FlinkServiceImpl) CountFlinkByGroupID(c context.Context, groupID int) (int, error) {
	return s.client.FLink.Query().Where(flink.GroupIDEQ(groupID)).Count(c)
}

func (s *FlinkServiceImpl) ListFlink(c context.Context, listPage model.FlinkListReq) ([]*ent.FLink, error) {
	var preds []predicate.FLinkGroup
	if listPage.GroupId != nil {
		preds = append(preds, flinkgroup.IDEQ(*listPage.GroupId))
	}
	if listPage.GroupName != nil {
		preds = append(preds, flinkgroup.NameContains(*listPage.GroupName))
	}
	flinks, err := s.client.FLinkGroup.
		Query().
		Where(preds...).
		QueryLinks().
		WithGroup().
		Order(ent.Desc(flink.FieldID)).
		All(c)
	if err != nil {
		return nil, err
	}
	return flinks, nil
}

func (s *FlinkServiceImpl) ListFlinkPage(c context.Context, pageQuery model.FlinkPageReq) ([]*ent.FLink, int, error) {
	var preds []predicate.FLink
	if pageQuery.GroupId != nil {
		preds = append(preds, flink.GroupIDEQ(*pageQuery.GroupId))
	}
	count, err := s.client.FLink.Query().Where(preds...).Count(c)
	if err != nil {
		return nil, 0, err
	}
	flinks, err := s.client.FLink.Query().
		Where(preds...).
		Order(ent.Desc(flink.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).All(c)
	if err != nil {
		return nil, 0, err
	}
	return flinks, count, nil
}
