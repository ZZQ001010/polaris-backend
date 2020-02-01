package handler

import (
	"errors"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/facade/resourcefacade"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const POS = "."

const DIR = "/"

const PATH = "resource"

var log = logger.GetDefaultLogger()

type Result struct {
	Code    int32      `json:"code"`
	Message string     `json:"message"`
	Data    UploadResp `json:"data"`
}

type UploadResp struct {
	FileList map[string]FileList `json:"fileList"`
}

type FileList struct {
	Url      string `json:"url"`
	Size     string `json:"size"`
	SourceId int64  `json:"sourceId"`
	FileName string `json:"fileName"`
	DistPath string `json:"distPath"`
}

func FileUploadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		runMode := config.GetConfig().Application.RunMode
		if runMode == 1 {
			errorHandle(errs.BuildSystemErrorInfo(errs.RunModeUnsupportUpload), c.Writer)
			return
		}
		r := c.Request
		// 整体限制 100Mb
		r.Body = http.MaxBytesReader(c.Writer, r.Body, 10 << 20)
		_, err := c.MultipartForm()
		if err != nil {
			log.Error(err)
			errorHandle(errs.BuildSystemErrorInfo(errs.FileTooLarge, err), c.Writer)
			return
		}

		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
		if userErr != nil {
			log.Error(userErr)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
			return
		}

		projectId, issueId, policyType := baseParam(r)
		result, err := uploadFile(c, cacheUserInfo, projectId, issueId, policyType)
		if err != nil {
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		} else {
			//jsonStr, _ := json.Marshal(result)
			jsonStr := json.ToJsonIgnoreError(result)
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Write([]byte(jsonStr))
			return
		}
	}
}

func FileReadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		readFile(c)
	}
}

func baseParam(r *http.Request) (int64, int64, int) {
	var issueId, projectId int64
	var policyType int
	queries := r.URL.Query()
	form := r.Form
	if len(queries["projectId"]) > 0 {
		projectId, _ = strconv.ParseInt(queries["projectId"][0], 10, 64)
	}
	if len(queries["issueId"]) > 0 {
		issueId, _ = strconv.ParseInt(queries["issueId"][0], 10, 64)
	}
	if len(queries["policyType"]) > 0 {
		policyType, _ = strconv.Atoi(queries["policyType"][0])
	}
	if len(form["projectId"]) > 0 {
		projectId, _ = strconv.ParseInt(form["projectId"][0], 10, 64)
	}
	if len(form["issueId"]) > 0 {
		issueId, _ = strconv.ParseInt(form["issueId"][0], 10, 64)
	}
	if len(form["policyType"]) > 0 {
		policyType, _ = strconv.Atoi(form["policyType"][0])
	}
	return projectId, issueId, policyType
}

func ImportDataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		_, err := c.MultipartForm()
		if err != nil {
			log.Error(err)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		}

		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
		if userErr != nil {
			log.Error(userErr)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
			return
		}

		projectId, issueId, policyType := baseParam(r)
		result, err := uploadFile(c, cacheUserInfo, projectId, issueId, policyType)
		if err != nil {
			log.Error(err)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		} else {
			log.Infof("上传文件成功，准备导入数据，result： %s", json.ToJsonIgnoreError(result))

			if len(result.Data.FileList) == 0 {
				errorHandle(errs.ImportFileNotExist, c.Writer)
				return
			}
			url := ""
			urlType := consts.UrlTypeDistPath
			for _, v := range result.Data.FileList {
				//url = v.Url
				//改为本地路径
				url = v.DistPath
				break
			}

			importInput := projectvo.ImportIssuesReqVo{
				UserId: cacheUserInfo.UserId,
				OrgId:  cacheUserInfo.OrgId,
				Input: vo.ImportIssuesReq{
					URL:       url,
					ProjectID: projectId,
					URLType:   urlType,
				},
			}

			log.Infof("批量导入任务请求结构体 %s", json.ToJsonIgnoreError(importInput))
			resp := projectfacade.ImportIssues(importInput)
			if resp.Failure() {
				log.Error(resp.Error())
				var errMsg string
				if resp.Error().Code() == errs.FileParseFail.Code() {
					data := resp.Message[len(errs.FileParseFail.Message()):len(resp.Message)]
					errMsg = "{\"code\":" + strconv.Itoa(resp.Error().Code()) + ",\"message\":\"" + errs.FileParseFail.Message() + "\", \"data\":" + data + "}"
				} else {
					errMsg = "{\"code\":" + strconv.Itoa(resp.Error().Code()) + ",\"message\":\"" + resp.Message + "\"}"
				}
				errorHandle(errors.New(errMsg), c.Writer)
				return
			}

			result := vo.Err{
				Code:    200,
				Message: fmt.Sprintf("共%d条任务数据上传成功！", resp.Data),
			}
			//jsonStr, _ := json.Marshal(result)
			jsonStr := json.ToJsonIgnoreError(result)
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Write([]byte(jsonStr))
			return
		}
	}
}

