package restapi

import (
	"distributed-configuration-management/shared/infrastructure/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Controller) authenticatedAdmin() gin.HandlerFunc {

	return func(c *gin.Context) {
		adminCredential := r.Config.Credentials.AdminCredential.SecretKey
		if adminCredential == "" {
			adminCredential = "AdminOnly"
		}

		// get header
		headerAuth := c.Request.Header.Get("Authorization")

		isAuthorized := validateHash(adminCredential, headerAuth)

		if !isAuthorized {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}

func (r *Controller) authenticatedAgent() gin.HandlerFunc {

	return func(c *gin.Context) {
		agentCredential := r.Config.Credentials.AgentCredential.SecretKey
		if agentCredential == "" {
			agentCredential = "AgentCredential"
		}

		// get header
		headerAuth := c.Request.Header.Get("Authorization")

		isAuthorized := validateHash(agentCredential, headerAuth)

		if !isAuthorized {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}

// authorized is an interceptor
func (r *Controller) authorized() gin.HandlerFunc {

	return func(c *gin.Context) {

		authorized := true

		if !authorized {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}

func validateHash(credential string, headerAuth string) bool {
	isAuthorized := true

	if len(headerAuth) <= 7 || headerAuth == "" {
		isAuthorized = false

		return isAuthorized
	}

	hashCredential := util.CreateSHA256Signature(credential)

	headerAuthWithoutBearer := headerAuth[7:]

	if hashCredential != headerAuthWithoutBearer {
		isAuthorized = false
	}

	return isAuthorized
}
