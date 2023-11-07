module github.com/arangodb-managed/oasisctl

go 1.20

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3

require (
	github.com/araddon/dateparse v0.0.0-20200409225146-d820a6159ab1
	github.com/arangodb-managed/apis v0.83.1
	github.com/arangodb-managed/arangocopy v0.0.0-20230330143258-9e03ba080b35
	github.com/coreos/go-semver v0.3.0
	github.com/dchest/uniuri v0.0.0-20160212164326-8902c56451e9
	github.com/dustin/go-humanize v1.0.0
	github.com/gizak/termui/v3 v3.1.0
	github.com/gogo/protobuf v1.3.2
	github.com/rs/zerolog v1.19.0
	github.com/ryanuber/columnize v2.1.0+incompatible
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.29.1
)

require (
	github.com/VividCortex/ewma v1.1.1 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/arangodb/go-driver v1.5.2 // indirect
	github.com/arangodb/go-velocypack v0.0.0-20200318135517-5af53c29c67e // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.6 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/nsf/termbox-go v0.0.0-20190121233118-02980233997d // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/vbauerster/mpb/v5 v5.2.2 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.37.0

replace github.com/hashicorp/vault/api => github.com/hashicorp/vault/api v1.9.1

replace github.com/hashicorp/vault/sdk => github.com/hashicorp/vault/sdk v0.9.0

replace github.com/hashicorp/vault => github.com/hashicorp/vault v1.10.11

replace github.com/kamilsk/retry => github.com/kamilsk/retry/v3 v3.4.4

replace github.com/nats-io/go-nats-streaming => github.com/nats-io/go-nats-streaming v0.4.4

replace github.com/nats-io/go-nats => github.com/nats-io/go-nats v1.7.2

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20190108154635-47c0da630f72

replace github.com/ugorji/go => github.com/ugorji/go v1.2.11

replace google.golang.org/api => google.golang.org/api v0.36.0

replace google.golang.org/grpc => google.golang.org/grpc v1.36.0

replace k8s.io/api => k8s.io/api v0.27.5

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.27.5

replace k8s.io/apimachinery => k8s.io/apimachinery v0.27.5

replace k8s.io/apiserver => k8s.io/apiserver v0.27.5

replace k8s.io/client-go => k8s.io/client-go v0.27.5

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.27.5

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.27.5

replace k8s.io/code-generator => k8s.io/code-generator v0.27.5

replace k8s.io/component-base => k8s.io/component-base v0.27.5

replace k8s.io/kubernetes => k8s.io/kubernetes v1.27.5

replace k8s.io/metrics => k8s.io/metrics v0.27.5

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.15.2

replace sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.2

replace github.com/arangodb/kube-arangodb => github.com/arangodb/kube-arangodb v0.0.0-20230616093551-2a4a07da3442

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.1-0.20210602190049-10e0b31633f1+incompatible

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6

replace github.com/cilium/cilium => github.com/cilium/cilium v1.13.4

replace github.com/optiopay/kafka => github.com/optiopay/kafka v2.0.4+incompatible

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.27.5

replace k8s.io/cri-api => k8s.io/cri-api v0.27.5

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.27.5

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.27.5

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.27.5

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.27.5

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.27.5

replace k8s.io/kubelet => k8s.io/kubelet v0.27.5

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.27.5

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.27.5

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20201211151036-40ec1c210f7a

replace k8s.io/kubectl => k8s.io/kubectl v0.27.5

replace github.com/nats-io/nats.go => github.com/nats-io/nats.go v1.24.0

replace github.com/nats-io/stan.go => github.com/nats-io/stan.go v0.10.4

replace github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring => github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.66.0

replace github.com/prometheus-operator/prometheus-operator/pkg/client => github.com/prometheus-operator/prometheus-operator/pkg/client v0.66.0

replace github.com/prometheus-operator/prometheus-operator => github.com/prometheus-operator/prometheus-operator v0.66.0

replace go.uber.org/multierr => go.uber.org/multierr v1.11.0

replace k8s.io/component-helpers => k8s.io/component-helpers v0.27.5

replace k8s.io/controller-manager => k8s.io/controller-manager v0.27.5

replace k8s.io/mount-utils => k8s.io/mount-utils v0.27.5

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.27.5

replace helm.sh/helm/v3 => helm.sh/helm/v3 v3.12.3

replace github.com/coreos/go-systemd/v22 => github.com/coreos/go-systemd/v22 v22.5.0

replace k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.27.5

replace k8s.io/kms => k8s.io/kms v0.27.5
