package helpers

import (
	"context"
	"errors"

	"github.com/nartim88/urlshortener/pkg/logger"
)

func GetUserIDFromCtx(ctx context.Context) (string, error) {
	ctxKey := "userID"
	userID := ctx.Value(ctxKey)
	if userID == nil {
		return "", errors.New("user id is not found in the request context")
	}
	val, ok := userID.(string)
	if !ok {
		return "", errors.New("user id is not valid, must be string")
	}
	logger.Log.Info().Str(ctxKey, val).Send()
	return val, nil
}
