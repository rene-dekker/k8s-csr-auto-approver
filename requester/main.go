package main

import (
	"context"
	"flag"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"k8s.io/api/certificates/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/goombaio/namegenerator"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	ctx := context.Background()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	certV1Client := clientset.CertificatesV1beta1()

	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()

	log.Println(name)

	signer := "example.com/signer"
	csr := &v1beta1.CertificateSigningRequest{
		TypeMeta: metaV1.TypeMeta{Kind: "CertificateSigningRequest", APIVersion: "certificates.k8s.io/v1beta1"},
		ObjectMeta: metaV1.ObjectMeta{Name: name},
		Spec: v1beta1.CertificateSigningRequestSpec{
			Request: []byte(certRequestPem),
			Username: "user",
			SignerName: &signer,
			Usages: []v1beta1.KeyUsage{v1beta1.UsageCodeSigning},
		},
	}
	created, err := certV1Client.CertificateSigningRequests().Create(ctx, csr, metaV1.CreateOptions{})
	if err != nil {
		log.Fatal("crashed while trying to create csr", err)
	}
	log.Printf("Created CSR: %v", created)

}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/rd/bzprofiles/eks/.local/kubeconfig", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "127.0.0.1:8001", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}


const certRequestPem = `
-----BEGIN CERTIFICATE REQUEST-----
MIIBhTCCASsCAQAwUzEVMBMGA1UEChMMc3lzdGVtOm5vZGVzMTowOAYDVQQDEzFz
eXN0ZW06bm9kZTpteS1wb2QubXktbmFtZXNwYWNlLnBvZC5jbHVzdGVyLmxvY2Fs
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEMX0LcAL1jmE4OTnqDe2k2Cd446dB
dSvTOOOHuXQ9arlai0GsmGrIqarrUfmKzfgGzf83HUaRCSbkJBb9OhDVmKB2MHQG
CSqGSIb3DQEJDjFnMGUwYwYDVR0RBFwwWoIlbXktc3ZjLm15LW5hbWVzcGFjZS5z
dmMuY2x1c3Rlci5sb2NhbIIlbXktcG9kLm15LW5hbWVzcGFjZS5wb2QuY2x1c3Rl
ci5sb2NhbIcEwAACGIcECgAiAjAKBggqhkjOPQQDAgNIADBFAiB8WJT2Moa5ZhnN
fWdz1BaAzQJ/lKo8dtyFO6AIJxLY4QIhAKSvCdgIwQ5oCPHPmRCpVkns8mKDbIv1
nyd0T4ebtcYu
-----END CERTIFICATE REQUEST-----
`
