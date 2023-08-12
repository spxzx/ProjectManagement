package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}

func ParseGrpcError(err error) BError {
	e, _ := status.FromError(err)
	return BError{Code: ErrorCode(e.Code()), Msg: e.Message()}
}
