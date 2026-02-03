package handlers

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/maxjove/defi-yield-aggregator/internal/api/middleware"
)

// MetricsResponse contains application metrics
type MetricsResponse struct {
	Timestamp     string         `json:"timestamp"`
	Uptime        string         `json:"uptime"`
	Go            GoMetrics      `json:"go"`
	HTTP          HTTPMetrics    `json:"http"`
	Memory        MemoryMetrics  `json:"memory"`
}

// GoMetrics contains Go runtime metrics
type GoMetrics struct {
	Version     string `json:"version"`
	NumGoroutine int   `json:"numGoroutine"`
	NumCPU      int    `json:"numCpu"`
}

// HTTPMetrics contains HTTP request metrics
type HTTPMetrics struct {
	TotalRequests    int64            `json:"totalRequests"`
	SuccessRequests  int64            `json:"successRequests"`
	ErrorRequests    int64            `json:"errorRequests"`
	AvgLatencyMs     float64          `json:"avgLatencyMs"`
	RequestsByStatus map[int]int64    `json:"requestsByStatus"`
}

// MemoryMetrics contains memory usage metrics
type MemoryMetrics struct {
	Alloc      string `json:"alloc"`
	TotalAlloc string `json:"totalAlloc"`
	Sys        string `json:"sys"`
	NumGC      uint32 `json:"numGc"`
}

// GetMetrics returns application metrics
// GET /api/v1/metrics
func (h *Handler) GetMetrics(c *fiber.Ctx) error {
	// Get runtime metrics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Get HTTP metrics
	httpMetrics := middleware.GetMetrics()

	// Calculate average latency
	var avgLatency float64
	if httpMetrics.TotalRequests > 0 {
		avgLatency = float64(httpMetrics.TotalLatencyMs) / float64(httpMetrics.TotalRequests)
	}

	response := MetricsResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(h.startTime).String(),
		Go: GoMetrics{
			Version:      runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
		},
		HTTP: HTTPMetrics{
			TotalRequests:    httpMetrics.TotalRequests,
			SuccessRequests:  httpMetrics.SuccessRequests,
			ErrorRequests:    httpMetrics.ErrorRequests,
			AvgLatencyMs:     avgLatency,
			RequestsByStatus: httpMetrics.RequestsByStatus,
		},
		Memory: MemoryMetrics{
			Alloc:      formatBytes(memStats.Alloc),
			TotalAlloc: formatBytes(memStats.TotalAlloc),
			Sys:        formatBytes(memStats.Sys),
			NumGC:      memStats.NumGC,
		},
	}

	return c.JSON(response)
}

// GetPrometheusMetrics returns metrics in Prometheus format
// GET /metrics
func (h *Handler) GetPrometheusMetrics(c *fiber.Ctx) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	httpMetrics := middleware.GetMetrics()

	// Build Prometheus format output
	output := fmt.Sprintf(`# HELP defi_http_requests_total Total number of HTTP requests
# TYPE defi_http_requests_total counter
defi_http_requests_total %d

# HELP defi_http_requests_success_total Total number of successful HTTP requests
# TYPE defi_http_requests_success_total counter
defi_http_requests_success_total %d

# HELP defi_http_requests_error_total Total number of failed HTTP requests
# TYPE defi_http_requests_error_total counter
defi_http_requests_error_total %d

# HELP defi_http_request_latency_ms_total Total HTTP request latency in milliseconds
# TYPE defi_http_request_latency_ms_total counter
defi_http_request_latency_ms_total %d

# HELP defi_go_goroutines Number of goroutines
# TYPE defi_go_goroutines gauge
defi_go_goroutines %d

# HELP defi_go_memory_alloc_bytes Current memory allocation in bytes
# TYPE defi_go_memory_alloc_bytes gauge
defi_go_memory_alloc_bytes %d

# HELP defi_go_memory_sys_bytes Total memory obtained from system
# TYPE defi_go_memory_sys_bytes gauge
defi_go_memory_sys_bytes %d

# HELP defi_go_gc_runs_total Total number of GC runs
# TYPE defi_go_gc_runs_total counter
defi_go_gc_runs_total %d

# HELP defi_uptime_seconds Service uptime in seconds
# TYPE defi_uptime_seconds gauge
defi_uptime_seconds %.0f
`,
		httpMetrics.TotalRequests,
		httpMetrics.SuccessRequests,
		httpMetrics.ErrorRequests,
		httpMetrics.TotalLatencyMs,
		runtime.NumGoroutine(),
		memStats.Alloc,
		memStats.Sys,
		memStats.NumGC,
		time.Since(h.startTime).Seconds(),
	)

	c.Set("Content-Type", "text/plain; charset=utf-8")
	return c.SendString(output)
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
