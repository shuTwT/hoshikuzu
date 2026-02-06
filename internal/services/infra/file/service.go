package file

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/file"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type FileService interface {
	ListFile(ctx context.Context) ([]*ent.File, error)
	ListFilePage(ctx context.Context, page, size int) (int, []*ent.File, error)
	ListFilePageWithQuery(ctx context.Context, req model.FilePageReq) (int, []*ent.File, error)
	QueryFile(ctx context.Context, id int) (*ent.File, error)
	DeleteFile(ctx context.Context, id int) error
	CreateFile(ctx context.Context, strategyID int, name, path, url, fileType, size string) (*ent.File, error)
}

type FileServiceImpl struct {
	client *ent.Client
}

func NewFileServiceImpl(client *ent.Client) *FileServiceImpl {
	return &FileServiceImpl{client: client}
}

func (s *FileServiceImpl) ListFile(ctx context.Context) ([]*ent.File, error) {
	files, err := s.client.File.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *FileServiceImpl) ListFilePage(ctx context.Context, page, size int) (int, []*ent.File, error) {
	count, err := s.client.File.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	files, err := s.client.File.Query().
		WithStorageStrategy().
		Order(ent.Desc(file.FieldID)).
		Offset((page - 1) * size).
		Limit(size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, files, nil
}

func (s *FileServiceImpl) ListFilePageWithQuery(ctx context.Context, req model.FilePageReq) (int, []*ent.File, error) {
	query := s.client.File.Query()

	if req.Name != "" {
		query.Where(file.NameContains(req.Name))
	}

	if req.Type != "" {
		query.Where(file.Type(req.Type))
	}

	if req.StorageStrategyID != nil {
		query.Where(file.StorageStrategyID(*req.StorageStrategyID))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	files, err := query.
		WithStorageStrategy().
		Order(ent.Desc(file.FieldID)).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, files, nil
}

func (s *FileServiceImpl) QueryFile(ctx context.Context, id int) (*ent.File, error) {
	file, err := s.client.File.Query().
		WithStorageStrategy().
		Where(file.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *FileServiceImpl) DeleteFile(ctx context.Context, id int) error {
	err := s.client.File.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileServiceImpl) CreateFile(ctx context.Context, strategyID int, name, path, url, fileType, size string) (*ent.File, error) {
	newFile, err := s.client.File.Create().
		SetName(name).
		SetPath(path).
		SetURL(url).
		SetType(fileType).
		SetSize(size).
		SetStorageStrategyID(strategyID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return newFile, nil
}
