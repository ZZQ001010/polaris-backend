package str

import (
	"fmt"
	"testing"
)

func TestUrlParse(t *testing.T) {
	fmt.Println(UrlParse("http://www.baidu.com/a.jpg"))
}

func TestParseOssKey(t *testing.T) {
	t.Log(ParseOssKey("http://www.baidu.com/a.jpg"))
	t.Log(ParseOssKey("/a.jpg"))
	t.Log(ParseOssKey("a.jpg"))
}

func TestCountStrByGBK(t *testing.T) {
	t.Log(CountStrByGBK("aa啊a12啊a"))
}