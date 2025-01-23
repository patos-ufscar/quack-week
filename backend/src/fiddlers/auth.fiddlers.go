package fiddlers

import (
	"github.com/gin-gonic/gin"
	"github.com/patos-ufscar/quack-week/common"
	"github.com/patos-ufscar/quack-week/models"
)

func GetClaimsFromGinCtx(ctx *gin.Context) (models.JwtClaims, error) {
	claims, ok := ctx.Get(common.GIN_CTX_JWT_CLAIM_KEY_NAME)
	if !ok {
		return models.JwtClaims{}, common.ErrAuth
	}

	parsedClaims, ok := claims.(models.JwtClaims)
	if !ok {
		return models.JwtClaims{}, common.ErrAuth
	}

	return parsedClaims, nil
}
