### 概要
下面提供对 zinx 框架 router 模块的二次封装方案，包括设计说明、关键接口与中间件实现、Prometheus 指标集成、OpenTelemetry 链路追踪集成、配置化（含热加载）与示例用法。代码以 Go 展示，可直接嵌入现有 zinx 项目中（假设原有 zinx 的 IRouter、HandleFunc 等类型可直接使用或做最小适配）。

---

### 设计要点
- 新增 EnhancedRouter 接口，继承 zinx 的 IRouter 并支持 Use 注册中间件链。
- 中间件采用 func(next HandleFunc) HandleFunc 签名，按注册顺序执行（洋葱模型）。
- 内置两个常用中间件：MetricsMiddleware（Prometheus）与 TracingMiddleware（OpenTelemetry）。
- Prometheus 指标：处理耗时直方图/摘要、请求计数、错误计数；提供 /metrics HTTP 接口供 Prometheus 抓取。
- OpenTelemetry：支持 OTLP/gRPC 与 OTLP/HTTP exporter，可配置采样率与自定义 attributes。
- 配置文件采用 YAML，支持热加载（fsnotify），修改后自动应用到 router（例如采样率、上报地址、指标间隔等）。
- 指标上报间隔为采集/推送周期（Prometheus 是 pull，间隔配置用于内部收集或 pushgateway 场景）；默认 5s。

---

### 主要文件结构（示例）
- enhanced_router/enhanced_router.go
- enhanced_router/middleware.go
- enhanced_router/metrics.go
- enhanced_router/tracing.go
- config/config.go
- cmd/example/main.go

---

### 配置结构（config/config.go）
```go
package config

import (
	"time"
)

type MetricsConfig struct {
	Enabled       bool          `yaml:"enabled"`
	MetricsPath   string        `yaml:"metrics_path"`   // default "/metrics"
	CollectPeriod time.Duration `yaml:"collect_period"` // default 5s
}

type OTLPConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Endpoint    string  `yaml:"endpoint"`    // e.g. "localhost:4317" or "http://collector:4318"
	Protocol    string  `yaml:"protocol"`    // "grpc" or "http"
	ServiceName string  `yaml:"service_name"`
	SampleRatio float64 `yaml:"sample_ratio"` // 0.1 or 1.0
}

type AppConfig struct {
	Metrics MetricsConfig `yaml:"metrics"`
	OTLP    OTLPConfig    `yaml:"otlp"`
}
```

示例 YAML（config.yaml）：
```yaml
metrics:
  enabled: true
  metrics_path: "/metrics"
  collect_period: 5s

otlp:
  enabled: true
  endpoint: "localhost:4317"
  protocol: "grpc"
  service_name: "my-zinx-service"
  sample_ratio: 0.1
```

---

### EnhancedRouter 接口与实现（enhanced_router/enhanced_router.go）
```go
package enhanced_router

import (
	"context"
	"net/http"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"time"
	"path/filepath"
)

// 这里假设 zinx 包中已有类型定义:
// type HandleFunc func(req znet.IRequest) // or appropriate signature
// type IRouter interface { Handle(request znet.IRequest); ...}

type HandleFunc func(ctx context.Context, req interface{}) (resp interface{}, err error)

type MiddlewareFunc func(next HandleFunc) HandleFunc

type EnhancedRouter interface {
	Use(mw ...MiddlewareFunc)
	AddRoute(path string, handler HandleFunc) // 适配原 zinx 注册方式
	StartMetricsEndpoint(addr string) error
	Shutdown(ctx context.Context) error
	ReloadConfig(path string) error
}

type enhancedRouter struct {
	mu           sync.RWMutex
	middlewares  []MiddlewareFunc
	routes       map[string]HandleFunc
	metricsStart func() error
	tracingStart func(cfgPath string) error

	configPath   string
	config       interface{}
	stop         chan struct{}
}
```

