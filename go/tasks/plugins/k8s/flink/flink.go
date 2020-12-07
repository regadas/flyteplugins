package flink

import (
	"context"
	"math/rand"
	"strconv"
	"time"
	"unsafe"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/flytek8s"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/flytek8s/config"

	"github.com/lyft/flyteplugins/go/tasks/errors"
	pluginsCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/k8s"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	flinkOp "github.com/regadas/flink-on-k8s-operator/api/v1beta1"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/plugins"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pluginsConfig "github.com/lyft/flyteplugins/go/tasks/config"
	"github.com/lyft/flytestdlib/logger"
)

const KindFlinkCluster = "FlinkCluster"

var (
	flinkTaskType       = "flink"
	jobManagerUiPort    = int32(8081)
	taskManagerReplicas = int32(1)
	jobParallelism      = int32(1)
	flinkImage          = "flink:1.10.1-scala_2.12"
	flinkConfigSection  = pluginsConfig.MustRegisterSubSection("flink", &Config{})
)

// Config ... Flink-specific configs
type Config struct {
	DefaultFlinkConfig map[string]string `json:"flink-config-default" pflag:",Key value pairs of default flink configuration that should be applied to every FlinkJob"`
}

func GetFlinkConfig() *Config {
	return flinkConfigSection.GetConfig().(*Config)
}

// This method should be used for unit testing only
func setFlinkConfig(cfg *Config) error {
	return flinkConfigSection.SetConfig(cfg)
}

type flinkResourceHandler struct{}

func (flinkResourceHandler) GetProperties() pluginsCore.PluginProperties {
	return pluginsCore.PluginProperties{}
}

// Creates a new Job that will execute the main container as well as any generated types the result from the execution.
func (flinkResourceHandler) BuildResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionContext) (k8s.Resource, error) {
	taskTemplate, err := taskCtx.TaskReader().Read(ctx)
	if err != nil {
		return nil, errors.Errorf(errors.BadTaskSpecification, "unable to fetch task specification [%v]", err.Error())
	} else if taskTemplate == nil {
		return nil, errors.Errorf(errors.BadTaskSpecification, "nil task specification")
	}

	flinkJob := plugins.FlinkJob{}
	err = utils.UnmarshalStruct(taskTemplate.GetCustom(), &flinkJob)
	if err != nil {
		return nil, errors.Wrapf(errors.BadTaskSpecification, err, "invalid TaskSpecification [%v], failed to unmarshal", taskTemplate.GetCustom())
	}

	annotations := utils.UnionMaps(
		config.GetK8sPluginConfig().DefaultAnnotations,
		utils.CopyMap(taskCtx.TaskExecutionMetadata().GetAnnotations()),
	)
	labels := utils.UnionMaps(
		config.GetK8sPluginConfig().DefaultLabels,
		utils.CopyMap(taskCtx.TaskExecutionMetadata().GetLabels()),
	)
	container := taskTemplate.GetContainer()
	envVars := flytek8s.DecorateEnvVars(
		ctx,
		flytek8s.ToK8sEnvVar(container.GetEnv()),
		taskCtx.TaskExecutionMetadata().GetTaskExecutionID(),
	)

	flinkEnvVars := make(map[string]string)
	for _, envVar := range envVars {
		flinkEnvVars[envVar.Name] = envVar.Value
	}
	flinkEnvVars["FLYTE_MAX_ATTEMPTS"] = strconv.Itoa(int(taskCtx.TaskExecutionMetadata().GetMaxAttempts()))

	logger.Debugf(ctx, "flinkEnvVars: %#v", flinkEnvVars)

	jobManager := flinkOp.JobManagerSpec{
		PodAnnotations: annotations,
		PodLabels:      labels,
		Ports: flinkOp.JobManagerPorts{
			UI: &jobManagerUiPort,
		},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
		},
	}

	taskManager := flinkOp.TaskManagerSpec{
		PodAnnotations: annotations,
		PodLabels:      labels,
		Replicas:       taskManagerReplicas,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
	}
	// // Add Tolerations/NodeSelector to TaskManagers
	// if taskCtx.TaskExecutionMetadata().IsInterruptible() {
	// 	taskManager.Tolerations = config.GetK8sPluginConfig().InterruptibleTolerations
	// 	taskManager.NodeSelector = config.GetK8sPluginConfig().InterruptibleNodeSelector
	// }

	job := flinkOp.JobSpec{
		JarFile:      "/cache/job.jar",
		ClassName:    &flinkJob.MainClass,
		Args:         flinkJob.Args,
		Parallelism:  &jobParallelism,
		Volumes:      []corev1.Volume{{Name: "cache-volume"}},
		VolumeMounts: []corev1.VolumeMount{{Name: "cache-volume", MountPath: "/cache"}},
		InitContainers: []corev1.Container{
			{
				Name:  "gcs-downloader",
				Image: "google/cloud-sdk",
				Args: []string{
					"cp",
					flinkJob.JarFile,
					"/cache/job.jar",
				},
			},
		},
	}

	// Start with default config values.
	flinkProperties := make(map[string]string)
	for k, v := range GetFlinkConfig().DefaultFlinkConfig {
		flinkProperties[k] = v
	}

	for k, v := range flinkJob.GetFlinkProperties() {
		flinkProperties[k] = v
	}

	logger.Debugf(ctx, "flinkProperties: %#v", flinkEnvVars)

	fc := &flinkOp.FlinkCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindFlinkCluster,
			APIVersion: flinkOp.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			// FIXME(regadas): find a proper way to set Name
			Name:        randStringBytesMaskImprSrcUnsafe(10),
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: flinkOp.FlinkClusterSpec{
			Image: flinkOp.ImageSpec{
				Name:       flinkImage,
				PullPolicy: corev1.PullAlways,
			},
			JobManager:      jobManager,
			TaskManager:     taskManager,
			Job:             &job,
			FlinkProperties: flinkProperties,
		},
	}

	return fc, nil
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func (flinkResourceHandler) BuildIdentityResource(ctx context.Context, taskCtx pluginsCore.TaskExecutionMetadata) (k8s.Resource, error) {
	logger.Infof(ctx, "Flink: BuildIdentityResource %v", taskCtx)
	return &flinkOp.FlinkCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindFlinkCluster,
			APIVersion: flinkOp.GroupVersion.String(),
		},
	}, nil
}

