package mailtemplate

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

func TestCreateMailTemplate(t *testing.T) {
	client := newTestClient(t)
	svc := NewMailTemplateServiceImpl(client)

	ctx := context.Background()
	req := model.MailTemplateCreateReq{
		Name:    "注册欢迎邮件",
		Code:    "register_welcome",
		Subject: "欢迎加入星屑，{nickname}",
		Content: "<h1>你好 {nickname}</h1><p>感谢你的注册。</p>",
		Remark:  ptr("注册场景使用"),
	}

	tmpl, err := svc.CreateMailTemplate(ctx, req)
	if err != nil {
		t.Fatalf("CreateMailTemplate failed: %v", err)
	}
	if tmpl.Code != req.Code {
		t.Errorf("code = %s, want %s", tmpl.Code, req.Code)
	}
	if tmpl.ContentType != "html" {
		t.Errorf("content_type = %s, want 默认 html", tmpl.ContentType)
	}
	if !tmpl.Enabled {
		t.Errorf("enabled = false, want 默认 true")
	}
}

func TestCreateMailTemplateCustomContentType(t *testing.T) {
	client := newTestClient(t)
	svc := NewMailTemplateServiceImpl(client)

	ctx := context.Background()
	contentType := "text"
	tmpl, err := svc.CreateMailTemplate(ctx, model.MailTemplateCreateReq{
		Name:        "纯文本通知",
		Code:        "plain_notice",
		Subject:     "通知",
		Content:     "你好 {nickname}",
		ContentType: &contentType,
	})
	if err != nil {
		t.Fatalf("CreateMailTemplate failed: %v", err)
	}
	if tmpl.ContentType != "text" {
		t.Errorf("content_type = %s, want text", tmpl.ContentType)
	}
}

func TestCreateMailTemplateDuplicateCode(t *testing.T) {
	client := newTestClient(t)
	svc := NewMailTemplateServiceImpl(client)

	ctx := context.Background()
	req := model.MailTemplateCreateReq{
		Name:    "注册欢迎邮件",
		Code:    "register_welcome",
		Subject: "欢迎加入",
		Content: "感谢注册",
	}
	if _, err := svc.CreateMailTemplate(ctx, req); err != nil {
		t.Fatalf("first CreateMailTemplate failed: %v", err)
	}

	_, err := svc.CreateMailTemplate(ctx, model.MailTemplateCreateReq{
		Name:    "重复编码",
		Code:    "register_welcome",
		Subject: "重复",
		Content: "重复",
	})
	if !errors.Is(err, model.ErrMailTemplateCodeExists) {
		t.Fatalf("err = %v, want ErrMailTemplateCodeExists", err)
	}
}

func TestUpdateMailTemplate(t *testing.T) {
	client := newTestClient(t)
	svc := NewMailTemplateServiceImpl(client)

	ctx := context.Background()
	tmpl, err := svc.CreateMailTemplate(ctx, model.MailTemplateCreateReq{
		Name:    "注册欢迎邮件",
		Code:    "register_welcome",
		Subject: "欢迎加入星屑，{nickname}",
		Content: "<p>感谢你的注册。</p>",
	})
	if err != nil {
		t.Fatalf("CreateMailTemplate failed: %v", err)
	}

	newSubject := "欢迎加入星屑！{nickname}"
	newContentType := "text"
	updated, err := svc.UpdateMailTemplate(ctx, tmpl.ID, model.MailTemplateUpdateReq{
		Subject:     &newSubject,
		ContentType: &newContentType,
	})
	if err != nil {
		t.Fatalf("UpdateMailTemplate failed: %v", err)
	}
	if updated.Subject != newSubject {
		t.Errorf("subject = %s, want %s", updated.Subject, newSubject)
	}
	if updated.ContentType != "text" {
		t.Errorf("content_type = %s, want text", updated.ContentType)
	}
	if updated.Code != tmpl.Code {
		t.Errorf("code = %s, want %s（未更新的字段应保持不变）", updated.Code, tmpl.Code)
	}
}

