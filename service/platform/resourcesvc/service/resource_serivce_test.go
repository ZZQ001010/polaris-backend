package service

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/po"
	"github.com/galaxy-book/polaris-backend/service/platform/resourcesvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TestGetIdByPath(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetIdByPath(1029, "https://polaris-hd2.oss-cn-shanghai.aliyuncs.com/project/undraw_Projectpicture_update_jjgk.png", consts.OssResource))
	}))
}

func TestCreateResource(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("获取任务bo", func() {
			//var folderId int64 = 0
			var sourcetype int = 2
			createResourceBo := bo.CreateResourceBo{
				OrgId:      10109,
				ProjectId:  10115,
				Bucket:     "file",
				Path:       "/path",
				Name:       "北极星项目管理_20191206_R3_Bug list.xlsx",
				Size:       100,
				Suffix:     "xlsx",
				Type:       1,
				Md5:        "",
				OperatorId: 10209,
				//新增folderId用户文件管理创建资源	2019/12/12
				//FolderId: &folderId,
				//文件本地路径, 用于图片压缩
				DistPath:   "test",
				SourceType: &sourcetype,
			}
			_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
				_, err := CreateResource(createResourceBo, tx)
				if err != nil {
					log.Error(err)
					return err
				}
				return nil
			})
		})
	}))
}

func TestUpdateResourceName(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		var filename = "folder2-file2(gai)"
		var suffix = ".png"
		updateResourceBo := bo.UpdateResourceInfoBo{
			UserId:       1046,
			OrgId:        1016,
			ProjectId:    10116,
			ResourceId:   10587,
			FileName:     &filename,
			FileSuffix:   &suffix,
			UpdateFields: []string{"fileName"},
		}
		res, err := UpdateResourceInfo(updateResourceBo)
		fmt.Println(res)
		fmt.Println(err)
	}))
}

func TestUpdateResourceFolder(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		updateResourceBo := bo.UpdateResourceFolderBo{
			UserId:          10209,
			OrgId:           10109,
			ProjectId:       10115,
			ResourceIds:     []int64{10553, 10554, 10555},
			CurrentFolderId: 1022,
			TargetFolderID:  1021,
		}
		res, err := UpdateResourceFolder(updateResourceBo)
		fmt.Println(res)
		fmt.Println(err)
	}))
}

func TestDeleteResource(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("获取任务bo", func() {
			//var folderId int64 = 1021
			deleteResourceBo := bo.DeleteResourceBo{
				UserId:      1003,
				OrgId:       1003,
				ProjectId:   1006,
				ResourceIds: []int64{10505, 10503},
				//FolderId:    nil,
			}
			res, err := DeleteResource(deleteResourceBo)
			fmt.Println(res)
			fmt.Println(err)
		})
	}))
}

func TestGetResource(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("获取任务bo", func() {
			var folderId int64 = 0
			getResourceBo := bo.GetResourceBo{
				UserId:    10209,
				OrgId:     10109,
				ProjectId: 10115,
				FolderId:  &folderId,
				Page:      1,
				Size:      5,
			}
			res, err := GetResource(getResourceBo)
			fmt.Println(res.Total)
			for _, value := range res.List {
				fmt.Println(*value)
			}
			fmt.Println(err)
		})
	}))
}

//sql注入
func TestSql(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		conn, _ := mysql.GetConnect()
		conn.SetLogging(true)
		defer func() {
			if conn != nil {
				if err := conn.Close(); err != nil {
					logger.GetDefaultLogger().Info(strs.ObjectToString(err))
				}
			}
		}()
		size := 5
		page := 1
		objs := &[]po.PpmResResource{}
		cond := db.Cond{
			consts.TcIsDelete:  consts.AppIsNoDelete,
			consts.TcProjectId: 1006,
			consts.TcOrgId:     1003,
			//新增获取的文件类型 2019/12/27
			consts.TcSourceType: consts.OssPolicyTypeIssueResource,
			consts.TcName:       db.Like("%好%"),
		}
		order := "id desc"
		table := consts.TableResource
		mid := conn.Collection(table).Find(cond)
		mid = mid.Page(uint(page)).Paginate(uint(size))
		mid = mid.OrderBy(order)
		err := mid.All(objs)
		fmt.Println(err)
	}))
}
