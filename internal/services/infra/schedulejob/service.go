package schedulejob

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/schedulejob"
	schedule_manager "github.com/shuTwT/hoshikuzu/internal/infra/schedule/manager"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type ScheduleJobService interface {
	ListScheduleJob(ctx context.Context) ([]*ent.ScheduleJob, error)
	ListScheduleJobPage(ctx context.Context, page, size int) (int, []*ent.ScheduleJob, error)
	QueryScheduleJob(ctx context.Context, id int) (*ent.ScheduleJob, error)
	CreateScheduleJob(ctx context.Context, req *model.CreateScheduleJobReq) (*ent.ScheduleJob, error)
	UpdateScheduleJob(ctx context.Context, id int, req *model.UpdateScheduleJobReq) (*ent.ScheduleJob, error)
	DeleteScheduleJob(ctx context.Context, id int) error
	ExecuteScheduleJobNow(ctx context.Context, id int) error
}

type ScheduleJobServiceImpl struct {
	client          *ent.Client
	scheduleManager *schedule_manager.ScheduleManager
}

func NewScheduleJobServiceImpl(client *ent.Client, scheduleManager *schedule_manager.ScheduleManager) *ScheduleJobServiceImpl {
	return &ScheduleJobServiceImpl{client: client, scheduleManager: scheduleManager}
}

func (s *ScheduleJobServiceImpl) ListScheduleJob(ctx context.Context) ([]*ent.ScheduleJob, error) {
	jobs, err := s.client.ScheduleJob.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *ScheduleJobServiceImpl) ListScheduleJobPage(ctx context.Context, page, size int) (int, []*ent.ScheduleJob, error) {
	count, err := s.client.ScheduleJob.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	jobs, err := s.client.ScheduleJob.Query().
		Order(ent.Desc(schedulejob.FieldID)).
		Offset((page - 1) * size).
		Limit(size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, jobs, nil
}

func (s *ScheduleJobServiceImpl) QueryScheduleJob(ctx context.Context, id int) (*ent.ScheduleJob, error) {
	job, err := s.client.ScheduleJob.Query().
		Where(schedulejob.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *ScheduleJobServiceImpl) CreateScheduleJob(ctx context.Context, req *model.CreateScheduleJobReq) (*ent.ScheduleJob, error) {
	if err := validateCreateScheduleJobReq(req); err != nil {
		return nil, err
	}

	exists, err := s.client.ScheduleJob.Query().
		Where(schedulejob.JobName(req.JobName)).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("内部任务名称已存在")
	}

	builder := s.client.ScheduleJob.Create().
		SetName(req.Name).
		SetType(req.Type).
		SetExpression(req.Expression).
		SetJobName(req.JobName).
		SetEnabled(req.Enabled).
		SetMaxRetries(req.MaxRetries).
		SetFailureNotification(req.FailureNotification)

	if req.Description != nil {
		builder.SetDescription(*req.Description)
	}

	job, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	if job.Enabled {
		if err := s.scheduleManager.AddJobToScheduler(job); err != nil {
			return nil, err
		}
	}

	return job, nil
}

func (s *ScheduleJobServiceImpl) UpdateScheduleJob(ctx context.Context, id int, req *model.UpdateScheduleJobReq) (*ent.ScheduleJob, error) {
	if err := validateUpdateScheduleJobReq(req); err != nil {
		return nil, err
	}

	exists, err := s.client.ScheduleJob.Query().Where(schedulejob.ID(id)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("定时任务不存在")
	}

	if req.JobName != nil {
		jobNameExists, err := s.client.ScheduleJob.Query().
			Where(
				schedulejob.JobName(*req.JobName),
				schedulejob.IDNEQ(id),
			).
			Exist(ctx)
		if err != nil {
			return nil, err
		}
		if jobNameExists {
			return nil, errors.New("内部任务名称已存在")
		}
	}

	builder := s.client.ScheduleJob.UpdateOneID(id)

	if req.Name != nil {
		builder.SetName(*req.Name)
	}

	if req.Type != nil {
		builder.SetType(*req.Type)
	}

	if req.Expression != nil {
		builder.SetExpression(*req.Expression)
	}

	if req.Description != nil {
		builder.SetDescription(*req.Description)
	}

	if req.Enabled != nil {
		builder.SetEnabled(*req.Enabled)
	}

	if req.JobName != nil {
		builder.SetJobName(*req.JobName)
	}

	if req.MaxRetries != nil {
		builder.SetMaxRetries(*req.MaxRetries)
	}

	if req.FailureNotification != nil {
		builder.SetFailureNotification(*req.FailureNotification)
	}

	job, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.scheduleManager.UpdateJobInScheduler(job); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *ScheduleJobServiceImpl) DeleteScheduleJob(ctx context.Context, id int) error {
	err := s.client.ScheduleJob.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	if err := s.scheduleManager.RemoveJobFromScheduler(id); err != nil {
		return err
	}

	return nil
}

func (s *ScheduleJobServiceImpl) ExecuteScheduleJobNow(ctx context.Context, id int) error {
	job, err := s.client.ScheduleJob.Query().
		Where(schedulejob.ID(id)).
		First(ctx)
	if err != nil {
		return err
	}

	go func() {
		if err := s.scheduleManager.ExecuteJobNow(job); err != nil {
			log.Printf("异步执行任务失败: %v", err)
		} else {
			now := time.Now()
			_, err := s.client.ScheduleJob.UpdateOneID(id).
				SetLastRunTime(now).
				Save(context.Background())
			if err != nil {
				log.Printf("更新任务执行时间失败: %v", err)
			}
		}
	}()

	return nil
}

func validateCreateScheduleJobReq(req *model.CreateScheduleJobReq) error {
	if req.Name == "" {
		return errors.New("任务名称不能为空")
	}
	if req.Type == "" {
		return errors.New("任务类型不能为空")
	}
	if req.Expression == "" {
		return errors.New("调度表达式不能为空")
	}
	if req.JobName == "" {
		return errors.New("内部任务名称不能为空")
	}
	return nil
}

func validateUpdateScheduleJobReq(_ *model.UpdateScheduleJobReq) error {
	return nil
}
