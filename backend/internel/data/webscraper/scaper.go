package webscraper

import (
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type EducationalResource struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConceptID       string             `bson:"concept_id" json:"concept_id"`
	ConceptName     string             `bson:"concept_name" json:"concept_name"`
	Title           string             `bson:"title" json:"title"`
	URL             string             `bson:"url" json:"url"`
	Description     string             `bson:"description" json:"description"`
	ResourceType    string             `bson:"resource_type" json:"resource_type"` // video, article, tutorial, example, practice
	SourceDomain    string             `bson:"source_domain" json:"source_domain"`
	DifficultyLevel string             `bson:"difficulty_level" json:"difficulty_level"` // beginner, intermediate, advanced
	QualityScore    float64            `bson:"quality_score" json:"quality_score"`       // 0.0 to 1.0
	ContentPreview  string             `bson:"content_preview" json:"content_preview"`
	ScrapedAt       time.Time          `bson:"scraped_at" json:"scraped_at"`
	Language        string             `bson:"language" json:"language"`
	Duration        *string            `bson:"duration,omitempty" json:"duration,omitempty"`           // For videos
	ThumbnailURL    *string            `bson:"thumbnail_url,omitempty" json:"thumbnail_url,omitempty"` // For videos
	ViewCount       *int64             `bson:"view_count,omitempty" json:"view_count,omitempty"`
	Rating          *float64           `bson:"rating,omitempty" json:"rating,omitempty"`
	AuthorChannel   *string            `bson:"author_channel,omitempty" json:"author_channel,omitempty"`
	PublishedAt     *time.Time         `bson:"published_at,omitempty" json:"published_at,omitempty"`
	Tags            []string           `bson:"tags" json:"tags"`
	IsVerified      bool               `bson:"is_verified" json:"is_verified"`
}

type ScraperConfig struct {
	MaxConcurrentRequests int           `json:"max_concurrent_requests"`
	RequestTimeout        time.Duration `json:"request_timeout"`
	RateLimit             float64       `json:"rate_limit"`
	UserAgent             string        `json:"user_agent"`
	MongoURI              string        `json:"mongo_uri"`
	DatabaseName          string        `json:"database_name"`
	CollectionName        string        `json:"collection_name"`
	MaxRetries            int           `json:"max_retries"`
	RetryDelay            time.Duration `json:"retry_delay"`
}

type EducationalWebScraper struct {
	config       ScraperConfig
	httpClient   *http.Client
	limiter      *rate.Limiter
	mongoClient  *mongo.Client
	collection   *mongo.Collection
	logger       *zap.Logger
	scrapedURLs  sync.Map // Thread-safe cache of scraped URLs
	sharedClient bool     // Whether we're using a shared MongoDB client

	// Educational domains to target
	educationalDomains []string
}

type YouTubeVideoData struct {
	VideoID       string `json:"videoId"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      string `json:"duration"`
	ViewCount     string `json:"views"`
	Channel       string `json:"channel"`
	ThumbnailURL  string `json:"thumbnail"`
	PublishedTime string `json:"publishedTime"`
}
