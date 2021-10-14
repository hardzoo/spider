package config

import (
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Spider SpiderConfig `gcfg:"spider"`
}

type SpiderConfig struct {
	UrlListFile     string `gcfg:"urlListFile"`     //种子文件路径
	OutputDirectory string `gcfg:"outputDirectory"` //抓取结果存储目录
	MaxDepth        int64  `gcfg:"maxDepth"`        //最大抓取深度(种子为0级)
	CrawlInterval   int64  `gcfg:"crawlInterval"`   //抓取间隔. 单位: 秒
	CrawlTimeout    int64  `gcfg:"crawlTimeout"`    //抓取超时. 单位: 秒
	TargetUrl       string `gcfg:"targetUrl"`       //需要存储的目标网页URL pattern(正则表达式)
	ThreadCount     int64  `gcfg:"threadCount"`     //抓取routine数
}

func LoadConfig(filename string) (*Config, error) {
	var config Config
	err := gcfg.ReadFileInto(&config, filename)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
