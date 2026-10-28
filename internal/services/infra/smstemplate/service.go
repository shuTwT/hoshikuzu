package smstemplate

import (
	"context"
	"fmt"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/smstemplate"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type SmsTemplateService interface {
	QuerySmsTemplate(ctx context.Context, id int) (*ent.SmsTemplate, error)
	QuerySmsTemplatePage(ctx context.Context, pageQuery model.SmsTemplatePageQuery) (int, []*ent.SmsTemplate, error)
	CreateSmsTemplate(ctx context.Context, createReq model.SmsTemplateCreateReq) (*ent.SmsTemplate, error)
	UpdateSmsTemplate(ctx context.Context, id int, updateReq model.SmsTemplateUpdateReq) (*ent.SmsTemplate, error)
	DeleteSmsTemplate(ctx context.Context, id int) error
}

type SmsTemplateServiceImpl struct {
	client *ent.Client
}

func NewSmsTemplateServiceImpl(client *ent.Client) *SmsTemplateServiceImpl {
	return &SmsTemplateServiceImpl{client: client}
}

func (s *SmsTemplateServiceImpl) QuerySmsTemplate(ctx context.Context, id int) (*ent.SmsTemplate, error) {
	t, err := s.client.SmsTemplate.Query().
		Where(smstemplate.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, model.ErrSmsTemplateNotFound
		}
		return nil, fmt.Errorf("查询短信模板失败: %w", err)
	}
	return t, nil
}

func (s *SmsTemplateServiceImpl) QuerySmsTemplatePage(ctx context.Context, pageQuery model.SmsTemplatePageQuery) (int, []*ent.SmsTemplate, error) {
	query := s.client.SmsTemplate.Query()
	if pageQuery.Name != "" {
		query = query.Where(smstemplate.NameContains(pageQuery.Name))
	}
	if pageQuery.Code != "" {
		query = query.Where(smstemplate.CodeContains(pageQuery.Code))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("统计短信模板失败: %w", err)
	}

	templates, err := query.
		Order(ent.Desc(smstemplate.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("查询短信模板列表失败: %w", err)
	}

	return count, templates, nil
}

func (s *SmsTemplateServiceImpl) CreateSmsTemplate(ctx context.Context, createReq model.SmsTemplateCreateReq) (*ent.SmsTemplate, error) {
	exists, err := s.client.SmsTemplate.Query().
		Where(smstemplate.CodeEQ(createReq.Code)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询短信模板编码失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("%w: %s", model.ErrSmsTemplateCodeExists, createReq.Code)
	}

	t, err := s.client.SmsTemplate.Create().
		SetName(createReq.Name).
		SetCode(createReq.Code).
		SetContent(createReq.Content).
		SetNillableSignName(createReq.SignName).
		SetNillableProviderTemplateID(createReq.ProviderTemplateID).
		SetNillableEnabled(createReq.Enabled).
		SetNillableRemark(createReq.Remark).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fmt.Errorf("%w: %s", model.ErrSmsTemplateCodeExists, createReq.Code)
		}
		return nil, fmt.Errorf("创建短信模板失败: %w", err)
	}
	return t, nil
}

func (s *SmsTemplateServiceImpl) UpdateSmsTemplate(ctx context.Context, id int, updateReq model.SmsTemplateUpdateReq) (*ent.SmsTemplate, error) {
	if updateReq.Code != nil && *updateReq.Code != "" {
		exists, err := s.client.SmsTemplate.Query().
			Where(smstemplate.CodeEQ(*updateReq.Code), smstemplate.IDNEQ(id)).
			Exist(ctx)
		if err != nil {
			return nil, fmt.Errorf("查询短信模板编码失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("%w: %s", model.ErrSmsTemplateCodeExists, *updateReq.Code)
		}
	}

	update := s.client.SmsTemplate.UpdateOneID(id)

	if updateReq.Name != nil {
		update.SetNillableName(updateReq.Name)
	}
	if updateReq.Code != nil {
		update.SetNillableCode(updateReq.Code)
	}
	if updateReq.Content != nil {
		update.SetNillableContent(updateReq.Content)
	}
	if updateReq.SignName != nil {
		update.SetNillableSignName(updateReq.SignName)
	}
	if updateReq.ProviderTemplateID != nil {
		update.SetNillableProviderTemplateID(updateReq.ProviderTemplateID)
	}
	if updateReq.Enabled != nil {
		update.SetNillableEnabled(updateReq.Enabled)
	}
	if updateReq.Remark != nil {
		update.SetNillableRemark(updateReq.Remark)
	}

	t, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, model.ErrSmsTemplateNotFound
		}
		if ent.IsConstraintError(err) {
			code := ""
			if updateReq.Code != nil {
				code = *updateReq.Code
			}
			return nil, fmt.Errorf("%w: %s", model.ErrSmsTemplateCodeExists, code)
		}
		return nil, fmt.Errorf("更新短信模板失败: %w", err)
	}
	return t, nil
}

func (s *SmsTemplateServiceImpl) DeleteSmsTemplate(ctx context.Context, id int) error {
	err := s.client.SmsTemplate.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.ErrSmsTemplateNotFound
		}
		return fmt.Errorf("删除短信模板失败: %w", err)
	}
	return nil
}
