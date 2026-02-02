package model

import (
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
)

type PostPageReq struct {
	Page         int    `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size         int    `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
	CategoryID   *int   `json:"category_id" query:"category_id" form:"category_id"`
	TagID        *int   `json:"tag_id" query:"tag_id" form:"tag_id"`
	Title        string `json:"title" query:"title" form:"title"`
	CategoryName string `json:"category_name" query:"category_name" form:"category_name"`
	TagName      string `json:"tag_name" query:"tag_name" form:"tag_name"`
	Year         *int   `json:"year" query:"year" form:"year"`
	Month        *int   `json:"month" query:"month" form:"month"`
}

type PostListReq struct {
	CategoryName string `json:"category_name" query:"category_name" form:"category_name"`
	TagName      string `json:"tag_name" query:"tag_name" form:"tag_name"`
	Year         *int   `json:"year" query:"year" form:"year"`
	Month        *int   `json:"month" query:"month" form:"month"`
	Limit        *int   `json:"limit" query:"limit" form:"limit" validate:"omitempty,min=1,max=100"`
	IsPinToTop   *bool  `json:"is_pin_to_top" query:"is_pin_to_top" form:"is_pin_to_top"`
}

// PostCreateReq represents the request body for creating a post.
type PostCreateReq struct {
	Title                 string  `json:"title" validate:"required"`                                //文章标题
	Slug                  *string `json:"slug,omitempty"`                                           //文章别名
	Content               string  `json:"content" validate:"required"`                              //文章内容
	MdContent             *string `json:"md_content,omitempty"`                                     //md文章内容
	HtmlContent           *string `json:"html_content,omitempty"`                                   //html文章内容
	ContentType           *string `json:"content_type" validate:"required,enum=markdown,html"`      //内容类型
	Status                *string `json:"status" validate:"required,enum=draft,published,archived"` //状态
	IsAutogenSummary      bool    `json:"is_autogen_summary"`                                       //是否自动生成摘要
	IsVisible             bool    `json:"is_visible"`                                               //是否可见
	IsPinToTop            bool    `json:"is_pin_to_top"`                                            //是否置顶
	IsAllowComment        bool    `json:"is_allow_comment"`                                         //是否允许评论
	IsVisibleAfterComment bool    `json:"is_visible_after_comment"`                                 //是否评论后可见
	IsVisibleAfterPay     bool    `json:"is_visible_after_pay"`                                     //是否支付后可见
	Price                 float32 `json:"price" validate:"required,min=0"`                          //文章价格
	Cover                 *string `json:"cover,omitempty"`                                          //文章封面
	Keywords              *string `json:"keywords,omitempty"`                                       //文章关键词
	Copyright             *string `json:"copyright,omitempty"`                                      //文章版权
	Author                string  `json:"author"`                                                   //作者
	Summary               *string `json:"summary,omitempty"`                                        //文章摘要
	Categories            []int   `json:"categories,omitempty"`                                     //分类ID列表
	Tags                  []int   `json:"tags,omitempty"`                                           //标签ID列表
}

