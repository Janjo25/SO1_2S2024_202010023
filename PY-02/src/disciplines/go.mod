module github.com/janjo25/src/disciplines

go 1.23

replace github.com/janjo25/proto => ../../proto

require (
	github.com/janjo25/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.67.1
)

require (
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
