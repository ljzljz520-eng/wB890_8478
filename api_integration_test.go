package memorialstation

import (
	"memorialstation/api"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPRecordEntryPoint(t *testing.T) {
	services := openTestServices(t)
	server := api.New(services.store)
	request := httptest.NewRequest("POST", "/records", strings.NewReader(`{"id":"http-1","batch_id":"http-batch","student_name":"吴凡","graduation_year":2024,"message":"网络入口登记","visibility":"class","tags":["memory"]}`))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != 201 {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, httptest.NewRequest("GET", "/records", nil))
	if response.Code != 200 || !strings.Contains(response.Body.String(), "http-1") {
		t.Fatalf("list response=%d %s", response.Code, response.Body.String())
	}
}
