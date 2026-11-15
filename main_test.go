// Copyright (c) 2026 Jody Scott
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// testCA is a self-signed CA that issues the child's server certificates.
type testCA struct {
	cert *x509.Certificate
	key  *ecdsa.PrivateKey
	pem  []byte
}

func newTestCA(t *testing.T, name string) *testCA {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return &testCA{cert: cert, key: key, pem: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})}
}

// serverCert issues a server certificate for the DNS name dnsName.
func (ca *testCA) serverCert(t *testing.T, dnsName string) tls.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: dnsName},
		DNSNames:     []string{dnsName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca.cert, &key.PublicKey, ca.key)
	if err != nil {
		t.Fatal(err)
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

// setHealthCA writes the CA's PEM to a file and points HEALTH_CA at it.
func (ca *testCA) setHealthCA(t *testing.T) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(path, ca.pem, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HEALTH_CA", path)
}

// startHealthServer serves grpc.health.v1 (SERVING) on a loopback port, over
// TLS with cert when cert is non-nil, plaintext otherwise, and returns the port.
func startHealthServer(t *testing.T, cert *tls.Certificate) string {
	t.Helper()
	var opts []grpc.ServerOption
	if cert != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{*cert},
			MinVersion:   tls.VersionTLS13,
		})))
	}
	srv := grpc.NewServer(opts...)
	healthpb.RegisterHealthServer(srv, health.NewServer())
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go srv.Serve(lis) //nolint:errcheck
	t.Cleanup(srv.Stop)
	_, port, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	return port
}

// probe runs one health check against the port, as healthMonitor does.
func probe(t *testing.T, port string) (healthpb.HealthCheckResponse_ServingStatus, error) {
	t.Helper()
	conn, err := newHealthConn("127.0.0.1:" + port)
	if err != nil {
		return healthpb.HealthCheckResponse_UNKNOWN, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{})
	return resp.GetStatus(), err
}

func clearHealthEnv(t *testing.T) {
	t.Helper()
	t.Setenv("HEALTH_CA", "")
	t.Setenv("HEALTH_SERVER_NAME", "")
}

func TestHealthProbeTLSVerifiedAgainstHealthCA(t *testing.T) {
	clearHealthEnv(t)
	ca := newTestCA(t, "service CA")
	cert := ca.serverCert(t, "localhost")
	port := startHealthServer(t, &cert)
	ca.setHealthCA(t)

	status, err := probe(t, port)
	if err != nil || status != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("probe = %v, %v; want SERVING", status, err)
	}
}

func TestHealthProbeWrongCertificateFails(t *testing.T) {
	clearHealthEnv(t)
	trusted := newTestCA(t, "service CA")
	other := newTestCA(t, "another CA")
	cert := other.serverCert(t, "localhost")
	port := startHealthServer(t, &cert)
	trusted.setHealthCA(t)

	if status, err := probe(t, port); err == nil {
		t.Fatalf("probe of a certificate HEALTH_CA did not issue = %v, nil; want an error", status)
	}
}

func TestHealthProbeWrongServerNameFails(t *testing.T) {
	clearHealthEnv(t)
	ca := newTestCA(t, "service CA")
	cert := ca.serverCert(t, "someone-else")
	port := startHealthServer(t, &cert)
	ca.setHealthCA(t)

	if status, err := probe(t, port); err == nil {
		t.Fatalf("probe of a certificate not for localhost = %v, nil; want an error", status)
	}

	t.Setenv("HEALTH_SERVER_NAME", "someone-else")
	status, err := probe(t, port)
	if err != nil || status != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("probe with HEALTH_SERVER_NAME=someone-else = %v, %v; want SERVING", status, err)
	}
}

func TestHealthProbePlaintextWithoutHealthCA(t *testing.T) {
	clearHealthEnv(t)
	port := startHealthServer(t, nil)

	status, err := probe(t, port)
	if err != nil || status != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("plaintext probe = %v, %v; want SERVING", status, err)
	}
}

func TestHealthProbePlaintextAgainstTLSPortFails(t *testing.T) {
	clearHealthEnv(t)
	ca := newTestCA(t, "service CA")
	cert := ca.serverCert(t, "localhost")
	port := startHealthServer(t, &cert)

	if status, err := probe(t, port); err == nil {
		t.Fatalf("plaintext probe of a TLS port = %v, nil; want an error", status)
	}
}

func TestHealthCredentialsRejectsBadHealthCA(t *testing.T) {
	clearHealthEnv(t)
	t.Setenv("HEALTH_CA", filepath.Join(t.TempDir(), "missing.pem"))
	if _, err := healthCredentials(); err == nil {
		t.Fatal("healthCredentials with a missing HEALTH_CA: want an error")
	}

	path := filepath.Join(t.TempDir(), "empty.pem")
	if err := os.WriteFile(path, []byte("not a certificate\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HEALTH_CA", path)
	if _, err := healthCredentials(); err == nil {
		t.Fatal("healthCredentials with no certificate in HEALTH_CA: want an error")
	}
}
