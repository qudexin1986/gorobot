//web weixin client
package main

import (
	"fmt"
	"net/http"
	// "log"
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	// "strings"
)

type wxweb struct {
	uuid string
}

type jsonData map[string]interface{}

func (self *wxweb) getUuid(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/jslogin"
	params := make(map[string]interface{})
	params["appid"] = "wx782c26e4c19acffb"
	params["fun"] = "new"
	params["lang"] = "zh_CN"
	params["_"] = self.unixStr()
	data, _ := self._get(urlstr, params, false)
	re := regexp.MustCompile("\"(\\S+?)\"")
	find := re.FindAllString(data, -1)
	if len(find) > 0 {
		self.uuid = strings.Replace(find[0], "\"", "", -1)
		return true
	} else {
		return false
	}
}

func (self *wxweb) _run(desc string, f func(...interface{}) bool, args ...interface{}) {
	start := time.Now().UnixNano()
	fmt.Print(desc)
	var result bool
	if len(args) > 1 {
		result = f(args)
	} else if len(args) == 1 {
		result = f(args[0])
	} else {
		result = f()
	}
	useTime := fmt.Sprintf("%.5f", (float64(time.Now().UnixNano()-start) / 1000000000))
	if result {
		fmt.Println("成功,用时" + useTime + "秒")
	} else {
		fmt.Println("失败\n[*] 退出程序")
		os.Exit(1)
	}
}

func (self *wxweb) _post(urlstr string, params jsonData, jsonFmt bool) (string, error) {
	var err error
	res := ""
	v := url.Values{}
	for key, value := range params {
		v.Add(key, value.(string))
	}
	resp, err := http.PostForm(urlstr, v)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return string(body), nil
}

func (self *wxweb) _get(urlstr string, params jsonData, jsonFmt bool) (string, error) {
	var err error
	res := ""
	v := url.Values{}
	for key, value := range params {
		v.Add(key, value.(string))
	}
	resp, err := http.Get(urlstr + "?" + v.Encode())
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return string(body), nil
}

func (self *wxweb) unixStr() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

//TODO support linux
func (self *wxweb) genQRcode(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/qrcode/" + self.uuid
	urlstr += "?t=webwx"
	urlstr += "&_=" + self.unixStr()
	path := "qrcode.jpg"
	out, err := os.Create(path)
	defer out.Close()
	resp, err := http.Get(urlstr)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	_, err = io.Copy(out, bytes.NewReader(pix))
	if err != nil {
		return false
	} else {
		exec.Command("open", path).Run()
		return true
	}
}

func (self *wxweb) start() {
	fmt.Println("[*] 微信网页版 ... 开动")
	self._run("[*] 正在获取 uuid ... ", self.getUuid)
	self._run("[*] 正在获取 二维码 ... ", self.genQRcode)
}

func forgeheadget(urlstr string) string {

	client := &http.Client{}

	reqest, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}

	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate, sdch")
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Host", "login.weixin.qq.com")
	reqest.Header.Add("Referer", "https://wx.qq.com/")
	reqest.Header.Add("Upgrade-Insecure-Requests", "1")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	response, err := client.Do(reqest)
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}