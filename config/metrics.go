package config

import (
	"github.com/leandro-lugaresi/hub"
	"github.com/prometheus/client_golang/prometheus"
)

const metricsPrefix = "telepush_"

const (
	labelTotalMessagesOrigin  = "origin"
	labelTotalMessagesType    = "type"
	labelTotalRequestsSuccess = "success"
)

var (
	counterTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests to this bot",
		},
		[]string{labelTotalRequestsSuccess},
	)
	counterTotalMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_total",
			Help: "Total number of messages delivered",
		},
		[]string{labelTotalMessagesOrigin, labelTotalMessagesType},
	)
)

func init() {
	registerer := prometheus.WrapRegistererWithPrefix(metricsPrefix, prometheus.DefaultRegisterer)
	registerer.MustRegister(counterTotalMessages)
	registerer.MustRegister(counterTotalRequests)

	h := GetHub()
	sub := h.NonBlockingSubscribe(0, AllEvents()...)

	go func(s hub.Subscription) {
		for event := range s.Receiver {
			if event.Name == EventOnRequestSuccessful {
				counterTotalRequests.With(prometheus.Labels{
					labelTotalRequestsSuccess: "1",
				}).Inc()
				continue
			}
			if event.Name == EventOnRequestFailed {
				counterTotalRequests.With(prometheus.Labels{
					labelTotalRequestsSuccess: "0",
				}).Inc()
				continue
			}
			if event.Name == EventOnMessageDelivered {
				counterTotalMessages.With(prometheus.Labels{
					labelTotalMessagesOrigin: event.Fields[FieldMessageOrigin].(string),
					labelTotalMessagesType:   event.Fields[FieldMessageType].(string),
				}).Inc()
				continue
			}
		}
	}(sub)
}
