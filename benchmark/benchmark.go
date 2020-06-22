package benchmark

import (
	context "context"
	"regexp"
	"runtime"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Options controls different options for the benchmarks
type Options struct {
	Concurrent int           // Number of concurrent calls
	Duration   time.Duration // Duration of benchmark
	MaxAsync   int           // Max async calls at once
	NumNodes   int           // Number of nodes to include in configuration
	Payload    int           // Size of message payload
	QuorumSize int           // Number of messages to wait for
	Warmup     time.Duration // Warmup time
	Remote     bool          // Whether the servers are remote (true) or local (false)
}

type Bench struct {
	Name        string
	Description string
	runBench    benchFunc
}

type benchFunc func(Options) (*Result, error)
type qcFunc func(context.Context, *Echo, ...grpc.CallOption) (*Echo, error)
type serverFunc func(*TimedMsg) error

func runQCBenchmark(opts Options, cfg *Configuration, f qcFunc) (*Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	msg := &Echo{Payload: make([]byte, opts.Payload)}
	s := &Stats{}
	var g errgroup.Group

	for n := 0; n < opts.Concurrent; n++ {
		g.Go(func() error {
			warmupEnd := time.Now().Add(opts.Warmup)
			for !time.Now().After(warmupEnd) {
				_, err := f(ctx, msg)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return nil, err
	}

	if opts.Remote {
		_, err := cfg.StartBenchmark(ctx, &StartRequest{})
		if err != nil {
			return nil, err
		}
	}

	s.Start()
	for n := 0; n < opts.Concurrent; n++ {
		g.Go(func() error {
			endTime := time.Now().Add(opts.Duration)
			for !time.Now().After(endTime) {
				start := time.Now()
				_, err := f(ctx, msg)
				if err != nil {
					return err
				}
				s.AddLatency(time.Since(start))
			}
			return nil
		})
	}

	err = g.Wait()
	s.End()
	if err != nil {
		return nil, err
	}

	result := s.GetResult()
	if opts.Remote {
		memStats, err := cfg.StopBenchmark(ctx, &StopRequest{})
		if err != nil {
			return nil, err
		}
		result.ServerStats = memStats.MemoryStats
	}

	return result, nil
}

func runServerBenchmark(opts Options, cfg *Configuration, f serverFunc) (*Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	payload := make([]byte, opts.Payload)
	var g errgroup.Group
	var start runtime.MemStats
	var end runtime.MemStats

	benchmarkFunc := func(stopTime time.Time) error {
		for !time.Now().After(stopTime) {
			msg := &TimedMsg{SendTime: time.Now().UnixNano(), Payload: payload}
			err := f(msg)
			if err != nil {
				return err
			}
		}
		return nil
	}

	warmupEnd := time.Now().Add(opts.Warmup)
	for n := 0; n < opts.Concurrent; n++ {
		g.Go(func() error { return benchmarkFunc(warmupEnd) })
	}
	g.Wait()

	_, err := cfg.StartServerBenchmark(ctx, &StartRequest{})
	if err != nil {
		return nil, err
	}

	runtime.ReadMemStats(&start)
	endTime := time.Now().Add(opts.Duration)
	for n := 0; n < opts.Concurrent; n++ {
		g.Go(func() error { return benchmarkFunc(endTime) })
	}
	g.Wait()
	runtime.ReadMemStats(&end)

	resp, err := cfg.StopServerBenchmark(ctx, &StopRequest{})
	if err != nil {
		return nil, err
	}

	clientAllocs := (end.Mallocs - start.Mallocs) / resp.TotalOps
	clientMem := (end.TotalAlloc - start.TotalAlloc) / resp.TotalOps

	resp.AllocsPerOp = clientAllocs
	resp.MemPerOp = clientMem
	return resp, nil
}

func GetBenchmarks(cfg *Configuration) []Bench {
	m := []Bench{
		{
			Name:        "UnorderedQC",
			Description: "Unary RPC based quorum call implementation without FIFO ordering",
			runBench: func(opts Options) (*Result, error) {
				return runQCBenchmark(opts, cfg, cfg.UnorderedQC)
			},
		},
		{
			Name:        "OrderedQC",
			Description: "NodeStream based quorum call implementation with FIFO ordering",
			runBench:    func(opts Options) (*Result, error) { return runQCBenchmark(opts, cfg, cfg.OrderedQC) },
		},
		{
			Name:        "UnorderedSlowServer",
			Description: "UnorderedQC with a 10ms processing time on the server",
			runBench: func(opts Options) (*Result, error) {
				return runQCBenchmark(opts, cfg, cfg.UnorderedSlowServer)
			},
		},
		{
			Name:        "OrderedSlowServer",
			Description: "OrderedQC with a 10s processing time on the server",
			runBench:    func(opts Options) (*Result, error) { return runQCBenchmark(opts, cfg, cfg.OrderedSlowServer) },
		},
	}
	return m
}

// RunBenchmarks runs all the benchmarks that match the given regex with the given options
func RunBenchmarks(benchRegex *regexp.Regexp, options Options, manager *Manager) ([]*Result, error) {
	nodeIDs := manager.NodeIDs()
	cfg, err := manager.NewConfiguration(nodeIDs[:options.NumNodes], &QSpec{QSize: options.QuorumSize, CfgSize: options.NumNodes})
	if err != nil {
		return nil, err
	}
	benchmarks := GetBenchmarks(cfg)
	var results []*Result
	for _, b := range benchmarks {
		if benchRegex.MatchString(b.Name) {
			result, err := b.runBench(options)
			if err != nil {
				return nil, err
			}
			result.Name = b.Name
			i := sort.Search(len(results), func(i int) bool {
				return results[i].Name >= result.Name
			})
			results = append(results, nil)
			copy(results[i+1:], results[i:])
			results[i] = result
		}
	}
	return results, nil
}
