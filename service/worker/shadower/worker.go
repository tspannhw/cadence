// Copyright (c) 2017-2021 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package shadower

import (
	"context"

	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"

	"github.com/uber/cadence/.gen/go/shadower"
	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/cache"
)

type (
	// BootstrapParams contains the set of params needed to bootstrap workflow shadower worker
	BootstrapParams struct {
		ServiceClient workflowserviceclient.Interface
		DomainCache   cache.DomainCache
	}

	// Worker is for executing decision task generated by shadowing workflows
	Worker struct {
		decisionWorker worker.Worker
		domainCache    cache.DomainCache
	}

	contextKey string
)

const (
	workerContextKey contextKey = "shadower-worker-context"
)

func New(params *BootstrapParams) *Worker {
	w := &Worker{
		domainCache: params.DomainCache,
	}
	ctx := context.WithValue(context.Background(), workerContextKey, w)
	w.decisionWorker = worker.New(
		params.ServiceClient,
		common.ShadowerLocalDomainName,
		shadower.TaskList,
		worker.Options{
			BackgroundActivityContext: ctx,
		},
	)
	register(w.decisionWorker)
	return w
}

// Start starts the decision worker
func (w *Worker) Start() error {
	if err := w.decisionWorker.Start(); err != nil {
		w.decisionWorker.Stop()
		return err
	}
	return nil
}
