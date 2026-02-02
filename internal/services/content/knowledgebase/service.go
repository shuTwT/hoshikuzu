package knowledgebase

import (
	"context"
	"fmt"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/knowledgebase"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type KnowledgeBaseService interface {
	CreateKnowledgeBase(ctx context.Context, req model.KnowledgeBaseCreateReq) (*ent.KnowledgeBase, error)
	UpdateKnowledgeBase(ctx context.Context, id int, req model.KnowledgeBaseUpdateReq) error
	GetKnowledgeBase(ctx context.Context, id int) (*ent.KnowledgeBase, error)
	GetKnowledgeBasePage(ctx context.Context, req model.KnowledgeBaseQueryReq) ([]*ent.KnowledgeBase, int, error)
	DeleteKnowledgeBase(ctx context.Context, id int) error
	GetKnowledgeBaseList(ctx context.Context) ([]*ent.KnowledgeBase, error)
}

type KnowledgeBaseServiceImpl struct {
	client *ent.Client
}

func NewKnowledgeBaseServiceImpl(client *ent.Client) KnowledgeBaseService {
	return &KnowledgeBaseServiceImpl{client: client}
}

func (s *KnowledgeBaseServiceImpl) CreateKnowledgeBase(ctx context.Context, req model.KnowledgeBaseCreateReq) (*ent.KnowledgeBase, error) {
	kb := s.client.KnowledgeBase.Create().
		SetName(req.Name).
		SetModelProvider(knowledgebase.ModelProvider(req.ModelProvider)).
		SetModel(req.Model).
		SetVectorDimension(req.VectorDimension).
		SetMaxBatchDocumentCount(req.MaxBatchDocumentCount)

	createdKb, err := kb.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create knowledge base: %w", err)
	}

	return createdKb, nil
}

func (s *KnowledgeBaseServiceImpl) UpdateKnowledgeBase(ctx context.Context, id int, req model.KnowledgeBaseUpdateReq) error {
	_, err := s.client.KnowledgeBase.UpdateOneID(id).
		SetName(req.Name).
		SetModelProvider(knowledgebase.ModelProvider(req.ModelProvider)).
		SetModel(req.Model).
		SetVectorDimension(req.VectorDimension).
		SetMaxBatchDocumentCount(req.MaxBatchDocumentCount).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update knowledge base: %w", err)
	}

	return nil
}

func (s *KnowledgeBaseServiceImpl) GetKnowledgeBase(ctx context.Context, id int) (*ent.KnowledgeBase, error) {
	kb, err := s.client.KnowledgeBase.Query().
		Where(knowledgebase.ID(id)).
		First(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}

	return kb, nil
}

func (s *KnowledgeBaseServiceImpl) GetKnowledgeBasePage(ctx context.Context, req model.KnowledgeBaseQueryReq) ([]*ent.KnowledgeBase, int, error) {
	query := s.client.KnowledgeBase.Query()

	if req.Name != "" {
		query = query.Where(knowledgebase.NameContains(req.Name))
	}

	if req.ModelProvider != "" {
		query = query.Where(knowledgebase.ModelProviderEQ(knowledgebase.ModelProvider(req.ModelProvider)))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count knowledge bases: %w", err)
	}

	kbs, err := query.
		Order(ent.Desc(knowledgebase.FieldID)).
		Limit(req.Size).
		Offset((req.Page - 1) * req.Size).
		All(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get knowledge bases: %w", err)
	}

	return kbs, count, nil
}

func (s *KnowledgeBaseServiceImpl) DeleteKnowledgeBase(ctx context.Context, id int) error {
	err := s.client.KnowledgeBase.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete knowledge base: %w", err)
	}
	return nil
}

func (s *KnowledgeBaseServiceImpl) GetKnowledgeBaseList(ctx context.Context) ([]*ent.KnowledgeBase, error) {
	kbs, err := s.client.KnowledgeBase.Query().
		Order(ent.Asc(knowledgebase.FieldName)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge bases: %w", err)
	}

	return kbs, nil
}
