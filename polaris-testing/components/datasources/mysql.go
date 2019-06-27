package datasources

import (
	"github.com/golang/glog"

	"upper.io/db.v3/lib/sqlbuilder"
	mysql "upper.io/db.v3/mysql"
)

var settings = mysql.ConnectionURL{
	"root",
	"mysqldev",
	"allstar",
	"192.168.1.170",
	"",
	nil,
}

func GetMysqlConnect() sqlbuilder.Database {
	sess, err := mysql.Open(settings)
	if err != nil {
		glog.Fatalf("Error getting mysql connection: %s", err)
		return nil
	}
	return sess
}
