package services

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/jinzhu/gorm"
)

/***
 es 配置
type Config struct {
	Addresses []string // A list of Elasticsearch nodes to use.
	Username  string   // Username for HTTP Basic Authentication.
	Password  string   // Password for HTTP Basic Authentication.

	CloudID string // Endpoint for the Elastic Service (https://elastic.co/cloud).
	APIKey  string // Base64-encoded token for authorization; if set, overrides username and password.

	RetryOnStatus        []int // List of status codes for retry. Default: 502, 503, 504.
	DisableRetry         bool  // Default: false.
	EnableRetryOnTimeout bool  // Default: false.
	MaxRetries           int   // Default: 3.

	RetryBackoff func(attempt int) time.Duration // Optional backoff duration. Default: nil.

	Transport http.RoundTripper  // The HTTP transport object.
	Logger    estransport.Logger // The logger object.
}
*/

type EsService struct {
	app   interfaces.App
	db    *gorm.DB
	redis *redis.Client
	es    *elasticsearch.Client
}

func NewEsService(app interfaces.App) (*EsService, error) {
	db, err := app.GetDefaultDb()
	if err != nil {
		return nil, err
	}

	redisClient, err := app.GetDefaultRedis()
	if err != nil {
		return nil, err
	}

	addresses := app.GetConfig().GetStringSlice("es.Address")
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &EsService{
		app:   app,
		db:    db,
		redis: redisClient,
		es:    es,
	}, nil
}

/***
es.Delete("test", "1")
es.Exists("test", "1")

es.Index(
	"test",
	strings.NewReader(`{"title" : "logging"}`),
	es.Index.WithDocumentID("1"),
	es.Index.WithRefresh("true"),
	es.Index.WithPretty(),
	es.Index.WithFilterPath("result", "_id"),
)

es.Search(
	es.Search.WithIndex("test"),
	es.Search.WithQuery("{FAIL"))

res, err := es.Search(
	es.Search.WithIndex("test"),
	es.Search.WithBody(strings.NewReader(`{"query" : {"match" : { "title" : "logging" } } }`)),
	es.Search.WithSize(1),
	es.Search.WithPretty(),
	es.Search.WithFilterPath("took", "hits.hits"),
)
*/
func (s *EsService) GetEs() *elasticsearch.Client {
	return s.es
}

func (s *EsService) GetServiceName() string {
	return "EsService"
}
