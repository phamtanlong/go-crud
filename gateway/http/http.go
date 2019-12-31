package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/phamtanlong/go-crud/gateway/endpoints"
)

func NewHTTPHandler(endpoints endpoints.Endpoints) http.Handler {
	m := mux.NewRouter()

	m.Methods("POST").Path("/service/alive").Handler(httptransport.NewServer(endpoints.ServiceAlive, DecodeServiceAliveRequest, EncodeServiceAliveResponse))
	m.Methods("POST").Path("/service/die").Handler(httptransport.NewServer(endpoints.ServiceDie, DecodeServiceDieRequest, EncodeServiceDieResponse))

	// news
	m.Methods("POST").Path("/news/create").Handler(httptransport.NewServer(endpoints.CreateNews,
		DecodeCreateNewsRequest, EncodeCreateNewsResponse, httptransport.ServerBefore(httptransport.PopulateRequestContext)))
	m.Methods("POST").Path("/news/update").Handler(httptransport.NewServer(endpoints.UpdateNews,
		DecodeUpdateNewsRequest, EncodeUpdateNewsResponse, httptransport.ServerBefore(httptransport.PopulateRequestContext)))
	m.Methods("DELETE").Path("/news/delete").Handler(httptransport.NewServer(endpoints.DeleteNews,
		DecodeDeleteNewsRequest, EncodeDeleteNewsResponse, httptransport.ServerBefore(httptransport.PopulateRequestContext)))
	m.Path("/news/read").Handler(httptransport.NewServer(endpoints.ReadNews,
		DecodeReadNewsRequest, EncodeReadNewsResponse))

	// auth

	m.Methods("POST").Path("/register").Handler(httptransport.NewServer(endpoints.Register,
		DecodeRegisterRequest, EncodeRegisterResponse))
	m.Methods("POST").Path("/login").Handler(httptransport.NewServer(endpoints.Login,
		DecodeLoginRequest, EncodeLoginResponse))

	return m
}
func DecodeServiceAliveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ServiceAliveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeServiceAliveResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeServiceDieRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ServiceDieRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeServiceDieResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeCreateNewsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateNewsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeCreateNewsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeUpdateNewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateNewsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeUpdateNewsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeDeleteNewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DeleteNewsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeDeleteNewsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeReadNewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ReadNewsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeReadNewsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeRegisterResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
