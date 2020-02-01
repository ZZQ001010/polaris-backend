package service

import (
	"fmt"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/maps"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/asyn"
	"github.com/galaxy-book/polaris-backend/common/core/util/pinyin"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
	"sort"
	"strings"
	"upper.io/db.v3"
)

type TagStyle struct {
	BgStyle   string
	FontStyle string
}

var TagDefaultStyle = []TagStyle{
	{BgStyle: "#E0E6E8", FontStyle: "#222B37"},
	{BgStyle: "#FFE1DC", FontStyle: "#801C0E"},
	{BgStyle: "#FFE3C4", FontStyle: "#813E00"},
	{BgStyle: "#FDEBB3", FontStyle: "#573D00"},
	{BgStyle: "#FDF5B8", FontStyle: "#6D6100"},
	{BgStyle: "#C5EDFF", FontStyle: "#015481"},
	{BgStyle: "#C5DCFF", FontStyle: "#002F82"},
	{BgStyle: "#EDF8C5", FontStyle: "#3F4F01"},
	{BgStyle: "#BDF1E9", FontStyle: "#004337"},
}

func judgeTagName(tagName string) (string, errs.SystemErrorInfo) {
	nameLengthLimit := 10
	tagName = strings.Trim(tagName, " ")
	if tagName == consts.BlankString || strs.Len(tagName) > nameLengthLimit {
		return "", errs.LengthOutOfLimit
	}

	return tagName, nil
}

func judgeTagStyle(bgStyle string) (string, string, errs.SystemErrorInfo) {
	isDefaultColor := false
	var fontStyle string
	for _, style := range TagDefaultStyle {
		if style.BgStyle == bgStyle {
			fontStyle = style.FontStyle
			isDefaultColor = true
			break
		}
	}
	if !isDefaultColor {
		return "", "", errs.NotDefaultStyle
	}

	return bgStyle, fontStyle, nil
}

func CreateTag(orgId, currentUserId int64, input vo.CreateTagReq) (*vo.Void, errs.SystemErrorInfo) {
	//用户角色权限校验
	authErr := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProTag, consts.RoleOperationCreate)
	if authErr != nil {
		log.Error(authErr)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, authErr)
	}

	newUUID := uuid.NewUuid()
	lockKey := fmt.Sprintf("%s%d", consts.CreateProjectTagLock, input.ProjectID)
	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
	if lockErr != nil{
		log.Error(lockErr)
		return nil, errs.TryDistributedLockError
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil{
				log.Error(err)
			}
		}()
	}else{
		//未获取到锁，直接响应错误信息
		return nil, errs.CreateTagFail
	}

	name, judgeErr := judgeTagName(input.Name)
	if judgeErr != nil {
		log.Error(judgeErr)
		return nil, judgeErr
	}
	bgStyle, fontStyle, judgeStyleErr := judgeTagStyle(input.BgStyle)
	if judgeStyleErr != nil {
		log.Error(judgeStyleErr)
		return nil, judgeStyleErr
	}

	//id, err := domain.InsertTag(orgId, currentUserId, name, bgStyle, fontStyle, input.ProjectID)
	insert := []bo.TagBo{}
	insert = append(insert, bo.TagBo{Name:name, BgStyle:bgStyle, FontStyle:fontStyle})
	tagInfo, err := domain.InsertTag(orgId, currentUserId, input.ProjectID, insert)
	if err != nil {
		return nil, err
	}

	asyn.Execute(func() {
		PushAddTagNotice(orgId, input.ProjectID, tagInfo)
	})

	return &vo.Void{ID: tagInfo[0].Id}, nil
}

func TagList(orgId int64, page int, size int, input vo.TagListReq) (*vo.TagList, errs.SystemErrorInfo) {
	cond := db.Cond{}
	cond[consts.TcOrgId] = orgId
	cond[consts.TcProjectId] = input.ProjectID
	cond[consts.TcIsDelete] = consts.AppIsNoDelete

	if input.Name != nil && *input.Name != "" {
		cond[consts.TcName] = db.Like(*input.Name + "%")
	}
	if input.NamePinyin != nil && *input.NamePinyin != "" {
		cond[consts.TcNamePinyin] = db.Like(*input.NamePinyin + "%")
	}

	total, list, err := domain.GetTagList(cond, page, size)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	tagVo := &[]*vo.Tag{}
	copyErr := copyer.Copy(list, tagVo)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	tagIds := []int64{}
	for _, tag := range *tagVo {
		tagIds = append(tagIds, tag.ID)
	}
	tagStat, err := domain.IssueTagStat(tagIds)

	tagStatMap := maps.NewMap("TagId", tagStat)
	for i, tag := range *tagVo {
		if _, ok := tagStatMap[tag.ID]; ok {
			stat := tagStatMap[tag.ID].(bo.IssueTagStatBo)
			(*tagVo)[i].UsedNum = stat.Total
		}
	}
	return &vo.TagList{
		Total: int64(total),
		List:  *tagVo,
	}, nil
}

