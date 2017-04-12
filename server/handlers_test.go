package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/stretchr/testify/assert"
)

type mockOptions struct {
	publishResult  error
	consumeResults [][]byte
}

func (s *mockStream) Publish(body []byte, opts interface{}) error {
	return opts.(*mockOptions).publishResult
}

func (s *mockStream) Consume(opts interface{}) ([][]byte, error) {
	if opts.(*mockOptions).consumeResults == nil {
		return nil, errors.New("Consume failed!")
	}

	return opts.(*mockOptions).consumeResults, nil
}

func TestReceiveMetricsStoresInStream(t *testing.T) {
	ctx := NewContext(context.Background(), &mockStream{})
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/send", strings.NewReader(""))

	h := ReceiveMetrics(&mockOptions{})
	h(ctx, wrt, req)

	fmt.Println(wrt.Body.String())
	assert.Equal(t, wrt.Code, 200)
}

func TestReceiveMetricsFailsWithInvalidContext(t *testing.T) {
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/send", strings.NewReader(""))

	h := ReceiveMetrics(&mockOptions{})
	h(context.Background(), wrt, req)

	fmt.Println(wrt.Body.String())
	assert.Equal(t, wrt.Code, 500)
}

func TestReceiveMetricsFailsWithStreamError(t *testing.T) {
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/send", strings.NewReader(""))
	ctx := NewContext(context.Background(), &mockStream{})

	h := ReceiveMetrics(&mockOptions{publishResult: errors.New("Publish failed!")})
	h(ctx, wrt, req)

	fmt.Println(wrt.Body.String())
	assert.Equal(t, wrt.Code, 503)
}

// =======

func TestSendMetricsWithInvalidContext(t *testing.T) {
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/metrics", strings.NewReader(""))

	h := SendMetrics(&mockOptions{})
	h(context.Background(), wrt, req)

	fmt.Println(wrt.Body.String())
	assert.Equal(t, wrt.Code, 500)
}

func TestSendMetricsMessageRetreivalFails(t *testing.T) {
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/metrics", strings.NewReader(""))
	ctx := NewContext(context.Background(), &mockStream{})

	h := SendMetrics(&mockOptions{})
	h(ctx, wrt, req)

	fmt.Println(wrt.Body.String())
	assert.Equal(t, wrt.Code, 503)
}

func TestSendMetricsSuccess(t *testing.T) {
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/metrics", strings.NewReader(""))
	ctx := NewContext(context.Background(), &mockStream{})

	ts := remote.WriteRequest{
		Timeseries: []*remote.TimeSeries{
			&remote.TimeSeries{},
		},
	}
	result, _ := proto.Marshal(&ts)
	h := SendMetrics(&mockOptions{
		consumeResults: [][]byte{result},
	})
	h(ctx, wrt, req)

	assert.Equal(t, wrt.Code, 200)
}

func TestSendMetricsConstructsValidResponse(t *testing.T) {
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/metrics", strings.NewReader(""))
	ctx := NewContext(context.Background(), &mockStream{})

	ts := remote.WriteRequest{
		Timeseries: []*remote.TimeSeries{
			&remote.TimeSeries{
				Labels: []*remote.LabelPair{
					{Name: "test", Value: "test"},
				},
				Samples: []*remote.Sample{
					{Value: 10, TimestampMs: 1111111111},
				},
			},
		},
	}
	result, _ := proto.Marshal(&ts)
	h := SendMetrics(&mockOptions{
		consumeResults: [][]byte{result},
	})
	h(ctx, wrt, req)

	metrics, _ := ioutil.ReadAll(wrt.Body)

	assert.Equal(t, wrt.Code, 200)
	assert.Equal(t, string(metrics), "{test=\"test\"} 10.000000 1111111111\n")
}

func TestSendMetricsMessageFormatInvalid(t *testing.T) {
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/metrics", strings.NewReader(""))
	ctx := NewContext(context.Background(), &mockStream{})

	h := SendMetrics(&mockOptions{consumeResults: [][]byte{[]byte("Hello :)")}})
	h(ctx, wrt, req)

	fmt.Println(wrt.Body.String())
	assert.Equal(t, wrt.Code, 500)
}
