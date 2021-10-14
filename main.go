package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/golang/glog"

	"github.com/zhuzeyu/spider/config"
	"github.com/zhuzeyu/spider/spider"
)

var (
	configFile  string // 配置文件
	logPath     string // 日志路径
	showVersion bool
	showHelp    bool
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "show spider help")
	flag.BoolVar(&showVersion, "version", false, "show spider version")
	flag.StringVar(&configFile, "c", "./conf/spider.conf", "assign config file path")
}

func main() {
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		version, err := GetVersion()
		if err != nil {
			fmt.Printf("show version of spider failed: %v\n", err)
		}
		fmt.Printf("the version of spider is v%s\n", version)
		os.Exit(0)
	}

	//设置默认日志路径
	flag.Lookup("log_dir").Value.Set("./log")

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Errorf("load config file failed: %v", err)
		Exit(1)
	}

	//创建output文件保存路径,如果路径存在，则不会创建
	if err = os.MkdirAll(conf.Spider.OutputDirectory, 0777); err != nil {
		log.Errorf("can not create output directory: %v", err)
	}

	//初始化Spider
	spiderConfig := conf.Spider
	sp := spider.NewSpider(&spiderConfig)
	if err = sp.InitCrawlTask(); err != nil {
		log.Errorf("spider init  task err: %s", err)
		Exit(1)
	}

	//最大并发
	for i := int64(0); i < sp.SpiderConfig.ThreadCount; i++ {
		go sp.Crawling()
	}

	for i := 0; i < len(sp.InitTask); i++ {
		sp.WaitGroup.Add(1)
		go func(ct spider.CrawlTask) {
			sp.Tasks <- ct
		}(sp.InitTask[i])
	}

	sp.WaitGroup.Wait()
	Exit(0)
}

func GetVersion() (string, error) {
	file, err := os.Open("./version")
	if err != nil {
		return "", fmt.Errorf("open version file failed: %v", err)
	}
	defer file.Close()
	version, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read content from version file failed: %v", err)
	}
	return string(version), nil
}

func Exit(num int) {
	log.Flush()
	time.Sleep(100 * time.Millisecond)
	os.Exit(num)
}
