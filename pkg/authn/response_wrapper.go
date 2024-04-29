package authn

import (
	"context"
	"net/http"
)

// See this: https://www.alexedwards.net/blog/how-to-use-the-http-responsecontroller-type

type responseWrapper struct {
	http.ResponseWriter
	requestCtx    context.Context
	statusCode    int
	headerWritten bool
}

func newResponseWrapper(ctx context.Context, w http.ResponseWriter) *responseWrapper {
	return &responseWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		requestCtx:     ctx,
	}
}

func (mw *responseWrapper) WriteHeader(statusCode int) {

	if !mw.headerWritten {
		mw.checkUserToken()
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
	mw.ResponseWriter.WriteHeader(statusCode)
}

func (mw *responseWrapper) Write(b []byte) (int, error) {
	if !mw.headerWritten {
		mw.checkUserToken()
		mw.headerWritten = true
	}
	return mw.ResponseWriter.Write(b)
}

func (mw *responseWrapper) Unwrap() http.ResponseWriter {
	return mw.ResponseWriter
}

func (mw *responseWrapper) checkUserToken() {
	authnData := getAuthData(mw.requestCtx)
	if authnData == nil {
		return
	}
	if authnData.changed {
		mw.ResponseWriter.Header().Set("X-User-Token", "true")
	}
}
