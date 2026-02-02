package friendcircle

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/friendcirclerecord"
)

type FriendCircleService interface {
	ListFriendCircleRecord(ctx context.Context) ([]*ent.FriendCircleRecord, error)
	ListFriendCircleRecordPage(ctx context.Context, page, size int) (int, []*ent.FriendCircleRecord, error)
	CreateFriendCircleRecord(ctx context.Context, author, title, linkURL, avatarURL, siteURL, publishedAt string) (*ent.FriendCircleRecord, error)
	UpdateFriendCircleRecord(ctx context.Context, id int, author, title, linkURL, avatarURL, siteURL, publishedAt string) (*ent.FriendCircleRecord, error)
	QueryFriendCircleRecord(ctx context.Context, id int) (*ent.FriendCircleRecord, error)
	DeleteFriendCircleRecord(ctx context.Context, id int) error
	ExistsRecord(c context.Context, linkUrl string) (bool, error)
	InsertRecord(c context.Context, author string, avatarUrl string, title string, linkUrl string, publishedAt string) error
}

type FriendCircleServiceImpl struct {
	client *ent.Client
}

func NewFriendCircleServiceImpl(client *ent.Client) *FriendCircleServiceImpl {
	return &FriendCircleServiceImpl{client: client}
}

func (s *FriendCircleServiceImpl) ListFriendCircleRecord(ctx context.Context) ([]*ent.FriendCircleRecord, error) {
	records, err := s.client.FriendCircleRecord.Query().
		Order(ent.Desc(friendcirclerecord.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *FriendCircleServiceImpl) ListFriendCircleRecordPage(ctx context.Context, page, size int) (int, []*ent.FriendCircleRecord, error) {
	count, err := s.client.FriendCircleRecord.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	records, err := s.client.FriendCircleRecord.Query().
		Order(ent.Desc(friendcirclerecord.FieldID)).
		Limit(size).
		Offset((page - 1) * size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, records, nil
}

func (s *FriendCircleServiceImpl) CreateFriendCircleRecord(ctx context.Context, author, title, linkURL, avatarURL, siteURL, publishedAt string) (*ent.FriendCircleRecord, error) {
	builder := s.client.FriendCircleRecord.Create().
		SetAuthor(author).
		SetTitle(title).
		SetLinkURL(linkURL).
		SetAvatarURL(avatarURL)

	if siteURL != "" {
		builder.SetSiteURL(siteURL)
	}
	if publishedAt != "" {
		builder.SetPublishedAt(publishedAt)
	}

	newRecord, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return newRecord, nil
}

func (s *FriendCircleServiceImpl) UpdateFriendCircleRecord(ctx context.Context, id int, author, title, linkURL, avatarURL, siteURL, publishedAt string) (*ent.FriendCircleRecord, error) {
	builder := s.client.FriendCircleRecord.UpdateOneID(id).
		SetAuthor(author).
		SetTitle(title).
		SetLinkURL(linkURL).
		SetAvatarURL(avatarURL)

	if siteURL != "" {
		builder.SetSiteURL(siteURL)
	} else {
		builder.ClearSiteURL()
	}
	if publishedAt != "" {
		builder.SetPublishedAt(publishedAt)
	} else {
		builder.ClearPublishedAt()
	}

	updatedRecord, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedRecord, nil
}

func (s *FriendCircleServiceImpl) QueryFriendCircleRecord(ctx context.Context, id int) (*ent.FriendCircleRecord, error) {
	record, err := s.client.FriendCircleRecord.Query().
		Where(friendcirclerecord.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *FriendCircleServiceImpl) DeleteFriendCircleRecord(ctx context.Context, id int) error {
	err := s.client.FriendCircleRecord.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 判断是否已存在
func (s *FriendCircleServiceImpl) ExistsRecord(c context.Context, linkUrl string) (bool, error) {
	client := s.client
	return client.FriendCircleRecord.Query().Where(friendcirclerecord.LinkURLEQ(linkUrl)).Exist(c)
}

// 插入文章
func (s *FriendCircleServiceImpl) InsertRecord(c context.Context, author string, avatarUrl string, title string, linkUrl string, publishedAt string) error {
	client := s.client
	_, err := client.FriendCircleRecord.Create().
		SetAuthor(author).
		SetAvatarURL(avatarUrl).
		SetTitle(title).
		SetLinkURL(linkUrl).
		SetPublishedAt(publishedAt).
		Save(c)
	return err
}
