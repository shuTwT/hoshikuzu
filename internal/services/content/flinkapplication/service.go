package flinkapplication

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/flink"
	"github.com/shuTwT/hoshikuzu/ent/flinkapplication"
	"github.com/shuTwT/hoshikuzu/ent/predicate"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type FlinkApplicationService interface {
	CreateFlinkApplication(c context.Context, createReq *model.FlinkApplicationCreateReq) (*ent.FLinkApplication, error)
	ListFlinkApplicationPage(c context.Context, pageQuery model.FlinkApplicationPageReq) ([]*ent.FLinkApplication, int, error)
	QueryFlinkApplication(c context.Context, id int) (*ent.FLinkApplication, error)
	ApproveFlinkApplication(c context.Context, id int, status int, rejectReason string) (*ent.FLinkApplication, error)
}

type FlinkApplicationServiceImpl struct {
	client *ent.Client
}

func NewFlinkApplicationServiceImpl(client *ent.Client) *FlinkApplicationServiceImpl {
	return &FlinkApplicationServiceImpl{
		client: client,
	}
}

func (s *FlinkApplicationServiceImpl) CreateFlinkApplication(c context.Context, createReq *model.FlinkApplicationCreateReq) (*ent.FLinkApplication, error) {
	return s.client.FLinkApplication.Create().
		SetWebsiteURL(createReq.WebsiteURL).
		SetApplicationType(createReq.ApplicationType).
		SetWebsiteName(createReq.WebsiteName).
		SetWebsiteLogo(createReq.WebsiteLogo).
		SetWebsiteDescription(createReq.WebsiteDescription).
		SetContactEmail(createReq.ContactEmail).
		SetNillableSnapshotURL(&createReq.SnapshotURL).
		SetNillableOriginalWebsiteURL(&createReq.OriginalWebsiteURL).
		SetNillableModificationReason(&createReq.ModificationReason).
		SetStatus(0).
		Save(c)
}

func (s *FlinkApplicationServiceImpl) ListFlinkApplicationPage(c context.Context, pageQuery model.FlinkApplicationPageReq) ([]*ent.FLinkApplication, int, error) {
	var preds []predicate.FLinkApplication
	if pageQuery.Status != nil {
		preds = append(preds, flinkapplication.StatusEQ(*pageQuery.Status))
	}
	if pageQuery.ApplicationType != nil {
		preds = append(preds, flinkapplication.ApplicationTypeEQ(*pageQuery.ApplicationType))
	}
	count, err := s.client.FLinkApplication.Query().Where(preds...).Count(c)
	if err != nil {
		return nil, 0, err
	}
	applications, err := s.client.FLinkApplication.Query().
		Where(preds...).
		Order(ent.Desc(flinkapplication.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).All(c)
	if err != nil {
		return nil, 0, err
	}
	return applications, count, nil
}

func (s *FlinkApplicationServiceImpl) QueryFlinkApplication(c context.Context, id int) (*ent.FLinkApplication, error) {
	return s.client.FLinkApplication.Query().
		Where(flinkapplication.ID(id)).
		First(c)
}

func (s *FlinkApplicationServiceImpl) ApproveFlinkApplication(c context.Context, id int, status int, rejectReason string) (*ent.FLinkApplication, error) {
	application, err := s.client.FLinkApplication.Query().
		Where(flinkapplication.ID(id)).
		First(c)
	if err != nil {
		return nil, err
	}

	application, err = s.client.FLinkApplication.UpdateOneID(id).
		SetStatus(status).
		SetNillableRejectReason(&rejectReason).
		Save(c)
	if err != nil {
		return nil, err
	}

	if status == 1 {
		if application.ApplicationType == "create" {
			_, err := s.client.FLink.Create().
				SetName(application.WebsiteName).
				SetURL(application.WebsiteURL).
				SetAvatarURL(application.WebsiteLogo).
				SetDescription(application.WebsiteDescription).
				SetSnapshotURL(application.SnapshotURL).
				SetEmail(application.ContactEmail).
				SetStatus(1).
				Save(c)
			if err != nil {
				return nil, err
			}
		} else if application.ApplicationType == "update" && application.OriginalWebsiteURL != "" {
			existingFlink, err := s.client.FLink.Query().
				Where(flink.URLEQ(application.OriginalWebsiteURL)).
				First(c)
			if err == nil {
				_, err := s.client.FLink.UpdateOneID(existingFlink.ID).
					SetName(application.WebsiteName).
					SetURL(application.WebsiteURL).
					SetAvatarURL(application.WebsiteLogo).
					SetDescription(application.WebsiteDescription).
					SetSnapshotURL(application.SnapshotURL).
					SetEmail(application.ContactEmail).
					Save(c)
				if err != nil {
					return nil, err
				}
			}
		}
	} else if status == 2 {
		s.sendRejectEmail(application.ContactEmail, application.WebsiteName, rejectReason)
	}

	return application, nil
}

func (s *FlinkApplicationServiceImpl) sendRejectEmail(email, websiteName, reason string) {
}
