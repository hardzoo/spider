package spider

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/axgle/mahonia"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"

	log "github.com/golang/glog"

	"github.com/withzoo/spider/config"
	"github.com/withzoo/spider/utils"
)

const (
	MAX_TASK_NUM = 500
)

type CrawlTask struct {
	Url   string
	Depth int64
}

type Spider struct {
	SpiderConfig *config.SpiderConfig
	InitTask     []CrawlTask
	Tasks        chan CrawlTask
	WaitGroup    sync.WaitGroup
}

func NewSpider(spiderConfig *config.SpiderConfig) *Spider {
	var initTask []CrawlTask
	spider := &Spider{
		SpiderConfig: spiderConfig,
		InitTask:     initTask,
		Tasks:        make(chan CrawlTask, MAX_TASK_NUM),
	}

	return spider
}

/*
* InitCrawlTask - 将url.data文件中目标seeds放入初始化爬虫任务
*
* RETURNS:
*   - error: 成功返回nil
 */
func (s *Spider) InitCrawlTask() error {
	seeds, err := utils.LoadSeedFromFile(s.SpiderConfig.UrlListFile)
	if err != nil {
		return err
	}

	for _, v := range seeds {
		s.InitTask = append(s.InitTask, CrawlTask{Url: v, Depth: 0})
	}

	return nil
}

func (s *Spider) Crawling() {
	for {
		select {
		case crawlTask := <-s.Tasks:
			s.Start(crawlTask)
			//抓取间隔
			time.Sleep(time.Duration(s.SpiderConfig.CrawlInterval) * time.Second)
		}
	}
}

/*
* AddCrawlTask - 添加新的爬虫任务，防止阻塞调用此方法的协程，用gorountine触发新任务
*
* PARAMS:
*   - task: 爬虫任务
 */
func (s *Spider) AddCrawlTask(task CrawlTask) {
	s.WaitGroup.Add(1)
	go func(ct CrawlTask) {
		s.Tasks <- ct
	}(task)
}

func (s *Spider) Start(crawlTask CrawlTask) {
	defer s.WaitGroup.Done()
	log.Infof("[%s]开始抓取, 当前深度[%d] ", crawlTask.Url, crawlTask.Depth)

	if ctx, err := s.CrawlHTML(crawlTask.Url); err != nil {
		log.Warningf("[%s]抓取失败, 当前深度[%d], 错误: %v", crawlTask.Url, crawlTask.Depth, err)
	} else {
		if err = s.SaveHtmlToFile(crawlTask, ctx); err != nil {
			log.Warningf("[%s]抓取结果存储为文件失败, 当前深度[%d], 错误: %v", crawlTask.Url,
				crawlTask.Depth, err)
		}

		if s.SpiderConfig.MaxDepth > crawlTask.Depth {
			if err = s.ParseHTML(crawlTask, ctx); err != nil {
				log.Warningf("[%s]解析抓取结果失败, 当前深度[%d], 错误: %v", crawlTask.Url,
					crawlTask.Depth, err)
			}
		}
	}
}

/*
* CrawlHTML - 抓取url网页内容
*
* PARAMS:
*   - url: 网页url地址
*
* RETURNS:
*   - string: 抓取的内容
*   - error: 成功返回nil
 */
func (s *Spider) CrawlHTML(url string) (string, error) {
	var timeout = time.Duration(s.SpiderConfig.CrawlTimeout) * time.Second
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, timeout) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(timeout)) //设置发送接受数据超时
			return c, nil
		},
		ResponseHeaderTimeout: timeout,
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	bodyStr := string(body.Bytes())

	//判断网页编码
	_, encodeName, _ := charset.DetermineEncoding(body.Bytes(), "")
	//将其他字符编码转换成UTF-8
	if encodeName != "" {
		decoder := mahonia.NewDecoder(encodeName)
		bodyStr = decoder.ConvertString(bodyStr)
	}

	return bodyStr, nil
}

/*
* ParseHTML - 解析HTML页面，匹配页面内符合条件的url
*
* PARAMS:
*   - crawlTask: 爬虫任务
*   - ctx: 当前爬虫任务抓取的页面内容
*
* RETURNS:
*   - error: 成功返回nil
 */
func (s *Spider) ParseHTML(crawlTask CrawlTask, ctx string) error {
	re, err := regexp.Compile(s.SpiderConfig.TargetUrl)
	if err != nil {
		return err
	}

	htmlNode, err := html.Parse(strings.NewReader(ctx))
	if err != nil {
		return err
	}

	urls := utils.ParseHTML(htmlNode, crawlTask.Url, re)

	//根据解析结果开始新的爬虫任务
	for _, url := range urls {
		s.AddCrawlTask(CrawlTask{Url: url, Depth: crawlTask.Depth + 1})
	}

	log.Infof("解析抓取的页面[%s]成功, 当前深度[%d]", crawlTask.Url, crawlTask.Depth)
	return nil
}

/*
* SaveHtmlToFile - 保存抓取的结果到文件
*
* PARAMS:
*   - crawlTask: 爬虫任务
*   - ctx: 当前爬虫任务抓取的页面内容
*
* RETURNS:
*   - error: 成功返回nil
 */
func (s *Spider) SaveHtmlToFile(crawlTask CrawlTask, ctx string) error {
	filename := s.SpiderConfig.OutputDirectory + "/" + utils.UrlToFilename(crawlTask.Url)

	if utils.IsFileExist(filename) {
		if err := os.Remove(filename); err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	if _, err = w.WriteString(ctx); err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}

	log.Infof("[%s]抓取成功, 当前深度[%d], 抓取结果存储为文件[%s]", crawlTask.Url, crawlTask.Depth, filename)
	return nil
}
