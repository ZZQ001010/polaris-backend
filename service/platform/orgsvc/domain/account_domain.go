package domain

import (
	"github.com/galaxy-book/common/core/util/uuid"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"upper.io/db.v3"
)

func SetUserPassword(userId int64, password string, salt string, operatorId int64) errs.SystemErrorInfo{
	log.Infof("密码: %s, Slat: %s", password, salt)
	_, err := mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcId: userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcPassword: password,
		consts.TcPasswordSalt: salt,
		consts.TcUpdator: operatorId,
	})
	if err != nil{
		log.Error(err)
		return errs.SetUserPasswordError
	}
	return nil
}

func UnbindUserName(userId int64, addressType int) errs.SystemErrorInfo{
	upd := mysql.Upd{
		consts.TcUpdator: userId,
	}
	if addressType == consts.ContactAddressTypeEmail{
		upd[consts.TcEmail] = consts.BlankString
	}else if addressType == consts.ContactAddressTypeMobile{
		upd[consts.TcMobile] = consts.BlankString
	}

	_, err := mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcId: userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if err != nil{
		log.Error(err)
		return errs.UnBindLoginNameFail
	}
	return nil
}

func BindUserName(userId int64, addressType int, username string) errs.SystemErrorInfo{
	lockKey := consts.UserBindLoginNameLock + username
	uid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uid)
	if lockErr != nil{
		log.Error(lockErr)
		return errs.BindLoginNameFail
	}
	if suc{
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil{
				log.Error(err)
			}
		}()
	}else{
		return errs.BindLoginNameFail
	}

	//判断是否被其他账户绑定过
	_, err := GetUserInfoByLoginName(username)
	if err != nil{
		if err.Code() != errs.UserNotExist.Code(){
			return err
		}
	}else{
		if addressType == consts.ContactAddressTypeEmail{
			return errs.EmailAlreadyBindByOtherAccountError
		}else{
			return errs.MobileAlreadyBindOtherAccountError
		}
	}


	upd := mysql.Upd{
		consts.TcUpdator: userId,
	}
	if addressType == consts.ContactAddressTypeEmail{
		upd[consts.TcEmail] = username
	}else if addressType == consts.ContactAddressTypeMobile{
		upd[consts.TcMobile] = username
	}

	_, dbErr := mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcId: userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if dbErr != nil{
		log.Error(dbErr)
		return errs.UnBindLoginNameFail
	}
	return nil
}