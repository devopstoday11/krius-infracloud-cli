module github.com/infracloudio/krius

go 1.15

require (
	github.com/gofrs/flock v0.8.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.5.1
	rsc.io/letsencrypt v0.0.3 // indirect
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309