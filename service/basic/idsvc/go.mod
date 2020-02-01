module github.com/galaxy-book/polaris-backend/service/basic/idsvc

go 1.13

replace github.com/galaxy-book/polaris-backend/common => ./../../../common

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/ainilili/go2sky v0.1.2
	github.com/galaxy-book/common v1.6.8
	github.com/galaxy-book/polaris-backend/common v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/polaris-team/dingtalk-sdk-golang v0.0.7
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
	upper.io/db.v3 v3.6.3+incompatible
)
