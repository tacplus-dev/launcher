// Copyright (c) 2026 Jody Scott
// SPDX-License-Identifier: Apache-2.0

// launcher — process supervisor for distroless containers.
//
// Forks and supervises a child process defined by LAUNCH_CMD. Forwards
// signals so the child receives SIGTERM/SIGINT cleanly, and propagates the
// child's exit code so the container runtime sees the right status.
//
// Optionally monitors the child's gRPC health endpoint to detect deadlocks.
// A deadlocked process stays alive but stops serving; the health check catches
// that and kills the container so it can be restarted.
//
// Environment variables:
//
//	LAUNCH_CMD       Command to run, whitespace-split (required)
//	HEALTH_PORT      gRPC port to poll for health (default: 5000; set to "0" to disable)
//	HEALTH_INTERVAL  How often to poll, e.g. "30s" (default: 30s)
//	HEALTH_TIMEOUT   Kill child after this long unhealthy, e.g. "5m" (default: 5m)
//	HEALTH_SERVICE   gRPC service name to check (default: "")
//
// Usage in a Dockerfile:
//
//	COPY launcher /usr/local/bin/stratum
//	ENV LAUNCH_CMD="java -jar /app/stratum.jar"
//	ENTRYPOINT ["stratum"]
//
// LAUNCH_CMD is whitespace-split only — no shell quoting or glob expansion.
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	defaultHealthPort     = "5000"
	defaultHealthInterval = 30 * time.Second
	defaultHealthTimeout  = 5 * time.Minute
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(self() + ": ")

	cmd := strings.TrimSpace(os.Getenv("LAUNCH_CMD"))
	if cmd == "" {
		log.Fatalf("LAUNCH_CMD is not set")
	}

	parts := strings.Fields(cmd)
	args := append(parts, os.Args[1:]...)

	binary, err := findBinary(args[0])
	if err != nil {
		log.Fatal(err)
	}

	child := exec.Command(binary, args[1:]...)
	child.Stdin = os.Stdin
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr

	if err := child.Start(); err != nil {
		log.Fatalf("start %s: %v", binary, err)
	}

	// Forward signals to the child.
	sigs := make(chan os.Signal, 8)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP,
		syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGWINCH)
	go func() {
		for s := range sigs {
			child.Process.Signal(s) //nolint:errcheck
		}
	}()

	// Start the health monitor unless explicitly disabled (HEALTH_PORT=0).
	port := strings.TrimSpace(os.Getenv("HEALTH_PORT"))
	if port == "" {
		port = defaultHealthPort
	}
	if port != "0" {
		go healthMonitor(child, port)
	}

	// Wait and propagate exit code.
	if err := child.Wait(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
		log.Fatalf("wait: %v", err)
	}
}

// healthMonitor polls the child's gRPC health endpoint. If the child is
// unhealthy for longer than HEALTH_TIMEOUT, it is killed so the container
// can be restarted.
func healthMonitor(child *exec.Cmd, port string) {
	interval := envDuration("HEALTH_INTERVAL", defaultHealthInterval)
	timeout := envDuration("HEALTH_TIMEOUT", defaultHealthTimeout)
	service := os.Getenv("HEALTH_SERVICE") // default ""

	if _, err := strconv.Atoi(port); err != nil {
		log.Printf("health: invalid HEALTH_PORT %q: %v — monitoring disabled", port, err)
		return
	}

	target := "127.0.0.1:" + port
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true, ServerName: "localhost", MinVersion: tls.VersionTLS13}))) //nolint:gosec
	if err != nil {
		log.Printf("health: failed to create client for %s: %v — monitoring disabled", target, err)
		return
	}
	defer conn.Close()

	client := healthpb.NewHealthClient(conn)

	// firstFailure holds the epoch-ms of the first consecutive unhealthy result;
	// zero means currently healthy.
	var firstFailure atomic.Int64

	log.Printf("health: monitoring %s every %s (timeout %s)", target, interval, timeout)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: service})
		cancel()

		serving := err == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING

		if serving {
			if prev := firstFailure.Swap(0); prev != 0 {
				log.Printf("health: %s healthy again after %s",
					target, time.Since(time.UnixMilli(prev)).Round(time.Second))
			}
			continue
		}

		now := time.Now().UnixMilli()
		prev := firstFailure.CompareAndSwap(0, now)
		firstMs := now
		if !prev {
			firstMs = firstFailure.Load()
		}
		elapsed := time.Since(time.UnixMilli(firstMs)).Round(time.Second)

		if elapsed >= timeout {
			if err != nil {
				log.Printf("health: %s unreachable for %s (%v) — killing child", target, elapsed, err)
			} else {
				log.Printf("health: %s not serving for %s — killing child", target, elapsed)
			}
			child.Process.Kill() //nolint:errcheck
			return
		}

		if err != nil {
			log.Printf("health: %s check failed (%s elapsed): %v", target, elapsed, err)
		} else {
			log.Printf("health: %s status=%s (%s elapsed)", target, resp.GetStatus(), elapsed)
		}
	}
}

// findBinary resolves the binary path, searching PATH when not absolute.
func findBinary(name string) (string, error) {
	if filepath.IsAbs(name) {
		return name, nil
	}
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		candidate := filepath.Join(dir, name)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("%q not found in PATH", name)
}

func self() string {
	if len(os.Args) > 0 {
		return filepath.Base(os.Args[0])
	}
	return "launcher"
}

func envDuration(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("health: invalid %s=%q: %v — using default %s", key, v, err, def)
		return def
	}
	return d
}
