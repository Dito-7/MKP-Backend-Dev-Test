package handler

import (
	"errors"
	"net/http"
	"strconv"

	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/usecase"
	"mkp-cinema-ticketing/pkg/response"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleUsecase usecase.ScheduleUsecase
}

func NewScheduleHandler(scheduleUsecase usecase.ScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{scheduleUsecase: scheduleUsecase}
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	var req entity.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Data jadwal tayang tidak valid", err.Error())
		return
	}

	result, err := h.scheduleUsecase.Create(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrScheduleConflict):
			response.Error(c, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, usecase.ErrMovieNotFound), errors.Is(err, usecase.ErrStudioNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, usecase.ErrInvalidTimeRange), errors.Is(err, usecase.ErrInvalidTimeFormat):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, "Jadwal tayang berhasil ditambahkan", result)
}

func (h *ScheduleHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := entity.ScheduleFilter{
		MovieID:  c.Query("movie_id"),
		StudioID: c.Query("studio_id"),
		CinemaID: c.Query("cinema_id"),
		CityID:   c.Query("city_id"),
		Date:     c.Query("date"),
		Status:   c.Query("status"),
		Page:     page,
		Limit:    limit,
	}

	items, total, err := h.scheduleUsecase.GetAll(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, "Gagal mengambil daftar jadwal tayang: "+err.Error())
		return
	}

	meta := gin.H{
		"page":        page,
		"limit":       limit,
		"total_items": total,
		"total_pages": (total + limit - 1) / limit,
	}

	response.SuccessWithMeta(c, http.StatusOK, "Daftar jadwal tayang berhasil diambil", items, meta)
}

func (h *ScheduleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	item, err := h.scheduleUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrScheduleNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Detail jadwal tayang berhasil didapatkan", item)
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req entity.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Data update tidak valid", err.Error())
		return
	}

	item, err := h.scheduleUsecase.Update(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrScheduleConflict):
			response.Error(c, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, usecase.ErrScheduleNotFound), errors.Is(err, usecase.ErrMovieNotFound), errors.Is(err, usecase.ErrStudioNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, usecase.ErrInvalidTimeRange), errors.Is(err, usecase.ErrInvalidTimeFormat):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Jadwal tayang berhasil diperbarui", item)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.scheduleUsecase.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrScheduleNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "Gagal menghapus jadwal tayang: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Jadwal tayang berhasil dihapus", nil)
}
