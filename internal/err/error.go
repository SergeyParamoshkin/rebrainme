package fwerr

import (
	"context"
	"errors"
	"fmt"

	fwctx "github.com/SergeyParamoshkin/alerts/internal/ctx"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e Error) String() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Err)
}

func (e Error) LogCtx(ctx context.Context, logMessage string, fields ...zapcore.Field) {
	switch e.Code {
	case
		CodeAuthError,
		CodeBadRequest,
		CodeNotFound:
		// do not log
		return
	case CodeDBError, CodeInternalError:
		fallthrough
	default:
		if e.Err != nil {
			fields = append(fields, zap.Error(e.Err))
		}

		if logger := fwctx.LoggerFromCtx(ctx, nil); logger != nil {
			logger.Error(logMessage, fields...)
		}
	}
}

func New(code Code, message string) Error {
	return Error{
		Code:    code,
		Message: message,
		Err:     nil,
	}
}

func Wrap(code Code, message string, err error) Error {
	return Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func FromError(err error) Error {
	var fwErr Error

	if errors.As(err, &fwErr) {
		return fwErr
	}

	return Wrap(CodeInternalError, "", err)
}

func APIErrorWrap(message string, err error) Error {
	return Wrap(CodeBadRequest, message, err)
}
