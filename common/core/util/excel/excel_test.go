package excel

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"testing"
)

func TestGenerateCSVFromXLSXFile(t *testing.T) {
	t.Log(GenerateCSVFromXLSXFile("issue.xlsx", consts.UrlTypeDistPath, 0, 0, []int64{1, 2}))
}

func TestGenerateCSVFromXLSXFile2(t *testing.T) {
	a := []int64{1, 2}
	fmt.Println(slice.Contain(a, 1))
}
