package statusCode

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	NotAuth          codes.Code = 101
	GroupNotFound    codes.Code = 102
	ChapterNotFound  codes.Code = 103
	ResourceNotFound codes.Code = 104
)

var Code2StrMap map[codes.Code]string = map[codes.Code]string{
	NotAuth:          "没有授权访问",
	GroupNotFound:    "未找到资源组",
	ChapterNotFound:  "未找章节",
	ResourceNotFound: "未找到资源",
}

var GrpcCodeStr2code = map[string]codes.Code{
	"Canceled":           codes.Canceled,
	"Unknown":            codes.Unknown,
	"InvalidArgument":    codes.InvalidArgument,
	"DeadlineExceeded":   codes.DeadlineExceeded,
	"NotFound":           codes.NotFound,
	"AlreadyExists":      codes.AlreadyExists,
	"PermissionDenied":   codes.PermissionDenied,
	"ResourceExhausted":  codes.ResourceExhausted,
	"FailedPrecondition": codes.FailedPrecondition,
	"Aborted":            codes.Aborted,
	"OutOfRange":         codes.OutOfRange,
	"Unimplemented":      codes.Unimplemented,
	"Internal":           codes.Internal,
	"Unavailable":        codes.Unavailable,
	"DataLoss":           codes.DataLoss,
	"Unauthenticated":    codes.Unauthenticated,
}

type Status struct {
	code    codes.Code
	message string
}

func NewStatus(code codes.Code, message string) *Status {
	if message == "" {
		if msg, ok := Code2StrMap[code]; ok {
			message = msg
		}
	}
	return &Status{
		code:    code,
		message: message,
	}
}

func (s *Status) Error() string {
	return s.message
}

func (s *Status) GrpcError() error {
	return status.Error(s.code, s.message)
}

func (s *Status) WithMessage(msg string) {
	s.message = msg
	return
}

func (s *Status) GetMessage() string {
	return s.message
}

//
func (s *Status) String() string {

	e := status.Error(s.code, s.message)
	if e == nil {
		return "success"
	}
	return e.Error()
}


func PackGrpcError(code codes.Code) error {
	return NewStatus(code, "").GrpcError()
}