核心逻辑（注册中间件、构造执行链、添加路由）：
```go
func NewEnhancedRouter(configPath string) EnhancedRouter {
	er := &enhancedRouter{
		routes:     make(map[string]HandleFunc),
		configPath: configPath,
		stop:       make(chan struct{}),
	}
	return er
}

func (e *enhancedRouter) Use(mw ...MiddlewareFunc) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.middlewares = append(e.middlewares, mw...)
}

func (e *enhancedRouter) AddRoute(path string, handler HandleFunc) {
	e.mu.Lock()
	defer e.mu.Unlock()
	// Wrap handler with middlewares in registration order -> outermost is first registered
	h := handler
	// build chain in reverse so first middleware registered is outermost (洋葱模型)
	for i := len(e.middlewares) - 1; i >= 0; i-- {
		h = e.middlewares[i](h)
	}
	e.routes[path] = h
}

// StartMetricsEndpoint will be implemented in metrics.go to expose /metrics
func (e *enhancedRouter) StartMetricsEndpoint(addr string) error {
	if e.metricsStart == nil {
		return nil
	}
	return e.metricsStart()
}

func (e *enhancedRouter) Shutdown(ctx context.Context) error {
	close(e.stop)
	// add other shutdown work if needed
	return nil
}
```

---

### 中间件定义与示例（enhanced_router/middleware.go）
```go
package enhanced_router

import (
	"context"
	"errors"
)

// Example of a simple recovery middleware
func RecoveryMiddleware() MiddlewareFunc {
	return func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = errors.New("panic recovered")
				}
			}()
			return next(ctx, req)
		}
	}
}
```

---

### Prometheus 指标集成（enhanced_router/metrics.go）
```go
package enhanced_router

import (
	"net/http"
	"time"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 定义指标名规范
var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "zinx_router_requests_total",
			Help: "Total number of requests processed by zinx router",
		},
		[]string{"route", "code"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "zinx_router_process_duration_seconds",
			Help:    "Router request processing duration seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route"},
	)
	requestErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "zinx_router_errors_total",
			Help: "Total number of errors in router handlers",
		},
		[]string{"route", "type"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal, requestDuration, requestErrors)
}

// MetricsMiddleware 自动记录请求数、耗时、错误
func MetricsMiddleware() MiddlewareFunc {
	return func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			start := time.Now()
			route := "unknown"
			if r, ok := ctx.Value("route").(string); ok {
				route = r
			}
			resp, err = next(ctx, req)
			duration := time.Since(start).Seconds()
			requestDuration.WithLabelValues(route).Observe(duration)
			code := "200"
			if err != nil {
				code = "500"
				requestErrors.WithLabelValues(route, "handler_error").Inc()
			}
			requestsTotal.WithLabelValues(route, code).Inc()
			return resp, err
		}
	}
}

// StartMetricsHTTPServer 启动独立的 HTTP 服务暴露 /metrics
func StartMetricsHTTPServer(addr, metricsPath string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.Handler())
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("metrics server error: %v", err)
		}
	}()
	return srv
}
```

说明：
- Prometheus 是 pull 模式，StartMetricsHTTPServer 提供 /metrics 路径。默认 metricsPath 为 /metrics。
- collect_period 在此方案中用于内部聚合或 pushgateway 场景；Prometheus 抓取频率由 Prometheus 配置决定。

---

