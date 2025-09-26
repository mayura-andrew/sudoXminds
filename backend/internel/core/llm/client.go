package llm

import (
	"context"
	"fmt"
	"mathprereq/internel/core/config"
	"mathprereq/internel/types"
	"mathprereq/pkg/logger"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genai"
)

type Client struct {
	genaiClient *genai.Client
	config      config.LLMConfig
	ctx         context.Context
	cancel      context.CancelFunc
	logger      *zap.Logger
}

const (
	DefaultModel      = "gemini-2.0-flash-exp"
	DefaultMaxTokens  = 4000
	DefaultTimeout    = 60 * time.Second
	HealthCheckPrompt = "Respond with 'OK' to confirm you are working."
)

type ExplanationRequest struct {
	Query            string          `json:"query"`
	PrerequisitePath []types.Concept `json:"prerequisite_path"`
	ContextChunks    []string        `json:"context_chunks"`
}

func NewClient(cfg config.LLMConfig) (*Client, error) {
	logger := logger.MustGetLogger()
	logger.Info("Initializing Gemini LLM Client",
		zap.String("model", cfg.Model),
		zap.Bool("api_key_provided", cfg.APIKey != ""))

	ctx, cancel := context.WithCancel(context.Background())

	// Get API key with fallback priority
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("MLF_LLM_API_KEY")
	}
	if apiKey == "" {
		cancel()
		return nil, fmt.Errorf("Gemini API key not found. Set GEMINI_API_KEY, GOOGLE_API_KEY, or MLF_LLM_API_KEY environment variable")
	}

	genaiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})

	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize Gemini client: %w", err)
	}

	client := &Client{
		genaiClient: genaiClient,
		config:      cfg,
		ctx:         ctx,
		cancel:      cancel,
		logger:      logger,
	}

	logger.Info("Gemini LLM client initialized successfully",
		zap.String("model", cfg.Model),
		zap.String("provider", "gemini"))

	return client, nil
}
