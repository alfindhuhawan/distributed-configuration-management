package workersdk

import (
	"context"
	"distributed-configuration-management/shared/model/errorenum"
	"distributed-configuration-management/shared/model/request"
	"distributed-configuration-management/shared/model/response"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

type Worker interface {
	UpdateConfig(ctx context.Context, cfgWorker request.ConfigWorker) error
	CurrentConfig(ctx context.Context) (request.ConfigWorker, bool)
	Hit(ctx context.Context) (*response.HitResponse, error)
}

type workerImpl struct {
	mu        sync.RWMutex
	cfgWorker request.ConfigWorker
	hasCfg    bool
	client    *http.Client
}

type Options struct {
	Timeout time.Duration
}

func New(opts Options) Worker {
	if opts.Timeout == 0 {
		opts.Timeout = 10 * time.Second
	}

	return &workerImpl{
		client: &http.Client{Timeout: opts.Timeout},
	}
}

func (w *workerImpl) UpdateConfig(ctx context.Context, cfgWorker request.ConfigWorker) error {
	if cfgWorker.URL == "" {
		return errorenum.ErrorURLWorkerRequired
	}

	w.mu.Lock()
	w.cfgWorker = cfgWorker
	w.hasCfg = true
	w.mu.Unlock()

	return nil
}

func (w *workerImpl) CurrentConfig(ctx context.Context) (request.ConfigWorker, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.cfgWorker, w.hasCfg
}

func (w *workerImpl) Hit(ctx context.Context) (*response.HitResponse, error) {
	w.mu.RLock()
	cfg := w.cfgWorker
	hasCfg := w.hasCfg
	w.mu.RUnlock()

	if !hasCfg || cfg.URL == "" {
		return nil, errorenum.ErrorWorker
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var externalResponse any
	json.Unmarshal(body, &externalResponse)
	return &response.HitResponse{
		URL:              cfg.URL,
		StatusCode:       resp.StatusCode,
		ExternalResponse: externalResponse,
	}, nil
}
