package domain

import (
	"fmt"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/date"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/idvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/po"
	"strconv"
	"strings"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()
var (
	_defaultPage      = int64(1)
	_defaultSize      = int64(10)
	_select_sql       = "SELECT a.`id`, a.`org_id`, a.`module1`, a.`module2_id`, a.`module2`, a.`module3_id`, a.`module3`, a.`oper_code`, a.`oper_obj_id`, a.`oper_obj_type`, a.`oper_obj_property`, a.`relation_obj_id`, a.`relation_type`, a.`new_value`, a.`old_value`, a.`ext`, a.`creator`, a.`create_time`, a.`is_delete` FROM `ppm_tre_trends` AS a WHERE "
	_select_count_sql = "SELECT count(*) AS `id` FROM `ppm_tre_trends` AS a WHERE "
)

/**
创建动态
*/
func CreateTrends(trendsBo *bo.TrendsBo, tx ...sqlbuilder.Tx) (*int64, errs.SystemErrorInfo) {
	trendsPo := &po.PpmTreTrends{}
	err1 := util.ConvertObject(&trendsBo, &trendsPo)
	if err1 != nil {
		return nil, err1
	}
	trendsPo.IsDelete = consts.AppIsNoDelete

	respVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: consts.TableTrends})
	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}

	trendsId := respVo.Id
	trendsPo.Id = trendsId

	err2 := dao.InsertTrends(*trendsPo, tx...)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.TrendsCreateError)
	}
	trendsBo.Id = trendsId
	return &trendsId, nil
}

func CreateTrendsBatch(trendBos []bo.TrendsBo, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	trendsPos := &[]po.PpmTreTrends{}
	_ = copyer.Copy(trendBos, trendsPos)

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableTrends, len(*trendsPos))
	if err != nil{
		log.Error(err)
		return err
	}

	for i, _ := range *trendsPos{
		(*trendsPos)[i].IsDelete = consts.AppIsNoDelete
		(*trendsPos)[i].Id = ids.Ids[i].Id
	}

	dbErr := dao.InsertTrendsBatch(*trendsPos, tx...)
	if dbErr != nil{
		log.Error(dbErr)
		return errs.MysqlOperateError
	}

	return nil
}

func QueryTrends(condBo *bo.TrendsQueryCondBo) (*bo.TrendsPageBo, errs.SystemErrorInfo) {
	errorParam := checkParam(condBo)

	if errorParam != nil {
		return nil, errorParam
	}

	if condBo.Page == nil {
		condBo.Page = &_defaultPage
	}
	if condBo.Size == nil {
		condBo.Size = &_defaultSize
	}

	//组装参数 和拼接 sql
	params, _sql := assemblyParamsAndSql(condBo)

	orderSql := " order by create_time asc"
	if condBo.OrderType != nil && *condBo.OrderType == 2 {
		orderSql = " order by create_time desc"
	}

	fCount := (*condBo.Page - 1) * *condBo.Size
	limitSql := ""
	extraSql := ""
	if *condBo.Page > 0 && *condBo.Size >= 0 {
		if condBo.LastTrendID != nil && *condBo.LastTrendID != 0 {
			if condBo.OrderType != nil && *condBo.OrderType == 2 {
				//降序
				extraSql += " AND a.`id` < " + strconv.FormatInt(*condBo.LastTrendID, 10) + "  "
			} else {
				//升序
				extraSql += " AND a.`id` < " + strconv.FormatInt(*condBo.LastTrendID, 10) + "  "
			}
		}
		limitSql = " limit " + strconv.FormatInt(fCount, 10) + " , " + strconv.FormatInt(*condBo.Size, 10) + " "
	}

	trendsPos := &[]po.PpmTreTrends{}
	trendsIdPos := &[]po.PpmTreTrends{}

	err := mysql.SelectByQuery(_select_sql+_sql+extraSql+orderSql+limitSql, trendsPos, *params...)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	err2 := mysql.SelectByQuery(_select_count_sql+_sql, trendsIdPos, *params...)
	if err2 != nil || len(*trendsIdPos) < 1 {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	trendsBos := make([]bo.TrendsBo, len(*trendsPos))
	for i, v := range *trendsPos {
		trendsBos[i] = bo.TrendsBo{}
		err2 := util.ConvertObject(&v, &trendsBos[i])
		if err2 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
		}
	}

	trendsPageBo := &bo.TrendsPageBo{
		Page:  *condBo.Page,
		Size:  *condBo.Size,
		List:  &trendsBos,
		Total: (*trendsIdPos)[0].Id,
	}

	return trendsPageBo, nil
}

