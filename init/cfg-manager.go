package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"text/template"
)

const (
	cfgPath        = "/etc/csphere-services.json"
	confdConfigDir = "/etc/confd/conf.d"
	confdTplDir    = "/etc/confd/templates"
	tplPath        = "/tpls"
	sslCertDir     = "/etc/nginx/ssl"
)

type Service struct {
	Domain            string `json:"domain_name"`
	App               string `json:"app"`
	Service           string `json:"service"`
	BackendPort       int    `json:"backend_port"`
	FrontendPort      int    `json:"frontend_port"`
	BackendRootPath   string `json:"backend_root_path"`
	SslCertificate    string `json:"ssl_certificate"`
	SslCertificateKey string `json:"ssl_certificate_key"`
	SslCertPath       string `json:"-"`
	SslKeyPath        string `json:"-"`
	EnableSsl         bool   `json:"-"`
	SslPort           int    `json:"ssl_port"`
}

func genConfig() {
	os.MkdirAll(confdTplDir, 0755)
	os.MkdirAll(confdConfigDir, 0755)
	os.MkdirAll(sslCertDir, 0755)
	services, err := parseConfig()
	if err != nil {
		log.Fatalf("Failed to parse config file %s: %v", cfgPath, err)
	}

	if err := genConfdToml(services); err != nil {
		log.Fatalf("Failed to generate confd toml: %v", err)
	}

	if err := genNginxTpl(services); err != nil {
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

func parseConfig() ([]*Service, error) {
	services := []*Service{}
	fp, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	err = json.NewDecoder(fp).Decode(&services)
	if err != nil {
		return nil, err
	}

	for _, s := range services {
		if s.BackendPort == 0 {
			s.BackendPort = 80
		}
		if s.FrontendPort == 0 {
			s.FrontendPort = 80
		}
		if s.Domain == "" {
			s.Domain = fmt.Sprintf("%s-%s", s.App, s.Service)
		}
		if s.BackendRootPath == "" {
			s.BackendRootPath = "/"
		}
		s.BackendRootPath = filepath.Clean(s.BackendRootPath) + "/"

		if s.SslCertificate != "" && s.SslCertificateKey != "" {
			s.EnableSsl = true
			s.SslCertPath = filepath.Join(sslCertDir, s.Domain+".pem")
			s.SslKeyPath = filepath.Join(sslCertDir, s.Domain+".key")
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
	}

	return services, nil
}

func genConfdToml(services []*Service) error {
	src := filepath.Join(tplPath, "confd.tpl")
	dst := filepath.Join(confdConfigDir, "nginx.toml")
	tpl := template.Must(getTplObj("confd.tpl").ParseFiles(src))
	keys := make([]string, len(services))
	for i, v := range services {
		keys[i] = fmt.Sprintf("/%s-%s/ips", v.App, v.Service)
	}

	fp, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	v, err := json.Marshal(keys)
	return tpl.Execute(fp, string(v))
}

func genNginxTpl(services []*Service) error {
	src := filepath.Join(tplPath, "nginx.tpl")
	dst := filepath.Join(confdTplDir, "nginx.tpl")
	tpl := template.Must(getTplObj("nginx.tpl").ParseFiles(src))

	fp, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	return tpl.Execute(fp, services)
}

func getTplObj(name string) *template.Template {
	tpl := template.New(name)
	tpl.Delims("(", ")")
	return tpl
}