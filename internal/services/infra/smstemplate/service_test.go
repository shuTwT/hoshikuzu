package smstemplate

import (
	"context"
	"errors"
	"testing"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/enttest"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	_ "github.com/mattn/go-sqlite3"
)

func newTestClient(t *testing.T) *ent.Client {
	t.Helper()
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() { client.Close() })
	return client
}

func TestCreateSmsTemplate(t *testing.T) {
	client := newTestClient(t)
	svc := NewSmsTemplateServiceImpl(client)

	ctx := context.Background()
	req := model.SmsTemplateCreateReq{
		Name:     "注册验证码",
		Code:     "register_verify_code",
		Content:  "您的验证码为{code}，{minutes}分钟内有效",
		SignName: ptr("星屑"),
		Enabled:  ptr(true),
	}

	tmpl, err := svc.CreateSmsTemplate(ctx, req)
	if err != nil {
		t.Fatalf("CreateSmsTemplate failed: %v", err)
	}
	if tmpl.Code != req.Code {
		t.Errorf("code = %s, want %s", tmpl.Code, req.Code)
	}
	if tmpl.SignName != "星屑" {
		t.Errorf("sign_name = %q, want 星屑", tmpl.SignName)
	}
	if !tmpl.Enabled {
		t.Errorf("enabled = false, want true")
	}
}

func TestCreateSmsTemplateDuplicateCode(t *testing.T) {
	client := newTestClient(t)
	svc := NewSmsTemplateServiceImpl(client)

	ctx := context.Background()
	req := model.SmsTemplateCreateReq{
		Name:    "注册验证码",
		Code:    "register_verify_code",
		Content: "您的验证码为{code}",
	}
	if _, err := svc.CreateSmsTemplate(ctx, req); err != nil {
		t.Fatalf("first CreateSmsTemplate failed: %v", err)
	}

	_, err := svc.CreateSmsTemplate(ctx, model.SmsTemplateCreateReq{
		Name:    "注册验证码2",
		Code:    "register_verify_code",
		Content: "重复编码",
	})
	if !errors.Is(err, model.ErrSmsTemplateCodeExists) {
		t.Fatalf("err = %v, want ErrSmsTemplateCodeExists", err)
	}
}

func TestUpdateSmsTemplate(t *testing.T) {
	client := newTestClient(t)
	svc := NewSmsTemplateServiceImpl(client)

	ctx := context.Background()
	tmpl, err := svc.CreateSmsTemplate(ctx, model.SmsTemplateCreateReq{
		Name:    "注册验证码",
		Code:    "register_verify_code",
		Content: "您的验证码为{code}",
	})
	if err != nil {
		t.Fatalf("CreateSmsTemplate failed: %v", err)
	}

	newName := "注册验证码(更新)"
	newContent := "您的验证码为{code}，{minutes}分钟内有效"
	updated, err := svc.UpdateSmsTemplate(ctx, tmpl.ID, model.SmsTemplateUpdateReq{
		Name:    &newName,
		Content: &newContent,
	})
	if err != nil {
		t.Fatalf("UpdateSmsTemplate failed: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("name = %s, want %s", updated.Name, newName)
	}
	if updated.Content != newContent {
		t.Errorf("content = %s, want %s", updated.Content, newContent)
	}
	if updated.Code != tmpl.Code {
		t.Errorf("code = %s, want %s（未更新的字段应保持不变）", updated.Code, tmpl.Code)
	}
}

func TestUpdateSmsTemplateDuplicateCode(t *testing.T) {
	client := newTestClient(t)
	svc := NewSmsTemplateServiceImpl(client)

	ctx := context.Background()
	tmplA, err := svc.CreateSmsTemplate(ctx, model.SmsTemplateCreateReq{
		Name:    "模板A",
		Code:    "code_a",
		Content: "内容A",
	})
	if err != nil {
		t.Fatalf("create template A failed: %v", err)
	}
	if _, err := svc.CreateSmsTemplate(ctx, model.SmsTemplateCreateReq{
		Name:    "模板B",
		Code:    "code_b",
		Content: "内容B",
	}); err != nil {
		t.Fatalf("create template B failed: %v", err)
	}

	dupCode := "code_b"
	_, err = svc.UpdateSmsTemplate(ctx, tmplA.ID, model.SmsTemplateUpdateReq{
		Code: &dupCode,
	})
	if !errors.Is(err, model.ErrSmsTemplateCodeExists) {
		t.Fatalf("err = %v, want ErrSmsTemplateCodeExists", err)
	}

	// 更新为自身已有编码不应报错
	ownCode := "code_a"
	if _, err := svc.UpdateSmsTemplate(ctx, tmplA.ID, model.SmsTemplateUpdateReq{
		Code: &ownCode,
	}); err != nil {
		t.Fatalf("update to own code failed: %v", err)
	}
}

func TestQueryAndDeleteSmsTemplate(t *testing.T) {
	client := newTestClient(t)
	svc := NewSmsTemplateServiceImpl(client)

	ctx := context.Background()
	tmpl, err := svc.CreateSmsTemplate(ctx, model.SmsTemplateCreateReq{
		Name:    "注册验证码",
		Code:    "register_verify_code",
		Content: "您的验证码为{code}",
	})
	if err != nil {
		t.Fatalf("CreateSmsTemplate failed: %v", err)
	}

	// 按 ID 查询
	got, err := svc.QuerySmsTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("QuerySmsTemplate failed: %v", err)
	}
	if got.ID != tmpl.ID {
		t.Errorf("id = %d, want %d", got.ID, tmpl.ID)
	}

	// 查询不存在的模板
	if _, err := svc.QuerySmsTemplate(ctx, 99999); !errors.Is(err, model.ErrSmsTemplateNotFound) {
		t.Fatalf("err = %v, want ErrSmsTemplateNotFound", err)
	}

	// 分页查询
	count, list, err := svc.QuerySmsTemplatePage(ctx, model.SmsTemplatePageQuery{Page: 1, Size: 10})
	if err != nil {
		t.Fatalf("QuerySmsTemplatePage failed: %v", err)
	}
	if count != 1 || len(list) != 1 {
		t.Errorf("count = %d, len = %d, want 1, 1", count, len(list))
	}

	// 按名称模糊搜索
	count, _, err = svc.QuerySmsTemplatePage(ctx, model.SmsTemplatePageQuery{Page: 1, Size: 10, Name: "注册"})
	if err != nil {
		t.Fatalf("QuerySmsTemplatePage with name filter failed: %v", err)
	}
	if count != 1 {
		t.Errorf("filtered count = %d, want 1", count)
	}

	// 按编码模糊搜索
	count, _, err = svc.QuerySmsTemplatePage(ctx, model.SmsTemplatePageQuery{Page: 1, Size: 10, Code: "register"})
	if err != nil {
		t.Fatalf("QuerySmsTemplatePage with code filter failed: %v", err)
	}
	if count != 1 {
		t.Errorf("filtered count = %d, want 1", count)
	}

	// 删除
	if err := svc.DeleteSmsTemplate(ctx, tmpl.ID); err != nil {
		t.Fatalf("DeleteSmsTemplate failed: %v", err)
	}
	if err := svc.DeleteSmsTemplate(ctx, tmpl.ID); !errors.Is(err, model.ErrSmsTemplateNotFound) {
		t.Fatalf("second delete err = %v, want ErrSmsTemplateNotFound", err)
	}
}

func ptr[T any](v T) *T {
	return &v
}
