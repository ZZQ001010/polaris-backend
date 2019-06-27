package main

import (
	"fmt"

	cache "github.com/allstar/coding-carefree/components/cache"
	mail "github.com/allstar/coding-carefree/components/mail"
	utils "github.com/allstar/coding-carefree/utils"

	utils1 "github.com/polaris-team/polaris-backend/polaris-common/"

	"strconv"
)

func testMail() {
	//定义收件人
	mailTo := []string{
		"ainililia@163.com",
	}
	//邮件主题为"Hello"
	subject := "Hello"
	// 邮件正文
	body := "Good"
	mail.SendMail(mailTo, subject, body)
}

func testCache() {
	rp := cache.GetProxy()
	key := "hello4"
	rp.SetEx(key, "abc", 60)
	fmt.Println(rp.Get(key))
	fmt.Println(rp.Exist(key))
	fmt.Println(rp.Del(key))
	fmt.Println(rp.Exist(key))

	fmt.Println(rp.Incrby("abc", 1))
}

func testUtils() {
	fmt.Println(utils.VerifyPwdFormat("123123"))
	fmt.Println(utils.VerifyPwdFormat("aaa"))
	fmt.Println(utils.VerifyPwdFormat("a123"))
	fmt.Println(utils.VerifyPwdFormat("1aaa"))
	fmt.Println(utils.VerifyPwdFormat("111111111"))
	fmt.Println(utils.VerifyPwdFormat("aaaaaaaaa"))
	fmt.Println(utils.VerifyPwdFormat("11111111a"))
}

func testUtils1(){

}

func main() {
	//testMail()
	//testCache()
	// testUtils()
	fmt.Println(utils1.Md5V("123"))
}
