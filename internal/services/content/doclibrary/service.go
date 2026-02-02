package doclibrary

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/doclibrary"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type DocLibraryService interface {
	CreateDocLibrary(ctx context.Context, req *model.DocLibraryCreateReq) (*ent.DocLibrary, error)
	UpdateDocLibrary(ctx context.Context, id int, req *model.DocLibraryUpdateReq) error
	GetDocLibrary(ctx context.Context, id int) (*ent.DocLibrary, error)
	GetDocLibraryPage(ctx context.Context, page, pageSize int) ([]*ent.DocLibrary, int, error)
	DeleteDocLibrary(ctx context.Context, id int) error
	GetDocLibraryList(ctx context.Context) ([]*ent.DocLibrary, error)
}

type DocLibraryServiceImpl struct {
	client *ent.Client
}

func NewDocLibraryServiceImpl(client *ent.Client) DocLibraryService {
	return &DocLibraryServiceImpl{client: client}
}

func (s *DocLibraryServiceImpl) CreateDocLibrary(ctx context.Context, req *model.DocLibraryCreateReq) (*ent.DocLibrary, error) {
	library := s.client.DocLibrary.Create().
		SetName(req.Name).
		SetAlias(req.Alias).
		SetDescription(req.Description).
		SetSource(doclibrary.Source(req.Source)).
		SetURL(req.URL).
		SaveX(ctx)
	return library, nil
}

func (s *DocLibraryServiceImpl) UpdateDocLibrary(ctx context.Context, id int, req *model.DocLibraryUpdateReq) error {
	return s.client.DocLibrary.UpdateOneID(id).
		SetName(req.Name).
		SetAlias(req.Alias).
		SetDescription(req.Description).
		SetSource(doclibrary.Source(req.Source)).
		SetURL(req.URL).
		Exec(ctx)
}

func (s *DocLibraryServiceImpl) GetDocLibrary(ctx context.Context, id int) (*ent.DocLibrary, error) {
	library, err := s.client.DocLibrary.Query().
		Where(doclibrary.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return library, nil
}

func (s *DocLibraryServiceImpl) GetDocLibraryPage(ctx context.Context, page, pageSize int) ([]*ent.DocLibrary, int, error) {
	count, err := s.client.DocLibrary.Query().
		Order(ent.Desc(doclibrary.FieldID)).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	libraries, err := s.client.DocLibrary.Query().
		Order(ent.Desc(doclibrary.FieldID)).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return libraries, count, nil
}

func (s *DocLibraryServiceImpl) DeleteDocLibrary(ctx context.Context, id int) error {
	return s.client.DocLibrary.DeleteOneID(id).Exec(ctx)
}

func (s *DocLibraryServiceImpl) GetDocLibraryList(ctx context.Context) ([]*ent.DocLibrary, error) {
	libraries, err := s.client.DocLibrary.Query().
		Order(ent.Asc(doclibrary.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return libraries, nil
}
