package doclibrarydetail

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/doclibrarydetail"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type DocLibraryDetailService interface {
	CreateDocLibraryDetail(ctx context.Context, req model.DocLibraryDetailCreateReq) (*ent.DocLibraryDetail, error)
	UpdateDocLibraryDetail(ctx context.Context, id int, req model.DocLibraryDetailUpdateReq) error
	GetDocLibraryDetail(ctx context.Context, id int) (*ent.DocLibraryDetail, error)
	GetDocLibraryDetailPage(ctx context.Context, page, pageSize int, libraryID int) ([]*ent.DocLibraryDetail, int, error)
	DeleteDocLibraryDetail(ctx context.Context, id int) error
	GetDocLibraryDetailTree(ctx context.Context, libraryID int) ([]*ent.DocLibraryDetail, error)
}

type DocLibraryDetailServiceImpl struct {
	client *ent.Client
}

func NewDocLibraryDetailServiceImpl(client *ent.Client) DocLibraryDetailService {
	return &DocLibraryDetailServiceImpl{client: client}
}

func (s *DocLibraryDetailServiceImpl) CreateDocLibraryDetail(ctx context.Context, req model.DocLibraryDetailCreateReq) (*ent.DocLibraryDetail, error) {
	detail, err := s.client.DocLibraryDetail.Create().
		SetLibraryID(req.LibraryID).
		SetTitle(req.Title).
		SetVersion(req.Version).
		SetContent(req.Content).
		SetParentID(req.ParentID).
		SetPath(req.Path).
		SetURL(req.URL).
		SetLanguage(req.Language).
		Save(ctx)
	return detail, err
}

func (s *DocLibraryDetailServiceImpl) UpdateDocLibraryDetail(ctx context.Context, id int, req model.DocLibraryDetailUpdateReq) error {
	update := s.client.DocLibraryDetail.UpdateOneID(id).
		SetTitle(req.Title).
		SetVersion(req.Version).
		SetContent(req.Content).
		SetPath(req.Path).
		SetURL(req.URL).
		SetLanguage(req.Language)

	if req.ParentID != 0 {
		update.SetParentID(req.ParentID)
	}

	return update.Exec(ctx)
}

func (s *DocLibraryDetailServiceImpl) GetDocLibraryDetail(ctx context.Context, id int) (*ent.DocLibraryDetail, error) {
	detail, err := s.client.DocLibraryDetail.Query().
		Where(doclibrarydetail.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

func (s *DocLibraryDetailServiceImpl) GetDocLibraryDetailPage(ctx context.Context, page, pageSize int, libraryID int) ([]*ent.DocLibraryDetail, int, error) {
	query := s.client.DocLibraryDetail.Query()
	if libraryID > 0 {
		query = query.Where(doclibrarydetail.LibraryID(libraryID))
	}

	count, err := query.
		Order(ent.Desc(doclibrarydetail.FieldID)).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	details, err := query.
		Order(ent.Desc(doclibrarydetail.FieldID)).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return details, count, nil
}

func (s *DocLibraryDetailServiceImpl) DeleteDocLibraryDetail(ctx context.Context, id int) error {
	return s.client.DocLibraryDetail.DeleteOneID(id).Exec(ctx)
}

func (s *DocLibraryDetailServiceImpl) GetDocLibraryDetailTree(ctx context.Context, libraryID int) ([]*ent.DocLibraryDetail, error) {
	query := s.client.DocLibraryDetail.Query()
	if libraryID > 0 {
		query = query.Where(doclibrarydetail.LibraryID(libraryID))
	}

	details, err := query.
		Order(ent.Asc(doclibrarydetail.FieldTitle)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return details, nil
}
