package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/domain"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) PostEstate(ctx echo.Context) error {
	body := &generated.EstateRequest{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&body)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	if body.Width < 1 || body.Width > 50000 || body.Length < 1 || body.Length > 50000 {
		var resp generated.ErrorResponse
		resp.Message = "invalid value or format"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	newEstate := &domain.Estate{
		Width:  body.Width,
		Length: body.Length,
	}

	createdEstate, err := s.Usecase.CreateEstate(ctx.Request().Context(), newEstate)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, resp)
	}

	var resp generated.EstateSuccessResponse
	resp.Id = createdEstate.ID
	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) PostEstateEstateIdTree(ctx echo.Context, estateId string) error {
	body := &generated.EstateTreeRequest{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&body)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	if body.Height < 1 || body.Height > 30 {
		var resp generated.ErrorResponse
		resp.Message = "invalid value or format"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	newTree := &domain.Tree{
		X:        body.X,
		Y:        body.Y,
		Height:   body.Height,
		EstateID: estateId,
	}

	createdTree, err := s.Usecase.CreateTree(ctx.Request().Context(), newTree)
	if err != nil {
		if err.Error() == "tree already exist on the coordinate" || err.Error() == "tree coordinate out of bound" {
			var resp generated.ErrorResponse
			resp.Message = err.Error()
			return ctx.JSON(http.StatusBadRequest, resp)
		} else if err.Error() == "estate id not found" {
			var resp generated.ErrorResponse
			resp.Message = err.Error()
			return ctx.JSON(http.StatusNotFound, resp)
		}

		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, resp)
	}

	var resp generated.EstateTreeSuccessResponse
	resp.Id = createdTree.ID
	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) GetEstateEstateIdStats(ctx echo.Context, estateId string) error {
	resp, err := s.Usecase.GetEstateStats(ctx.Request().Context(), estateId)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetEstateEstateIdDronePlan(ctx echo.Context, estateId string) error {
	distance, err := s.Usecase.GetDronePlan(ctx.Request().Context(), estateId)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	resp := &generated.DronePlanSuccessResponse{
		Distance: distance,
	}
	return ctx.JSON(http.StatusOK, resp)
}
