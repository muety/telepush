package inlets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/muety/telepush/config"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/resolvers"
	"github.com/muety/telepush/util"
	"io"
	"net/http"
	"strings"
	"text/template"
)

var templateFuncs = template.FuncMap{
	"escapemd": util.EscapeMarkdown,
	"div":      util.Div,
}

type InletConfig struct {
	Name        string            `yaml:"name,omitempty"`
	ContentType string            `yaml:"content_type,omitempty"`
	Template    string            `yaml:"template,omitempty"`
	HeaderVars  map[string]string `yaml:"header_vars"`
}

type ConfigInlet struct {
	Config *InletConfig
	tpl    *template.Template
}

func NewConfigInlet(inletConfig *InletConfig) (*ConfigInlet, error) {
	inletConfig.Name = strings.ToLower(inletConfig.Name)

	tpl, err := template.New(inletConfig.Name).Funcs(templateFuncs).Parse(inletConfig.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template for %s: %v\n", inletConfig.Name, err)
	}

	return &ConfigInlet{
		Config: inletConfig,
		tpl:    tpl,
	}, nil
}

func (c *ConfigInlet) SupportedMethods() []string {
	return []string{http.MethodPost} // config-based inlets only support POST atm.
}

type templateVars struct {
	Message interface{}
	Vars    map[string]string
}

func (c *ConfigInlet) Name() string {
	return c.Config.Name
}

func (c *ConfigInlet) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		responseText, err := c.getTextResponse(payload, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		message := &model.DefaultMessage{
			Text:   responseText,
			Type:   resolvers.TextType,
			Origin: c.getOrigin(payload, r),
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, config.KeyMessage, message)
		ctx = context.WithValue(ctx, config.KeyParams, &model.MessageParams{DisableLinkPreviews: true})

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *ConfigInlet) getTextResponse(bodyBytes []byte, r *http.Request) (string, error) {
	var payload interface{} = string(bodyBytes)
	if c.Config.ContentType == "application/json" {
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			return "", err
		}
	}

	var buf bytes.Buffer
	if err := c.tpl.Execute(&buf, templateVars{
		Message: payload,
		Vars:    c.getHeaderVars(r),
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *ConfigInlet) getOrigin(bodyBytes []byte, r *http.Request) string {
	if r.Header.Get("content-type") == "application/json" {
		var payload map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			return model.DefaultOrigin
		}
		if v, ok := payload["origin"]; ok {
			switch v.(type) {
			case string:
				return v.(string)
			}
		}
		return model.DefaultOrigin
	}
	return model.DefaultOrigin
}

func (c *ConfigInlet) getHeaderVars(r *http.Request) map[string]string {
	vars := make(map[string]string)
	for k, v := range c.Config.HeaderVars {
		vars[k] = r.Header.Get(v)
	}
	return vars
}
