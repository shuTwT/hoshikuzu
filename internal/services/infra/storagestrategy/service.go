package storagestrategy

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/storagestrategy"
)

type StorageStrategyService interface {
	ListStorageStrategy(ctx context.Context) ([]*ent.StorageStrategy, error)
	ListStorageStrategyPage(ctx context.Context, page, size int) (int, []*ent.StorageStrategy, error)
	CreateStorageStrategy(ctx context.Context, name, strategyType, nodeID, endpoint, region, bucket, accessKey, secretKey, basePath, domain string, master bool) (*ent.StorageStrategy, error)
	UpdateStorageStrategy(ctx context.Context, id int, name, strategyType, nodeID, endpoint, region, bucket, accessKey, secretKey, basePath, domain string, master bool) (*ent.StorageStrategy, error)
	QueryStorageStrategy(ctx context.Context, id int) (*ent.StorageStrategy, error)
	DeleteStorageStrategy(ctx context.Context, id int) error
	SetDefaultStorageStrategy(ctx context.Context, id int) error
	GetMasterStorageStrategy(ctx context.Context) (*ent.StorageStrategy, error)
	GetStorageStrategyByID(ctx context.Context, id int) (*ent.StorageStrategy, error)
}

type StorageStrategyServiceImpl struct {
	client *ent.Client
}

func NewStorageStrategyServiceImpl(client *ent.Client) *StorageStrategyServiceImpl {
	return &StorageStrategyServiceImpl{client: client}
}

func (s *StorageStrategyServiceImpl) ListStorageStrategy(ctx context.Context) ([]*ent.StorageStrategy, error) {
	strategies, err := s.client.StorageStrategy.Query().
		Order(ent.Desc(storagestrategy.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return strategies, nil
}

func (s *StorageStrategyServiceImpl) ListStorageStrategyPage(ctx context.Context, page, size int) (int, []*ent.StorageStrategy, error) {
	count, err := s.client.StorageStrategy.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	strategies, err := s.client.StorageStrategy.Query().
		Order(ent.Desc(storagestrategy.FieldID)).
		Limit(size).
		Offset((page - 1) * size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, strategies, nil
}

func (s *StorageStrategyServiceImpl) CreateStorageStrategy(ctx context.Context, name, strategyType, nodeID, endpoint, region, bucket, accessKey, secretKey, basePath, domain string, master bool) (*ent.StorageStrategy, error) {
	newStrategy, err := s.client.StorageStrategy.Create().
		SetName(name).
		SetType(storagestrategy.Type(strategyType)).
		SetNodeID(nodeID).
		SetEndpoint(endpoint).
		SetRegion(region).
		SetBucket(bucket).
		SetAccessKey(accessKey).
		SetSecretKey(secretKey).
		SetBasePath(basePath).
		SetDomain(domain).
		SetMaster(master).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return newStrategy, nil
}

func (s *StorageStrategyServiceImpl) UpdateStorageStrategy(ctx context.Context, id int, name, strategyType, nodeID, endpoint, region, bucket, accessKey, secretKey, basePath, domain string, master bool) (*ent.StorageStrategy, error) {
	updatedStrategy, err := s.client.StorageStrategy.UpdateOneID(id).
		SetName(name).
		SetType(storagestrategy.Type(strategyType)).
		SetNodeID(nodeID).
		SetEndpoint(endpoint).
		SetRegion(region).
		SetBucket(bucket).
		SetAccessKey(accessKey).
		SetSecretKey(secretKey).
		SetBasePath(basePath).
		SetDomain(domain).
		SetMaster(master).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedStrategy, nil
}

func (s *StorageStrategyServiceImpl) QueryStorageStrategy(ctx context.Context, id int) (*ent.StorageStrategy, error) {
	strategy, err := s.client.StorageStrategy.Query().
		Where(storagestrategy.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return strategy, nil
}

func (s *StorageStrategyServiceImpl) DeleteStorageStrategy(ctx context.Context, id int) error {
	err := s.client.StorageStrategy.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageStrategyServiceImpl) SetDefaultStorageStrategy(ctx context.Context, id int) error {
	err := s.client.StorageStrategy.Update().SetMaster(false).Exec(ctx)
	if err != nil {
		return err
	}

	err = s.client.StorageStrategy.UpdateOneID(id).SetMaster(true).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetMasterStorageStrategy 获取主存储策略
func (s *StorageStrategyServiceImpl) GetMasterStorageStrategy(ctx context.Context) (*ent.StorageStrategy, error) {
	strategy, err := s.client.StorageStrategy.Query().
		Where(storagestrategy.Master(true)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return strategy, nil
}

// GetStorageStrategyByID 根据ID获取存储策略
func (s *StorageStrategyServiceImpl) GetStorageStrategyByID(ctx context.Context, id int) (*ent.StorageStrategy, error) {
	strategy, err := s.client.StorageStrategy.Query().
		Where(storagestrategy.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return strategy, nil
}
