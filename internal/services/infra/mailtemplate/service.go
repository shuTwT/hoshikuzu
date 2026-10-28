package mailtemplate

import (
	"context"
	"fmt"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/mailtemplate"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type MailTemplateService interface {
	QueryMailTemplate(ctx context.Context, id int) (*ent.MailTemplate, error)
	QueryMailTemplatePage(ctx context.Context, pageQuery model.MailTemplatePageQuery) (int, []*ent.MailTemplate, error)
	CreateMailTemplate(ctx context.Context, createReq model.MailTemplateCreateReq) (*ent.MailTemplate, error)
	UpdateMailTemplate(ctx context.Context, id int, updateReq model.MailTemplateUpdateReq) (*ent.MailTemplate, error)
	DeleteMailTemplate(ctx context.Context, id int) error
}

type MailTemplateServiceImpl struct {
	client *ent.Client
}

func NewMailTemplateServiceImpl(client *ent.Client) *MailTemplateServiceImpl {
	return &MailTemplateServiceImpl{client: client}
}

func (s *MailTemplateServiceImpl) QueryMailTemplate(ctx context.Context, id int) (*ent.MailTemplate, error) {
	t, err := s.client.MailTemplate.Query().
		Where(mailtemplate.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, model.ErrMailTemplateNotFound
		}
		return nil, fmt.Errorf("查询邮件模板失败: %w", err)
	}
	return t, nil
}

func (s *MailTemplateServiceImpl) QueryMailTemplatePage(ctx context.Context, pageQuery model.MailTemplatePageQuery) (int, []*ent.MailTemplate, error) {
	query := s.client.MailTemplate.Query()
	if pageQuery.Name != "" {
		query = query.Where(mailtemplate.NameContains(pageQuery.Name))
	}
	if pageQuery.Code != "" {
		query = query.Where(mailtemplate.CodeContains(pageQuery.Code))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("统计邮件模板失败: %w", err)
	}

	templates, err := query.
		Order(ent.Desc(mailtemplate.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("查询邮件模板列表失败: %w", err)
	}

	return count, templates, nil
}

func (s *MailTemplateServiceImpl) CreateMailTemplate(ctx context.Context, createReq model.MailTemplateCreateReq) (*ent.MailTemplate, error) {
	exists, err := s.client.MailTemplate.Query().
		Where(mailtemplate.CodeEQ(createReq.Code)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询邮件模板编码失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("%w: %s", model.ErrMailTemplateCodeExists, createReq.Code)
	}

	contentType := "html"
	if createReq.ContentType != nil && *createReq.ContentType != "" {
		contentType = *createReq.ContentType
	}

	t, err := s.client.MailTemplate.Create().
		SetName(createReq.Name).
		SetCode(createReq.Code).
		SetSubject(createReq.Subject).
		SetContent(createReq.Content).
		SetContentType(contentType).
		SetNillableEnabled(createReq.Enabled).
		SetNillableRemark(createReq.Remark).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fmt.Errorf("%w: %s", model.ErrMailTemplateCodeExists, createReq.Code)
		}
		return nil, fmt.Errorf("创建邮件模板失败: %w", err)
	}
	return t, nil
}

func (s *MailTemplateServiceImpl) UpdateMailTemplate(ctx context.Context, id int, updateReq model.MailTemplateUpdateReq) (*ent.MailTemplate, error) {
	if updateReq.Code != nil && *updateReq.Code != "" {
		exists, err := s.client.MailTemplate.Query().
			Where(mailtemplate.CodeEQ(*updateReq.Code), mailtemplate.IDNEQ(id)).
			Exist(ctx)
		if err != nil {
			return nil, fmt.Errorf("查询邮件模板编码失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("%w: %s", model.ErrMailTemplateCodeExists, *updateReq.Code)
		}
	}

	update := s.client.MailTemplate.UpdateOneID(id)

	if updateReq.Name != nil {
		update.SetNillableName(updateReq.Name)
	}
	if updateReq.Code != nil {
		update.SetNillableCode(updateReq.Code)
	}
	if updateReq.Subject != nil {
		update.SetNillableSubject(updateReq.Subject)
	}
	if updateReq.Content != nil {
		update.SetNillableContent(updateReq.Content)
	}
	if updateReq.ContentType != nil {
		update.SetNillableContentType(updateReq.ContentType)
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
			return nil, model.ErrMailTemplateNotFound
		}
		if ent.IsConstraintError(err) {
			code := ""
			if updateReq.Code != nil {
				code = *updateReq.Code
			}
			return nil, fmt.Errorf("%w: %s", model.ErrMailTemplateCodeExists, code)
		}
		return nil, fmt.Errorf("更新邮件模板失败: %w", err)
	}
	return t, nil
}

func (s *MailTemplateServiceImpl) DeleteMailTemplate(ctx context.Context, id int) error {
	err := s.client.MailTemplate.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.ErrMailTemplateNotFound
		}
		return fmt.Errorf("删除邮件模板失败: %w", err)
	}
	return nil
}
