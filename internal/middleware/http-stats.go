package middleware

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type HTTPAPIStats struct {
	APIStats map[string]int
	sync.RWMutex
}

type HTTPStats struct {
	TotalRequests HTTPAPIStats
	TotalErrors   HTTPAPIStats
	TotalCanceled HTTPAPIStats
}

var GlobalHTTPStats *HTTPStats

func init() {
	GlobalHTTPStats = &HTTPStats{
		TotalRequests: HTTPAPIStats{APIStats: make(map[string]int)},
		TotalErrors:   HTTPAPIStats{APIStats: make(map[string]int)},
		TotalCanceled: HTTPAPIStats{APIStats: make(map[string]int)},
	}
}

func (stats *HTTPAPIStats) Inc(api string) {
	if stats == nil {
		return
	}
	stats.Lock()
	defer stats.Unlock()
	stats.APIStats[api]++
}

func (stats *HTTPAPIStats) Dec(api string) {
	if stats == nil {
		return
	}
	stats.Lock()
	defer stats.Unlock()
	if val, ok := stats.APIStats[api]; ok && val > 0 {
		stats.APIStats[api]--
	}
}

func (stats *HTTPAPIStats) Load() map[string]int {
	stats.Lock()
	defer stats.Unlock()
	apiStats := make(map[string]int, len(stats.APIStats))
	for key, value := range stats.APIStats {
		apiStats[key] = value
	}
	return apiStats
}

func (stats *HTTPStats) updateStats(api string, ctx *gin.Context) {
	successReq := ctx.Writer.Status() >= 200 && ctx.Writer.Status() < 400
	if !strings.HasSuffix(ctx.FullPath(), "/metrics") {
		if !successReq {
			switch ctx.Writer.Status() {
			case 0:
			case 499:
				stats.TotalCanceled.Inc(api)
			default:
				stats.TotalErrors.Inc(api)
			}
		}
	}
}

func CollectAPIStats(api string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		GlobalHTTPStats.TotalRequests.Inc(api)
		ctx.Next()
		GlobalHTTPStats.updateStats(api, ctx)
	}
}
