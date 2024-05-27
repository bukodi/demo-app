package server

import (
	"context"
	"net/http"
)

// See this: https://www.alexedwards.net/blog/how-to-use-the-http-responsecontroller-type

type responseWrapper struct {
	http.ResponseWriter
	requestCtx           context.Context
	statusCode           int
	headerWritten        bool
	beforeHeaderComplete func(requestCtx context.Context, w http.ResponseWriter)
}

func NewResponseWrapper(ctx context.Context, w http.ResponseWriter, beforeHeaderComplete func(requestCtx context.Context, w http.ResponseWriter)) http.ResponseWriter {
	return &responseWrapper{
		ResponseWriter:       w,
		statusCode:           http.StatusOK,
		requestCtx:           ctx,
		beforeHeaderComplete: beforeHeaderComplete,
	}
}

func (mw *responseWrapper) WriteHeader(statusCode int) {

	if !mw.headerWritten {
		mw.beforeHeaderComplete(mw.requestCtx, mw.ResponseWriter)
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
	mw.ResponseWriter.WriteHeader(statusCode)
}

func (mw *responseWrapper) Write(b []byte) (int, error) {
	if !mw.headerWritten {
		mw.beforeHeaderComplete(mw.requestCtx, mw.ResponseWriter)
		mw.headerWritten = true
	}
	return mw.ResponseWriter.Write(b)
}

func (mw *responseWrapper) Unwrap() http.ResponseWriter {
	return mw.ResponseWriter
}
