package runner

import (
	"context"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type Input struct {
	Method    string
	URL       string
	Duration  time.Duration
	ReqPerSec int
}

type Output struct {
	totalLatency time.Duration
	RequestCount int64
}

func (o Output) AverageLatency() time.Duration {
	if o.RequestCount == 0 {
		return 0
	}
	return o.totalLatency / time.Duration(o.RequestCount)
}

func Run(ctx context.Context, in Input) (out Output) {
	atk := vegeta.NewAttacker(vegeta.Timeout(time.Second))

	res := atk.Attack(
		vegeta.NewStaticTargeter(vegeta.Target{
			Method: in.Method,
			URL:    in.URL,
		}),
		vegeta.Rate{Freq: in.ReqPerSec, Per: time.Second},
		in.Duration,
		"",
	)

	for {
		select {
		case <-ctx.Done():
			atk.Stop()
			return
		case r, ok := <-res:
			if !ok {
				return
			}
			// TODO: Account for r.Code / r.Error
			out.RequestCount++
			out.totalLatency += r.Latency
		}
	}
}
