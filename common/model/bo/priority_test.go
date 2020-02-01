package bo

import (
	"github.com/galaxy-book/common/core/util/json"
	"testing"
)

func TestSortPriorityBo(t *testing.T) {

	priorities := []PriorityBo{
		{
			Name: "1",
			Sort: 1,
		},
		{
			Name: "2",
			Sort: 2,
		},
		{
			Name: "4",
			Sort: 4,
		},
		{
			Name: "3",
			Sort: 3,
		},
	}

	SortPriorityBo(priorities)
	t.Log(json.ToJsonIgnoreError(priorities))
}