package manager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/shuTwT/hoshikuzu/ent"
	schedule_model "github.com/shuTwT/hoshikuzu/pkg/domain/model/schedule"
)

type ScheduleManager struct {
	jobCache  map[string]schedule_model.Job
	scheduler gocron.Scheduler
	jobIDMap  map[int]gocron.Job
}

func NewScheduleManager() *ScheduleManager {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(fmt.Sprintf("创建调度器失败: %v", err))
	}
	m := &ScheduleManager{
		jobCache:  map[string]schedule_model.Job{},
		scheduler: scheduler,
		jobIDMap:  make(map[int]gocron.Job),
	}
	return m
}

func (m *ScheduleManager) ClearCache() {
	for key := range m.jobCache {
		delete(m.jobCache, key)
	}
}

func (m *ScheduleManager) AddJobToCache(jobName string, job schedule_model.Job) {
	m.jobCache[jobName] = job
}

func (m *ScheduleManager) Start() {
	m.scheduler.Start()
}

func (m *ScheduleManager) Shutdown() {
	_ = m.scheduler.Shutdown()
}

func (m *ScheduleManager) GetJob(jobName string) (schedule_model.Job, bool) {
	job, ok := m.jobCache[jobName]
	return job, ok
}

func (m *ScheduleManager) AddJobToScheduler(jobEntity *ent.ScheduleJob) error {
	job, ok := m.jobCache[jobEntity.JobName]
	if !ok {
		return fmt.Errorf("找不到任务 '%s' 的实现", jobEntity.JobName)
	}

	var taskJob gocron.Job
	var err error

	switch jobEntity.Type {
	case "interval":
		durationJob, ok := job.(schedule_model.DurationJob)
		if !ok {
			return fmt.Errorf("任务 '%s' 不是 DurationJob 类型", jobEntity.JobName)
		}

		duration, err := time.ParseDuration(jobEntity.Expression)
		if err != nil {
			return fmt.Errorf("任务 '%s' 的表达式 '%s' 解析失败: %w", jobEntity.JobName, jobEntity.Expression, err)
		}

		taskJob, err = m.scheduler.NewJob(
			gocron.DurationJob(duration),
			gocron.NewTask(
				func(ctx context.Context, j schedule_model.DurationJob) {
					if err := j.Execute(ctx); err != nil {
						log.Printf("任务 '%s' 执行失败: %v", jobEntity.JobName, err)
					}
				},
				context.Background(),
				durationJob,
			),
		)

	case "cron":
		cronJob, ok := job.(schedule_model.CronJob)
		if !ok {
			return fmt.Errorf("任务 '%s' 不是 CronJob 类型", jobEntity.JobName)
		}

		taskJob, err = m.scheduler.NewJob(
			gocron.CronJob(jobEntity.Expression, true),
			gocron.NewTask(
				func(ctx context.Context, j schedule_model.CronJob) {
					if err := j.Execute(ctx); err != nil {
						log.Printf("任务 '%s' 执行失败: %v", jobEntity.JobName, err)
					}
				},
				context.Background(),
				cronJob,
			),
		)

	default:
		return fmt.Errorf("任务 '%s' 的类型 '%s' 不支持", jobEntity.JobName, jobEntity.Type)
	}

	if err != nil {
		return fmt.Errorf("添加任务 '%s' 失败: %w", jobEntity.JobName, err)
	}

	m.jobIDMap[jobEntity.ID] = taskJob
	log.Printf("已添加任务: %s (ID: %d, 类型: %s, 表达式: %s)", jobEntity.Name, jobEntity.ID, jobEntity.Type, jobEntity.Expression)
	return nil
}

func (m *ScheduleManager) UpdateJobInScheduler(jobEntity *ent.ScheduleJob) error {
	if err := m.RemoveJobFromScheduler(jobEntity.ID); err != nil {
		return err
	}

	if jobEntity.Enabled {
		return m.AddJobToScheduler(jobEntity)
	}

	return nil
}

func (m *ScheduleManager) RemoveJobFromScheduler(jobID int) error {
	taskJob, ok := m.jobIDMap[jobID]
	if !ok {
		return fmt.Errorf("任务 ID %d 不存在于调度器中", jobID)
	}

	if err := m.scheduler.RemoveJob(taskJob.ID()); err != nil {
		return fmt.Errorf("移除任务 ID %d 失败: %w", jobID, err)
	}

	delete(m.jobIDMap, jobID)
	log.Printf("已移除任务 ID: %d", jobID)
	return nil
}

func (m *ScheduleManager) ExecuteJobNow(jobEntity *ent.ScheduleJob) error {
	job, ok := m.jobCache[jobEntity.JobName]
	if !ok {
		return fmt.Errorf("找不到任务 '%s' 的实现", jobEntity.JobName)
	}

	switch jobEntity.Type {
	case "interval":
		durationJob, ok := job.(schedule_model.DurationJob)
		if !ok {
			return fmt.Errorf("任务 '%s' 不是 DurationJob 类型", jobEntity.JobName)
		}
		if err := durationJob.Execute(context.Background()); err != nil {
			log.Printf("任务 '%s' 执行失败: %v", jobEntity.JobName, err)
			return err
		}
	case "cron":
		cronJob, ok := job.(schedule_model.CronJob)
		if !ok {
			return fmt.Errorf("任务 '%s' 不是 CronJob 类型", jobEntity.JobName)
		}
		if err := cronJob.Execute(context.Background()); err != nil {
			log.Printf("任务 '%s' 执行失败: %v", jobEntity.JobName, err)
			return err
		}
	default:
		return fmt.Errorf("任务 '%s' 的类型 '%s' 不支持", jobEntity.JobName, jobEntity.Type)
	}

	log.Printf("任务 '%s' (ID: %d) 已立即执行", jobEntity.Name, jobEntity.ID)
	return nil
}
