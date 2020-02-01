package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/random"
	"github.com/galaxy-book/common/sdk/oss"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/core/util/str"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/resourcevo"
	"github.com/galaxy-book/feishu-sdk-golang/core/util/encrypt"
	"time"
)

func GetOssSignURL(req resourcevo.OssApplySignURLReqVo) (*vo.OssApplySignURLResp, errs.SystemErrorInfo) {
	oc := config.GetOSSConfig()
	if oc == nil {
		log.Error("oss缺少配置")
		return nil, errs.OssError
	}

	input := req.Input
	orgId := req.OrgId
	//userId := req.UserId

	ossKey := str.ParseOssKey(input.URL)
	ossKeyInfo := util.GetOssKeyInfo(ossKey)

	log.Infof("oss key %s, oss key info %s", ossKey, json.ToJsonIgnoreError(ossKeyInfo))
	if ossKeyInfo.OrgId != orgId {
		return nil, errs.NoOperationPermissions
	}
	//TODO 权限

	signUrl, signErr := oss.GetObjectUrl(ossKey, 60*60)
	if signErr != nil {
		log.Error(signErr)
		return nil, errs.BuildSystemErrorInfo(errs.OssError, signErr)
	}

	_, path := str.UrlParse(signUrl)
	//使用了自定义域名，要替换
	signUrl = util.JointUrl(oc.EndPoint, path)

	return &vo.OssApplySignURLResp{
		SignURL: signUrl,
	}, nil
}

func GetOssPostPolicy(req resourcevo.GetOssPostPolicyReqVo) (*vo.OssPostPolicyResp, errs.SystemErrorInfo) {
	oc := config.GetOSSConfig()
	if oc == nil {
		log.Error("oss缺少配置")
		return nil, errs.OssError
	}

	input := req.Input
	orgId := req.OrgId
	userId := req.UserId

	projectId := int64(0)
	issueId := int64(0)
	resourceFolderId := int64(0)
	if input.ProjectID != nil {
		projectId = *input.ProjectID
	}
	if input.IssueID != nil {
		issueId = *input.IssueID
	}
	if input.FolderID != nil {
		resourceFolderId = *input.FolderID
	}

	var policyConfig *config.OSSPolicyInfo = nil

	policyType := input.PolicyType

	toDay := time.Now()
	year, month, day := toDay.Date()

	//文件名
	fileName := random.RandomFileName()

	//定义callback
	callbackJson := ""
	switch policyType {
	case consts.OssPolicyTypeProjectCover:
		c := config.GetProjectCoverPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})

		policyConfig = &c
	case consts.OssPolicyTypeIssueResource:
		c := config.GetIssueResourcePolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyIssueIdConstName:   issueId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c

		callbackBody := fmt.Sprintf("{\"userId\":%d,\"orgId\":%d,\"type\":%d,\"host\":\"%s\",\"path\":\"%s\",\"filename\":\"%s\",\"issueId\":%d,\"bucket\":\"%s\",\"size\":%s,\"format\":%s,\"object\":%s,\"realName\":%s}", userId, orgId, policyType, oc.EndPoint, c.Dir, fileName, issueId, oc.BucketName, "${size}", "${imageInfo.format}", "${object}", "${x:filename}")
		callbackBo := bo.OssCallbackBo{
			CallbackBody:     callbackBody,
			CallbackUrl:      policyConfig.CallbackUrl,
			CallbackBodyType: consts.OssCallbackBodyTypeApplicationJson,
		}

		callbackJson = json.ToJsonIgnoreError(callbackBo)
		callbackJson = encrypt.BASE64([]byte(callbackJson))
	case consts.OssPolicyTypeIssueInputFile:
		c := config.GetIssueInputFilePolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeProjectResource:
		c := config.GetProjectResourcePolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c

		callbackBody := fmt.Sprintf("{\"userId\":%d,\"orgId\":%d,\"type\":%d,\"host\":\"%s\",\"path\":\"%s\",\"filename\":\"%s\",\"folderId\":%d,\"projectId\":%d,\"bucket\":\"%s\",\"size\":%s,\"format\":%s,\"object\":%s,\"realName\":%s}", userId, orgId, policyType, oc.EndPoint, c.Dir, fileName, resourceFolderId, projectId, oc.BucketName, "${size}", "${imageInfo.format}", "${object}", "${x:filename}")
		callbackBo := bo.OssCallbackBo{
			CallbackBody:     callbackBody,
			CallbackUrl:      policyConfig.CallbackUrl,
			CallbackBodyType: consts.OssCallbackBodyTypeApplicationJson,
		}

		callbackJson = json.ToJsonIgnoreError(callbackBo)
		callbackJson = encrypt.BASE64([]byte(callbackJson))
	case consts.OssPolicyTypeCompatTest:
		c := config.GetCompatTestPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName: orgId,
			consts.CacheKeyYearConstName:  year,
			consts.CacheKeyMonthConstName: int(month),
			consts.CacheKeyDayConstName:   day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeUserAvatar:
		c := config.GetUserAvatarPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeFeedback:
		c := config.GetFeedbackPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName: orgId,
			consts.CacheKeyYearConstName:  year,
			consts.CacheKeyMonthConstName: int(month),
			consts.CacheKeyDayConstName:   day,
		})
		policyConfig = &c
	}

	if policyConfig == nil {
		return nil, errs.BuildSystemErrorInfo(errs.OssPolicyTypeError)
	}

	//policyBo := oss.PostPolicyWithCallback(policyConfig.Dir, policyConfig.Expire, policyConfig.MaxFileSize, callbackJson)
	policyBo := oss.PostPolicy(policyConfig.Dir, policyConfig.Expire, policyConfig.MaxFileSize)

	//oss绑定自定义域名，所以覆盖host
	policyBo.Host = oc.EndPoint

	resp := &vo.OssPostPolicyResp{}
	copyErr := copyer.Copy(policyBo, resp)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	resp.FileName = fileName
	resp.Callback = callbackJson
	return resp, nil
}
