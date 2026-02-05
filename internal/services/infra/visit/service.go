package visit

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/visitlog"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/user_agent"
)

type VisitService interface {
	CreateVisitLog(ctx context.Context, ip string, userAgent []byte, req model.VisitLogReq) error
	QueryVisitLogPage(c *fiber.Ctx, req model.VisitLogPageQuery) ([]*ent.VisitLog, int, error)
	QueryVisitLog(c *fiber.Ctx, id int) (*ent.VisitLog, error)
	DeleteVisitLog(c *fiber.Ctx, id int) error
	BatchDeleteVisitLog(c *fiber.Ctx, ids []int) error
}

type VisitServiceImpl struct {
	client *ent.Client
}

func NewVisitServiceImpl(client *ent.Client) VisitService {
	return &VisitServiceImpl{client: client}
}

func (s *VisitServiceImpl) CreateVisitLog(ctx context.Context, ip string, userAgent []byte, req model.VisitLogReq) error {
	uaString := string(userAgent)
	ua := user_agent.New(uaString)
	browserName, browserVersion := ua.Browser()
	var device string
	if ua.Mobile() {
		device = "Mobile"
	} else {
		device = "Desktop"
	}
	_, err := s.client.VisitLog.Create().
		SetIP(ip).
		SetUserAgent(uaString).
		SetPath(req.UrlPath).
		SetOs(ua.OS()).
		SetBrowser(browserName + " " + browserVersion).
		SetDevice(device).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *VisitServiceImpl) QueryVisitLogPage(c *fiber.Ctx, req model.VisitLogPageQuery) ([]*ent.VisitLog, int, error) {
	query := s.client.VisitLog.Query()

	if req.IP != "" {
		query.Where(visitlog.IPContains(req.IP))
	}

	if req.Path != "" {
		query.Where(visitlog.PathContains(req.Path))
	}

	count, err := query.Count(c.Context())
	if err != nil {
		return nil, 0, err
	}

	visitLogs, err := query.
		Order(ent.Desc(visitlog.FieldID)).
		Limit(req.Size).
		Offset((req.Page - 1) * req.Size).
		All(c.Context())
	if err != nil {
		return nil, 0, err
	}
	return visitLogs, count, nil
}

func (s *VisitServiceImpl) QueryVisitLog(c *fiber.Ctx, id int) (*ent.VisitLog, error) {
	return s.client.VisitLog.Query().Where(visitlog.IDEQ(id)).First(c.Context())
}

func (s *VisitServiceImpl) DeleteVisitLog(c *fiber.Ctx, id int) error {
	return s.client.VisitLog.DeleteOneID(id).Exec(c.Context())
}

func (s *VisitServiceImpl) BatchDeleteVisitLog(c *fiber.Ctx, ids []int) error {
	_, err := s.client.VisitLog.Delete().
		Where(visitlog.IDIn(ids...)).
		Exec(c.Context())
	return err
}
