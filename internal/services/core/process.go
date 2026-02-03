package core

import (
	"context"
	"errors"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
)

func LookupProcess(ctx context.Context, client core.Client, f func(p core.Process) bool) (*core.Process, error) {
	processes, err := client.GetProcesses(ctx, core.GetProcessesArgs{})
	if err != nil {
		return nil, err
	}
	if processes == nil {
		return nil, errors.New("unexpected null processes")
	}

	for _, process := range *processes {
		if f(process) {
			return &process, nil
		}
	}

	return nil, errors.New("process not found")
}
