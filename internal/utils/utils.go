package utils

import (
	"context"
	"errors"

	"github.com/nartim88/urlshortener/internal/models"
	"github.com/nartim88/urlshortener/pkg/logger"
)

func GetUserIDFromCtx(ctx context.Context) (string, error) {
	ctxKey := models.UserIDCtxKey("userID")
	userID := ctx.Value(ctxKey)
	if userID == nil {
		return "", errors.New("user id is not found in the request context")
	}
	val, ok := userID.(string)
	if !ok {
		return "", errors.New("user id is not valid, must be string")
	}
	logger.Log.Info().Str("userID", val).Msg("got user id from ctx")
	return val, nil
}
