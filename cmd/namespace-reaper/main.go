package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const expiryLabel string = "namespace-reaper.dwardu.com/expires-in-hours"
const reapingLabel string = "namespace-reaper.dwardu.com/reap"
const dateUpdatedLabel string = "namespace-reaper.dwardu.com/date-updated"

var expiresInHours float64 = 8
var gracePeriodSeconds int64 = 600 // 10 minutes gracePeriod
var checkInterval uint64 = 60
var excludedNamespaces []string = []string{"kube-node-lease", "kube-public", "kube-system"}

// Determine whether you are running inside or outside of a cluster.
func getConfiguration(kubeConfigPath string) (*rest.Config, error) {

	kubeConfigPath = os.Getenv("KUBE_CONFIG_PATH")

	var config *rest.Config
	var err error

	if len(kubeConfigPath) != 0 {
		// creates the config from kubeconfig.
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
	}
	return config, err

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func markReady() {
	f, err := os.Create("/tmp/namespace-reaper-ready")

	if err != nil {
		log.WithFields(log.Fields{
			"err":           err,
			"error_message": err.Error(),
		}).Panic("Error creating ready file.")
	}
	defer f.Close()
}

func cleanLive() {
	err := os.Remove("/tmp/namespace-reaper-live")
	if err != nil {
		log.WithFields(log.Fields{
			"err":           err,
			"error_message": err.Error(),
		}).Panic("Error removing liveness file.")
	}
}

func markLive() {
	f, err := os.Create("/tmp/namespace-reaper-live")
	if err != nil {
		log.WithFields(log.Fields{
			"err":           err,
			"error_message": err.Error(),
		}).Panic("Error creating ready file.")
	}
	defer f.Close()
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	logLevelEnv := getEnv("LOG_LEVEL", "info")
	logLevel, err := log.ParseLevel(logLevelEnv)
	if err != nil {
		log.WithFields(log.Fields{
			"err":           err,
			"error_message": err.Error(),
		}).Warningf("Log Level passed is invalid [%s]. Falling back to [INFO] log level.", os.Getenv("LOG_LEVEL"))
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	checkIntervalEnv := getEnv("CHECK_INTERVAL", "60")
	checkInterval, _ = strconv.ParseUint(checkIntervalEnv, 10, 64)

	excludedNamespacesEnv := getEnv("EXCLUDED_NAMESPACES", "kube-node-lease,kube-public,kube-system")
	excludedNamespaces = strings.Split(excludedNamespacesEnv, ",")
}

func main() {

	// Get the configuration of the cluster
	config, err := getConfiguration(os.Getenv("KUBE_CONFIG_PATH"))
	if err != nil {
		log.WithFields(log.Fields{
			"err":           err,
			"error_message": err.Error(),
		}).Panicf("Error loading kube config from path %s", os.Getenv("KUBE_CONFIG_PATH"))
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.WithFields(log.Fields{
			"err":           err,
			"error_message": err.Error(),
		}).Panic("Error loading kube config from in cluster config.")
	}

	markReady()
	for {
		namespaces := getNamespaces(clientset)
		cleanupNamespaces(clientset, namespaces)
		markLive()
		time.Sleep(time.Duration(checkInterval) * time.Second)
		cleanLive()
	}
}

func getNamespaces(clientset *kubernetes.Clientset) []v1.Namespace {
	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})

	log.WithFields(log.Fields{
		"namespaces_count": len(namespaces.Items),
	}).Infof("There are %d namespaces in the cluster", len(namespaces.Items))

	if err != nil {
		panic(err.Error())
	}
	return namespaces.Items
}

// func describeNamespace(namespace v1.Namespace) {
// 	if validUntil, ok := namespace.Annotations["namespace-reaper/valid-until"]; ok {
// 		log.WithFields(log.Fields{
// 			"namespace": namespace.Name,
// 			"eternal":   true,
// 		}).Infof("Namespace %s is Valid until %s.", namespace.Name, validUntil)
// 	} else {
// 		log.WithFields(log.Fields{
// 			"namespace": namespace.Name,
// 			"eternal":   true,
// 		}).Tracef("Namespace %s is Eternal.", namespace.Name)
// 	}
// }

func deleteNamespace(clientset *kubernetes.Clientset, namespace v1.Namespace) {
	deletePolicy := metav1.DeletePropagationBackground
	if err := clientset.CoreV1().Namespaces().Delete(namespace.GetName(),
		&metav1.DeleteOptions{
			PropagationPolicy:  &deletePolicy,
			GracePeriodSeconds: &gracePeriodSeconds,
			DryRun:             []string{}}); err != nil {
		log.WithFields(log.Fields{
			"err":           err,
			"error_message": err.Error(),
		}).Panicf("Error deleting the namespace [%s].", namespace.GetName())
	}
	log.WithFields(log.Fields{
		"namespace": namespace.Name,
	}).Infof("Namespace [%s] has been marked for deletion.", namespace.Name)
}

func isNamespaceExpired(namespace v1.Namespace) bool {
	now := time.Now()

	annotations := namespace.GetAnnotations()
	reapingAnnotation := annotations[reapingLabel]

	if len(reapingAnnotation) > 0 {
		reap, err := strconv.ParseBool(reapingAnnotation)
		if err != nil {
			log.WithFields(log.Fields{
				"err":           err,
				"error_message": err.Error(),
			}).Panicf("Reaping Annotation value is malformed, it should be true/false. Actual Value [%s].", reapingAnnotation)
		}
		if reap {
			namespaceAge := now.Sub(namespace.CreationTimestamp.Time)
			isExpired := (namespaceAge.Hours() > expiresInHours)
			log.WithFields(log.Fields{
				"namespace": namespace.Name,
				"eternal":   false,
				"expired":   isExpired,
			}).Debugf("Namespace [%s] is still valid? [%t].", namespace.Name, !isExpired)

			return isExpired
		}
	}
	log.WithFields(log.Fields{
		"namespace": namespace.Name,
	}).Debugf("Namespace [%s] doesnt have reaping annotation. It will be ignored.", namespace.Name)
	return false
}

func isNotExcluded(namespace v1.Namespace) bool {
	for _, excludedNamespace := range excludedNamespaces {
		if excludedNamespace == namespace.Name {
			log.WithFields(log.Fields{
				"namespace":          namespace.Name,
				"excludedNamespaces": excludedNamespaces,
			}).Debugf("Namespace [%s] is excluded from reaping.", namespace.Name)
			return true
		}
	}
	log.WithFields(log.Fields{
		"namespace":          namespace.Name,
		"excludedNamespaces": excludedNamespaces,
	}).Debugf("Namespace [%s] is not excluded from reaping.", namespace.Name)
	return false
}

func cleanupNamespaces(clientset *kubernetes.Clientset, namespaces []v1.Namespace) {
	for _, namespace := range namespaces {
		if isNotExcluded(namespace) && isNamespaceExpired(namespace) {
			deleteNamespace(clientset, namespace)
		}
	}
}
