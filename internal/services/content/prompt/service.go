package prompt

import (
	"context"
	"fmt"
	"strings"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/prompt"
	"github.com/shuTwT/hoshikuzu/ent/promptimage"
	"github.com/shuTwT/hoshikuzu/ent/prompttag"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type PromptService interface {
	ListPromptPage(ctx context.Context, query *model.PromptPageQuery) (int, []*ent.Prompt, error)
	CreatePrompt(ctx context.Context, req *model.PromptCreateReq) (*ent.Prompt, error)
	UpdatePrompt(ctx context.Context, id int, req *model.PromptUpdateReq) (*ent.Prompt, error)
	QueryPrompt(ctx context.Context, id int) (*ent.Prompt, error)
	DeletePrompt(ctx context.Context, id int) error
	ListPromptTag(ctx context.Context) ([]*ent.PromptTag, error)
}

type PromptServiceImpl struct {
	client *ent.Client
}

func NewPromptServiceImpl(client *ent.Client) *PromptServiceImpl {
	return &PromptServiceImpl{client: client}
}

// withImagesAndTags 预加载模型图片（按排序）与标签
func withImagesAndTags(q *ent.PromptQuery) *ent.PromptQuery {
	return q.
		WithImages(func(iq *ent.PromptImageQuery) {
			iq.Order(ent.Asc(promptimage.FieldSort), ent.Asc(promptimage.FieldID))
		}).
		WithTags(func(tq *ent.PromptTagQuery) {
			tq.Order(ent.Asc(prompttag.FieldID))
		})
}

func (s *PromptServiceImpl) ListPromptPage(ctx context.Context, query *model.PromptPageQuery) (int, []*ent.Prompt, error) {
	q := s.client.Prompt.Query()
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		q = q.Where(prompt.Or(
			prompt.TitleContains(keyword),
			prompt.ContentContains(keyword),
		))
	}
	if query.TagID > 0 {
		q = q.Where(prompt.HasTagsWith(prompttag.ID(query.TagID)))
	}

	count, err := q.Clone().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	prompts, err := withImagesAndTags(q).
		Order(ent.Asc(prompt.FieldSort), ent.Desc(prompt.FieldID)).
		Limit(query.Size).
		Offset((query.Page - 1) * query.Size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}
	return count, prompts, nil
}

func (s *PromptServiceImpl) CreatePrompt(ctx context.Context, req *model.PromptCreateReq) (*ent.Prompt, error) {
	if err := validateImages(req.Images); err != nil {
		return nil, err
	}

	var id int
	err := withTx(ctx, s.client, func(tx *ent.Tx) error {
		tagIDs, err := ensureTags(ctx, tx, req.Tags)
		if err != nil {
			return err
		}
		created, err := tx.Prompt.Create().
			SetTitle(req.Title).
			SetContent(req.Content).
			SetNegativePrompt(req.NegativePrompt).
			SetDescription(req.Description).
			SetSort(req.Sort).
			AddTagIDs(tagIDs...).
			Save(ctx)
		if err != nil {
			return err
		}
		id = created.ID
		return createImages(ctx, tx, id, req.Images)
	})
	if err != nil {
		return nil, err
	}
	return s.QueryPrompt(ctx, id)
}

func (s *PromptServiceImpl) UpdatePrompt(ctx context.Context, id int, req *model.PromptUpdateReq) (*ent.Prompt, error) {
	if err := validateImages(req.Images); err != nil {
		return nil, err
	}

	err := withTx(ctx, s.client, func(tx *ent.Tx) error {
		tagIDs, err := ensureTags(ctx, tx, req.Tags)
		if err != nil {
			return err
		}
		err = tx.Prompt.UpdateOneID(id).
			SetTitle(req.Title).
			SetContent(req.Content).
			SetNegativePrompt(req.NegativePrompt).
			SetDescription(req.Description).
			SetSort(req.Sort).
			ClearTags().
			AddTagIDs(tagIDs...).
			Exec(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return model.ErrPromptNotFound
			}
			return err
		}
		// 图片为全量覆盖：先删除旧图片再重建
		if _, err = tx.PromptImage.Delete().Where(promptimage.PromptID(id)).Exec(ctx); err != nil {
			return err
		}
		return createImages(ctx, tx, id, req.Images)
	})
	if err != nil {
		return nil, err
	}
	return s.QueryPrompt(ctx, id)
}

func (s *PromptServiceImpl) QueryPrompt(ctx context.Context, id int) (*ent.Prompt, error) {
	p, err := withImagesAndTags(s.client.Prompt.Query().Where(prompt.ID(id))).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, model.ErrPromptNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *PromptServiceImpl) DeletePrompt(ctx context.Context, id int) error {
	return withTx(ctx, s.client, func(tx *ent.Tx) error {
		if _, err := tx.PromptImage.Delete().Where(promptimage.PromptID(id)).Exec(ctx); err != nil {
			return err
		}
		if err := tx.Prompt.DeleteOneID(id).Exec(ctx); err != nil {
			if ent.IsNotFound(err) {
				return model.ErrPromptNotFound
			}
			return err
		}
		return nil
	})
}

func (s *PromptServiceImpl) ListPromptTag(ctx context.Context) ([]*ent.PromptTag, error) {
	return s.client.PromptTag.Query().
		Order(ent.Asc(prompttag.FieldName)).
		All(ctx)
}

func validateImages(images []model.PromptImageReq) error {
	for _, img := range images {
		if strings.TrimSpace(img.ImageURL) == "" {
			return model.ErrPromptImageURLEmpty
		}
	}
	return nil
}

func createImages(ctx context.Context, tx *ent.Tx, promptID int, images []model.PromptImageReq) error {
	if len(images) == 0 {
		return nil
	}
	builders := make([]*ent.PromptImageCreate, 0, len(images))
	for _, img := range images {
		builders = append(builders, tx.PromptImage.Create().
			SetPromptID(promptID).
			SetImageURL(strings.TrimSpace(img.ImageURL)).
			SetModelName(strings.TrimSpace(img.ModelName)).
			SetSort(img.Sort))
	}
	return tx.PromptImage.CreateBulk(builders...).Exec(ctx)
}

// ensureTags 按名称查找标签，不存在则创建，返回去重后的标签ID
func ensureTags(ctx context.Context, tx *ent.Tx, names []string) ([]int, error) {
	seen := make(map[string]struct{}, len(names))
	ids := make([]int, 0, len(names))
	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}

		tag, err := tx.PromptTag.Query().Where(prompttag.Name(name)).Only(ctx)
		if ent.IsNotFound(err) {
			tag, err = tx.PromptTag.Create().SetName(name).Save(ctx)
		}
		if err != nil {
			return nil, err
		}
		ids = append(ids, tag.ID)
	}
	return ids, nil
}

func withTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("%w: rollback failed: %v", err, rerr)
		}
		return err
	}
	return tx.Commit()
}
