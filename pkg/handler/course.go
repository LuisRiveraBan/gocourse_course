package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LuisRiveraBan/go_lib_response/response"
	course "github.com/LuisRiveraBan/gocourse_course/internal"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {
	r := mux.NewRouter()

	otps := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/courses").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.ListCourse),
		decodeGetAllCourse,
		encodeResponse,
		otps...,
	))

	r.Methods("GET").Path("/courses/{id}").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetCourseByID),
		decodeGetUser,
		encodeResponse,
		otps...,
	))

	r.Methods("POST").Path("/courses").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.CreateCourse),
		decodeCreateCourse,
		encodeResponse,
		otps...,
	))

	r.Methods("PATCH").Path("/courses/{id}").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.UpdateCourse),
		decodeUpdateCourse,
		encodeResponse,
		otps...,
	))

	r.Methods("DELETE").Path("/courses/{id}").Handler(httptransport.NewServer(
		endpoint.Endpoint(endpoints.DeleteCourse),
		decodeDeleteCourse,
		encodeResponse,
		otps...,
	))

	return r
}

func decodeCreateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req course.Create
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	id := params["id"]

	req := course.GetReq{
		ID: id,
	}

	return req, nil
}

func decodeGetAllCourse(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := course.GetAllReq{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}
	return req, nil
}

func decodeUpdateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req course.Update

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	params := mux.Vars(r)
	req.ID = params["id"]

	return req, nil
}

func decodeDeleteCourse(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	id := params["id"]

	req := course.DeleteReq{
		ID: id,
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(response.Response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
