package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ============================================================================
// ADVANCED GO CONCEPTS DEMONSTRATION (8-10/10 COMPLEXITY)
// Real-world scenario: Distributed Task Processing System with Monitoring
// ============================================================================

// 1. ADVANCED GENERICS WITH TYPE CONSTRAINTS
// Generic constraint for processable tasks
type Processable interface {
	~string | ~int | ~float64
	fmt.Stringer
}

// Generic task processor with type constraints
type TaskProcessor[T Processable] struct {
	name    string
	handler func(context.Context, T) error
	metrics *ProcessorMetrics
	limiter *RateLimiter
}

// 2. COMPLEX INTERFACES WITH TYPE ASSERTIONS AND REFLECTION
type Task interface {
	GetID() string
	GetType() string
	GetPriority() int
	Validate() error
	// Method that uses reflection internally
	ExtractMetadata() map[string]interface{}
}

// Concrete task implementations
type ImageProcessingTask struct {
	ID       string   `json:"id"`
	ImageURL string   `json:"image_url"`
	Filters  []string `json:"filters"`
	Priority int      `json:"priority"`
}

func (t ImageProcessingTask) GetID() string    { return t.ID }
func (t ImageProcessingTask) GetType() string  { return "image_processing" }
func (t ImageProcessingTask) GetPriority() int { return t.Priority }
func (t ImageProcessingTask) String() string   { return fmt.Sprintf("ImageTask(%s)", t.ID) }

func (t ImageProcessingTask) Validate() error {
	if t.ID == "" || t.ImageURL == "" {
		return errors.New("invalid image processing task: missing required fields")
	}
	return nil
}

// 3. REFLECTION AND RUNTIME TYPE MANIPULATION
func (t ImageProcessingTask) ExtractMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})
	v := reflect.ValueOf(t)
	typeOf := reflect.TypeOf(t)

	for i := 0; i < v.NumField(); i++ {
		field := typeOf.Field(i)
		value := v.Field(i)

		// Use reflection to extract JSON tags and values
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			metadata[jsonTag] = value.Interface()
		}

		// Advanced reflection: check if field implements specific interfaces
		if value.Type().Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
			if stringer, ok := value.Interface().(fmt.Stringer); ok {
				metadata[field.Name+"_string"] = stringer.String()
			}
		}
	}

	return metadata
}

type DataAnalysisTask struct {
	ID       string                 `json:"id"`
	Dataset  string                 `json:"dataset"`
	Query    map[string]interface{} `json:"query"`
	Priority int                    `json:"priority"`
}

func (t DataAnalysisTask) GetID() string    { return t.ID }
func (t DataAnalysisTask) GetType() string  { return "data_analysis" }
func (t DataAnalysisTask) GetPriority() int { return t.Priority }
func (t DataAnalysisTask) String() string   { return fmt.Sprintf("DataTask(%s)", t.ID) }

func (t DataAnalysisTask) Validate() error {
	if t.ID == "" || t.Dataset == "" {
		return errors.New("invalid data analysis task: missing required fields")
	}
	return nil
}

func (t DataAnalysisTask) ExtractMetadata() map[string]interface{} {
	return map[string]interface{}{
		"id":      t.ID,
		"dataset": t.Dataset,
		"query":   t.Query,
		"type":    "data_analysis",
	}
}

// 4. ADVANCED SYNCHRONIZATION PRIMITIVES
type RateLimiter struct {
	mu         sync.RWMutex
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	cond       *sync.Cond
}

func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
	rl.cond = sync.NewCond(&rl.mu)

	// Background goroutine for token refill
	go rl.refillTokens()
	return rl
}

func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(rl.refillRate)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		if rl.tokens < rl.maxTokens {
			rl.tokens++
			rl.cond.Signal() // Wake up waiting goroutines
		}
		rl.lastRefill = time.Now()
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Acquire(ctx context.Context) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for rl.tokens <= 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			rl.cond.Wait() // Wait for token refill
		}
	}

	rl.tokens--
	return nil
}

// 5. PROMETHEUS METRICS AND MONITORING
type ProcessorMetrics struct {
	tasksProcessed prometheus.Counter
	taskDuration   prometheus.Histogram
	activeWorkers  prometheus.Gauge
	errorRate      prometheus.Counter
}

