package notification

import (
	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/notification"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type NotificationService interface {
	QueryNotificationPage(c *fiber.Ctx, req model.NotificationPageQuery) ([]*ent.Notification, int, error)
	QueryNotification(c *fiber.Ctx, id int) (*ent.Notification, error)
	DeleteNotification(c *fiber.Ctx, id int) error
	BatchMarkAsRead(c *fiber.Ctx, ids []int) error
}

type NotificationServiceImpl struct {
	client *ent.Client
}

func NewNotificationServiceImpl(client *ent.Client) *NotificationServiceImpl {
	return &NotificationServiceImpl{client: client}
}

func (s *NotificationServiceImpl) QueryNotificationPage(c *fiber.Ctx, req model.NotificationPageQuery) ([]*ent.Notification, int, error) {
	query := s.client.Notification.Query()

	if req.IsRead != nil {
		query.Where(notification.IsReadEQ(*req.IsRead))
	}

	count, err := query.Count(c.Context())
	if err != nil {
		return nil, 0, err
	}

	notifications, err := query.
		Order(ent.Desc(notification.FieldID)).
		Limit(req.Size).
		Offset((req.Page - 1) * req.Size).
		All(c.Context())
	if err != nil {
		return nil, 0, err
	}
	return notifications, count, nil
}

func (s *NotificationServiceImpl) QueryNotification(c *fiber.Ctx, id int) (*ent.Notification, error) {
	return s.client.Notification.Query().Where(notification.IDEQ(id)).First(c.Context())
}

func (s *NotificationServiceImpl) DeleteNotification(c *fiber.Ctx, id int) error {
	return s.client.Notification.DeleteOneID(id).Exec(c.Context())
}

func (s *NotificationServiceImpl) BatchMarkAsRead(c *fiber.Ctx, ids []int) error {
	_, err := s.client.Notification.Update().
		Where(notification.IDIn(ids...)).
		SetIsRead(true).
		Save(c.Context())
	return err
}
