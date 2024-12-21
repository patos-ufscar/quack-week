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
// @Param   payload 	body 		schemas.CreateEvent true "event json"
// @Success 200 		{object} 	schemas.Id
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/events/organization/{orgId} [PUT]
func (c *EventController) CreateEvent(ctx *gin.Context) {
	var createEvent schemas.CreateEvent

	if err := ctx.ShouldBind(&createEvent); err != nil {
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

	event, err := c.eventService.CreateEvent(ctx, createEvent.EventName, claims.UserId, *claims.OrganizationId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	if createEvent.CoverBase64 != nil {
		picBytes, err := base64.StdEncoding.DecodeString(*createEvent.CoverBase64)
		if err != nil {
			slog.Error(err.Error())
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		if len(picBytes) > 500*1024 { // if > 500k
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
		err = c.objService.Upload(ctx, common.S3_BUCKET, objPath, bytes.NewReader(picBytes))
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
	}

	ctx.JSON(http.StatusOK, schemas.Id{Id: event.EventId})
}

func (c *EventController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/events")

	g.PUT("/organization/:orgId", authMiddleware.AuthorizeOrganization(true), c.CreateEvent)
}
