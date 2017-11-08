package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/prompb"
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
	req := httptest.NewRequest("POST", "/expose", nil)

	h := ReceiveMetrics(&mockOptions{})
	h(ctx, wrt, req)

	fmt.Println(wrt.Body.String())
	assert.Equal(t, wrt.Code, 200)
}

func TestReceiveMetricsFailsWithInvalidContext(t *testing.T) {
	wrt := httptest.NewRecorder()
	//req := httptest.NewRequest("POST", "/send", strings.NewReader(""))
	matcher1, err := labels.NewMatcher(labels.MatchEqual, "__name__", "test_metric1")
	if err != nil {
		t.Fatal(err)
	}
	matcher2, err := labels.NewMatcher(labels.MatchEqual, "d", "e")
	if err != nil {
		t.Fatal(err)
	}
	query, err := remote.ToQuery(0, 1, []*labels.Matcher{matcher1, matcher2})
	if err != nil {
		t.Fatal(err)
	}
	req := &prompb.ReadRequest{Queries: []*prompb.Query{query}}
	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	compressed := snappy.Encode(nil, data)
	req := httptest.NewRequest("POST", "/expose", bytes.NewBuffer(compressed))

	h := ReceiveMetrics(&mockOptions{})
	h(context.Background(), wrt, req)

	fmt.Println(wrt.Body.String())
	assert.Equal(t, wrt.Code, 500)
}

func TestReceiveMetricsFailsWithStreamError(t *testing.T) {
	wrt := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/expose", strings.NewReader(""))
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

	ts := prompb.WriteRequest{
		Timeseries: []*prompb.TimeSeries{
			&prompb.TimeSeries{},
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

	ts := prompb.WriteRequest{
		Timeseries: []*prompb.TimeSeries{
			&prompb.TimeSeries{
				Labels: []*prompb.Label{
					{Name: "test", Value: "test"},
				},
				Samples: []*prompb.Sample{
					{Value: 10, Timestamp: 1111111111},
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
