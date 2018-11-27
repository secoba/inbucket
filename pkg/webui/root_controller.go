package webui

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jhillyerd/inbucket/pkg/config"
	"github.com/jhillyerd/inbucket/pkg/server/web"
)

// RootGreeting serves the Inbucket greeting.
func RootGreeting(w http.ResponseWriter, req *http.Request, ctx *web.Context) (err error) {
	greeting, err := ioutil.ReadFile(ctx.RootConfig.Web.GreetingFile)
	if err != nil {
		return fmt.Errorf("Failed to load greeting: %v", err)
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(greeting)
	return err
}

// RootMonitor serves the Inbucket monitor page
func RootMonitor(w http.ResponseWriter, req *http.Request, ctx *web.Context) (err error) {
	if !ctx.RootConfig.Web.MonitorVisible {
		ctx.Session.AddFlash("Monitor is disabled in configuration", "errors")
		_ = ctx.Session.Save(req, w)
		http.Redirect(w, req, web.Reverse("RootIndex"), http.StatusSeeOther)
		return nil
	}
	// Get flash messages, save session
	errorFlash := ctx.Session.Flashes("errors")
	if err = ctx.Session.Save(req, w); err != nil {
		return err
	}
	// Render template
	return web.RenderTemplate("root/monitor.html", w, map[string]interface{}{
		"ctx":        ctx,
		"errorFlash": errorFlash,
	})
}

// RootMonitorMailbox serves the Inbucket monitor page for a particular mailbox
func RootMonitorMailbox(w http.ResponseWriter, req *http.Request, ctx *web.Context) (err error) {
	if !ctx.RootConfig.Web.MonitorVisible {
		ctx.Session.AddFlash("Monitor is disabled in configuration", "errors")
		_ = ctx.Session.Save(req, w)
		http.Redirect(w, req, web.Reverse("RootIndex"), http.StatusSeeOther)
		return nil
	}
	name, err := ctx.Manager.MailboxForAddress(ctx.Vars["name"])
	if err != nil {
		ctx.Session.AddFlash(err.Error(), "errors")
		_ = ctx.Session.Save(req, w)
		http.Redirect(w, req, web.Reverse("RootIndex"), http.StatusSeeOther)
		return nil
	}
	// Get flash messages, save session
	errorFlash := ctx.Session.Flashes("errors")
	if err = ctx.Session.Save(req, w); err != nil {
		return err
	}
	// Render template
	return web.RenderTemplate("root/monitor.html", w, map[string]interface{}{
		"ctx":        ctx,
		"errorFlash": errorFlash,
		"name":       name,
	})
}

type jsonStatus struct {
	Version      string         `json:"version"`
	BuildDate    string         `json:"build-date"`
	POP3Listener string         `json:"pop3-listener"`
	WebListener  string         `json:"web-listener"`
	SMTPConfig   jsonSMTPConfig `json:"smtp-config"`
}

type jsonSMTPConfig struct {
	Addr           string   `json:"addr"`
	DefaultAccept  bool     `json:"default-accept"`
	AcceptDomains  []string `json:"accept-domains"`
	RejectDomains  []string `json:"reject-domains"`
	DefaultStore   bool     `json:"default-store"`
	StoreDomains   []string `json:"store-domains"`
	DiscardDomains []string `json:"discard-domains"`
}

// RootStatus serves the Inbucket status page
func RootStatus(w http.ResponseWriter, req *http.Request, ctx *web.Context) (err error) {
	smtpConfig := ctx.RootConfig.SMTP

	return web.RenderJSON(w,
		&jsonStatus{
			Version:      config.Version,
			BuildDate:    config.BuildDate,
			POP3Listener: ctx.RootConfig.POP3.Addr,
			WebListener:  ctx.RootConfig.Web.Addr,
			SMTPConfig: jsonSMTPConfig{
				Addr:           smtpConfig.Addr,
				DefaultAccept:  smtpConfig.DefaultAccept,
				AcceptDomains:  smtpConfig.AcceptDomains,
				RejectDomains:  smtpConfig.RejectDomains,
				DefaultStore:   smtpConfig.DefaultStore,
				StoreDomains:   smtpConfig.StoreDomains,
				DiscardDomains: smtpConfig.DiscardDomains,
			},
		})
	// return web.RenderTemplate("root/status.html", w, map[string]interface{}{
	// 	"ctx":           ctx,
	// 	"errorFlash":    errorFlash,
	// 	"version":       config.Version,
	// 	"buildDate":     config.BuildDate,
	// 	"smtpListener":  ctx.RootConfig.SMTP.Addr,
	// 	"pop3Listener":  ctx.RootConfig.POP3.Addr,
	// 	"webListener":   ctx.RootConfig.Web.Addr,
	// 	"smtpConfig":    ctx.RootConfig.SMTP,
	// 	"storageConfig": ctx.RootConfig.Storage,
	// })
}
