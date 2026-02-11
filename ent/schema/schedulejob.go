package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// 定时任务
type ScheduleJob struct {
	ent.Schema
}

func (ScheduleJob) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the ScheduleJob.
func (ScheduleJob) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("任务名称"),
		field.String("type").Comment("任务类型: cron, interval"),
		field.String("expression").Comment("调度表达式: cron表达式或时间间隔"),
		field.Text("description").Optional().Comment("任务描述"),
		field.Bool("enabled").Default(true).Comment("是否启用"),
		field.Time("last_run_time").Optional().Comment("上次执行时间"),
		field.String("job_name").Comment("内部任务名称"),
		field.Int("max_retries").Default(3).Comment("最大重试次数"),
		field.Bool("failure_notification").Default(false).Comment("失败是否通知"),
	}
}

// Edges of the ScheduleJob.
func (ScheduleJob) Edges() []ent.Edge {
	return nil
}
