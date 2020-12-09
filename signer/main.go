package main

import (
	"context"
	"flag"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"k8s.io/api/certificates/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	watchers, err := certV1Client.CertificateSigningRequests().Watch(ctx, metaV1.ListOptions{})

	ch := watchers.ResultChan()

	for event := range ch {
		csr, ok := event.Object.(*v1beta1.CertificateSigningRequest)
		if !ok {
			log.Fatal("unexpected type in cert channel")
		}
		log.Printf("CSR: %v", csr)
		cert := csr.DeepCopy()
		if csr.Status.Certificate == nil {

			cert.Status.Certificate = []byte("pancake")
			r, err := certV1Client.CertificateSigningRequests().UpdateStatus(ctx, cert, metaV1.UpdateOptions{})
			if err != nil {
				log.Fatal("unexpected err when updating csr")
			}
			log.Printf("CSR status updated: %v", r.Status)
		} else if len(csr.Status.Conditions) == 0 {
			cert.Status.Conditions = []v1beta1.CertificateSigningRequestCondition{
				{
					Type: v1beta1.CertificateApproved,
					Message: "Approved",
					Reason: "Approved",
				},
			}
			certV1Client.CertificateSigningRequests().UpdateApproval(ctx, cert, metaV1.UpdateOptions{})
		}
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/rd/bzprofiles/kadm/.local/kubeconfig", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "127.0.0.1:8001", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
