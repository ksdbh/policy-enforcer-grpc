package main
import (
  "context"
  "crypto/tls"
  "crypto/x509"
  "io/ioutil"
  "log"
  "os"
  "time"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  policyv1 "github.com/ksdbh/policy-enforcer-grpc/proto/policyv1"
)
func main() {
  addr := "localhost:50051"
  ca := os.Getenv("TLS_CA"); cc := os.Getenv("CLIENT_CERT"); ck := os.Getenv("CLIENT_KEY")
  if ca == "" || cc == "" || ck == "" { log.Fatal("set TLS_CA, CLIENT_CERT, CLIENT_KEY") }
  cert, err := tls.LoadX509KeyPair(cc, ck); if err != nil { log.Fatalf("load client cert: %v", err) }
  caBytes, err := ioutil.ReadFile(ca); if err != nil { log.Fatalf("read ca: %v", err) }
  pool := x509.NewCertPool(); pool.AppendCertsFromPEM(caBytes)
  creds := credentials.NewTLS(&tls.Config{ Certificates: []tls.Certificate{cert}, RootCAs: pool })
  conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds)); if err != nil { log.Fatalf("dial: %v", err) }
  defer conn.Close()
  c := policyv1.NewPolicyServiceClient(conn)
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second); defer cancel()
  resp, err := c.Evaluate(ctx, &policyv1.EvaluateRequest{ Resource: "payments", Action: "delete", Attributes: map[string]string{"role":"user"} })
  if err != nil { log.Fatalf("rpc: %v", err) }
  log.Printf("allow=%v reasons=%v", resp.Allow, resp.Reasons)
}
