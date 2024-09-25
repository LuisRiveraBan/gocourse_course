package course

import (
	"context"
	"errors"
	"github.com/LuisRiveraBan/go_lib_response/response"
	"github.com/LuisRiveraBan/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		GetCourseByID Controller
		CreateCourse  Controller
		UpdateCourse  Controller
		DeleteCourse  Controller
		ListCourse    Controller
	}

	Create struct {
		Name     string `json:"name"`
		StarDate string `json:"star_date"`
		EndDate  string `json:"end_date"`
	}

	GetReq struct {
		ID string
	}

	DeleteReq struct {
		ID string
	}

	GetAllReq struct {
		Page  int
		Limit int
		Name  string
	}

	Update struct {
		ID       string
		Name     *string `json:"name"`
		StarDate *string `json:"star_date"`
		EndDate  *string `json:"end_date"`
	}

	// ErrResponse is a simple struct to represent an error response
	ErrResponse struct {
		Message string `json:"message"`
	}

	Response struct {
		Status int         `json:"status"`
		Err    string      `json:"error,omitempty"`
		Data   interface{} `json:"data,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		CreateCourse:  makeCreateEndpoint(s),
		GetCourseByID: makeToFindEndpoints(s),
		UpdateCourse:  makeUpdateEndpoints(s),
		DeleteCourse:  makeDeleteEndpoint(s),
		ListCourse:    makeListEndpoint(s, config),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(Create)

		if req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StarDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		course, err := s.Create(ctx, req.Name, req.StarDate, req.EndDate)

		if err != nil {

			if err == ErrEndLesserStart || err == ErrInvalidStartDate || err == ErrInvalidEndDate {
				return nil, response.BadRequest(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", course, nil), nil
	}
}

func makeListEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Extract query parameters from the request URL
		req := request.(GetAllReq)

		filters := Filters{
			Name: req.Name,
		}
		count, err := s.Count(ctx, filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		meta, err := meta.NewMeta(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		course, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", course, meta), nil
	}
}

func makeToFindEndpoints(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		course, err := s.GetByID(ctx, req.ID)

		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.NotFound(err.Error())
		}
		return response.OK("sucess", course, nil), nil

	}
}

func makeUpdateEndpoints(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(Update)

		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		if req.StarDate != nil && *req.StarDate == "" {
			return nil, response.BadRequest(ErrStartDateRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.Name, req.StarDate, req.EndDate)

		if err != nil {

			if err == ErrEndLesserStart || err == ErrInvalidStartDate || err == ErrInvalidEndDate {
				return nil, response.BadRequest(err.Error())
			}

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", nil, nil), nil

	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteReq)

		if err := s.Delete(ctx, req.ID); err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())

		}
		return response.OK("success", nil, nil), nil
	}
}
