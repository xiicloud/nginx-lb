package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	confdBin = "/bin/confd"
	nginxBin = "/usr/sbin/nginx"
	pidFile  = "/tmp/pm.pid"
)

var (
	confdCmd *exec.Cmd
	nginxCmd *exec.Cmd
)

func startCmds() {
	genConfig()
	if err := confdOnce(); err != nil {
		log.Fatalf("confdOnce: %v", err)
	}

	if err := startNginx(); err != nil {
		log.Fatalf("startNginx: %v", err)
	}

	if err := startConfd(); err != nil {
		log.Fatalf("startConfd: %v", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for sig := range ch {
			switch sig {
			case syscall.SIGUSR2:
				if err := restartConfd(); err != nil {
					log.Printf("restartConfd: %v\n", err)
					return
				}
			case syscall.SIGTERM, syscall.SIGINT:
				confdCmd.Process.Signal(syscall.SIGTERM)
				nginxCmd.Process.Signal(syscall.SIGTERM)
				confdCmd.Wait()
				nginxCmd.Wait()
				return
			case syscall.SIGHUP:
				reloadNginx()
			}
		}
	}()

	go func() {
		if err := nginxCmd.Wait(); err != nil {
			log.Printf("nginxCmd.Wait: %v\n", err)
		}
		ch <- syscall.SIGTERM
		wg.Done()
	}()
	wg.Wait()
}

func getEtcdUrl() string {
	etcdPort := getEnvWithDefault("ETCD_PORT", "2379")
	etcdHost := getEnvWithDefault("HOST_IP", "172.17.0.1")
	return fmt.Sprintf("%s:%s", etcdHost, etcdPort)
}

func getEnvWithDefault(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaultValue
	}
	return v
}

// nginx -g "daemon off;"
func startNginx() error {
	nginxCmd = exec.Command(nginxBin, "-g", "daemon off;")
	op, ep, err := getPipes(nginxCmd)
	if err != nil {
		return err
	}
	go copyStreams(op, ep)
	return nginxCmd.Start()
}

func reloadNginx() {
	nginxCmd.Process.Signal(syscall.SIGHUP)
}

// confd -watch -backend etcd -node $ETCD
func startConfd() error {
	confdCmd = exec.Command(confdBin, "-watch", "-backend", "etcd", "-node", getEtcdUrl())
	op, ep, err := getPipes(confdCmd)
	if err != nil {
		return err
	}
	go copyStreams(op, ep)
	return confdCmd.Start()
}

// confd -onetime -node ${ETCD}
func confdOnce() error {
	etcd := getEtcdUrl()
	for i := 0; i < 120; i++ {
		cmd := exec.Command(confdBin, "-onetime", "-backend", "etcd", "-node", etcd)
		out, err := cmd.CombinedOutput()
		if err == nil {
			return nil
		}
		fmt.Println("Waiting for confd to generate the initial nginx config: ", string(out))
		time.Sleep(time.Second)
	}

	return fmt.Errorf("failed to generate initial config for nginx")
}

// restart confd process
func restartConfd() error {
	confdCmd.Process.Signal(syscall.SIGTERM)
	confdCmd.Wait()
	return startConfd()
}

func getPipes(cmd *exec.Cmd) (op io.ReadCloser, ep io.ReadCloser, err error) {
	ep, err = cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	op, err = cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	return op, ep, nil
}

func copyStreams(op, ep io.ReadCloser) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	cp := func(dst io.Writer, src io.ReadCloser) {
		io.Copy(dst, src)
		src.Close()
		wg.Done()
	}
	go cp(os.Stdout, op)
	go cp(os.Stderr, ep)
	wg.Wait()
}