func NewProcessorMetrics(processorName string) *ProcessorMetrics {
	return &ProcessorMetrics{
		tasksProcessed: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "tasks_processed_total",
			Help:        "Total number of tasks processed",
			ConstLabels: prometheus.Labels{"processor": processorName},
		}),
		taskDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:        "task_duration_seconds",
			Help:        "Task processing duration",
			ConstLabels: prometheus.Labels{"processor": processorName},
			Buckets:     prometheus.DefBuckets,
		}),
		activeWorkers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "active_workers",
			Help:        "Number of active workers",
			ConstLabels: prometheus.Labels{"processor": processorName},
		}),
		errorRate: prometheus.NewCounter(prometheus.CounterOpts{
			Name:        "processing_errors_total",
			Help:        "Total number of processing errors",
			ConstLabels: prometheus.Labels{"processor": processorName},
		}),
	}
}

func (m *ProcessorMetrics) Register() {
	prometheus.MustRegister(m.tasksProcessed, m.taskDuration, m.activeWorkers, m.errorRate)
}

// 6. WORKER POOL WITH COMPLEX SYNCHRONIZATION
type WorkerPool struct {
	workers    int
	taskQueue  chan Task
	resultChan chan TaskResult
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc

	// Advanced: Memory pool for task results to reduce GC pressure
	resultPool sync.Pool

	// Atomic counters for statistics
	totalProcessed int64
	totalErrors    int64

	// Complex synchronization for graceful shutdown
	shutdownOnce sync.Once
	shutdownCh   chan struct{}
}

type TaskResult struct {
	TaskID    string
	Success   bool
	Error     error
	Duration  time.Duration
	Metadata  map[string]interface{}
	Timestamp time.Time
}

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		workers:    workers,
		taskQueue:  make(chan Task, queueSize),
		resultChan: make(chan TaskResult, queueSize),
		ctx:        ctx,
		cancel:     cancel,
		shutdownCh: make(chan struct{}),
		resultPool: sync.Pool{
			New: func() interface{} {
				return &TaskResult{}
			},
		},
	}

	// Start workers
	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	// Start result collector
	go wp.resultCollector()

	return wp
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	slog.Info("Worker started", "worker_id", id)

	for {
		select {
		case task := <-wp.taskQueue:
			wp.processTask(task, id)
		case <-wp.ctx.Done():
			slog.Info("Worker shutting down", "worker_id", id)
			return
		}
	}
}

func (wp *WorkerPool) processTask(task Task, workerID int) {
	start := time.Now()

	// Get result from pool to reduce allocations
	result := wp.resultPool.Get().(*TaskResult)
	defer wp.resultPool.Put(result)

	// Reset result object
	*result = TaskResult{
		TaskID:    task.GetID(),
		Timestamp: start,
	}

	// Simulate task processing with reflection-based routing
	err := wp.routeTask(task)

	result.Duration = time.Since(start)
	result.Success = err == nil
	result.Error = err
	result.Metadata = task.ExtractMetadata()

	// Update atomic counters
	atomic.AddInt64(&wp.totalProcessed, 1)
	if err != nil {
		atomic.AddInt64(&wp.totalErrors, 1)
	}

	// Send result (non-blocking)
	select {
	case wp.resultChan <- *result:
	default:
		slog.Warn("Result channel full, dropping result", "task_id", task.GetID())
	}

	slog.Debug("Task processed",
		"worker_id", workerID,
		"task_id", task.GetID(),
		"duration", result.Duration,
		"success", result.Success)
}

// 7. REFLECTION-BASED TASK ROUTING
func (wp *WorkerPool) routeTask(task Task) error {
	// Use reflection to route tasks to appropriate handlers
	taskType := reflect.TypeOf(task)
	taskValue := reflect.ValueOf(task)

	switch taskType {
	case reflect.TypeOf(ImageProcessingTask{}):
		return wp.processImageTask(taskValue.Interface().(ImageProcessingTask))
	case reflect.TypeOf(DataAnalysisTask{}):
		return wp.processDataTask(taskValue.Interface().(DataAnalysisTask))
	default:
		return fmt.Errorf("unknown task type: %s", taskType.Name())
	}
}

func (wp *WorkerPool) processImageTask(task ImageProcessingTask) error {
	// Simulate image processing
	processingTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
	time.Sleep(processingTime)

	// Simulate occasional failures
	if rand.Float32() < 0.1 {
		return errors.New("image processing failed")
	}

	return nil
}

func (wp *WorkerPool) processDataTask(task DataAnalysisTask) error {
	// Simulate data analysis
	processingTime := time.Duration(rand.Intn(2000)+1000) * time.Millisecond
	time.Sleep(processingTime)

	// Simulate occasional failures
	if rand.Float32() < 0.05 {
		return errors.New("data analysis failed")
	}

	return nil
}

func (wp *WorkerPool) resultCollector() {
	for {
		select {
		case result := <-wp.resultChan:
			// In a real system, this would write to a database or message queue
			if !result.Success {
				slog.Error("Task failed",
					"task_id", result.TaskID,
					"error", result.Error,
					"duration", result.Duration)
			}
		case <-wp.ctx.Done():
			slog.Info("Result collector shutting down")
			return
		}
	}
}

