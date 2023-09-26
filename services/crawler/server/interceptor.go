package server

import (
	"context"

	"connectrpc.com/connect"
)

func NewLoggerInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			procedure := req.Spec().Procedure
			logger.Info(procedure, "status", "run", "args", req.Any(), "query", req.Peer().Query)
			res, err := next(ctx, req)
			if err != nil {
				logger.Error(procedure, err, "status", "error")
				return res, err
			}
			logger.Info(procedure, "status", "done")
			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