func TestUpdateMailTemplateDuplicateCode(t *testing.T) {
	client := newTestClient(t)
	svc := NewMailTemplateServiceImpl(client)

	ctx := context.Background()
	tmplA, err := svc.CreateMailTemplate(ctx, model.MailTemplateCreateReq{
		Name:    "模板A",
		Code:    "code_a",
		Subject: "主题A",
		Content: "内容A",
	})
	if err != nil {
		t.Fatalf("create template A failed: %v", err)
	}
	if _, err := svc.CreateMailTemplate(ctx, model.MailTemplateCreateReq{
		Name:    "模板B",
		Code:    "code_b",
		Subject: "主题B",
		Content: "内容B",
	}); err != nil {
		t.Fatalf("create template B failed: %v", err)
	}

	dupCode := "code_b"
	_, err = svc.UpdateMailTemplate(ctx, tmplA.ID, model.MailTemplateUpdateReq{
		Code: &dupCode,
	})
	if !errors.Is(err, model.ErrMailTemplateCodeExists) {
		t.Fatalf("err = %v, want ErrMailTemplateCodeExists", err)
	}

	// 更新为自身已有编码不应报错
	ownCode := "code_a"
	if _, err := svc.UpdateMailTemplate(ctx, tmplA.ID, model.MailTemplateUpdateReq{
		Code: &ownCode,
	}); err != nil {
		t.Fatalf("update to own code failed: %v", err)
	}
}

func TestQueryAndDeleteMailTemplate(t *testing.T) {
	client := newTestClient(t)
	svc := NewMailTemplateServiceImpl(client)

	ctx := context.Background()
	tmpl, err := svc.CreateMailTemplate(ctx, model.MailTemplateCreateReq{
		Name:    "注册欢迎邮件",
		Code:    "register_welcome",
		Subject: "欢迎加入星屑，{nickname}",
		Content: "感谢你的注册。",
	})
	if err != nil {
		t.Fatalf("CreateMailTemplate failed: %v", err)
	}

	// 按 ID 查询
	got, err := svc.QueryMailTemplate(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("QueryMailTemplate failed: %v", err)
	}
	if got.ID != tmpl.ID {
		t.Errorf("id = %d, want %d", got.ID, tmpl.ID)
	}

	// 查询不存在的模板
	if _, err := svc.QueryMailTemplate(ctx, 99999); !errors.Is(err, model.ErrMailTemplateNotFound) {
		t.Fatalf("err = %v, want ErrMailTemplateNotFound", err)
	}

	// 分页查询
	count, list, err := svc.QueryMailTemplatePage(ctx, model.MailTemplatePageQuery{Page: 1, Size: 10})
	if err != nil {
		t.Fatalf("QueryMailTemplatePage failed: %v", err)
	}
	if count != 1 || len(list) != 1 {
		t.Errorf("count = %d, len = %d, want 1, 1", count, len(list))
	}

	// 按名称模糊搜索
	count, _, err = svc.QueryMailTemplatePage(ctx, model.MailTemplatePageQuery{Page: 1, Size: 10, Name: "注册"})
	if err != nil {
		t.Fatalf("QueryMailTemplatePage with name filter failed: %v", err)
	}
	if count != 1 {
		t.Errorf("filtered count = %d, want 1", count)
	}

	// 按编码模糊搜索
	count, _, err = svc.QueryMailTemplatePage(ctx, model.MailTemplatePageQuery{Page: 1, Size: 10, Code: "register"})
	if err != nil {
		t.Fatalf("QueryMailTemplatePage with code filter failed: %v", err)
	}
	if count != 1 {
		t.Errorf("filtered count = %d, want 1", count)
	}

	// 删除
	if err := svc.DeleteMailTemplate(ctx, tmpl.ID); err != nil {
		t.Fatalf("DeleteMailTemplate failed: %v", err)
	}
	if err := svc.DeleteMailTemplate(ctx, tmpl.ID); !errors.Is(err, model.ErrMailTemplateNotFound) {
		t.Fatalf("second delete err = %v, want ErrMailTemplateNotFound", err)
	}
}

func ptr[T any](v T) *T {
	return &v
}
