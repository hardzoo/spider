package spider

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/withzoo/spider/config"
)

func TestSpider(t *testing.T) {
	config, err := config.LoadConfig("../conf/spider.conf")
	if err != nil {
		t.Errorf("load config file failed: %v", err)
	}
	config.Spider.OutputDirectory = "./output"
	//创建output文件保存路径,如果路径存在，则不会创建
	if err = os.MkdirAll(config.Spider.OutputDirectory, 0777); err != nil {
		t.Errorf("can not create output directory: %v", err)
	}

	sipderConfig := config.Spider
	sp := NewSpider(&sipderConfig)

	//两个协程并发操作
	go sp.Crawling()
	go sp.Crawling()

	ct := CrawlTask{
		Url:   "http://www.baidu.com",
		Depth: 0,
	}
	sp.WaitGroup.Add(1)
	go func(ct CrawlTask) {
		sp.Tasks <- ct
	}(ct)
	sp.WaitGroup.Wait()

	files, _ := ioutil.ReadDir(config.Spider.OutputDirectory)
	if len(files) != 3 {
		t.Errorf("TestSpider failed, expect crawl 3 file, but got [%v] file", len(files))
	}

	os.RemoveAll(config.Spider.OutputDirectory)
}

func TestInitCrawlTask(t *testing.T) {
	config, err := config.LoadConfig("../conf/spider.conf")
	if err != nil {
		t.Errorf("load config file failed: %v", err)
	}
	config.Spider.UrlListFile = "../data/url.data"
	sipderConfig := config.Spider
	sp := NewSpider(&sipderConfig)
	sp.InitCrawlTask()

	if sp.InitTask[0].Url != "http://www.baidu.com" {
		t.Errorf("TestInitCrawlTask failed, expect seed is: %v, but got %s", "http://www.baidu.com", sp.InitTask[0].Url)
	}

	if sp.InitTask[1].Url != "http://www.sina.com.cn" {
		t.Errorf("TestInitCrawlTask failed, expect seed is: %v, but got %s", "http://www.sina.com.cn", sp.InitTask[1].Url)
	}
}

func TestCrawlHTML(t *testing.T) {
	config, err := config.LoadConfig("../conf/spider.conf")
	if err != nil {
		t.Errorf("load config file failed: %v", err)
	}
	config.Spider.UrlListFile = "../data/url.data"
	sipderConfig := config.Spider
	sp := NewSpider(&sipderConfig)

	url := "www.baidu123.com"
	_, err = sp.CrawlHTML(url)
	if err == nil {
		t.Errorf("TestCrawlHTML failed: %v", err)
	}

	url = "http://www.baidu.com"
	_, err = sp.CrawlHTML(url)
	if err != nil {
		t.Errorf("TestCrawlHTML failed: %v", err)
	}

}

func TestSaveHtmlToFileForErrPath(t *testing.T) {
	config, err := config.LoadConfig("../conf/spider.conf")
	if err != nil {
		t.Errorf("load config file failed: %v", err)
	}
	config.Spider.OutputDirectory = "./output"
	sipderConfig := config.Spider
	sp := NewSpider(&sipderConfig)

	content := "123"
	ct := CrawlTask{
		Url:   "http://www.baidu.com",
		Depth: 0,
	}

	err = sp.SaveHtmlToFile(ct, content)
	if err == nil {
		t.Errorf("TestSaveHtmlToFile failed: %v", err)
	}
}
