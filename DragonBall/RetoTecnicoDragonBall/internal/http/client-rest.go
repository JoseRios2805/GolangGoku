package internal

import (
	"RetoTecnicoDragonBall/internal/logs"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClientRest struct {
	Log logs.ILogger
}

func NewClientRest(log logs.ILogger) *ClientRest {
	return &ClientRest{
		Log: log,
	}
}

func (c *ClientRest) InvokeAPI(url string, httpMethod string, data interface{}) ([]byte, int, error) {
	var request *http.Request
	var err error

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			c.Log.Error("Error serializando JSON", err)
			return nil, http.StatusInternalServerError, err
		}
		request, err = http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonData))
	} else {
		request, err = http.NewRequest(httpMethod, url, nil)
	}

	if err != nil {
		c.Log.Error("Error creando request", err)
		return nil, http.StatusInternalServerError, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: LoggingRoundTripper{http.DefaultTransport},
		Timeout:   10 * time.Second, // Corrige el timeout (10 segundos)
	}

	resp, err := client.Do(request)
	if err != nil {
		c.Log.Error("Error en llamada HTTP", err)
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Log.Error("Error leyendo respuesta", err)
		return nil, http.StatusInternalServerError, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	return body, resp.StatusCode, nil
}