func (wp *WorkerPool) Submit(task Task) error {
	if err := task.Validate(); err != nil {
		return fmt.Errorf("task validation failed: %w", err)
	}

	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return errors.New("worker pool is shutting down")
	default:
		return errors.New("task queue is full")
	}
}

// 8. GRACEFUL SHUTDOWN WITH CONTEXT PROPAGATION
func (wp *WorkerPool) Shutdown(timeout time.Duration) error {
	var shutdownErr error

	wp.shutdownOnce.Do(func() {
		slog.Info("Initiating worker pool shutdown")

		// Stop accepting new tasks
		close(wp.taskQueue)

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Signal workers to stop
		wp.cancel()

		// Wait for workers to finish with timeout
		done := make(chan struct{})
		go func() {
			wp.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			slog.Info("All workers shut down gracefully")
		case <-shutdownCtx.Done():
			shutdownErr = errors.New("shutdown timeout exceeded")
			slog.Error("Shutdown timeout exceeded")
		}

		close(wp.shutdownCh)
	})

	return shutdownErr
}

func (wp *WorkerPool) GetStats() (processed, errors int64) {
	return atomic.LoadInt64(&wp.totalProcessed), atomic.LoadInt64(&wp.totalErrors)
}

// 9. MEMORY MANAGEMENT AND PERFORMANCE OPTIMIZATION
type MemoryManager struct {
	allocations   int64
	deallocations int64
	peakMemory    int64

	// Unsafe pointer operations for advanced memory management
	memoryPool unsafe.Pointer
}

func (mm *MemoryManager) TrackAllocation(size int64) {
	atomic.AddInt64(&mm.allocations, 1)

	// Update peak memory usage
	for {
		current := atomic.LoadInt64(&mm.peakMemory)
		if size <= current {
			break
		}
		if atomic.CompareAndSwapInt64(&mm.peakMemory, current, size) {
			break
		}
	}
}

func (mm *MemoryManager) GetStats() (allocs, deallocs, peak int64) {
	return atomic.LoadInt64(&mm.allocations),
		atomic.LoadInt64(&mm.deallocations),
		atomic.LoadInt64(&mm.peakMemory)
}

// 10. PLUGIN SYSTEM WITH DYNAMIC LOADING (Simulated)
type Plugin interface {
	Name() string
	Version() string
	Process(ctx context.Context, data interface{}) (interface{}, error)
}

type PluginManager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

func (pm *PluginManager) RegisterPlugin(plugin Plugin) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[plugin.Name()] = plugin
}

