module github.com/arangodb-managed/oasis

go 1.12

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3

require (
	github.com/arangodb-managed/apis v0.12.1
	github.com/dustin/go-humanize v1.0.0
	github.com/gogo/protobuf v1.2.1
	github.com/grpc-ecosystem/grpc-gateway v1.9.2 // indirect
	github.com/rs/zerolog v1.14.3
	github.com/ryanuber/columnize v2.1.0+incompatible
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	golang.org/x/net v0.0.0-20190619014844-b5b0513f8c1b // indirect
	golang.org/x/sys v0.0.0-20190618155005-516e3c20635f // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20190611190212-a7e196e89fd3 // indirect
	google.golang.org/grpc v1.21.1
)
