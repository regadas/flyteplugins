module github.com/lyft/flyteplugins

go 1.13

require (
	github.com/GoogleCloudPlatform/spark-on-k8s-operator v0.0.0-20200210170106-c8ae9e352035
	github.com/Masterminds/semver v1.5.0
	github.com/aws/amazon-sagemaker-operator-for-k8s v1.0.1-0.20200410212604-780c48ecb21a
	github.com/aws/aws-sdk-go v1.29.23
	github.com/coocood/freecache v1.1.0
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-test/deep v1.0.5
	github.com/golang/protobuf v1.4.3
	github.com/hashicorp/golang-lru v0.5.4
	github.com/kubeflow/pytorch-operator v0.6.0
	github.com/kubeflow/tf-operator v0.5.3
	github.com/lyft/flyteidl v0.18.9
	github.com/lyft/flytestdlib v0.3.9
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/mapstructure v1.3.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0

	// github.com/lyft/flinkk8soperator v0.5.0
	github.com/regadas/flink-on-k8s-operator v0.0.0-20201126113113-ee4a54861ce2
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	google.golang.org/grpc v1.33.2
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v11.0.1-0.20190918222721-c0e3722d5cf0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

// Pin the version of client-go to something that's compatible with katrogan's fork of api and apimachinery
// Type the following
//   replace k8s.io/client-go => k8s.io/client-go kubernetes-1.16.2
// and it will be replaced with the 'sha' variant of the version

replace (
	github.com/GoogleCloudPlatform/spark-on-k8s-operator => github.com/lyft/spark-on-k8s-operator v0.1.4-0.20201027003055-c76b67e3b6d0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1
	//FIXME(regadas):
	// github.com/lyft/flyteidl => /Users/regadas/projects/flyteidl
	github.com/lyft/flyteidl => github.com/regadas/flyteidl v0.18.11-0.20201207222255-cb6e9f38973d
	k8s.io/api => github.com/lyft/api v0.0.0-20191031200350-b49a72c274e0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.3
	k8s.io/apimachinery => github.com/lyft/apimachinery v0.0.0-20191031200210-047e3ea32d7f
	k8s.io/apiserver => k8s.io/apiserver v0.18.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.3
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.3
	k8s.io/code-generator => k8s.io/code-generator v0.16.10-beta.0
	k8s.io/component-base => k8s.io/component-base v0.18.3
	k8s.io/cri-api => k8s.io/cri-api v0.18.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.3
	k8s.io/kubectl => k8s.io/kubectl v0.18.3
	k8s.io/kubelet => k8s.io/kubelet v0.18.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.3
	k8s.io/metrics => k8s.io/metrics v0.18.3
	k8s.io/node-api => k8s.io/node-api v0.18.3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.3
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.18.3
	k8s.io/sample-controller => k8s.io/sample-controller v0.18.3
)
