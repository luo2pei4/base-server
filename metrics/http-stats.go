package metrics

import (
	"net/http"
	"strings"
	"sync"
)

type HTTPAPIStats struct {
	apiStats map[string]int
	sync.RWMutex
}

type HTTPStats struct {
	totalRequests HTTPAPIStats
	totalErrors   HTTPAPIStats
	totalCanceled HTTPAPIStats
}

func (stats *HTTPAPIStats) Inc(api string) {
	if stats == nil {
		return
	}
	stats.Lock()
	defer stats.Unlock()
	if stats.apiStats == nil {
		stats.apiStats = make(map[string]int)
	}
	stats.apiStats[api]++
}

func (stats *HTTPAPIStats) Dec(api string) {
	if stats == nil {
		return
	}
	stats.Lock()
	defer stats.Unlock()
	if val, ok := stats.apiStats[api]; ok && val > 0 {
		stats.apiStats[api]--
	}
}

func (stats *HTTPAPIStats) Load() map[string]int {
	stats.Lock()
	defer stats.Unlock()
	var apiStats = make(map[string]int, len(stats.apiStats))
	for k, v := range stats.apiStats {
		apiStats[k] = v
	}
	return apiStats
}

func NewHTTPStats() *HTTPStats {
	return &HTTPStats{}
}

func (st *HTTPStats) UpdateStats(api string, r *http.Request, w *http.Response) {
	successReq := w.StatusCode >= 200 && w.StatusCode < 400
	if !strings.HasPrefix(r.URL.Path, "/metrics") {
		st.totalRequests.Inc(api)
		if !successReq {
			switch w.StatusCode {
			case 0:
			case 499:
				st.totalCanceled.Inc(api)
			default:
				st.totalErrors.Inc(api)
			}
		}
	}
}
