package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Handlers interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func TimeoutHandlerServer(h Handlers, dt time.Duration, data interface{}) Handlers {
	return &timeoutHandlers{
		handler: h,
		dt:      dt,
		data:    data,
	}
}

func (h *timeoutHandlers) errorBodyData() []byte {
	result, _ := json.Marshal(h.data)
	return result
}

type timeoutHandlers struct {
	handler     Handlers
	dt          time.Duration
	data        interface{}
	testContext context.Context
}
type timeoutWriters struct {
	w    http.ResponseWriter
	h    http.Header
	wbuf bytes.Buffer
	req  *http.Request

	mu          sync.Mutex
	err         error
	wroteHeader bool
	code        int
}

func (tw *timeoutWriters) Header() http.Header {
	return tw.h
}

func (tw *timeoutWriters) Write(p []byte) (int, error) {

	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.err != nil {
		return 0, tw.err
	}
	if !tw.wroteHeader {
		tw.writeHeaderLocked(http.StatusOK)
	}
	return tw.wbuf.Write(p)
}

func (tw *timeoutWriters) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.writeHeaderLocked(code)
}

func checkWriteHeaderCode(code int) {

	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

func relevantCaller() runtime.Frame {
	pc := make([]uintptr, 16)
	n := runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc[:n])
	var frame runtime.Frame
	for {
		frame, more := frames.Next()
		if !strings.HasPrefix(frame.Function, "net/http.") {
			return frame
		}
		if !more {
			break
		}
	}
	return frame
}
func (tw *timeoutWriters) writeHeaderLocked(code int) {
	checkWriteHeaderCode(code)

	switch {
	case tw.err != nil:
		return
	case tw.wroteHeader:
		if tw.req != nil {
			caller := relevantCaller()
			logf(tw.req, "http: superfluous response.WriteHeader call from %s (%s:%d)", caller.Function, path.Base(caller.File), caller.Line)
		}
	default:
		tw.wroteHeader = true
		tw.code = code
	}
}

func logf(r *http.Request, format string, args ...any) {
	s, _ := r.Context().Value(http.ServerContextKey).(*http.Server)
	if s != nil && s.ErrorLog != nil {
		s.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (h *timeoutHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.testContext
	if ctx == nil {
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithTimeout(r.Context(), h.dt)
		defer cancelCtx()
	}
	r = r.WithContext(ctx)
	done := make(chan struct{})
	tw := &timeoutWriters{
		w:   w,
		h:   make(http.Header),
		req: r,
	}
	panicChan := make(chan any, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		h.handler.ServeHTTP(tw, r)
		close(done)
	}()
	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		tw.mu.Lock()
		defer tw.mu.Unlock()
		dst := w.Header()
		for k, vv := range tw.h {
			dst[k] = vv
		}
		if !tw.wroteHeader {
			tw.code = http.StatusOK
		}
		w.WriteHeader(tw.code)
		w.Write(tw.wbuf.Bytes())
	case <-ctx.Done():
		tw.mu.Lock()
		defer tw.mu.Unlock()
		switch err := ctx.Err(); err {
		case context.DeadlineExceeded:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write(h.errorBodyData())
			tw.err = http.ErrHandlerTimeout
		default:
			w.WriteHeader(http.StatusServiceUnavailable)
			tw.err = err
		}
	}
}
