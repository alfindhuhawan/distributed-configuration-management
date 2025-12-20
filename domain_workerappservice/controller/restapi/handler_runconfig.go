package restapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"distributed-configuration-management/domain_workerappservice/usecase/runconfig"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/model/payload"
	requestModel "distributed-configuration-management/shared/model/request"
)

// runConfigHandler ...
func (r *Controller) runConfigHandler(inputPort runconfig.Inport) gin.HandlerFunc {

	type request struct {
		requestModel.ConfigWorker
	}

	type response struct {
	}

	return func(c *gin.Context) {

		traceID := util.GenerateID(16)

		ctx := logger.SetTraceID(context.Background(), traceID)

		var jsonReq request
		if err := c.BindJSON(&jsonReq); err != nil {
			r.Log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		req := runconfig.InportRequest(jsonReq)

		r.Log.Info(ctx, util.MustJSON(req))

		res, err := inputPort.Execute(ctx, req)
		if err != nil {
			r.Log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		var jsonRes response
		_ = res

		r.Log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusOK, payload.NewSuccessResponse(jsonRes, traceID))

	}
}
