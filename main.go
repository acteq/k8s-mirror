package main

import (
	"encoding/json"
	// "strings"
	// "flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"crypto/tls"
	// TODO: try this library to see if it generates correct json patch
	// https://github.com/mattbaird/jsonpatch


)

// toAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

// admitFunc is the type we use for all of our validators and mutators
type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

// serve handles the http portion of a request prior to handing to an admit
// function
func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	klog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	// The AdmissionReview that was sent to the webhook
	requestedAdmissionReview := v1beta1.AdmissionReview{}

	// The AdmissionReview that will be returned
	responseAdmissionReview := v1beta1.AdmissionReview{}

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		klog.Error(err)
		responseAdmissionReview.Response = toAdmissionResponse(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = admit(requestedAdmissionReview)
	}

	// Return the same UID
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		klog.Error(err)
	}
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func configTLS(certFileName string, keyFileName string) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(certFileName, keyFileName)
	if err != nil {
		klog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}

func withEnv(env map[string]string, origin func(v1beta1.AdmissionReview, map[string]string) *v1beta1.AdmissionResponse) func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
    innerfunc := func(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse{
        return origin(ar, env)
    }
    return innerfunc
}

func handleWithEnv(origin func(v1beta1.AdmissionReview, map[string]string) *v1beta1.AdmissionResponse, env map[string]string) func(w http.ResponseWriter, r *http.Request) {
    handle := func(w http.ResponseWriter, r *http.Request){
		innerfunc := func(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse{
			return origin(ar, env)
		}
		serve(w, r, innerfunc)
    }
    return handle
}

func main() {

	env, certFile, keyFile := getConfigure()
	
	handleMutate := handleWithEnv(mutateResource, env)

	http.HandleFunc("/mutate", handleMutate)
	http.HandleFunc("/ping", pong)
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: configTLS(certFile, keyFile),
	}
	fmt.Println("start ......")
	server.ListenAndServeTLS("", "")
}
