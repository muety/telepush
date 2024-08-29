package config

import (
	"github.com/leandro-lugaresi/hub"
	"github.com/muety/telepush/services"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const metricsPrefix = "telepush_"

const (
	labelTotalMessagesOrigin  = "origin"
	labelTotalMessagesType    = "type"
	labelTotalRequestsSuccess = "success"
)

const usersActiveTimeout = 24 * time.Hour

var (
	counterTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricsPrefix + "requests_total",
			Help: "Total number of requests to this bot",
		},
		[]string{labelTotalRequestsSuccess},
	)
	counterTotalMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricsPrefix + "messages_total",
			Help: "Total number of messages delivered",
		},
		[]string{labelTotalMessagesOrigin, labelTotalMessagesType},
	)
	gaugeTotalUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: metricsPrefix + "users_total",
			Help: "Total number of registered users",
		},
	)
	gaugeActiveUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: metricsPrefix + "users_active",
			Help: "Total number of users who received notifications within past 24 hours",
		},
	)
)

var userService *services.UserService
var activeUsers map[int64]time.Time

func countActive() (count int) {
	for _, t := range activeUsers {
		if time.Since(t) < usersActiveTimeout {
			count++
		}
	}
	return count
}

func init() {
	prometheus.MustRegister(counterTotalMessages)
	prometheus.MustRegister(counterTotalRequests)
	prometheus.MustRegister(gaugeTotalUsers)
	prometheus.MustRegister(gaugeActiveUsers)

	userService = services.NewUserService(GetStore())
	activeUsers = make(map[int64]time.Time)

	h := GetHub()
	sub := h.NonBlockingSubscribe(0, AllEvents()...)

	// init active users
	for _, i := range userService.GetUsers() {
		activeUsers[i] = time.Time{}
	}
	gaugeTotalUsers.Set(float64(len(activeUsers)))
	gaugeActiveUsers.Set(float64(countActive()))

	// event hub subscription
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

				recipientUsers := userService.GetUsersByRecipient(event.Fields[FieldMessageRecipient].(string))
				for _, u := range recipientUsers {
					activeUsers[u] = time.Now()
				}
				continue
			}
			if event.Name == EventOnTokenIssued {
				activeUsers[event.Fields[FieldTokenUser].(int64)] = time.Now()
				gaugeTotalUsers.Set(float64(len(userService.GetUsers())))
				gaugeActiveUsers.Set(float64(countActive()))
			}
		}
	}(sub)

	// background jobs
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			gaugeActiveUsers.Set(float64(countActive()))
		}
	}()
}
