package schedulejob

import (
	"log/slog"
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/services/infra/schedulejob"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type ScheduleJobHandler interface {
	CreateScheduleJob(c *fiber.Ctx) error
	ListScheduleJobPage(c *fiber.Ctx) error
	QueryScheduleJob(c *fiber.Ctx) error
	UpdateScheduleJob(c *fiber.Ctx) error
	DeleteScheduleJob(c *fiber.Ctx) error
	ExecuteScheduleJobNow(c *fiber.Ctx) error
}

type ScheduleJobHandlerImpl struct {
	scheduleJobService schedulejob.ScheduleJobService
}

func NewScheduleJobHandlerImpl(scheduleJobService schedulejob.ScheduleJobService) *ScheduleJobHandlerImpl {
	return &ScheduleJobHandlerImpl{scheduleJobService: scheduleJobService}
}

func (h *ScheduleJobHandlerImpl) CreateScheduleJob(c *fiber.Ctx) error {
	var req model.CreateScheduleJobReq
	if err := c.BodyParser(&req); err != nil {
		slog.Error("Failed to parse request body", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "请求参数解析失败"))
	}

	job, err := h.scheduleJobService.CreateScheduleJob(c.Context(), &req)
	if err != nil {
		slog.Error("Failed to create schedule job", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	resp := h.buildScheduleJobResp(job)
	slog.Info("Schedule job created successfully", "job_id", job.ID, "job_name", job.Name)
	return c.JSON(model.NewSuccess("定时任务创建成功", resp))
}

func (h *ScheduleJobHandlerImpl) ListScheduleJobPage(c *fiber.Ctx) error {
	var pageQuery model.PageQuery
	if err := c.QueryParser(&pageQuery); err != nil {
		slog.Error("Failed to parse query parameters", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "查询参数解析失败"))
	}

	count, jobs, err := h.scheduleJobService.ListScheduleJobPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		slog.Error("Failed to list schedule jobs", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	jobResps := make([]*model.ScheduleJobResp, 0, len(jobs))
	for _, job := range jobs {
		jobResps = append(jobResps, h.buildScheduleJobResp(job))
	}

	pageResult := model.PageResult[*model.ScheduleJobResp]{
		Total:   int64(count),
		Records: jobResps,
	}
	return c.JSON(model.NewSuccess("定时任务列表获取成功", pageResult))
}

func (h *ScheduleJobHandlerImpl) QueryScheduleJob(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid job ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的任务ID"))
	}

	job, err := h.scheduleJobService.QueryScheduleJob(c.Context(), id)
	if err != nil {
		slog.Error("Failed to query schedule job", "job_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	resp := h.buildScheduleJobResp(job)
	return c.JSON(model.NewSuccess("定时任务查询成功", resp))
}

func (h *ScheduleJobHandlerImpl) UpdateScheduleJob(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid job ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的任务ID"))
	}

	var req model.UpdateScheduleJobReq
	if err := c.BodyParser(&req); err != nil {
		slog.Error("Failed to parse request body", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "请求参数解析失败"))
	}

	job, err := h.scheduleJobService.UpdateScheduleJob(c.Context(), id, &req)
	if err != nil {
		slog.Error("Failed to update schedule job", "job_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	resp := h.buildScheduleJobResp(job)
	slog.Info("Schedule job updated successfully", "job_id", job.ID, "job_name", job.Name)
	return c.JSON(model.NewSuccess("定时任务更新成功", resp))
}

func (h *ScheduleJobHandlerImpl) DeleteScheduleJob(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid job ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的任务ID"))
	}

	err = h.scheduleJobService.DeleteScheduleJob(c.Context(), id)
	if err != nil {
		slog.Error("Failed to delete schedule job", "job_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Schedule job deleted successfully", "job_id", id)
	return c.JSON(model.NewSuccess("定时任务删除成功", nil))
}

func (h *ScheduleJobHandlerImpl) ExecuteScheduleJobNow(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid job ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的任务ID"))
	}

	err = h.scheduleJobService.ExecuteScheduleJobNow(c.Context(), id)
	if err != nil {
		slog.Error("Failed to execute schedule job", "job_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Schedule job executed successfully", "job_id", id)
	return c.JSON(model.NewSuccess("定时任务执行成功", nil))
}

func (*ScheduleJobHandlerImpl) buildScheduleJobResp(job *ent.ScheduleJob) *model.ScheduleJobResp {
	return &model.ScheduleJobResp{
		ID:                  job.ID,
		CreatedAt:           job.CreatedAt,
		UpdatedAt:           job.UpdatedAt,
		Name:                job.Name,
		Type:                job.Type,
		Expression:          job.Expression,
		Description:         job.Description,
		Enabled:             job.Enabled,
		NextRunTime:         job.NextRunTime,
		LastRunTime:         job.LastRunTime,
		JobName:             job.JobName,
		MaxRetries:          job.MaxRetries,
		FailureNotification: job.FailureNotification,
	}
}
