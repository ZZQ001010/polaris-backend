module github.com/galaxy-book/polaris-backend/service/basic/appsvc

replace github.com/galaxy-book/polaris-backend/common => ./../../../common

replace github.com/galaxy-book/polaris-backend/facade => ./../../../facade

go 1.13

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/galaxy-book/common v1.6.8
	github.com/galaxy-book/polaris-backend/common v0.0.0-00010101000000-000000000000
	github.com/galaxy-book/polaris-backend/facade v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
	upper.io/db.v3 v3.6.3+incompatible
)
