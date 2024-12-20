package newsapi

import (
	"net/http"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/barthr/newsapi"
)

func GenerateNewsApi() *newsapi.Client {
	newsapi_key := env.GetEnv("NEWS_API_KEY", "my-news-api-key")
	c := newsapi.NewClient(newsapi_key, newsapi.WithHTTPClient(http.DefaultClient), newsapi.WithUserAgent("GoNews/1.0 (+https://gonews-example.com)"))

	return c
}
