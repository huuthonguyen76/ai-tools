package main

import (
	"net/http"

	_ "example.com/m/v2/docs"
	"go.uber.org/zap"

	contextualizelink "example.com/m/v2/internal/components/contextualize_link"
	redirecturl "example.com/m/v2/internal/components/redirect_url"
	"example.com/m/v2/internal/pkg"
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
func ContextualizeLinkHandler(c *gin.Context, contextualHandler *contextualizelink.ContextualizeLinkHandler) {
	ctx := c.Request.Context()

	link := c.Query("link")

	if link == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "link query parameter is required",
		})
		return
	}

	logger.Info(ctx, "link: ", zap.String("link", link))

	contextualizedLink, err := contextualHandler.Handler(link)
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

// RedirectURLHandler is the HTTP handler for redirecting contextualized links to original URLs.
// @Summary Redirect to original URL
// @Description Receives a contextualized link and redirects to the original URL
// @Tags redirect
// @Accept json
// @Produce json
// @Param contextualizedLink query string true "Contextualized link to redirect"
// @Success 302 {string} string "Redirect to original URL"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /redirect [get]
func RedirectURLHandler(c *gin.Context, redirectHandler *redirecturl.RedirectURLHandler) {
	ctx := c.Request.Context()

	responseFormat := pkg.NewResponseFormat(ctx, http.StatusInternalServerError, "", nil)

	contextualizedLink := c.Query("contextualizedLink")

	if contextualizedLink == "" {
		responseFormat.ErrorMsg = "contextualizedLink query parameter is required"
		c.JSON(http.StatusBadRequest, responseFormat)
		return
	}

	logger.Info(ctx, "contextualizedLink: ", zap.String("contextualizedLink", contextualizedLink))

	originalLink, err := redirectHandler.Handler(contextualizedLink)
	if err != nil {
		responseFormat.ErrorMsg = err.Error()
		c.JSON(http.StatusInternalServerError, responseFormat)
		return
	}

	logger.Info(ctx, "redirecting to original link: ", zap.String("originalLink", originalLink))

	c.Redirect(http.StatusFound, originalLink)
}
