package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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

/*
ParseCPU parses the string containing cores specified in nano, milli, micro and gives the number of cores as floating point value
*/
func ParseCPU(val string) (float64, error) {
	result := 0.0
	var errorE error
	n := len(val)
	re, err := regexp.Compile(`^([0-9]+)((n)|(u)|(m))$`)
	if err != nil {
		log.Err(err)
		errorE = err
	} else {
		matches := re.FindStringSubmatch(val)
		if len(matches) != 0 {
			value, err := strconv.Atoi(string(matches[1]))
			if err != nil {
				log.Err(err)
				errorE = err
			} else {
				var unit string = string(val[n-1])
				switch unit {
				case "n":
					result = float64(value) / 1000000000.0
				case "u":
					result = float64(value) / 1000000.0
				case "m":
					result = float64(value) / 1000.0
				default:
					result = float64(value)
				}
			}
		} else {
			var value int
			runes := []rune(val)
			value, err = strconv.Atoi(string(runes[0 : n-1]))
			if err != nil {
				log.Err(err)
				errorE = err
			} else {
				result = float64(value)
			}
		}
	}
	return result, errorE
}

/*
ParseMemory parses the string containing memory specified in Gi, Mi, Ki and gives the memory in bytes as floating point value
*/
func ParseMemory(val string) (float64, error) {
	result := 0.0
	var errorE error
	n := len(val)
	re, err := regexp.Compile(`^([0-9]+)((Ki)|(Mi)|(Gi)|(Ti))$`)
	if err != nil {
		log.Err(err)
		errorE = err
	} else {
		matches := re.FindStringSubmatch(val)
		if len(matches) != 0 {
			value, err := strconv.Atoi(string(matches[1]))
			if err != nil {
				log.Err(err)
				errorE = err
			} else {
				var unit string = strings.Trim(fmt.Sprintf("%s%s", string(val[n-2]), string(val[n-1])), " \t")
				switch unit {
				case "Ki":
					result = float64(value) * 1024.0
				case "Mi":
					result = float64(value) * 1024.0 * 1024.0
				case "Gi":
					result = float64(value) * 1024.0 * 1024.0 * 1024.0
				case "Ti":
					result = float64(value) * 1024.0 * 1024.0 * 1024.0 * 1024.0
				default:
					result = float64(value)
				}
			}
		} else {
			var value int
			runes := []rune(val)
			value, err = strconv.Atoi(string(runes[0 : n-1]))
			if err != nil {
				log.Err(err)
				errorE = err
			} else {
				result = float64(value)
			}
		}
	}
	return result, errorE
}
