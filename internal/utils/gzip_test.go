package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"testing"
)

func Test_compressReader_Close(t *testing.T) {
	type fields struct {
		r  io.ReadCloser
		zr *gzip.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressReader{
				r:  tt.fields.r,
				zr: tt.fields.zr,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_compressReader_Read(t *testing.T) {
	type fields struct {
		r  io.ReadCloser
		zr *gzip.Reader
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := compressReader{
				r:  tt.fields.r,
				zr: tt.fields.zr,
			}
			gotN, err := c.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Read() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_compressWriter_Close(t *testing.T) {
	type fields struct {
		w  http.ResponseWriter
		zw *gzip.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressWriter{
				w:  tt.fields.w,
				zw: tt.fields.zw,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_compressWriter_Write(t *testing.T) {
	type fields struct {
		w  http.ResponseWriter
		zw *gzip.Writer
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressWriter{
				w:  tt.fields.w,
				zw: tt.fields.zw,
			}
			got, err := c.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressWriter_WriteHeader(t *testing.T) {
	type fields struct {
		w  http.ResponseWriter
		zw *gzip.Writer
	}
	type args struct {
		statusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressWriter{
				w:  tt.fields.w,
				zw: tt.fields.zw,
			}
			c.WriteHeader(tt.args.statusCode)
		})
	}
}
