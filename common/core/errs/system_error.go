package errs

import (
	"github.com/galaxy-book/common/core/errors"
)

type SystemErrorInfo errors.SystemErrorInfo

func BuildSystemErrorInfo(resultCodeInfo errors.ResultCodeInfo, e ...error) errors.SystemErrorInfo {
	return errors.BuildSystemErrorInfo(resultCodeInfo, e...)
}

func BuildSystemErrorInfoWithMessage(resultCodeInfo errors.ResultCodeInfo, message string) errors.SystemErrorInfo {
	resultCodeInfo.SetMessage(resultCodeInfo.Message() + message)
	return resultCodeInfo
}

//add system ResultCodeInfo
func AddResultCodeInfo(code int, message string, langCode string) errors.ResultCodeInfo {
	return errors.AddResultCodeInfo(code, message, langCode)
}

func GetResultCodeInfoByCode(code int) errors.ResultCodeInfo {
	return errors.GetResultCodeInfoByCode(code)
}
