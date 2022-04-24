package service

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock "github.com/Asymmetriq/shortener/internal/app/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService_getHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
		wantErr     bool
	}
	tests := []struct {
		name    string
		request *http.Request
		want    want
	}{
		{
			name:    "test positive 1",
			request: httptest.NewRequest(http.MethodGet, "/short-url-mock", nil),
			want: want{
				contentType: "text/html; charset=utf-8",
				statusCode:  http.StatusTemporaryRedirect,
				body:        "<a href=\"https://www.google.com\">Temporary Redirect</a>.\n\n",
			},
		},
		{
			name:    "test negative 1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("wow-url")),
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        fmt.Sprintf("no original url found with shortcut %q\n", "wow-url"),
				wantErr:     true,
			},
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockRepository(ctrl)
		if !tt.want.wantErr {
			repo.EXPECT().Get("short-url-mock").Return("https://www.google.com", nil)
		} else {
			repo.EXPECT().Get(gomock.Any()).Return("", fmt.Errorf("no original url found with shortcut %q", "wow-url"))
		}

		recoder := httptest.NewRecorder()

		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: repo,
			}
			s.getHandler(recoder, tt.request)
			resp := recoder.Result()

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "error reading resp body")

			require.Equal(t, tt.want.statusCode, resp.StatusCode, "status codes don't match")
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "content types don't match")
			require.Equal(t, tt.want.body, string(body), "shortened urls don't match")
		})
	}
}

func TestService_postHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
		wantErr     bool
	}
	tests := []struct {
		name    string
		request *http.Request
		want    want
	}{
		{
			name:    "test positive 1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com")),
			want: want{
				contentType: "application/text",
				statusCode:  http.StatusCreated,
				body:        "http://example.com/short-url-mock",
			},
		},
		{
			name:    "test positive 2",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("GLKGDSL;FG;DLSKFGLSDF;GK")),
			want: want{
				contentType: "application/text",
				statusCode:  http.StatusCreated,
				body:        "http://example.com/short-url-mock",
			},
		},
		{
			name:    "test negative 1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("")),
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "no request body\n",
				wantErr:     true,
			},
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockRepository(ctrl)
		if !tt.want.wantErr {
			repo.EXPECT().Set(gomock.Any()).Return("short-url-mock")
		}
		recoder := httptest.NewRecorder()

		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: repo,
			}
			s.postHandler(recoder, tt.request)
			resp := recoder.Result()

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "error reading resp body")

			require.Equal(t, tt.want.statusCode, resp.StatusCode, "status codes don't match")
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "content types don't match")
			require.Equal(t, tt.want.body, string(body), "shortened urls don't match")
		})
	}
}
