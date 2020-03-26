package alertmanager_webhook

type Message struct {
	Alerts []*Alert
}

type Alert struct {
	Status      string            `json:"status"`
	Url         string            `json:"generatorURL"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}
