module github.com/galaxy-book/polaris-backend/schedule

go 1.13

replace github.com/galaxy-book/polaris-backend/common => ./../common

replace github.com/galaxy-book/polaris-backend/facade => ./../facade

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/Jeffail/tunny v0.0.0-20190930221602-f13eb662a36a
	github.com/galaxy-book/common v1.6.8
	github.com/galaxy-book/polaris-backend/common v0.0.0-00010101000000-000000000000
	github.com/galaxy-book/polaris-backend/facade v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.4.0
	github.com/robfig/cron v1.2.0
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337
)
