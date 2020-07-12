package utils

import (
	"os"

	// "github.com/rs/zerolog/log"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // To allow clusters running behind OIDC proxy
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// KubeConfig represents kubernetes config to connect with kube-api-server
var KubeConfig *rest.Config = nil

/*
CheckKubeAPI checks if kubeconfig exists and is working
*/
func CheckKubeAPI() bool {
	status := false

	if KubeConfig == nil {
		InitKubeConfig()
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(KubeConfig)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// check running pods
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})

	if err != nil {
		log.Error().Msg(err.Error())
	} else {
		log.Info().Msgf("Pods running in cluster: %d", len(pods.Items))
		status = true
	}
	return status

}

/*
HomeDir provides you the path to home directory from environment variable
*/
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") //windows
}

// InitKubeConfig returns the kube config or nil in case of error
func InitKubeConfig() {
	var err error
	KubeConfig, err = config.GetConfig()
	if err != nil {
		log.Fatal().Msg("Unable to find kube config. Either use service account for access within container or ensure .kube/config in your HOME")
	}
}
