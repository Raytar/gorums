.PHONY: all qf ordering future

all: qf ordering future

qf: qf/qf.pb.go qf/qf_grpc.pb.go qf/qf_gorums.pb.go

ordering: ordering/order.pb.go ordering/order_gorums.pb.go

future: future/future.pb.go future/future_gorums.pb.go

%.pb.go : %.proto
	@protoc --go_out=paths=source_relative:. -I=../:. $^

%_grpc.pb.go : %.proto
	@protoc --go-grpc_out=paths=source_relative:. -I=../:. $^

%_gorums.pb.go : %.proto
	@protoc --gorums_out=trace=true,paths=source_relative:. -I=../:. $^
