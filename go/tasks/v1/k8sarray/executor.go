package k8sarray

import (
	"context"
	"fmt"
	pluginMachinery "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/v1"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/v1/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/v1/k8s"
	"github.com/lyft/flyteplugins/go/tasks/v1/errors"
	v1 "k8s.io/api/core/v1"
)

const executorName = "k8s-array-executor"
const arrayTaskType = "container_array"
const pluginStateVersion = 0

type Executor struct {
	catalogReader DoOnceWorkerInterface
	catalogWriter DoOnceWorkerInterface
	kubeClient core.KubeClient
}

func NewExecutor() Executor {
	return Executor{}
}

func (e Executor) GetID() string {
	return executorName
}

func (Executor) GetProperties() core.PluginProperties {
	return core.PluginProperties{}
}

func (e Executor) Handle(ctx context.Context, tCtx core.TaskExecutionContext) (core.Transition, error) {

	pluginState := State{}

	if _, err := tCtx.PluginStateReader().Get(&pluginState); err != nil {
		return core.UnknownTransition, errors.Wrapf(errors.CorruptedPluginState, err, "Failed to read unmarshal custom state")
	}

	switch pluginState.currentPhase {
	case NotStarted:
		taskTemplate, err := tCtx.TaskReader().Read(ctx)
		if err != nil {
			return core.UnknownTransition, err
		}
		taskTemplate.GetCustom()

		for item := range tCtx.TaskReader

	case MappingFileCreated:
	case JobSubmitted:
	case JobsFinished:
	}


}

func (Executor) Abort(ctx context.Context, tCtx core.TaskExecutionContext) error {
	panic("implement me")
}

func (Executor) Finalize(ctx context.Context, tCtx core.TaskExecutionContext) error {
	panic("implement me")
}

func init() {
	pluginMachinery.PluginRegistry().RegisterCorePlugin(
		core.PluginEntry{
			ID:                  executorName,
			RegisteredTaskTypes: []core.TaskType{arrayTaskType},
			LoadPlugin:          GetNewExecutorPlugin,
			IsDefault:           false,
		})
}

func GetNewExecutorPlugin(ctx context.Context, iCtx core.SetupContext) (core.Plugin, error) {
	return NewExecutor(), nil
}