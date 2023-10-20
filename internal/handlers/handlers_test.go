package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMiddleware(t *testing.T) {
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Middleware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Middleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
