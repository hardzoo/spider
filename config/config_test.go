package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	fileName := "../conf/spider.conf"
	config, err := LoadConfig(fileName)
	if err != nil {
		t.Errorf("TestLoadConfig err :%v", err)
	}

	if config.Spider.UrlListFile != "./data/url.data" {
		t.Errorf("TestLoadConfig failed, expect UrlListFile is :%s, but got:%s", "./data/url.data", config.Spider.UrlListFile)
	}

	if config.Spider.OutputDirectory != "./output" {
		t.Errorf("TestLoadConfig failed, expect outputDirectory is :%s, but got:%s", "./output", config.Spider.OutputDirectory)
	}

	if config.Spider.MaxDepth != 1 {
		t.Errorf("TestLoadConfig failed, expect MaxDepth is :%d, but got:%d", 1, config.Spider.MaxDepth)
	}

	if config.Spider.CrawlInterval != 1 {
		t.Errorf("TestLoadConfig failed, expect CrawlInterval is :%d, but got:%d", 1, config.Spider.CrawlInterval)
	}

	if config.Spider.CrawlTimeout != 1 {
		t.Errorf("TestLoadConfig failed, expect CrawlTimeout is :%d, but got:%d", 1, config.Spider.CrawlTimeout)
	}

	if config.Spider.TargetUrl != ".*.(htm|html)$" {
		t.Errorf("TestLoadConfig failed, expect threadCount is :%s, but got:%s", ".*.(htm|html)$", config.Spider.TargetUrl)
	}

	if config.Spider.ThreadCount != 8 {
		t.Errorf("TestLoadConfig failed, expect threadCount is :%d, but got:%d", 8, config.Spider.ThreadCount)
	}

}
