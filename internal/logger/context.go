package logger

import (
	"context"
	"time"
)

type RequestStartTime string

const RequestStartTimeKey = RequestStartTime("requestStartTime")

func GetRequestStartTimeFromCtx(ctx context.Context) (time.Time, bool) {
	v := ctx.Value(RequestStartTimeKey)
	if v == nil {
		return time.Time{}, false
	}
	data, ok := v.(time.Time)
	if !ok {
		return time.Time{}, false
	}

	return data, true
}

func PutRequestStartTimeInCtx(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, RequestStartTimeKey, t)
}