// PostUpdateReq represents the request body for updating a post.
type PostUpdateReq struct {
	Title string `json:"title,omitempty"` //文章标题
	//Slug                  *string `json:"slug,omitempty"`                                           //文章别名
	Content               string  `json:"content"`                                                  //文章内容
	MdContent             *string `json:"md_content,omitempty"`                                     //md文章内容
	HtmlContent           *string `json:"html_content,omitempty"`                                   //html文章内容
	ContentType           *string `json:"content_type" validate:"required,enum=markdown,html"`      //内容类型
	Status                *string `json:"status" validate:"required,enum=draft,published,archived"` //状态
	IsAutogenSummary      bool    `json:"is_autogen_summary,omitempty"`                             //是否自动生成摘要
	IsVisible             bool    `json:"is_visible,omitempty"`                                     //是否可见
	IsPinToTop            bool    `json:"is_pin_to_top,omitempty"`                                  //是否置顶
	IsAllowComment        bool    `json:"is_allow_comment,omitempty"`                               //是否允许评论
	IsVisibleAfterComment bool    `json:"is_visible_after_comment,omitempty"`                       //是否评论后可见
	IsVisibleAfterPay     bool    `json:"is_visible_after_pay,omitempty"`                           //是否支付后可见
	Price                 float32 `json:"price" validate:"required,min=0"`                          //文章价格
	Cover                 string  `json:"cover,omitempty"`                                          //文章封面
	Keywords              string  `json:"keywords,omitempty"`                                       //文章关键词
	Copyright             string  `json:"copyright,omitempty"`                                      //文章版权
	Author                string  `json:"author"`                                                   //作者
	Summary               string  `json:"summary,omitempty"`                                        //文章摘要
	Categories            []int   `json:"categories,omitempty"`                                     //分类ID列表
	Tags                  []int   `json:"tags,omitempty"`                                           //标签ID列表
}

// PostResp represents the response body for a post.
type PostResp struct {
	ID                    int             `json:"id"`                       //文章ID
	CreatedAt             time.Time       `json:"created_at"`               //创建时间
	Title                 string          `json:"title"`                    //文章标题
	Slug                  *string         `json:"slug"`                     //文章别名
	Content               string          `json:"content"`                  //文章内容
	MdContent             *string         `json:"md_content,omitempty"`     //md文章内容
	HtmlContent           *string         `json:"html_content,omitempty"`   //html文章内容
	ContentType           string          `json:"content_type"`             //内容类型
	Status                string          `json:"status"`                   //状态
	IsAutogenSummary      bool            `json:"is_autogen_summary"`       //是否自动生成摘要
	IsVisible             bool            `json:"is_visible"`               //是否可见
	IsPinToTop            bool            `json:"is_pin_to_top"`            //是否置顶
	IsAllowComment        bool            `json:"is_allow_comment"`         //是否允许评论
	IsVisibleAfterComment bool            `json:"is_visible_after_comment"` //是否评论后可见
	IsVisibleAfterPay     bool            `json:"is_visible_after_pay"`     //是否支付后可见
	Price                 float32         `json:"price"`                    //文章价格
	PublishedAt           *time.Time      `json:"published_at,omitempty"`   //发布时间
	ViewCount             int             `json:"view_count"`               //浏览次数
	CommentCount          int             `json:"comment_count"`            //评论次数
	Cover                 string          `json:"cover"`                    //文章封面
	Keywords              string          `json:"keywords"`                 //文章关键词
	Copyright             string          `json:"copyright"`                //文章版权
	Author                string          `json:"author"`                   //作者
	Summary               string          `json:"summary"`                  //文章摘要
	Categories            []*ent.Category `json:"categories"`               //分类ID列表
	CategoryIds           []int           `json:"category_ids,omitempty"`   //分类ID列表
	Tags                  []*ent.Tag      `json:"tags"`                     //标签ID列表
	TagIds                []int           `json:"tag_ids,omitempty"`        //标签ID列表
}

type PostMonthStat struct {
	Month string `json:"month"` // MM-YYYY 格式
	Count int    `json:"count"` // 文章数量
}

type PostMonthStatsReq struct {
	Limit int `json:"limit" query:"limit" form:"limit" validate:"min=1,max=100"`
}

type PostSearchReq struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" validate:"required"`
	Page    int    `json:"page" query:"page" form:"page" validate:"required,min=1"`
	Size    int    `json:"page_size" query:"page_size" form:"page_size" validate:"required,min=1,max=100"`
}

type PostSearchResp struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content"`
	Slug        *string    `json:"slug"`
	Cover       string     `json:"cover"`
	Author      string     `json:"author"`
	PublishedAt *time.Time `json:"published_at"`
	ViewCount   int        `json:"view_count"`
	Relevance   float64    `json:"relevance"`
}
