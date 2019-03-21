module github.com/arangodb-managed/oasis

go 1.12

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3

require (
	github.com/arangodb-managed/apis v0.3.0
	github.com/dustin/go-humanize v1.0.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.8.5 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/rs/zerolog v1.13.0
	github.com/ryanuber/columnize v2.1.0+incompatible
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	golang.org/x/net v0.0.0-20190320064053-1272bf9dcd53 // indirect
	golang.org/x/sys v0.0.0-20190321052220-f7bb7a8bee54 // indirect
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2 // indirect
	google.golang.org/genproto v0.0.0-20190307195333-5fe7a883aa19 // indirect
	google.golang.org/grpc v1.19.1
)
