package aggendpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/sony/gobreaker"
	"github.com/telematics/go-kit-example/aggsvc/aggservice"
	"github.com/telematics/types"
	"golang.org/x/time/rate"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

func New(svc aggservice.Service, logger log.Logger) Set {
	duration := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "toll_calculator",
		Subsystem: "aggservice",
		Name:      "request_duration_seconds",
		Help:      "Request duration in seconds.",
	}, []string{"method", "success"})

	var aggregateEndpoint endpoint.Endpoint
	{
		aggregateEndpoint = MakeAggregateEndpoint(svc)
		aggregateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(aggregateEndpoint)
		aggregateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(aggregateEndpoint)
		aggregateEndpoint = LoggingMiddleware(log.With(logger, "method", "Aggregate"))(aggregateEndpoint)
		aggregateEndpoint = InstrumentingMiddleware(duration.With("method", "Aggregate"))(aggregateEndpoint)
	}
	var calculateEndpoint endpoint.Endpoint
	{
		calculateEndpoint = MakeConcatEndpoint(svc)
		calculateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(calculateEndpoint)
		calculateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(calculateEndpoint)
		calculateEndpoint = LoggingMiddleware(log.With(logger, "method", "Invoice"))(calculateEndpoint)
		calculateEndpoint = InstrumentingMiddleware(duration.With("method", "Invoice"))(calculateEndpoint)
	}
	return Set{
		AggregateEndpoint: aggregateEndpoint,
		CalculateEndpoint: calculateEndpoint,
	}
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type AggregateResponse struct {
	Err error `json:"err"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"err"`
}

func (s Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{OBUID: obuID})
	if err != nil {
		return nil, err
	}
	result := resp.(CalculateResponse)
	return &types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, nil
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		OBUID: dist.OBUID,
		Value: dist.Value,
		Unix:  dist.Unix,
	})
	return err
}

func MakeAggregateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			OBUID: req.OBUID,
			Value: req.Value,
			Unix:  req.Unix,
		})
		return AggregateResponse{Err: err}, nil
	}
}

func MakeConcatEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		inv, err := s.Calculate(ctx, req.OBUID)
		return CalculateResponse{
			Err:           err,
			OBUID:         inv.OBUID,
			TotalDistance: inv.TotalDistance,
			TotalAmount:   inv.TotalAmount,
		}, nil
	}
}
