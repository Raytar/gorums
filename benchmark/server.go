package benchmark

import (
	context "context"
	fmt "fmt"
	log "log"
	net "net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	stats *Stats
	UnimplementedBenchmarkServer
}

func (srv *server) StartServerBenchmark(_ context.Context, req *StartRequest) (resp *StartResponse, _ error) {
	srv.stats.Clear()
	srv.stats.Start()
	resp = &StartResponse{}
	return
}

func (srv *server) StopServerBenchmark(_ context.Context, req *StopRequest) (resp *Result, _ error) {
	srv.stats.End()
	resp = srv.stats.GetResult()
	return
}

func (srv *server) StartBenchmark(_ context.Context, req *StartRequest) (resp *StartResponse, _ error) {
	srv.stats.Clear()
	srv.stats.Start()
	resp = &StartResponse{}
	return
}

func (srv *server) StopBenchmark(_ context.Context, req *StopRequest) (resp *MemoryStat, _ error) {
	srv.stats.End()
	resp = &MemoryStat{
		Allocs: srv.stats.endMs.Mallocs - srv.stats.startMs.Mallocs,
		Memory: srv.stats.endMs.TotalAlloc - srv.stats.startMs.TotalAlloc,
	}
	return
}

func (srv *server) UnorderedQC(_ context.Context, in *Echo) (out *Echo, _ error) {
	out = in
	return
}

func (srv *server) OrderedQC(stream Benchmark_OrderedQCServer) (err error) {
	return OrderedQCServerLoop(stream, func(in *Echo) (out *Echo) {
		out = in
		return
	})
}

func (srv *server) UnorderedSlowServer(_ context.Context, in *Echo) (out *Echo, _ error) {
	time.Sleep(10 * time.Millisecond)
	out = in
	return
}

func (srv *server) OrderedSlowServer(stream Benchmark_OrderedSlowServerServer) (err error) {
	return OrderedSlowServerServerLoop(stream, func(in *Echo) (out *Echo) {
		time.Sleep(10 * time.Millisecond)
		out = in
		return
	})
}

// Server is a unified server for both ordered and unordered methods
type Server struct {
	*grpc.Server
	impl  server
	stats Stats
}

// NewServer returns a new Server
func NewServer() *Server {
	srv := &Server{Server: grpc.NewServer()}
	srv.impl.stats = &srv.stats

	RegisterBenchmarkServer(srv.Server, &srv.impl)
	return srv
}

// StartLocalServers starts benchmark servers locally
func StartLocalServers(ctx context.Context, n int) []string {
	var ports []string
	basePort := 40000
	var servers []*Server
	for p := basePort; p < basePort+n; p++ {
		port := fmt.Sprintf(":%d", p)
		ports = append(ports, port)
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("Failed to start local server: %v\n", err)
		}
		srv := NewServer()
		servers = append(servers, srv)
		go srv.Serve(lis)
	}
	go func() {
		<-ctx.Done()
		for _, srv := range servers {
			srv.Stop()
		}
	}()
	return ports
}
