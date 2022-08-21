package alertmanager

import "time"

// See https://grafana.com/docs/grafana/latest/alerting/contact-points/notifiers/webhook-notifier/

type Message struct {
	Receiver    string   `json:"receiver"`
	Status      string   `json:"status"`
	OrgId       int      `json:"orgId"`
	Alerts      []*Alert `json:"alerts"`
	ExternalURL string   `json:"externalURL"`
	Version     string   `json:"version"`
	GroupKey    string   `json:"groupKey"`
	Title       string   `json:"title"`
	State       string   `json:"state"`
	Message     string   `json:"message"`
}

type Alert struct {
	Status       string            `json:"status"`
	Url          string            `json:"url"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	SilenceURL   string            `json:"silenceURL"`
	DashboardURL string            `json:"dashboardURL"`
	PanelURL     string            `json:"panelURL"`
	Fingerprint  string            `json:"fingerprint"`
	ValueString  string            `json:"valueString"`
}
