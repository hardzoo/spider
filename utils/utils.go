package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	log "github.com/golang/glog"
)

/*
* ParseHTML - 解析HTML
*
* PARAMS:
*   - node: html特面节点
*   - url: 当前页面基础路径
*   - re: 匹配正则表达式
*
* RETURNS:
*   - []string: 当前页面符合条件的url地址集合
 */
func ParseHTML(node *html.Node, url string, re *regexp.Regexp) []string {
	var urls []string
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "src" || attr.Key == "href" {
				if matched, err := regexp.MatchString("javascript:location", attr.Val); err == nil && matched {
					continue
				}

				if !re.MatchString(attr.Val) {
					continue
				}

				absoluteURL, err := GetAbsoluteAddress(attr.Val, url)
				if err != nil {
					log.Warningf("GetAbsoluteAddress failed: %s", err)
					continue
				}

				urls = append(urls, absoluteURL)
			}
		}
	}

	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		urls = append(urls, ParseHTML(childNode, url, re)...)
	}

	return urls
}

/*
* GetAbsoluteAddress - 获取url绝对地址
*
* PARAMS:
*   - relativePath: 抓取的相对地址
*   - basePath: 当前页面基础地址
*
* RETURNS:
*   - string: 返回组合之后的绝对地址url
*   - error: 成功返回nil
 */
func GetAbsoluteAddress(relativePath string, basePath string) (string, error) {
	relativePath = strings.TrimSpace(relativePath)
	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		log.Warningf("GetAbsoluteAddress parse url[%s] failed: %s", relativePath, err)
		return "", err
	}

	baseURL, err := url.Parse(basePath)
	if err != nil {
		log.Warningf("GetAbsoluteAddress parse url[%s] failed: %s", basePath, err)
		return "", err
	}

	res := baseURL.ResolveReference(relativeURL).String()
	return res, nil
}

/*
* LoadSeedFromFile - 从url.data里读取目标seed
*
* PARAMS:
*   - filename: 目标文件地址
*
* RETURNS:
*   - []string: 返回目标文件中列出的url
*   - error: 成功返回nil
 */
func LoadSeedFromFile(filename string) ([]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Warningf("ReadFile failed: %s ", err)
		return nil, err
	}

	var seeds []string
	if err := json.Unmarshal(bytes, &seeds); err != nil {
		log.Warningf("Unmarshal failed: %s", err)
		return nil, err
	}

	return seeds, nil
}

/*
* UrlToFilename - 将url用md5转换成文件名
*
* PARAMS:
*   - targetUrl: 目标地址
*
* RETURNS:
*   - string: 返回转换的文件名
 */
func UrlToFilename(targetUrl string) string {
	name := strings.Replace(targetUrl, "/", "%20F", -1)
	h := md5.New()
	h.Write([]byte(name))
	fileName := hex.EncodeToString(h.Sum(nil))
	return fileName
}

func IsFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
