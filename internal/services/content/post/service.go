package post

import (
	"context"
	"fmt"

	"math/rand/v2"
	"sort"
	"strings"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/category"
	"github.com/shuTwT/hoshikuzu/ent/post"
	"github.com/shuTwT/hoshikuzu/ent/tag"
	"github.com/shuTwT/hoshikuzu/pkg/cache"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
	"github.com/shuTwT/hoshikuzu/pkg/utils"
)

type PostService interface {
	QueryPostList(c context.Context, req model.PostListReq) ([]*ent.Post, error)
	QueryPostPage(c context.Context, req model.PostPageReq) ([]*ent.Post, int, error)
	QueryPostBySlug(c context.Context, slug string) (*ent.Post, error)
	QueryPostById(c context.Context, id int) (*ent.Post, error)
	CreatePost(c context.Context, title string, content string) (*ent.Post, error)
	UpdatePostContent(c context.Context, id int, content string, htmlContent *string, mdContent *string) (*ent.Post, error)
	UpdatePostSetting(c context.Context, id int, post model.PostUpdateReq) (*ent.Post, error)
	DeletePost(c context.Context, id int) error
	GetPostCount(c context.Context) (int, error)
	GetPostMonthStats(c context.Context, req model.PostMonthStatsReq) ([]model.PostMonthStat, error)
	GetRandomPost(c context.Context) (*ent.Post, error)
	SearchPosts(c context.Context, req model.PostSearchReq) ([]*model.PostSearchResp, int, error)
	PublishPost(c context.Context, id int) (*ent.Post, error)
	UnpublishPost(c context.Context, id int) (*ent.Post, error)
}

type PostServiceImpl struct {
	client *ent.Client
}

func NewPostServiceImpl(client *ent.Client) *PostServiceImpl {
	return &PostServiceImpl{client: client}
}

