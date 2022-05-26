module github.com/shestakovda/fdbx/v2

go 1.13

require (
	github.com/apple/foundationdb/bindings/go v0.0.0-20201222225940-f3aef311ccfb
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/flatbuffers v1.12.0
	github.com/kr/text v0.2.0 // indirect
	github.com/shestakovda/errx v1.2.0
	github.com/shestakovda/typex v1.0.0
	github.com/stretchr/testify v1.6.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/shestakovda/fdbx/v2 v2.0.3-dev.3 => github.com/triumphpc/fdbx/v2 v2.0.4
