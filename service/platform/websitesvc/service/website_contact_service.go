package service

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/websitevo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/websitesvc/domain"
	"github.com/prometheus/common/log"
	"strings"
)

func RegisterWebSiteContact(reqVo websitevo.RegisterWebSiteContactReqVo) (*vo.Void, errs.SystemErrorInfo) {
	reqInput := reqVo.Input
	userId := reqVo.UserId
	orgId := reqVo.OrgId

	sex := 3
	if reqInput.Sex != nil {
		sex = *reqInput.Sex
		if sex != 1 && sex != 2 {
			log.Error("无效的性别")
			return nil, errs.BuildSystemErrorInfo(errs.InvalidSex)
		}
	}

	name := ""
	if reqInput.Name != nil {
		name = strings.TrimSpace(*reqInput.Name)
		//err := util.CheckUserNameLen(name)
		//if err != nil{
		//	log.Error(err)
		//	return nil, err
		//}
		isNameRight := format.VerifyUserNameFormat(name)
		if !isNameRight {
			return nil, errs.UserNameLenError
		}
	} else {
		userInfo, err := orgfacade.GetBaseUserInfoRelaxed("", orgId, userId)
		if err == nil && userInfo != nil {
			name = userInfo.Name
		}
	}

	//不考虑是否重复
	//isRepet, err := domain.CheckContactRepetition(reqInput.ContactInfo, consts.ContactStatusWait)
	//if isRepet{
	//	//重复就不插入了
	//	return &vo.Void{ID: 0,}, nil
	//}
	source := consts.AppSourceChannelWeb
	if reqInput.Source != nil {
		source = *reqInput.Source
	}

	remark := ""
	if reqInput.Remark != nil {
		remark = *reqInput.Remark
		if strs.Len(remark) > 512 {
			return nil, errs.ContactRemarkLenErr
		}
	}

	resourceInfo := ""
	if reqInput.ResourceUrls != nil {
		if len(reqInput.ResourceUrls) > 5 {
			return nil, errs.ContactResourceSizeErr
		}
		resourceInfo = json.ToJsonIgnoreError(reqInput.ResourceUrls)
		if strs.Len(resourceInfo) > 2048 {
			return nil, errs.ContactResourceInfoLenErr
		}
	}

	contactBo := bo.ContactBo{
		Name:         name,
		Mobile:       reqInput.ContactInfo,
		Sex:          sex,
		Status:       consts.ContactStatusWait,
		Source:       source,
		ResourceInfo: resourceInfo,
		Remark:       remark,
		Creator:      userId,
		Updator:      userId,
	}

	contactId, err := domain.RegisterWebSiteContact(contactBo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &vo.Void{ID: contactId}, nil
}
