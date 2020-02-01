module github.com/galaxy-book/polaris-backend/facade

go 1.13

replace github.com/galaxy-book/polaris-backend/common => ./../common

require (
	github.com/galaxy-book/common v1.6.8
	github.com/galaxy-book/polaris-backend/common v0.0.0-00010101000000-000000000000
	github.com/polaris-team/dingtalk-sdk-golang v0.0.7
	gopkg.in/fatih/set.v0 v0.2.1
	upper.io/db.v3 v3.6.3+incompatible
)