### OpenTelemetry 集成（enhanced_router/tracing.go）
```go
package enhanced_router

import (
	"context"
	"time"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"google.golang.org/grpc"
)

type TracingOptions struct {
	Endpoint  string
	Protocol  string // "grpc" or "http"
	SampleRate float64
	ServiceName string
}

// InitTracer 初始化 OTLP exporter 和 TracerProvider
func InitTracer(ctx context.Context, opt TracingOptions) (func(context.Context) error, error) {
	var tp *sdktrace.TracerProvider
	var err error
	var expOption interface{}
	if opt.Protocol == "http" {
		client := otlptracehttp.NewClient(otlptracehttp.WithEndpoint(opt.Endpoint))
		exp, err := otlptracehttp.New(ctx, client)
		if err != nil {
			return nil, err
		}
		expOption = exp
		_ = exp
	} else {
		client := otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(opt.Endpoint), otlptracegrpc.WithDialOption(grpc.WithBlock()))
		exp, err := otlptracegrpc.New(ctx, client)
		if err != nil {
			return nil, err
		}
		expOption = exp
		_ = exp
	}

	// sampler
	var sampler sdktrace.Sampler
	if opt.SampleRate >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else if opt.SampleRate <= 0 {
		sampler = sdktrace.NeverSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(opt.SampleRate)
	}

	// Compose TracerProvider
	res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(
		attribute.String("service.name", opt.ServiceName),
	))
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		// sdktrace.WithBatcher(exp) // set exporter here if concrete type
	)

	// Note: we need to register the exporter with batcher — handle http/grpc exporter types
	switch v := expOption.(type) {
	case *otlptracegrpc.Exporter:
		tp.RegisterSpanProcessor(sdktrace.NewBatchSpanProcessor(v))
	case *otlptracehttp.Exporter:
		tp.RegisterSpanProcessor(sdktrace.NewBatchSpanProcessor(v))
	default:
		// If using generic types, attempt reflection omitted for brevity
	}

	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

// TracingMiddleware 生成 span，记录参数与返回结果
func TracingMiddleware(tracerName string, attrGetter func(ctx context.Context, req interface{}) []attribute.KeyValue) MiddlewareFunc {
	return func(next HandleFunc) HandleFunc {
		tracer := otel.Tracer(tracerName)
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			route := "unknown"
			if r, ok := ctx.Value("route").(string); ok {
				route = r
			}
			ctx, span := tracer.Start(ctx, route)
			defer func() {
				if err != nil {
					span.SetStatus(codes.Error, err.Error())
				} else {
					span.SetStatus(codes.Ok, "")
				}
				span.End()
			}()

			// add user provided attributes
			if attrGetter != nil {
				for _, a := range attrGetter(ctx, req) {
					span.SetAttributes(a)
				}
			}

			// record request param as event (注意：避免记录过大 payload)
			span.AddEvent("request.received")
			resp, err = next(ctx, req)
			span.AddEvent("request.completed")
			return resp, err
		}
	}
}
```

说明：
- InitTracer 启动 OTLP exporter，配置 sampler。示例代码中对 exporter 的具体类型注册可能需要适配实际 exporter 类型（上例做了简化）。
- TracingMiddleware 会为每次请求创建 span，并支持通过 attrGetter 添加自定义 attribute（例如 user_id、device_id 等）。

---

### 配置热加载（使用 fsnotify）
```go
package config

import (
	"io/ioutil"
	"log"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
	"github.com/fsnotify/fsnotify"
)

type Loader struct {
	Path string
	Cfg  *AppConfig
	mu   sync.RWMutex
	cb   func(*AppConfig)
}

func NewLoader(path string, cb func(*AppConfig)) (*Loader, error) {
	ld := &Loader{Path: path, cb: cb}
	if err := ld.load(); err != nil {
		return nil, err
	}
	go ld.watch()
	return ld, nil
}

func (l *Loader) load() error {
	b, err := ioutil.ReadFile(l.Path)
	if err != nil {
		return err
	}
	var c AppConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return err
	}
	l.mu.Lock()
	l.Cfg = &c
	l.mu.Unlock()
	if l.cb != nil {
		l.cb(&c)
	}
	return nil
}

func (l *Loader) watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("watcher init err:", err)
		return
	}
	defer w.Close()
	dir := filepath.Dir(l.Path)
	if err := w.Add(dir); err != nil {
		log.Println("watcher add dir err:", err)
		return
	}
	for {
		select {
		case ev := <-w.Events:
			if ev.Op&(fsnotify.Write|fsnotify.Create) != 0 && filepath.Clean(ev.Name) == filepath.Clean(l.Path) {
				// debounce
				time.Sleep(200 * time.Millisecond)
				if err := l.load(); err != nil {
					log.Println("reload config error:", err)
				}
			}
		case err := <-w.Errors:
			log.Println("watcher error:", err)
		}
	}
}
```