func GetTagDefaultStyle() []string {
	res := []string{}
	for _, v := range TagDefaultStyle {
		res = append(res, v.BgStyle)
	}
	return res
}

func UpdateTag(orgId, currentUserId int64, input vo.UpdateTagReq) (*vo.Void, errs.SystemErrorInfo) {
	tagInfo, err := domain.GetTagInfo(input.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//用户角色权限校验
	authErr := domain.AuthProject(orgId, currentUserId, tagInfo.ProjectId, consts.RoleOperationPathOrgProTag, consts.RoleOperationModify)
	if authErr != nil {
		log.Error(authErr)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, authErr)
	}

	newUUID := uuid.NewUuid()
	lockKey := fmt.Sprintf("%s%d", consts.CreateTagLock, tagInfo.ProjectId)
	suc, lockErr := cache.TryGetDistributedLock(lockKey, newUUID)
	if lockErr != nil{
		log.Error(lockErr)
		return nil, errs.TryDistributedLockError
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, newUUID); err != nil{
				log.Error(err)
			}
		}()
	}else{
		//未获取到锁，直接响应错误信息
		return nil, errs.RepeatTag
	}

	upd := mysql.Upd{}
	if input.Name != nil {
		name, judgeErr := judgeTagName(*input.Name)
		if judgeErr != nil {
			log.Error(judgeErr)
			return nil, judgeErr
		}
		upd[consts.TcName] = name
		upd[consts.TcNamePinyin] = pinyin.ConvertToPinyin(name)
	}
	if input.BgStyle != nil {
		bgStyle, fontStyle, judgeStyleErr := judgeTagStyle(*input.BgStyle)
		if judgeStyleErr != nil {
			log.Error(judgeStyleErr)
			return nil, judgeStyleErr
		}
		upd[consts.TcBgStyle] = bgStyle
		upd[consts.TcFontStyle] = fontStyle
	}
	if len(upd) > 0 {
		updateErr := domain.UpdateTag(orgId, currentUserId, input.ID, upd)
		if updateErr != nil {
			log.Error(updateErr)
			return nil, updateErr
		}
	}

	asyn.Execute(func() {
		PushModifyTagNotice(orgId, tagInfo.ProjectId, tagInfo.Id)
	})
	return &vo.Void{ID:input.ID}, nil
}

func DeleteTag(orgId, currentUserId int64, input vo.DeleteTagReq) (*vo.Void, errs.SystemErrorInfo) {
	//用户角色权限校验
	authErr := domain.AuthProject(orgId, currentUserId, input.ProjectID, consts.RoleOperationPathOrgProTag, consts.RoleOperationDelete)
	if authErr != nil {
		log.Error(authErr)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, authErr)
	}

	if len(input.Ids) == 0 {
		return &vo.Void{ID:0}, nil
	}

	count, err := domain.DeleteTag(orgId, currentUserId, input.ProjectID, input.Ids)
	if err != nil {
		return nil, err
	}

	asyn.Execute(func() {
		PushRemoveTagNotice(orgId, input.ProjectID, input.Ids)
	})
	return &vo.Void{ID:count}, nil
}

func HotTagList(orgId int64, projectId int64) (*vo.TagList, errs.SystemErrorInfo) {
	cond := db.Cond{}
	cond[consts.TcOrgId] = orgId
	cond[consts.TcProjectId] = projectId
	cond[consts.TcIsDelete] = consts.AppIsNoDelete

	//获取所有项目标签
	total, list, err := domain.GetTagList(cond, -1, -1)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	tagIds := []int64{}
	for _, tagBo := range *list {
		tagIds = append(tagIds, tagBo.Id)
	}
	//从缓存获取标签使用数量
	tagStat, err := domain.IssueTagStatByCache(orgId, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	tagVo := &[]*vo.Tag{}
	copyErr := copyer.Copy(list, tagVo)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	tagStatMap := maps.NewMap("TagId", tagStat)
	for i, tag := range *tagVo {
		if _, ok := tagStatMap[tag.ID]; ok {
			stat := tagStatMap[tag.ID].(bo.IssueTagStatBo)
			(*tagVo)[i].UsedNum = stat.Total
		}
	}

	//根据热度从高到低排序
	sort.SliceStable(*tagVo, func(i, j int) bool {
		if (*tagVo)[i].UsedNum > (*tagVo)[j].UsedNum {
			return true
		}
		return false
	})

	return &vo.TagList{
		Total: int64(total),
		List:  *tagVo,
	}, nil
}