package main

import (
	"net/http"

	_ "example.com/m/v2/docs"
	"example.com/m/v2/internal/services"
	"go.uber.org/zap"

	logger "example.com/m/v2/internal/pkg"
	"github.com/gin-gonic/gin"
)

// @title AI Tools API
// @version 1.0
// @description This is the AI Tools API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthz [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

// Handler is the HTTP handler for contextualizing links.
// @Summary Contextualize a link
// @Description Get contextualized version of a link using AI
// @Tags contextualize
// @Accept json
// @Produce json
// @Param link query string true "URL to contextualize"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /contextualize-link [get]
func ContextualizeLinkHandler(c *gin.Context, difyService *services.DifyService) {
	ctx := c.Request.Context()

	link := c.Query("link")

	if link == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "link query parameter is required",
		})
		return
	}

	logger.Info(ctx, "link: ", zap.String("link", link))

	contextualizedLink, err := difyService.GetContextualLink(link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Info(ctx, "contextualizedLink: ", zap.String("contextualizedLink", contextualizedLink))

	c.JSON(http.StatusOK, gin.H{
		"contextualized_link": contextualizedLink,
	})

}
