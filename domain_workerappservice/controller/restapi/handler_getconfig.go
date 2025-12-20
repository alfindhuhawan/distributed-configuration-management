package restapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"distributed-configuration-management/domain_workerappservice/usecase/getconfig"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/model/payload"
	responseHit "distributed-configuration-management/shared/model/response"
)

// getConfigHandler ...
func (r *Controller) getConfigHandler(inputPort getconfig.Inport) gin.HandlerFunc {

	type request struct {
	}

	type response struct {
		*responseHit.HitResponse
	}

	return func(c *gin.Context) {

		traceID := util.GenerateID(16)

		ctx := logger.SetTraceID(context.Background(), traceID)

		var jsonReq request
		if err := c.Bind(&jsonReq); err != nil {
			r.Log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		var req getconfig.InportRequest

		r.Log.Info(ctx, util.MustJSON(req))

		res, err := inputPort.Execute(ctx, req)
		if err != nil {
			r.Log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		// var jsonRes response
		jsonRes := response(*res)

		r.Log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusOK, payload.NewSuccessResponse(jsonRes, traceID))

	}
}
