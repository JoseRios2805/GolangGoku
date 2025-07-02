package internal

import (
	"RetoTecnicoDragonBall/internal/utils/helpers"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

/**
  *Start: Interceptor peticiones http-client
  *- Imprime los logs de request y response de las invocaciones a apis
**/
type LoggingRoundTripper struct {
	proxy http.RoundTripper
}

func (l LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	//agregamos en la cabecera un identificador de logs

	var bodyInBytes []byte
	if req.Body != nil {
		defer req.Body.Close()
		bodyInBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyInBytes))
	}

	//
	respChannel := make(chan loggerResponse)

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go l.saveLogsRequest(bodyInBytes, req, wg)
	go l.clientProxy(req, wg, respChannel)
	go l.saveLogsResponse(req, wg, respChannel)

	outs := <-respChannel
	wg.Wait()

	return outs.res, outs.err
}

func (l LoggingRoundTripper) clientProxy(req *http.Request, wg *sync.WaitGroup, respChannel chan<- loggerResponse) {
	defer wg.Done()
	timeStart := time.Now()
	res, e := l.proxy.RoundTrip(req)
	timeEnd := time.Now()

	out := new(loggerResponse)
	out.err = e
	out.res = res
	out.latency = timeEnd.Sub(timeStart)

	var bodyOutBytes []byte
	if res != nil && res.Body != nil {
		defer res.Body.Close()
		bodyOutBytes, _ = io.ReadAll(res.Body)
		res.Body = io.NopCloser(bytes.NewBuffer(bodyOutBytes))
	}

	out.body = bodyOutBytes

	respChannel <- *out
	respChannel <- *out
	close(respChannel)
}

func (l LoggingRoundTripper) saveLogsRequest(bodyInBytes []byte, req *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	path := req.URL.Path
	raw := req.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}

	var bodyIn interface{}
	if bodyInBytes != nil {
		var sizeBody = req.ContentLength
		//logs cloud logging solo permite escribir maximo 256 KB
		if sizeBody > 255900 {
			bodyIn = "[body (" + strconv.FormatInt(sizeBody, 10) + ") greater than 256 KB]"
		} else {
			json.Unmarshal(bodyInBytes, &bodyIn)
			jsonObjdata := helpers.SerializeLogsStruct(bodyIn)
			bodyIn = string(jsonObjdata)
		}
	} else {
		bodyIn = ""
	}

}

func (l LoggingRoundTripper) saveLogsResponse(req *http.Request, wg *sync.WaitGroup, respChannel <-chan loggerResponse) {
	defer wg.Done()

	path := req.URL.Path
	raw := req.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}
	resProxy := <-respChannel
	var bodyOut interface{}
	if resProxy.body != nil {
		json.Unmarshal(resProxy.body, &bodyOut)
		jsonObjdata := helpers.SerializeLogsStruct(bodyOut)
		bodyOut = string(jsonObjdata)
	} else {
		bodyOut = ""
	}

	if resProxy.res.StatusCode >= http.StatusOK && resProxy.res.StatusCode < http.StatusBadRequest {
		//l.logger.Info(loggerOut)
	} else {
		//l.logger.Error(loggerOut, resProxy.err)
	}

}

type loggerResponse struct {
	res     *http.Response
	err     error
	body    []byte
	latency time.Duration
}

//End: Interceptor peticiones http-client