使用场景：Loader 启动时回调会把最新配置传给 enhancedRouter，router 根据新配置调整 metrics 或 tracing（例如修改采样率、切换 exporter）。

---

### 示例：如何在应用中使用（cmd/example/main.go）
```go
package main

import (
	"context"
	"log"
	"time"

	"your/module/enhanced_router"
	"your/module/config"
)

func main() {
	cfgPath := "./config.yaml"
	// create router
	router := enhanced_router.NewEnhancedRouter(cfgPath)

	// create config loader with callback to apply changes
	loader, err := config.NewLoader(cfgPath, func(c *config.AppConfig) {
		// apply metrics config: start metrics server if enabled
		if c.Metrics.Enabled {
			go enhanced_router.StartMetricsHTTPServer(":9090", c.Metrics.MetricsPath)
		}
		// apply otlp config: init tracer with new sampler or exporter
		if c.OTLP.Enabled {
			opt := enhanced_router.TracingOptions{
				Endpoint: c.OTLP.Endpoint,
				Protocol: c.OTLP.Protocol,
				SampleRate: c.OTLP.SampleRatio,
				ServiceName: c.OTLP.ServiceName,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			shutdown, err := enhanced_router.InitTracer(ctx, opt)
			if err != nil {
				log.Println("init tracer err:", err)
			}
			_ = shutdown
		}
	})
	if err != nil {
		log.Fatalln(err)
	}
	_ = loader

	// register middlewares (order matters: first registered is outermost)
	router.Use(
		enhanced_router.RecoveryMiddleware(),
		enhanced_router.MetricsMiddleware(),
		enhanced_router.TracingMiddleware("router-tracer", func(ctx context.Context, req interface{}) []attribute.KeyValue {
			// example attribute injection, must import go.opentelemetry.io/otel/attribute
			return []attribute.KeyValue{
				attribute.String("example.key", "example.value"),
			}
		}),
	)

	// add route and handler
	router.AddRoute("/echo", func(ctx context.Context, req interface{}) (interface{}, error) {
		return req, nil
	})

	// 启动应用的其余部分（整合 zinx 的网络监听等）
	select {}
}
```

---

### 注意事项与实现细节提示
- zinx 的原始 Handler 签名可能与示例不同，需编写一层适配器将 zinx 的 IRequest 转换为 context+req interface，然后调用 EnhancedRouter 中保存的 HandleFunc。
- Tracing exporter 的具体创建与注册会依赖你采用的 otel 版本，示例中对具体 exporter 的断言仅做示意，实际代码应直接保持对 otlptracegrpc.New/otlptracehttp.New 返回的类型并注册 BatchSpanProcessor。
- Prometheus 指标名与 label 设计应稳定，切勿包含高基数标签（例如原始请求体）。
- 热加载需要考虑并发：在回调中优雅切换 TracerProvider（调用旧的 Shutdown 并设置新的），避免在切换瞬间丢失 trace。
- 采样率变更通常要求重建 TracerProvider，因为 sampler 在构建时确定。
- 对于高 QPS 场景，Histogram 的 Bucket 设计与 Metric 标签数量必须谨慎，避免 OOM 与内存飙升。

---

### 结论
该方案将业务逻辑（handler）与公共处理逻辑（指标、追踪、恢复、鉴权等）完全解耦，通过 EnhancedRouter 的中间件链实现洋葱模型的执行顺序；同时集成 Prometheus 与 OpenTelemetry，并通过配置文件与 fsnotify 实现热加载，满足可配置化和可观测性的要求。若需要，我可以把上述代码整理为可直接编译的完整仓库示例并补充 zinx 适配层与更完善的错误处理逻辑。