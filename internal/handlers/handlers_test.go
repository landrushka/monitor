package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/metrics"
	"github.com/landrushka/monitor.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetValueHandle(t *testing.T) {
	type fields struct {
		memStorage storage.MemStorage
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "positive test #1",
			fields: fields{
				memStorage: storage.MemStorage{GaugeMetric: map[string]float64{"gauge_test_name": 0.001}},
			},
			want: want{
				code:     200,
				response: `0.001`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &Handler{
				memStorage: test.fields.memStorage,
			}

			request := httptest.NewRequest(http.MethodGet, "/value/{type}/{name}", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", "gauge")
			rctx.URLParams.Add("name", "gauge_test_name")
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.GetValueHandle(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}

func TestHandler_GetAllNamesHandle(t *testing.T) {
	type fields struct {
		memStorage storage.MemStorage
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "positive test #1",
			fields: fields{
				memStorage: storage.MemStorage{GaugeMetric: map[string]float64{"gauge_test_name": 0.001}},
			},
			want: want{
				code:     200,
				response: "\n<h1>Metric Names</h1>\n<dl>\n[gauge_test_name]\n</dl>\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &Handler{
				memStorage: test.fields.memStorage,
			}

			request := httptest.NewRequest(http.MethodGet, "/", nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.GetAllNamesHandle(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}

func TestHandler_UpdateHandle(t *testing.T) {
	type fields struct {
		memStorage storage.MemStorage
	}
	type want struct {
		code        int
		value       float64
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "positive test #1",
			fields: fields{
				memStorage: storage.MemStorage{GaugeMetric: map[string]float64{"gauge_test_name": 0.001}},
			},
			want: want{
				code:  200,
				value: 100,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &Handler{
				memStorage: test.fields.memStorage,
			}
			var buf bytes.Buffer
			var val = 100.00

			m := metrics.Metrics{ID: "gauge_test_name", MType: "gauge", Value: &val}
			json.NewEncoder(&buf).Encode(m)

			request := httptest.NewRequest(http.MethodPost, "/update", &buf)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.UpdateHandle(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.value, h.memStorage.GaugeMetric["gauge_test_name"])
		})
	}
}