//读取文件
func readFile(c *gin.Context) {
	w := c.Writer
	path := c.Param("path")
	file, err := os.Open(config.GetConfig().OSS.RootPath + path)
	if err != nil {
		log.Error(err)
		errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), w)
		return
	}

	defer file.Close()
	buff, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		errorHandle(errs.FileReadFail, w)
		return
	}
	w.Write(buff)
}

// 统一错误输出接口
func errorHandle(err error, w http.ResponseWriter) {
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func getFileName(ext string, r *http.Request, projectId, issueId int64, policyType int) (string, string, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(r.Context())
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return consts.BlankString, consts.BlankString, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	reqVo := resourcevo.GetOssPostPolicyReqVo{
		Input: vo.OssPostPolicyReq{
			ProjectID:  &projectId,
			IssueID:    &issueId,
			PolicyType: policyType,
		},
		OrgId: cacheUserInfo.OrgId,
	}

	resp := resourcefacade.GetOssPostPolicy(reqVo)

	if resp.Failure() {
		return consts.BlankString, consts.BlankString, resp.Error()
	}
	commonPath := resp.GetOssPostPolicy.Dir

	rootPath := config.GetOSSConfig().RootPath
	dstPath := rootPath + DIR + commonPath
	suffixSplit := strings.Split(ext, POS)
	suffix := suffixSplit[len(suffixSplit)-1]
	fileName := resp.GetOssPostPolicy.FileName + POS + suffix
	dstFile := dstPath + DIR + fileName

	_, err1 := os.Stat(dstPath)
	res := os.IsNotExist(err1)
	if res == true {
		os.MkdirAll(dstPath, os.ModePerm)
	}

	relatePath := commonPath + DIR + fileName
	return dstFile, relatePath, nil
}

//上传文件
func uploadFile(c *gin.Context, cacheUserInfo *bo.CacheUserInfoBo, projectId, issueId int64, policyType int) (Result, error) {
	result := Result{
		Code:    200,
		Message: "success",
		Data: UploadResp{
			FileList: map[string]FileList{},
		},
	}
	r := c.Request

	//支持多文件上传
	for k, _ := range r.MultipartForm.File {
		file, handler, err := r.FormFile(k)
		if err != nil {
			log.Error(err)
			return result, err
		}
		defer file.Close()

		dstFile, relatePath, err := getFileName(handler.Filename, r, projectId, issueId, policyType)
		if err != nil {
			log.Error(err)
			return result, err
		}

		fp, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Error(err)
			return result, err
		}
		defer fp.Close()

		size, err := io.Copy(fp, file)
		if err != nil {
			log.Error(err)
			return result, err
		}

		url := config.GetOSSConfig().LocalDomain + "/" + relatePath
		////单机部署
		//if runMode == 2 {
		//	url = fmt.Sprintf("%s/read/%s", config.GetOSSConfig().EndPoint, relatePath)
		//}
		respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{
			CreateResourceBo: bo.CreateResourceBo{
				Name:       handler.Filename,
				Suffix:     util.ParseFileSuffix(handler.Filename),
				Size:       handler.Size,
				Path:       url,
				OrgId:      cacheUserInfo.OrgId,
				OperatorId: cacheUserInfo.UserId,
				Type:       consts.LocalResource,
				DistPath:   dstFile,
			},
		})
		if respVo.Failure() {
			log.Error(respVo.Error())
		}

		result.Data.FileList[k] = FileList{
			Url:      url,
			Size:     strconv.FormatInt(size/1024, 10) + "KB",
			SourceId: respVo.ResourceId,
			FileName: handler.Filename,
			DistPath: dstFile,
		}
	}

	return result, nil
}
