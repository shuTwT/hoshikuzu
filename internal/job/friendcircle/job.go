package job

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/flink"
	friend_circle_service "github.com/shuTwT/hoshikuzu/internal/services/content/friendcircle"
	schedule_model "github.com/shuTwT/hoshikuzu/pkg/domain/model/schedule"

	"io"
	"log"
	"net/http"

	"github.com/shuTwT/hoshikuzu/pkg/domain/model/rss"
)

type FriendCircleJob struct {
	DbClient            *ent.Client
	FriendCircleService friend_circle_service.FriendCircleService
}

func NewFriendCircleJob(dbClient *ent.Client, friendCircleService friend_circle_service.FriendCircleService) schedule_model.CronJob {
	return &FriendCircleJob{
		DbClient:            dbClient,
		FriendCircleService: friendCircleService,
	}
}

func (job FriendCircleJob) Execute(c context.Context) error {
	// 朋友圈规则
	// rules, err := dbClient.FriendCircleRule.Query().All(c)
	// if err != nil {
	// 	return err
	// }
	// 遍历友链
	flinks, err := job.DbClient.FLink.Query().Where(
		flink.EnableFriendCircleEQ(true),
		flink.FriendCircleRuleIDNotNil(),
	).All(c)
	if err != nil {
		return err
	}
	httpClient := &http.Client{}
	for _, flink := range flinks {
		req, err := http.NewRequest("GET", flink.URL, nil)
		if err != nil {
			log.Println("创建请求失败:%w", err)
			continue
		}
		// 设置User-Agent，避免被识别为爬虫
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Println("请求失败:%w", err)
			continue
		}
		// 必须关闭响应体，避免内存泄漏
		defer resp.Body.Close()

		// 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			log.Println("响应状态码错误: &d", resp.StatusCode)
			continue
		}

		// 状态码为 200，表示可以访问
		// 接下来先访问/rss.xml
		err = job.VisitRss(httpClient, flink.URL, flink)
		if err != nil {
			log.Println("%w", err)
			// 尝试访问/atom.xml
			err = job.VisitAtom(httpClient, flink.URL, flink)
			if err != nil {
				log.Println("%w", err)

			}
		}

	}

	return nil
}

func (job FriendCircleJob) VisitRss(httpClient *http.Client, baseUrl string, flink *ent.FLink) error {
	rssReq, err := http.NewRequest("GET", baseUrl+"/rss.xml", nil)
	if err != nil {
		log.Println("创建请求失败:%w", err)
		return err
	}

	rssReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	rssResp, err := httpClient.Do(rssReq)
	if err != nil {
		log.Println("请求失败:%w", err)
		return err
	}

	defer rssResp.Body.Close()
	// 判断状态码
	if rssResp.StatusCode != http.StatusOK {
		return fmt.Errorf("响应状态码错误: %d", rssResp.StatusCode)
	}

	// 状态码是 200，有东西
	// 一次性读取响应体到字节数组（核心优化点）
	log.Println("开始读取响应体")
	body, err := io.ReadAll(rssResp.Body)
	if err != nil {
		log.Println("读取响应体失败: %w", err)
		return err
	}

	// 空响应体处理
	if len(body) == 0 {
		log.Printf("响应体为空 (URL: %s)", baseUrl+"/rss.xml")
		return err
	}

	log.Println("响应体不为空")
	rssRes, err := job.parseRss(body)
	if err != nil {
		return err
	}
	switch v := rssRes.(type) {
	case rss.RSS2:
		fmt.Printf("\n===== RSS 2.0 解析结果 =====\n")
		fmt.Printf("频道标题：%s\n", v.Channel.Title)
		fmt.Printf("频道链接：%s\n", v.Channel.Link)
		fmt.Printf("频道描述：%s\n", v.Channel.Description)
		fmt.Printf("文章总数：%d\n", len(v.Channel.Items))
		for _, item := range v.Channel.Items[:5] { // 仅打印前5篇，避免输出过长
			var exist bool
			exist, err = job.FriendCircleService.ExistsRecord(context.TODO(), item.Link)
			if err == nil && !exist {
				_ = job.FriendCircleService.InsertRecord(context.TODO(), flink.Name, flink.AvatarURL, item.Title, item.Link, item.PubDate)
			}
		}
	case rss.AtomFeed:
		fmt.Printf("\n===== Atom 1.0 解析结果 =====\n")
		fmt.Printf("Feed标题：%s\n", v.Title)
		// 提取主链接
		mainLink := ""
		for _, link := range v.Link {
			if link.Rel == "alternate" || link.Rel == "" {
				mainLink = link.Href
				break
			}
		}
		fmt.Printf("Feed链接：%s\n", mainLink)
		fmt.Printf("文章总数：%d\n", len(v.Entries))
		for _, entry := range v.Entries[:5] { // 仅打印前5篇
			// 提取文章链接
			entryLink := ""
			for _, link := range entry.Link {
				if link.Rel == "alternate" || link.Rel == "" {
					entryLink = link.Href
					break
				}
			}
			var exist bool
			exist, err = job.FriendCircleService.ExistsRecord(context.TODO(), entryLink)
			if err == nil && !exist {
				job.FriendCircleService.InsertRecord(context.TODO(), flink.Name, flink.AvatarURL, entry.Title, entryLink, entry.Updated)
			}
		}

	default:
		log.Println("不支持的RSS/Atom格式")
	}

	return nil
}

