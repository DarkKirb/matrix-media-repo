package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var HttpRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_http_requests_total",
}, []string{"host", "action", "method"})
var InvalidHttpRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_invalid_http_requests_total",
}, []string{"action", "method"})
var HttpResponses = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_http_responses_total",
}, []string{"host", "action", "method", "statusCode"})
var HttpResponseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "media_http_response_time_seconds",
}, []string{"host", "action", "method"})
var CacheHits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_cache_hits_total",
}, []string{"cache"})
var CacheMisses = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_cache_misses_total",
}, []string{"cache"})
var ThumbnailsGenerated = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_thumbnails_generated_total",
}, []string{"width", "height", "method", "animated", "origin"})
var MediaDownloaded = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_downloaded_total",
}, []string{"origin"})
var UrlPreviewsGenerated = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_url_previews_generated_total",
}, []string{"type"})
var S3Operations = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "media_s3_operations_total",
}, []string{"operation"})

func init() {
	prometheus.MustRegister(HttpRequests)
	prometheus.MustRegister(InvalidHttpRequests)
	prometheus.MustRegister(HttpResponses)
	prometheus.MustRegister(HttpResponseTime)
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(CacheMisses)
	prometheus.MustRegister(ThumbnailsGenerated)
	prometheus.MustRegister(MediaDownloaded)
	prometheus.MustRegister(UrlPreviewsGenerated)
	prometheus.MustRegister(S3Operations)
}
