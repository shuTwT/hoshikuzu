package prompt

import (
	"context"
	"errors"
	"testing"

	"github.com/shuTwT/hoshikuzu/ent/enttest"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	_ "github.com/mattn/go-sqlite3"
)

func TestPromptServiceCRUD(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:prompt?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	ctx := context.Background()
	svc := NewPromptServiceImpl(client)

	created, err := svc.CreatePrompt(ctx, &model.PromptCreateReq{
		Title:   "猫",
		Content: "a cat",
		Images: []model.PromptImageReq{
			{ImageURL: "https://img/b.png", ModelName: "sdxl", Sort: 1},
			{ImageURL: "https://img/a.png", ModelName: "flux", Sort: 0},
		},
		Tags: []string{"动物", " 动物 ", "写实", ""},
	})
	if err != nil {
		t.Fatalf("CreatePrompt: %v", err)
	}
	if got := len(created.Edges.Images); got != 2 {
		t.Fatalf("images = %d, want 2", got)
	}
	if created.Edges.Images[0].ModelName != "flux" {
		t.Errorf("images not ordered by sort: %+v", created.Edges.Images)
	}
	if got := len(created.Edges.Tags); got != 2 {
		t.Fatalf("tags = %d, want 2 (deduplicated)", got)
	}

	updated, err := svc.UpdatePrompt(ctx, created.ID, &model.PromptUpdateReq{
		Title:   "猫2",
		Content: "a cat, 4k",
		Images:  []model.PromptImageReq{{ImageURL: "https://img/c.png", ModelName: "mj"}},
		Tags:    []string{"写实", "新标签"},
	})
	if err != nil {
		t.Fatalf("UpdatePrompt: %v", err)
	}
	if updated.Title != "猫2" || len(updated.Edges.Images) != 1 || len(updated.Edges.Tags) != 2 {
		t.Fatalf("unexpected update result: %+v", updated)
	}

	tags, err := svc.ListPromptTag(ctx)
	if err != nil || len(tags) != 3 {
		t.Fatalf("ListPromptTag = %d, %v; want 3 tags", len(tags), err)
	}

	var xieShiID int
	for _, tag := range tags {
		if tag.Name == "写实" {
			xieShiID = tag.ID
		}
	}
	count, list, err := svc.ListPromptPage(ctx, &model.PromptPageQuery{Page: 1, Size: 10, Keyword: "4k", TagID: xieShiID})
	if err != nil || count != 1 || len(list) != 1 {
		t.Fatalf("ListPromptPage = %d/%d, %v; want 1", count, len(list), err)
	}

	if _, err := svc.CreatePrompt(ctx, &model.PromptCreateReq{
		Title: "bad", Content: "x", Images: []model.PromptImageReq{{ImageURL: " "}},
	}); !errors.Is(err, model.ErrPromptImageURLEmpty) {
		t.Errorf("empty image url err = %v", err)
	}

	if err := svc.DeletePrompt(ctx, created.ID); err != nil {
		t.Fatalf("DeletePrompt: %v", err)
	}
	if n := client.PromptImage.Query().CountX(ctx); n != 0 {
		t.Errorf("images left after delete = %d", n)
	}
	if _, err := svc.QueryPrompt(ctx, created.ID); !errors.Is(err, model.ErrPromptNotFound) {
		t.Errorf("QueryPrompt after delete err = %v", err)
	}
	if err := svc.DeletePrompt(ctx, created.ID); !errors.Is(err, model.ErrPromptNotFound) {
		t.Errorf("DeletePrompt twice err = %v", err)
	}
}
