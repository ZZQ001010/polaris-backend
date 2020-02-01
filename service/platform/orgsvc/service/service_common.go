package service

import "github.com/galaxy-book/common/core/util/slice"

func NeedUpdate(updateFields []string, field string) bool {
	if updateFields == nil || len(updateFields) == 0 {
		return true
	}
	bol, err := slice.Contain(updateFields, field)
	if err != nil {
		return false
	}
	return bol
}
