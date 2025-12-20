package restapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"distributed-configuration-management/domain_controllerappservice/usecase/registeragent"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/model/payload"
	"distributed-configuration-management/shared/model/repository"
)

// registerAgentHandler ...
func (r *Controller) registerAgentHandler(inputPort registeragent.Inport) gin.HandlerFunc {

	type request struct {
		repository.RegisterAgentRequest
	}

	type response struct {
		repository.RegisterAgentResponse
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

		// var req registeragent.InportRequest
		req := registeragent.InportRequest(jsonReq)
		req.Name = strings.TrimSpace(jsonReq.Name)
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
