package controllers

import (
	"bytes"
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patos-ufscar/quack-week/common"
	"github.com/patos-ufscar/quack-week/fiddlers"
	"github.com/patos-ufscar/quack-week/fiddlers/storage"
	"github.com/patos-ufscar/quack-week/middlewares"
	"github.com/patos-ufscar/quack-week/schemas"
	"github.com/patos-ufscar/quack-week/services"
)

type EventController struct {
	userService  services.UserService
	emailService services.EmailService
	orgService   services.OrganizationService
	eventService services.EventService
	objService   services.ObjectService
}

func NewEventController(
	userService services.UserService,
	emailService services.EmailService,
	orgService services.OrganizationService,
	eventService services.EventService,
	objService services.ObjectService,
) EventController {
	return EventController{
		userService:  userService,
		emailService: emailService,
		orgService:   orgService,
		eventService: eventService,
		objService:   objService,
	}
}

// @Summary CreateEvent
// @Security JWT
// @Tags Event
// @Description Creates an Event
// @Consume application/json
// @Accept json
// @Produce plain
// @Param	orgId 		path string true "Organization Id"
// @Param   payload 	body 		schemas.Name true "name json"
// @Success 200 		{object} 	schemas.Id
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/events/organization/{orgId} [PUT]
func (c *EventController) CreateEvent(ctx *gin.Context) {
	var name schemas.Name

	if err := ctx.ShouldBind(&name); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	claims, err := fiddlers.GetClaimsFromGinCtx(ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	event, err := c.eventService.CreateEvent(ctx, name.Name, claims.UserId, *claims.OrganizationId, "")
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, schemas.Id{Id: event.EventId})
}

// @Summary GetEvent
// @Tags Event
// @Description Gets an Event
// @Consume application/json
// @Accept json
// @Produce plain
// @Param	eventId 	path string true "Event Id"
// @Success 200 		{object} 	models.Event
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/events/{eventId} [GET]
func (c *EventController) GetEvent(ctx *gin.Context) {
	eventId := ctx.Param("eventId")
	event, err := c.eventService.GetEvent(ctx, eventId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// @Summary SetBanner
// @Security JWT
// @Tags Event
// @Description Gets an Event
// @Consume application/json
// @Accept json
// @Produce plain
// @Param	eventId 	path string true "Event Id"
// @Param	orgId 		path string true "Organization Id"
// @Param   payload 	body 		schemas.UploadPicture true "picture json"
// @Success 200 		{object} 	schemas.Url
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/events/{eventId}/organization/{orgId}/banner [PUT]
func (c *EventController) SetBanner(ctx *gin.Context) {
	eventId := ctx.Param("eventId")
	orgId := ctx.Param("orgId")
	var uploadPicture schemas.UploadPicture

	event, err := c.eventService.GetEvent(ctx, eventId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusConflict, "Conflict")
		return
	}

	if event.OwnerOrganizationId != orgId {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := ctx.ShouldBind(&uploadPicture); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	picBytes, err := base64.StdEncoding.DecodeString(uploadPicture.Content)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if len(picBytes) > 5*1024*1024 { // if > 5MB
		ctx.String(http.StatusBadGateway, "ImgTooLarge")
		return
	}

	imgFmt, err := common.GetImageFormat(picBytes)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	if imgFmt != "png" && imgFmt != "jpeg" {
		ctx.String(http.StatusUnsupportedMediaType, "UnsupportedMediaType")
		return
	}

	objPath := storage.GetPublicPath(storage.EVENT_BANNERS, event.EventId)
	err = c.objService.Upload(ctx, common.S3_BUCKET, objPath, int64(len(picBytes)), bytes.NewReader(picBytes))
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	objUrl, err := storage.GetFullObjUrl(objPath)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.eventService.SetCover(ctx, event.EventId, objUrl)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, schemas.Url{Url: objUrl})
}

func (c *EventController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/events")

	g.GET("/:eventId", c.GetEvent)
	g.PUT("/organization/:orgId", authMiddleware.AuthorizeOrganization(true), c.CreateEvent)
	g.PUT("/:eventId/organization/:orgId/banner", authMiddleware.AuthorizeOrganization(true), c.SetBanner)
}
