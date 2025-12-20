package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"distributed-configuration-management/domain_controllerappservice/usecase/updateconfig"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/model/payload"
	"distributed-configuration-management/shared/model/repository"
)

// updateConfigHandler ...
func (r *Controller) updateConfigHandler(inputPort updateconfig.Inport) gin.HandlerFunc {

	type request struct {
		repository.GlobalConfigRequest
	}

	type response struct {
		repository.GlobalConfigResponse
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

		// var req updateconfig.InportRequest
		req := updateconfig.InportRequest(jsonReq)
		req.TimeNow = time.Now()

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
