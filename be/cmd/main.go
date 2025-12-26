package main

import (
	"context"
	"os"
	"time"

	_ "example.com/m/v2/docs"
	contextualizelink "example.com/m/v2/internal/components/contextualize_link"
	"example.com/m/v2/internal/components/cronjob/classify"
	syncdatabase "example.com/m/v2/internal/components/cronjob/sync_data"
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
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		logger.Error(ctx, "Warning: .env file could not be loaded:", zap.Error(err))
	}

	difyBaseURL := "https://api.dify.ai/v1/workflows/run"
	difyContextualAPIKey := os.Getenv("DIFY_CONTEXTUAL_API_KEY")

	difyService := services.NewDifyService(difyBaseURL, difyContextualAPIKey)

	supabaseClient, err := services.NewSupabaseClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_PRIVATE_API_KEY"))
	if err != nil {
		logger.Error(ctx, "Error: failed to initialize supabase client:", zap.Error(err))
		os.Exit(1)
	}

	contextualLinkRepository := repositories.NewContextualLinkRepository(supabaseClient)

	contextualizeLinkHandler := contextualizelink.NewContextualizeLinkHandler(difyService, contextualLinkRepository)
	redirectURLHandler := redirecturl.NewRedirectURLHandler(contextualLinkRepository)

	// Initialize Apify service
	apifyAPIToken := os.Getenv("APIFY_API_KEY")
	apifyService, err := services.NewApifyService(apifyAPIToken)
	if err != nil {
		logger.Error(ctx, "Error: failed to initialize apify service:", zap.Error(err))
		os.Exit(1)
	}

	socialPostRawRepository := repositories.NewSocialPostRawRepository(supabaseClient)
	syncDatasetCronjob, err := syncdatabase.NewSyncService(apifyService, supabaseClient, socialPostRawRepository)
	if err != nil {
		logger.Error(ctx, "Error: failed to initialize sync dataset cronjob:", zap.Error(err))
		os.Exit(1)
	}

	classifySocialPostCronjob := classify.NewClassifySocialPost(socialPostRawRepository)

	go func(ctx context.Context) {
		for {
			time.Sleep(1 * time.Hour)
			syncDatasetCronjob.SyncAllDatasets(ctx)
		}
	}(ctx)

	go func(ctx context.Context) {
		for {
			classifySocialPostCronjob.ClassifySocialPost(ctx)
			time.Sleep(1 * time.Hour)
		}
	}(ctx)

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
	router.GET("/redirect/:client/:link", func(c *gin.Context) {
		client := c.Param("client")
		link := c.Param("link")
		RedirectURLHandler(c, redirectURLHandler, client, link)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":2000")
}
