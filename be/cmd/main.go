package main

import (
	"context"
	"os"

	_ "example.com/m/v2/docs"
	contextualizelink "example.com/m/v2/internal/components/contextualize_link"
	redirecturl "example.com/m/v2/internal/components/redirect_url"
	"example.com/m/v2/internal/repositories"
	"example.com/m/v2/internal/services"
	"go.uber.org/zap"

	logger "example.com/m/v2/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRequestID(c *gin.Context) {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	ctx := context.WithValue(c.Request.Context(), logger.RequestIDKey, requestID)
	c.Request = c.Request.WithContext(ctx)
}

func main() {
	logger.Init(false)

	if err := godotenv.Load(); err != nil {
		logger.Error(context.Background(), "Warning: .env file could not be loaded:", zap.Error(err))
	}

	difyBaseURL := "https://api.dify.ai/v1/workflows/run"
	difyContextualAPIKey := os.Getenv("DIFY_CONTEXTUAL_API_KEY")

	difyService := services.NewDifyService(difyBaseURL, difyContextualAPIKey)

	supabaseClient, err := services.NewSupabaseClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_PRIVATE_API_KEY"))
	if err != nil {
		logger.Error(context.Background(), "Error: failed to initialize supabase client:", zap.Error(err))
		os.Exit(1)
	}

	contextualLinkRepository := repositories.NewContextualLinkRepository(supabaseClient)

	contextualizeLinkHandler := contextualizelink.NewContextualizeLinkHandler(difyService, contextualLinkRepository)
	redirectURLHandler := redirecturl.NewRedirectURLHandler(contextualLinkRepository)

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(setupRequestID)

	// Health check endpoint
	router.GET("/healthz", HealthCheck)

	// Contextualize link endpoint
	router.GET("/contextualize-link", func(c *gin.Context) {
		ContextualizeLinkHandler(c, contextualizeLinkHandler)
	})

	// Redirect URL endpoint
	router.GET("/redirect", func(c *gin.Context) {
		RedirectURLHandler(c, redirectURLHandler)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")
}
