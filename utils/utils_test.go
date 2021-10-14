package utils

import (
	"regexp"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParseHTML(t *testing.T) {
	ctx := `<!DOCTYPE html>
<html>
    <head>
        <meta charset=utf8>
        <title>TestParseHTML</title>
    </head>
    <body>
        <ul>
            <li><a href=test1.html>test1</a></li>
            <li><a href='test2/test2.html'>test2</a></li>
            <li><a href='javascript:location.href="test3.html"'>test3</a></li>
        </ul>
    </body>
</html>
`

	re, err := regexp.Compile(".*.(htm|html)$")
	if err != nil {
		t.Errorf("failed to compile regexp: %s", err)
	}

	htmlNode, err := html.Parse(strings.NewReader(ctx))
	if err != nil {
		t.Errorf("failed to parse content: %s", err)
	}

	urls := ParseHTML(htmlNode, "http://localhost/index.html", re)

	if len(urls) != 2 {
		t.Errorf("TestParseHTML failed, Expect 2 urls, but got %d urls", len(urls))
	}

	if urls[1] != "http://localhost/test2/test2.html" {
		t.Errorf("TestParseHTML failed, Expect url [%s], but got [%v]",
			"http://localhost/test2/test2.html", urls[1])
	}
}

func TestGetAbsoluteAddress(t *testing.T) {
	relativePath := "//www.baidu.com/cache/sethelp/help.html"
	basePath := "http://www.baidu.com"
	res, err := GetAbsoluteAddress(relativePath, basePath)
	if err != nil {
		t.Errorf("GetAbsoluteAddress failed, err: %s", err)
	}

	if res != "http://www.baidu.com/cache/sethelp/help.html" {
		t.Errorf("TestGetAbsoluteAddress failed, Expect res [%s], but got [%s]",
			"http://www.baidu.com/cache/sethelp/help.html", res)
	}

	basePath = " http://www.baidu.com"
	res, err = GetAbsoluteAddress(relativePath, basePath)
	if err == nil {
		t.Errorf("TestGetAbsoluteAddress failed, Expect GetAbsoluteAddress failed, but got nothing")
	}
}

func TestLoadSeedFromFile(t *testing.T) {
	fileName := "./data/url.data"
	seeds, err := LoadSeedFromFile(fileName)
	if err == nil {
		t.Errorf("TestLoadSeedFromFile err, fileName[%s] is not exist, but LoadSeedFromFile return no err", fileName)
	}

	fileName = "../data/url.data"
	seeds, err = LoadSeedFromFile(fileName)
	if err != nil {
		t.Errorf("TestLoadSeedFromFile err :%v", err)
	}
	if len(seeds) != 2 {
		t.Errorf("TestLoadSeedFromFile failed, expect 2 seeds, but got %d", len(seeds))
	}

	if seeds[0] != "http://www.baidu.com" {
		t.Errorf("TestLoadSeedFromFile failed, expect seed is :%s, but got:%s",
			"http://www.baidu.com", seeds[0])
	}
}

func TestUrlToFilename(t *testing.T) {
	targetUrl := "http://www.baidu.com/cache/sethelp/help.html"
	filename := UrlToFilename(targetUrl)
	if filename != "93a27aab7c56d581f91841bcd4bfc83c" {
		t.Errorf("TestUrlToFilename failed, expect filename is :%s, but got:%s",
			"93a27aab7c56d581f91841bcd4bfc83c", filename)
	}
}

func TestIsFileExist(t *testing.T) {
	fileName := "../conf/spider.conf"
	if !IsFileExist(fileName) {
		t.Errorf("TestIsFileExist failed, file [%s] is exist, but got false",
			"../conf/spider.conf")
	}
	fileName = "./conf/spider.conf"
	if IsFileExist(fileName) {
		t.Errorf("TestIsFileExist failed, file [%s] is not exist, but got true",
			"./conf/spider.conf")
	}
}
