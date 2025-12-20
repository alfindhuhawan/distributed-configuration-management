package restapi

import (
	"context"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	"distributed-configuration-management/domain_controllerappservice/usecase/getconfig"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/model/payload"
	"distributed-configuration-management/shared/model/repository"
)

// getConfigHandler ...
func (r *Controller) getConfigHandler(inputPort getconfig.Inport) gin.HandlerFunc {

	type request struct {
		TypeConfig string `form:"type"`
	}

	type response struct {
		repository.GlobalConfigResponse
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
		req.Type = jsonReq.TypeConfig

		var latestVersion float64
		latestVersionStr := c.Request.Header.Get("ETag")
		if reflect.TypeOf(latestVersionStr).Kind() == reflect.String {
			latestVersion, _ = strconv.ParseFloat(latestVersionStr, 64)
		}
		req.Version = latestVersion

		r.Log.Info(ctx, util.MustJSON(req))

		res, err := inputPort.Execute(ctx, req)
		if err != nil {
			r.Log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		jsonRes := response(*res)

		r.Log.Info(ctx, util.MustJSON(jsonRes))

		if !jsonRes.IsDataChanged {
			c.JSON(http.StatusOK, payload.NoChangeResponse(jsonRes, traceID))
		} else {
			c.JSON(http.StatusOK, payload.NewSuccessResponse(jsonRes, traceID))
		}
	}
}
