package client

import (
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	k "github.com/infracloudio/krius/pkg/kubeClient"
	"gopkg.in/yaml.v2"
)

func NewThanosClient(thanosCluster *Cluster) (*Thanos, error) {
	thanosConfig, err := GetConfig(thanosCluster.Data, "thanos")
	if err != nil {
		log.Printf("Error getting config %s", err)
		return nil, err
	}
	spec, _ := yaml.Marshal(thanosConfig)
	var thanos Thanos
	err = yaml.Unmarshal(spec, &thanos)
	if err != nil {
		log.Printf("Error unmarshaling %s", err)
		return nil, err
	}
	return &thanos, nil
}

func (t *Thanos) PreflightChecks(clusterConfig *Config, clusterName string) ([]string, error) {
	kubeClient, err := k.GetKubeClient(t.Namespace, clusterName)
	if err != nil {
		return nil, err
	}
	thanosErrs := []string{}
	err = kubeClient.CreateNSIfNotExist()
	if err != nil {
		e := fmt.Sprintf("cluster.%s: %s,", clusterName, err)
		thanosErrs = append(thanosErrs, e)
		return thanosErrs, nil // don't create secret, if error in creating namespace
	}

	found := false
	for _, v := range clusterConfig.ObjStoreConfigslist {
		if v.Name == t.ObjStoreConfig {
			found = true
			secretSpec, err := createSecretforObjStore(v.Type, v.Config)
			if err != nil {
				return nil, err
			}
			err = kubeClient.CreateSecret(secretSpec, t.ObjStoreConfig)

			if err != nil {
				e := fmt.Sprintf("cluster.%s: %s,", t.Name, err)
				thanosErrs = append(thanosErrs, e)
			}
			break
		}
	}
	if !found {
		e := fmt.Sprintf("cluster.%s: Bucket config doesn't exist,", clusterName)
		thanosErrs = append(thanosErrs, e)
	}
	return thanosErrs, nil
}

func (t *Thanos) InstallClient(clusterName string) (string, error) {

	chartConfiguration := &helm.Config{
		Repo: "bitnami",
		Name: "thanos",
		URL:  "https://charts.bitnami.com/bitnami",
	}

	helmClient, err := createHelmClientObject(clusterName, t.Namespace, chartConfiguration)
	if err != nil {
		return "", err
	}
	helmClient.ChartName = "thanos"
	helmClient.ReleaseName = "thanos"
	var extraFlags []string
	if t.Querier.AutoDownsample {
		extraFlags = append(extraFlags, "--query.auto-downsampling")
	}
	if t.Querier.PartialResponse {
		extraFlags = append(extraFlags, "--query.partial-response")
	}
	t.Querier.ExtraFlags = extraFlags
	Values := createThanosValuesMap(*t)
	_, err = helmClient.InstallChart(Values)
	log.Println("error installing Thanos", err)
	return "", err

}
