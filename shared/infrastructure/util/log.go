package util

import (
	"context"
	"distributed-configuration-management/shared/infrastructure/logger"
)

func InsertLog(ctx context.Context, log logger.Logger, err error) {
	log.Error(ctx, err.Error())
}
