package sdk

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"time"
)

type ApiOrderer struct {
	Url           string `yaml:"url"`
	Tls           bool   `yaml:"tls"`
	TlsCertFile   string `yaml:"tlsCertFile"`
	TlsServerName string `yaml:"tlsServerName"`
}

func (ap *ApiOrderer) GrpcConn(ctx context.Context) (conn *grpc.ClientConn, err error) {

	// TODO keep alive, timeout etc
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	peerCertPEM, err := ioutil.ReadFile(ap.TlsCertFile)
	if err != nil {
		fmt.Println(err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(peerCertPEM)
	tlsConfig := &tls.Config{RootCAs: certPool, ServerName: ap.TlsServerName}

	if ap.Tls {
		conn, err = grpc.DialContext(ctx, ap.Url, grpc.WithBlock(), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		conn, err = grpc.DialContext(ctx, ap.Url, grpc.WithBlock(), grpc.WithInsecure())
	}
	return
}