func getEventInfoForFlink(fc *flinkOp.FlinkCluster) (*pluginsCore.TaskInfo, error) {
	var taskLogs []*core.TaskLog
	customInfoMap := make(map[string]string)

	customInfo, err := utils.MarshalObjToStruct(customInfoMap)
	if err != nil {
		return nil, err
	}

	return &pluginsCore.TaskInfo{
		Logs:       taskLogs,
		CustomInfo: customInfo,
	}, nil
}

func (r flinkResourceHandler) GetTaskPhase(ctx context.Context, pluginContext k8s.PluginContext, resource k8s.Resource) (pluginsCore.PhaseInfo, error) {
	logger.Info(ctx, "Flink: GetTaskPhase")

	app := resource.(*flinkOp.FlinkCluster)
	info, err := getEventInfoForFlink(app)
	if err != nil {
		return pluginsCore.PhaseInfoUndefined, err
	}

	occurredAt := time.Now()
	// TODO(regadas): we should probably cover cluster state
	// case flinkOp.JobStatePending, flinkOp.JobStateUpdating:
	//  return pluginsCore.PhaseInfoQueued(occurredAt, pluginsCore.DefaultPhaseVersion, "job queued"), nil
	// case flinkOp.JobStateRunning:
	// 	return pluginsCore.PhaseInfoInitializing(occurredAt, pluginsCore.DefaultPhaseVersion, "job submitted", info), nil
	// case flinkOp.JobStateFailed:
	// 	reason := fmt.Sprintln("Flink Job  Submission Failed with Error")
	// 	return pluginsCore.PhaseInfoRetryableFailure(errors.DownstreamSystemError, reason, info), nil
	// case flinkOp.JobStateCancelled:
	// 	reason := fmt.Sprintf("Flink Job cancelled")
	// 	return pluginsCore.PhaseInfoRetryableFailure(errors.DownstreamSystemError, reason, info), nil
	// case flinkOp.JobStateSucceeded:
	// 	return pluginsCore.PhaseInfoSuccess(info), nil

	// FIXME(regadas):💣

	logger.Infof(ctx, "State %s", app.Status.State)
	switch app.Status.State {
	case flinkOp.ClusterStateCreating, flinkOp.ClusterStateReconciling, flinkOp.ClusterStateUpdating:
		return pluginsCore.PhaseInfoQueued(occurredAt, pluginsCore.DefaultPhaseVersion, "job queued"), nil
	case flinkOp.ClusterStateRunning:
		return pluginsCore.PhaseInfoInitializing(occurredAt, pluginsCore.DefaultPhaseVersion, "job submitted", info), nil
	case flinkOp.ClusterStateStopped, flinkOp.ClusterStateStopping, flinkOp.ClusterStatePartiallyStopped:
		return pluginsCore.PhaseInfoSuccess(info), nil
	}

	return pluginsCore.PhaseInfoRunning(pluginsCore.DefaultPhaseVersion, info), nil
}

func init() {
	if err := flinkOp.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}

	pluginmachinery.PluginRegistry().RegisterK8sPlugin(
		k8s.PluginEntry{
			ID:                  flinkTaskType,
			RegisteredTaskTypes: []pluginsCore.TaskType{flinkTaskType},
			ResourceToWatch:     &flinkOp.FlinkCluster{},
			Plugin:              flinkResourceHandler{},
			IsDefault:           false,
			DefaultForTaskTypes: []pluginsCore.TaskType{flinkTaskType},
		})
}
