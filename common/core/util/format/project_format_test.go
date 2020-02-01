package format

import "testing"
import "github.com/magiconair/properties/assert"

func TestVerifyProjectNameFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyProjectNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyProjectNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "/"
	}
	suc = VerifyProjectNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyProjectNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "好"
	}
	testStr = testStr + "1"
	suc = VerifyProjectNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "asd/*你好04"
	suc = VerifyProjectNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)
}

func TestVerifyProjectPreviousCodeFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 10; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyProjectPreviousCodeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 10; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyProjectPreviousCodeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 10; i++ {
		testStr = testStr + "/"
	}
	suc = VerifyProjectPreviousCodeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = ""
	for i := 0; i < 10; i++ {
		testStr = testStr + "你"
	}
	suc = VerifyProjectPreviousCodeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "/*你01"
	suc = VerifyProjectPreviousCodeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "adsa01A"
	suc = VerifyProjectPreviousCodeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)
}

func TestVerifyProjectRemarkFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyProjectRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyProjectRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "/"
	}
	suc = VerifyProjectRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyProjectRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "好"
	}
	testStr = testStr + "101"
	suc = VerifyProjectRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "0sd你/**"
	suc = VerifyProjectRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)
}

func TestVerifyIssueNameFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 50; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyIssueNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 50; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyIssueNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 50; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyIssueNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 50; i++ {
		testStr = testStr + "/"
	}
	suc = VerifyIssueNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 50; i++ {
		testStr = testStr + "/"
	}
	testStr = testStr + "你0156"
	suc = VerifyIssueNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "0asdsa/**好"
	suc = VerifyIssueNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)
}

func TestVerifyIssueRemarkFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 10000; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyIssueRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 10000; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyIssueRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 10000; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyIssueRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 10000; i++ {
		testStr = testStr + "*"
	}
	suc = VerifyIssueRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 10000; i++ {
		testStr = testStr + "1"
	}
	testStr = testStr + "你好"
	suc = VerifyIssueRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = ""
	for i := 0; i < 10000; i++ {
		testStr = testStr + "/"
	}
	testStr = testStr + "***"
	suc = VerifyIssueRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "0656你好*/.sadaAsas"
	suc = VerifyIssueRemarkFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)
}

func TestVerifyIssueCommenFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyIssueCommenFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyIssueCommenFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyIssueCommenFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "*"
	}
	suc = VerifyIssueCommenFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 500; i++ {
		testStr = testStr + "1"
	}
	testStr = testStr + "0595"
	suc = VerifyIssueCommenFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "0656你好*/.sadaAsas"
	suc = VerifyIssueCommenFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)
}

func TestVerifyProjectObjectTypeNameFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyProjectObjectTypeNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "a"
	}
	suc = VerifyProjectObjectTypeNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "/"
	}
	suc = VerifyProjectObjectTypeNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "号"
	}
	suc = VerifyProjectObjectTypeNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 30; i++ {
		testStr = testStr + "号"
	}
	testStr = testStr + "01ji"
	suc = VerifyProjectObjectTypeNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "你好01LSADlasda*/***"
	suc = VerifyProjectObjectTypeNameFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)
}


func TestVerifyProjectNoticeFormat(t *testing.T) {
	testStr := ""
	for i := 0; i < 2000; i++ {
		testStr = testStr + "1"
	}
	suc := VerifyProjectNoticeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 2000; i++ {
		testStr = testStr + "*"
	}
	suc = VerifyProjectNoticeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 2000; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyProjectNoticeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)

	testStr = ""
	for i := 0; i < 2001; i++ {
		testStr = testStr + "好"
	}
	suc = VerifyProjectNoticeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, false)

	testStr = "好121*/n\\"
	suc = VerifyProjectNoticeFormat(testStr)
	t.Log(suc)
	assert.Equal(t, suc, true)
}