func (s *PostServiceImpl) QueryPostList(c context.Context, req model.PostListReq) ([]*ent.Post, error) {
	query := s.client.Post.Query()

	if req.CategoryName != "" {
		query.Where(post.HasCategoriesWith(category.Name(req.CategoryName)))
	}

	if req.TagName != "" {
		query.Where(post.HasTagsWith(tag.Name(req.TagName)))
	}

	if req.Year != nil {
		startDate := time.Date(*req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(*req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		query.Where(post.PublishedAtGTE(startDate), post.PublishedAtLT(endDate))
	}

	if req.Month != nil && req.Year != nil {
		startDate := time.Date(*req.Year, time.Month(*req.Month), 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(*req.Year, time.Month(*req.Month)+1, 1, 0, 0, 0, 0, time.UTC)
		query.Where(post.PublishedAtGTE(startDate), post.PublishedAtLT(endDate))
	}

	if req.IsPinToTop != nil {
		query.Where(post.IsPinToTop(*req.IsPinToTop))
	}

	query = query.
		WithCategories().
		WithTags().
		Order(ent.Desc(post.FieldID))

	if req.Limit != nil {
		query.Limit(*req.Limit)
	}

	posts, err := query.All(c)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostServiceImpl) QueryPostPage(c context.Context, req model.PostPageReq) ([]*ent.Post, int, error) {
	query := s.client.Post.Query()

	if req.CategoryID != nil {
		query.Where(post.HasCategoriesWith(category.ID(*req.CategoryID)))
	}

	if req.CategoryName != "" {
		query.Where(post.HasCategoriesWith(category.Name(req.CategoryName)))
	}

	if req.TagID != nil {
		query.Where(post.HasTagsWith(tag.ID(*req.TagID)))
	}

	if req.TagName != "" {
		query.Where(post.HasTagsWith(tag.Name(req.TagName)))
	}

	if req.Title != "" {
		query.Where(post.TitleContains(req.Title))
	}

	if req.Year != nil {
		startDate := time.Date(*req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(*req.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		query.Where(post.PublishedAtGTE(startDate), post.PublishedAtLT(endDate))
	}

	if req.Month != nil && req.Year != nil {
		startDate := time.Date(*req.Year, time.Month(*req.Month), 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(*req.Year, time.Month(*req.Month)+1, 1, 0, 0, 0, 0, time.UTC)
		query.Where(post.PublishedAtGTE(startDate), post.PublishedAtLT(endDate))
	}

	count, err := query.Count(c)
	if err != nil {
		return nil, 0, err
	}
	posts, err := query.
		WithCategories().
		WithTags().
		Order(ent.Desc(post.FieldID)).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		All(c)
	if err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

func (s *PostServiceImpl) CreatePost(c context.Context, title string, content string) (*ent.Post, error) {

	newPost, err := s.client.Post.Create().
		SetTitle(title).
		SetContent(content).
		Save(c)
	slug, err := utils.GenerateSlug(title, newPost.CreatedAt.Unix())
	if err != nil {
		return newPost, err
	}
	newPost, err = s.client.Post.UpdateOneID(newPost.ID).
		SetSlug(slug).
		Save(c)
	return newPost, err
}

func (s *PostServiceImpl) QueryPostBySlug(c context.Context, slug string) (*ent.Post, error) {
	post, err := s.client.Post.Query().
		Where(post.Slug(slug)).
		WithCategories().
		WithTags().
		Only(c)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostServiceImpl) QueryPostById(c context.Context, id int) (*ent.Post, error) {
	post, err := s.client.Post.Query().
		Where(post.ID(id)).
		WithCategories().
		WithTags().
		Only(c)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostServiceImpl) UpdatePostContent(c context.Context, id int, content string, htmlContent *string, mdContent *string) (*ent.Post, error) {
	newPost, err := s.client.Post.UpdateOneID(id).
		SetContent(content).
		SetNillableHTMLContent(htmlContent).
		SetNillableMdContent(mdContent).
		Save(c)
	return newPost, err
}

func (s *PostServiceImpl) UpdatePostSetting(c context.Context, id int, updateReq model.PostUpdateReq) (*ent.Post, error) {
	client := s.client
	var summary string
	if updateReq.IsAutogenSummary {
		summary = "生成失败"
	} else {
		summary = updateReq.Summary
	}

	oldPost, err := client.Post.Query().Where(post.ID(id)).First(c)
	if err != nil {
		return nil, err
	}

	var slug string

	generatedSlug, err := utils.GenerateSlug(updateReq.Title, oldPost.CreatedAt.Unix())
	if err != nil {
		return nil, err
	}
	slug = generatedSlug

	newPost, err := client.Post.UpdateOneID(id).
		SetTitle(updateReq.Title).
		SetSlug(slug).
		SetCover(updateReq.Cover).
		SetKeywords(updateReq.Keywords).
		SetCopyright(updateReq.Copyright).
		SetAuthor(updateReq.Author).
		SetIsAutogenSummary(updateReq.IsAutogenSummary).
		SetIsVisible(updateReq.IsVisible).
		SetIsPinToTop(updateReq.IsPinToTop).
		SetIsAllowComment(updateReq.IsAllowComment).
		SetIsVisibleAfterComment(updateReq.IsVisibleAfterComment).
		SetIsVisibleAfterPay(updateReq.IsVisibleAfterPay).
		SetPrice(int(updateReq.Price * 100)).
		SetSummary(summary).
		AddCategoryIDs(updateReq.Categories...).
		AddTagIDs(updateReq.Tags...).
		Save(c)
	return newPost, err
}

func (s *PostServiceImpl) DeletePost(c context.Context, id int) error {
	return s.client.Post.DeleteOneID(id).Exec(c)
}

func (s *PostServiceImpl) GetPostCount(c context.Context) (int, error) {
	count, err := s.client.Post.Query().Count(c)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *PostServiceImpl) GetPostMonthStats(c context.Context, req model.PostMonthStatsReq) ([]model.PostMonthStat, error) {
	posts, err := s.client.Post.Query().
		Order(ent.Desc(post.FieldCreatedAt)).
		All(c)
	if err != nil {
		return nil, err
	}

	monthMap := make(map[string]int)
	for _, p := range posts {

		monthKey := fmt.Sprintf("%d-%02d", p.CreatedAt.Year(), p.CreatedAt.Month())
		monthMap[monthKey]++

	}

	var stats []model.PostMonthStat
	for month, count := range monthMap {
		stats = append(stats, model.PostMonthStat{
			Month: month,
			Count: count,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Month > stats[j].Month
	})

	if req.Limit > 0 && req.Limit < len(stats) {
		stats = stats[:req.Limit]
	}

	return stats, nil
}

func (s *PostServiceImpl) GetRandomPost(c context.Context) (*ent.Post, error) {
	count, err := s.client.Post.Query().Count(c)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	offset := rand.IntN(count)

	post, err := s.client.Post.Query().
		WithCategories().
		WithTags().
		Order(ent.Asc(post.FieldID)).
		Offset(offset).
		Limit(1).
		First(c)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostServiceImpl) SearchPosts(c context.Context, req model.PostSearchReq) ([]*model.PostSearchResp, int, error) {
	cacheKey := fmt.Sprintf("post:search:%s:%d:%d", req.Keyword, req.Page, req.Size)

	if cached, found := cache.GetCache().Get(cacheKey); found {
		if result, ok := cached.([]*model.PostSearchResp); ok {
			return result, len(result), nil
		}
	}

	keyword := strings.ToLower(req.Keyword)

	posts, err := s.client.Post.Query().
		Where(
			post.Or(
				post.TitleContains(keyword),
				post.ContentContains(keyword),
				post.SummaryContains(keyword),
				post.KeywordsContains(keyword),
			),
			post.StatusEQ("published"),
			post.IsVisible(true),
		).
		Order(ent.Desc(post.FieldPublishedAt)).
		All(c)

	if err != nil {
		return nil, 0, err
	}

	var results []*model.PostSearchResp

	for _, p := range posts {
		relevance := s.calculateRelevance(p, keyword)
		if relevance > 0 {
			results = append(results, &model.PostSearchResp{
				ID:          p.ID,
				Title:       p.Title,
				Summary:     p.Summary,
				Content:     p.Content,
				Slug:        p.Slug,
				Cover:       p.Cover,
				Author:      p.Author,
				PublishedAt: (*model.LocalTime)(p.PublishedAt),
				ViewCount:   p.ViewCount,
				Relevance:   relevance,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Relevance == results[j].Relevance {
			if results[i].PublishedAt != nil && results[j].PublishedAt != nil {
				return results[i].PublishedAt.Time().After(results[j].PublishedAt.Time())
			}
			return results[i].ID > results[j].ID
		}
		return results[i].Relevance > results[j].Relevance
	})

	total := len(results)

	start := (req.Page - 1) * req.Size
	end := start + req.Size

	if start >= total {
		return []*model.PostSearchResp{}, total, nil
	}
	if end > total {
		end = total
	}

	pagedResults := results[start:end]

	cache.GetCache().Set(cacheKey, pagedResults, 5*time.Minute)

	return pagedResults, total, nil
}

func (s *PostServiceImpl) calculateRelevance(p *ent.Post, keyword string) float64 {
	var relevance float64 = 0

	title := strings.ToLower(p.Title)
	content := strings.ToLower(p.Content)
	summary := strings.ToLower(p.Summary)
	keywords := strings.ToLower(p.Keywords)

	if strings.Contains(title, keyword) {
		if title == keyword {
			relevance += 10.0
		} else if strings.HasPrefix(title, keyword) {
			relevance += 8.0
		} else {
			relevance += 5.0
		}
	}

	if strings.Contains(summary, keyword) {
		relevance += 3.0
	}

	if strings.Contains(keywords, keyword) {
		relevance += 2.0
	}

	if strings.Contains(content, keyword) {
		count := strings.Count(content, keyword)
		relevance += float64(count) * 0.5
	}

	return relevance
}

func (s *PostServiceImpl) PublishPost(c context.Context, id int) (*ent.Post, error) {
	post, err := s.client.Post.UpdateOneID(id).
		SetStatus("published").
		SetPublishedAt(time.Now()).
		Save(c)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostServiceImpl) UnpublishPost(c context.Context, id int) (*ent.Post, error) {
	post, err := s.client.Post.UpdateOneID(id).
		SetStatus("draft").
		Save(c)
	if err != nil {
		return nil, err
	}
	return post, nil
}
