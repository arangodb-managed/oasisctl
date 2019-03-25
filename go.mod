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
	golang.org/x/net v0.0.0-20190324223953-e3b2ff56ed87 // indirect
	golang.org/x/sys v0.0.0-20190322080309-f49334f85ddc // indirect
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2 // indirect
	google.golang.org/genproto v0.0.0-20190321212433-e79c0c59cdb5 // indirect
	google.golang.org/grpc v1.19.1
)