func (job FriendCircleJob) VisitAtom(httpClient *http.Client, baseUrl string, flink *ent.FLink) error {
	rssReq, err := http.NewRequest("GET", baseUrl+"/atom.xml", nil)
	if err != nil {
		log.Println("创建请求失败:%w", err)
		return err
	}

	rssReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	rssResp, err := httpClient.Do(rssReq)
	if err != nil {
		log.Println("请求失败:%w", err)
		return err
	}
	defer rssResp.Body.Close()
	// 判断状态码
	if rssResp.StatusCode != http.StatusOK {
		return fmt.Errorf("响应状态码错误: %d", rssResp.StatusCode)
	}

	// 状态码是 200，有东西
	// 一次性读取响应体到字节数组（核心优化点）
	log.Println("开始读取响应体")
	body, err := io.ReadAll(rssResp.Body)
	if err != nil {
		log.Println("读取响应体失败: %w", err)
		return err
	}

	// 空响应体处理
	if len(body) == 0 {
		log.Printf("响应体为空 (URL: %s)", baseUrl+"/rss.xml")
		return err
	}

	log.Println("响应体不为空")
	atomRes, err := job.parseAtom(body)
	if err != nil {
		return err
	}
	fmt.Printf("\n===== Atom 1.0 解析结果 =====\n")
	fmt.Printf("Feed标题：%s\n", atomRes.Title)
	// 提取主链接
	mainLink := ""
	for _, link := range atomRes.Link {
		if link.Rel == "alternate" || link.Rel == "" {
			mainLink = link.Href
			break
		}
	}
	fmt.Printf("Feed链接：%s\n", mainLink)
	fmt.Printf("文章总数：%d\n", len(atomRes.Entries))
	for _, entry := range atomRes.Entries[:5] { // 仅打印前5篇
		// 提取文章链接
		entryLink := ""
		for _, link := range entry.Link {
			if link.Rel == "alternate" || link.Rel == "" {
				entryLink = link.Href
				break
			}
		}
		var exist bool
		exist, err = job.FriendCircleService.ExistsRecord(context.TODO(), entryLink)
		if err == nil && !exist {
			job.FriendCircleService.InsertRecord(context.TODO(), flink.Name, flink.AvatarURL, entry.Title, entryLink, entry.Updated)
		}
	}
	return nil
}

func (FriendCircleJob) parseRss(body []byte) (interface{}, error) {

	var rss2 rss.RSS2

	if err := xml.Unmarshal(body, &rss2); err == nil && rss2.Version == "2.0" {
		//解析到了
		return rss2, nil
	}

	return nil, fmt.Errorf("不支持的 RSS 格式")
}

func (FriendCircleJob) parseAtom(body []byte) (*rss.AtomFeed, error) {
	var atom rss.AtomFeed
	if err := xml.Unmarshal(body, &atom); err == nil {
		if atom.Title != "" && len(atom.Entries) > 0 {
			return &atom, nil
		}
	}

	return nil, fmt.Errorf("不支持的 Atom 格式")
}

func (FriendCircleJob) Type() schedule_model.JobType {
	return schedule_model.DurationJobType
}

func (FriendCircleJob) Duration() time.Duration {
	return 24 * time.Hour
}

func (FriendCircleJob) Description() string {
	return "朋友圈爬虫任务"
}
