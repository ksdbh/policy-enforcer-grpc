package main
import (
  "crypto/tls"
  "crypto/x509"
  "io/ioutil"
  "log"
  "net"
  "os"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  policyv1 "github.com/ksdbh/policy-enforcer-grpc/proto/policyv1"
)
type policyServer struct{ policyv1.UnimplementedPolicyServiceServer }
func (s *policyServer) Evaluate(_ interface{}, req *policyv1.EvaluateRequest) (*policyv1.EvaluateResponse, error) {
  allow := true; reasons := []string{}
  if req.Action == "delete" && req.Resource == "payments" {
    if role, ok := req.Attributes["role"]; !ok || role != "admin" {
      allow = false; reasons = append(reasons, "only admin can delete payments")
    }
  }
  return &policyv1.EvaluateResponse{Allow: allow, Reasons: reasons}, nil
}
func main() {
  addr := ":50051"
  certFile := get("TLS_CERT", "certs/server.pem")
  keyFile  := get("TLS_KEY",  "certs/server.key")
  caFile   := get("TLS_CA",   "certs/ca.pem")
  cert, err := tls.LoadX509KeyPair(certFile, keyFile); if err != nil { log.Fatalf("load keypair: %v", err) }
  caBytes, err := ioutil.ReadFile(caFile); if err != nil { log.Fatalf("read ca: %v", err) }
  pool := x509.NewCertPool(); pool.AppendCertsFromPEM(caBytes)
  creds := credentials.NewTLS(&tls.Config{ Certificates: []tls.Certificate{cert}, ClientCAs: pool, ClientAuth: tls.RequireAndVerifyClientCert })
  s := grpc.NewServer(grpc.Creds(creds))
  policyv1.RegisterPolicyServiceServer(s, &policyServer{})
  lis, err := net.Listen("tcp", addr); if err != nil { log.Fatalf("listen: %v", err) }
  log.Printf("policy server listening on %s (mTLS)", addr)
  if err := s.Serve(lis); err != nil { log.Fatalf("serve: %v", err) }
}
func get(k, d string) string { if v, ok := os.LookupEnv(k); ok { return v }; return d }
