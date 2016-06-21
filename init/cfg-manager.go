package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
)

const (
	cfgPath        = "/etc/csphere-services.json"
	confdConfigDir = "/etc/confd/conf.d"
	confdTplDir    = "/etc/confd/templates"
	tplPath        = "/tpls"
	sslCertDir     = "/etc/nginx/ssl"
	nginxConfDir   = "/etc/nginx"
)

type Route struct {
	App         string `json:"app"`
	Service     string `json:"service"`
	Port        int    `json:"port"`
	BackendPath string `json:"backend_path"`
	Backup      *Route `json:"backup"`
	// Any config directives passed to the `location` config.
	Opaque string `json:"opaque"`
}

func (b *Route) fixup() error {
	if b.Port == 0 {
		b.Port = 80
	}

	if b.BackendPath == "" {
		b.BackendPath = "/"
	} else {
		b.BackendPath = filepath.Clean(b.BackendPath) + "/"
	}

	if b.Backup != nil {
		return b.Backup.fixup()
	}
	return nil
}

type Server struct {
	DomainName        string `json:"domain_name"`
	FrontendPort      int    `json:"frontend_port"`
	SslCertificate    string `json:"ssl_certificate"`
	SslCertificateKey string `json:"ssl_certificate_key"`
	SslPort           int    `json:"ssl_port"`
	SslCertPath       string `json:"-"`
	SslKeyPath        string `json:"-"`
	EnableSsl         bool   `json:"-"`

	// The key is the URL that will be exposed to the frontend user.
	Routes map[string]*Route `json:"Routes"`

	// Any config directives passed to the `server` config.
	Opaque string `json:"opaque"`
}

type Config struct {
	Version string    `json:"version"`
	Servers []*Server `json:"servers"`
}

func genConfig() {
	os.MkdirAll(confdTplDir, 0755)
	os.MkdirAll(confdConfigDir, 0755)
	os.MkdirAll(sslCertDir, 0755)
	config, err := parseConfig()
	if err != nil {
		log.Fatalf("Failed to parse config file %s: %v", cfgPath, err)
	}

	if err := genConfdToml(config.Servers); err != nil {
		log.Fatalf("Failed to generate confd toml: %v", err)
	}

	if err := genNginxTpl(config.Servers); err != nil {
		log.Fatalf("Failed to generate nginx.tpl: %v", err)
	}
}

func reload() {
	proc, err := os.FindProcess(1)
	if err != nil {
		log.Fatalf("os.FindProcess: %v", err)
	}
	if err := proc.Signal(syscall.SIGUSR2); err != nil {
		log.Fatalf("Failed to restart confd process: %v", err)
	}
}

func parseConfig() (*Config, error) {
	cfg := &Config{}
	fp, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	err = json.NewDecoder(fp).Decode(cfg)
	if err != nil {
		return nil, err
	}

	for i, s := range cfg.Servers {
		if len(s.Routes) == 0 {
			log.Printf("Routes must contain at least 1 items")
			continue
		}

		if s.FrontendPort == 0 {
			s.FrontendPort = 80
		}
		if s.DomainName == "" {
			s.DomainName = fmt.Sprintf("%d.example.com", i)
		}

		if s.SslCertificate != "" && s.SslCertificateKey != "" {
			s.EnableSsl = true
			s.SslCertPath = filepath.Join(sslCertDir, s.DomainName+".pem")
			s.SslKeyPath = filepath.Join(sslCertDir, s.DomainName+".key")
			if err := ioutil.WriteFile(s.SslCertPath, []byte(s.SslCertificate), 0400); err != nil {
				return nil, err
			}
			if err := ioutil.WriteFile(s.SslKeyPath, []byte(s.SslCertificateKey), 0400); err != nil {
				return nil, err
			}
		}

		if s.SslPort == 0 {
			s.SslPort = 443
		}

		for _, b := range s.Routes {
			b.fixup()
		}
	}

	return cfg, nil
}

func genConfdToml(servers []*Server) error {
	src := filepath.Join(tplPath, "confd.tpl")
	dst := filepath.Join(confdConfigDir, "nginx.toml")
	tpl := template.Must(getTplObj("confd.tpl").ParseFiles(src))
	keys := make([]string, len(servers))
	for _, s := range servers {
		for _, b := range s.Routes {
			keys = append(keys, fmt.Sprintf("/%s-%s/ips", b.App, b.Service))
		}
	}

	fp, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	v, err := json.Marshal(keys)
	vars := map[string]string{
		"Keys":     string(v),
		"DestPath": getConfdNginxConfDestPath(),
	}

	return tpl.Execute(fp, vars)
}

func genNginxTpl(servers []*Server) error {
	src := filepath.Join(tplPath, getNginxTpl())
	srcUpstreams := filepath.Join(tplPath, "upstreams.tpl")
	dst := filepath.Join(confdTplDir, "nginx.tpl")
	tpl := template.Must(getTplObj("nginx.tpl").ParseFiles(src, srcUpstreams))

	fp, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	return tpl.Execute(fp, servers)
}

func getLbKey(b *Route) string {
	return fmt.Sprintf("/%s-%s/ips", b.App, b.Service)
}

// split is a version of strings.Split that can be piped
func split(sep, s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

func upstreamName(r *Route) string {
	names := []string{fmt.Sprintf("%s-%s", r.App, r.Service)}
	if r.Backup != nil {
		names = append(names, fmt.Sprintf("%s-%s", r.Backup.App, r.Backup.Service))
	}
	return strings.Join(names, "-")
}

func normalizeURI(uri string) string {
	if uri == "" {
		return "/"
	}
	return strings.TrimRight(filepath.Clean(uri), "/") + "/"
}

var tplFuncs = map[string]interface{}{
	"split":        split,
	"getLbKey":     getLbKey,
	"upstreamName": upstreamName,
	"normalizeURI": normalizeURI,
}

func getTplObj(name string) *template.Template {
	tpl := template.New(name)
	tpl.Delims("(", ")")
	tpl.Funcs(tplFuncs)
	return tpl
}

func getNginxTpl() string {
	return "nginx.tpl"
}

func getConfdNginxConfDestPath() string {
	return filepath.Join(nginxConfDir, "nginx.conf")
}