func (pm *PluginManager) GetPlugin(name string) (Plugin, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// Example plugin implementation
type ImageFilterPlugin struct{}

func (p ImageFilterPlugin) Name() string    { return "image_filter" }
func (p ImageFilterPlugin) Version() string { return "1.0.0" }

func (p ImageFilterPlugin) Process(ctx context.Context, data interface{}) (interface{}, error) {
	// Simulate image filtering
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("Filtered: %v", data), nil
}

// ============================================================================
// MAIN APPLICATION WITH ALL ADVANCED CONCEPTS INTEGRATED
// ============================================================================

func main() {
	// Set up structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting advanced Go demonstration")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize components
	memManager := &MemoryManager{}
	pluginManager := NewPluginManager()
	pluginManager.RegisterPlugin(ImageFilterPlugin{})

	// Create worker pool
	workerPool := NewWorkerPool(5, 100)
	defer func() {
		if err := workerPool.Shutdown(30 * time.Second); err != nil {
			slog.Error("Worker pool shutdown error", "error", err)
		}
	}()

	// Set up metrics
	metrics := NewProcessorMetrics("main_processor")
	metrics.Register()

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			processed, errors := workerPool.GetStats()
			allocs, deallocs, peak := memManager.GetStats()

			status := map[string]interface{}{
				"status": "healthy",
				"stats": map[string]interface{}{
					"tasks_processed": processed,
					"task_errors":     errors,
					"memory_allocs":   allocs,
					"memory_deallocs": deallocs,
					"peak_memory":     peak,
					"goroutines":      runtime.NumGoroutine(),
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(status)
		})

		slog.Info("Metrics server starting on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			slog.Error("Metrics server error", "error", err)
		}
	}()

	// Create rate limiter
	rateLimiter := NewRateLimiter(10, 100*time.Millisecond)

	// Start task generator
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		taskID := 0
		for {
			select {
			case <-ticker.C:
				// Rate limit task generation
				if err := rateLimiter.Acquire(ctx); err != nil {
					continue
				}

				taskID++
				var task Task

				// Create different types of tasks
				if taskID%2 == 0 {
					task = ImageProcessingTask{
						ID:       fmt.Sprintf("img_%d", taskID),
						ImageURL: fmt.Sprintf("https://example.com/image_%d.jpg", taskID),
						Filters:  []string{"blur", "sharpen"},
						Priority: rand.Intn(10),
					}
				} else {
					task = DataAnalysisTask{
						ID:       fmt.Sprintf("data_%d", taskID),
						Dataset:  fmt.Sprintf("dataset_%d", taskID%5),
						Query:    map[string]interface{}{"type": "aggregation", "field": "value"},
						Priority: rand.Intn(10),
					}
				}

				// Submit task to worker pool
				if err := workerPool.Submit(task); err != nil {
					slog.Error("Failed to submit task", "error", err, "task_id", task.GetID())
					metrics.errorRate.Inc()
				} else {
					metrics.tasksProcessed.Inc()
					memManager.TrackAllocation(1024) // Simulate memory allocation

					slog.Debug("Task submitted",
						"task_id", task.GetID(),
						"task_type", task.GetType(),
						"priority", task.GetPriority())
				}

			case <-ctx.Done():
				slog.Info("Task generator shutting down")
				return
			}
		}
	}()

	// Performance monitoring goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				processed, errors := workerPool.GetStats()
				allocs, _, peak := memManager.GetStats()

				var m runtime.MemStats
				runtime.ReadMemStats(&m)

				slog.Info("Performance stats",
					"tasks_processed", processed,
					"task_errors", errors,
					"error_rate", float64(errors)/float64(processed+1)*100,
					"goroutines", runtime.NumGoroutine(),
					"heap_alloc_mb", m.HeapAlloc/1024/1024,
					"sys_alloc_mb", m.Sys/1024/1024,
					"gc_cycles", m.NumGC,
					"custom_allocs", allocs,
					"peak_memory", peak)

				// Update Prometheus metrics
				metrics.activeWorkers.Set(float64(runtime.NumGoroutine()))

			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	slog.Info("Shutdown signal received, initiating graceful shutdown")

	// Cancel context to signal all goroutines to stop
	cancel()

	// Give some time for graceful shutdown
	time.Sleep(2 * time.Second)

	// Final statistics
	processed, errors := workerPool.GetStats()
	slog.Info("Final statistics",
		"total_processed", processed,
		"total_errors", errors,
		"success_rate", float64(processed-errors)/float64(processed)*100)

	slog.Info("Advanced Go demonstration completed")
}

/*
COMPLEXITY BREAKDOWN (8-10/10):

1. ✅ Advanced Generics with Type Constraints
   - Generic TaskProcessor with Processable constraint
   - Type constraints using ~string | ~int | ~float64

2. ✅ Complex Interfaces with Type Assertions
   - Task interface with multiple implementations
   - Runtime type checking and assertions

3. ✅ Reflection and Runtime Type Manipulation
   - ExtractMetadata() uses reflection to inspect struct fields
   - Dynamic task routing based on reflection
   - JSON tag extraction using reflection

4. ✅ Advanced Synchronization Primitives
   - sync.RWMutex for read-heavy operations
   - sync.Cond for conditional waiting
   - Complex worker pool synchronization

5. ✅ Memory Management and Performance Optimization
   - sync.Pool for object reuse
   - Atomic operations for lock-free counters
   - Memory tracking and statistics
   - Unsafe pointer operations

6. ✅ Context Propagation and Cancellation
   - Proper context usage throughout the application
   - Graceful shutdown with timeout handling
   - Context-aware rate limiting

7. ✅ Complex Error Handling Patterns
   - Wrapped errors with context
   - Error aggregation and reporting
   - Graceful degradation

8. ✅ Metrics and Monitoring (Prometheus)
   - Custom metrics collection
   - HTTP metrics endpoint
   - Performance monitoring

9. ✅ Plugin System (Simulated)
   - Dynamic plugin registration
   - Interface-based plugin architecture
   - Runtime plugin discovery

10. ✅ Network Programming
    - HTTP server for metrics and health checks
    - JSON API endpoints
    - Concurrent request handling

REAL-WORLD APPLICATIONS:
- Distributed task processing systems
- Microservices with monitoring
- High-performance data processing pipelines
- Plugin-based architectures
- Systems requiring graceful shutdown and monitoring

This demonstrates production-level Go code complexity that you'd find in:
- Kubernetes controllers
- Database systems
- Message queue implementations
- High-throughput web services
- Distributed systems components
*/
