package interceptors

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonades/logger"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"os"
)

// interceptor basically used for intercepting entire traffic of the web app
// we are going to push data points from here to our statsd deamons
type HttpInterceptor struct {
	mux         *bone.Mux
	gApiCounter metrics.Counter
}

func NewInterceptor(mux *bone.Mux) *HttpInterceptor {
	i := HttpInterceptor{}
	i.mux = mux
	i.gApiCounter = metrics.NewCounter()
	i.gApiCounter = metrics.GetOrRegisterCounter(os.Getenv("ENV")+".api.served", nil)
	return &i
}

// for following the Handler interface (http://golang.org/pkg/net/http/#Handler)
func (i HttpInterceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.gApiCounter.Inc(1)
	apiTimer := metrics.GetOrRegisterTimer(os.Getenv("ENV")+".api.time", nil)
	// before
	logger.Get().Debug("Start => " + r.URL.String())
	apiTimer.Time(func() { i.mux.ServeHTTP(w, r) })
	// after
	logger.Get().Debug("End => " + r.URL.String())
}
