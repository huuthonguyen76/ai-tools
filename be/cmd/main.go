package main

import (
	"context"
	"os"

	_ "example.com/m/v2/docs"
	contextualizelink "example.com/m/v2/internal/components/contextualize_link"
	"example.com/m/v2/internal/repositories"
	"example.com/m/v2/internal/services"
	"go.uber.org/zap"

	logger "example.com/m/v2/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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

	router := gin.Default()

	// Health check endpoint
	router.GET("/healthz", HealthCheck)

	// Contextualize link endpoint
	router.GET("/contextualize-link", func(c *gin.Context) {
		ContextualizeLinkHandler(c, contextualizeLinkHandler)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")
}
