package nts

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

type KEServer struct {
	Addr            string
	CertPEM, KeyPEM []byte
	OnConn          func(net.Conn)
}

func SelfSigned() ([]byte, []byte, error) {
	k, e := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if e != nil {
		return nil, nil, e
	}
	der, e := x509.CreateCertificate(rand.Reader, &x509.Certificate{SerialNumber: bigOne(), NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(24 * time.Hour), DNSNames: []string{"chronos-mesh"}}, &x509.Certificate{SerialNumber: bigOne()}, &k.PublicKey, k)
	if e != nil {
		return nil, nil, e
	}
	kb, _ := x509.MarshalECPrivateKey(k)
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: kb}), nil
}
func bigOne() *big.Int { return big.NewInt(1) }
func (s *KEServer) Start(ctx context.Context) error {
	if len(s.CertPEM) == 0 {
		return fmt.Errorf("certificate required")
	}
	ln, e := net.Listen("tcp", s.Addr)
	if e != nil {
		return e
	}
	go func() { <-ctx.Done(); ln.Close() }()
	for {
		c, e := ln.Accept()
		if e != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return e
			}
		}
		if s.OnConn != nil {
			s.OnConn(c)
		} else {
			c.Close()
		}
	}
}