func assemblyParamsAndSql(condBo *bo.TrendsQueryCondBo) (*[]interface{}, string) {
	params := &[]interface{}{}

	_sql := " "

	_sql = _sql + " a.`org_id` = ? and a.`is_delete` = ? "
	*params = append(*params, condBo.OrgId, 2)

	if condBo.ObjType != nil {
		_sql = _sql + " AND ( ( a.`module2_id` = ? AND a.`module2` = ? ) OR ( a.`module3_id` = ? AND a.`module3` = ? ) OR ( a.`oper_obj_id` = ? AND a.`oper_obj_type` = ? ) OR ( a.`relation_obj_id` = ? AND a.`relation_obj_type` = ? ) )  "
		//_sql = _sql + " AND ( ( a.`module2_id` = ? AND a.`module2` = ? ) OR ( a.`module3_id` = ? AND a.`module3` = ? ) )  "
		*params = append(*params, condBo.ObjId, condBo.ObjType, condBo.ObjId, condBo.ObjType, condBo.ObjId, condBo.ObjType, condBo.ObjId, condBo.ObjType)
		//*params = append(*params, condBo.ObjId, condBo.ObjType, condBo.ObjId, condBo.ObjType)
	}

	if condBo.Type != nil {
		if *condBo.Type == 2 {
			//任务评论
			_sql = _sql + " AND a.`relation_type` = ?"
			*params = append(*params, consts.TrendsRelationTypeCreateIssueComment)
		} else if *condBo.Type == 1 {
			//任务动态（不包括评论）
			_sql = _sql + " AND a.`relation_type` != ?"
			*params = append(*params, consts.TrendsRelationTypeCreateIssueComment)
		} else if *condBo.Type == 3 {
			//项目动态
			//_sql = _sql + " AND a.`relation_type` in ?"
			//*params = append(*params, consts.ValidRelationTypesOfProject)
			relationStr := strings.Replace(strings.Trim(fmt.Sprint(consts.ValidRelationTypesOfProject), "[]"), " ", "\",\"", -1)
			_sql = _sql + " AND a.`relation_type` in (\"" + relationStr + "\")"
		}
	}

	if condBo.StartTime != nil {
		_sql = _sql + " AND a.`create_time` >= ?  "
		*params = append(*params, date.FormatTime(*condBo.StartTime))
	}

	if condBo.EndTime != nil {
		_sql = _sql + " AND a.`create_time` <= ?  "
		*params = append(*params, date.FormatTime(*condBo.EndTime))
	}

	if condBo.OperId != nil {
		_sql = _sql + " AND a.`creator` = ?  "
		*params = append(*params, condBo.OperId)
	}

	return params, _sql
}

func checkParam(condBo *bo.TrendsQueryCondBo) errs.SystemErrorInfo {
	if condBo.OrgId <= 0 {
		return errs.BuildSystemErrorInfo(errs.IllegalityOrg)
	}

	if condBo.ObjType != nil && (condBo.ObjId == nil || *condBo.ObjId <= 0) {
		return errs.BuildSystemErrorInfo(errs.TrendsObjIdNilError)
	}

	if condBo.ObjId != nil && condBo.ObjType == nil {
		return errs.BuildSystemErrorInfo(errs.TrendsObjTypeNilError)
	}

	return nil
}
