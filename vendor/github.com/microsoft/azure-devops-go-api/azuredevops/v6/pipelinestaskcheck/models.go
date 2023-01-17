// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package pipelinestaskcheck

import (
	"github.com/google/uuid"
)

// Config to facilitate task check
type TaskCheckConfig struct {
	DefinitionRef       *TaskCheckDefinitionReference `json:"definitionRef,omitempty"`
	DisplayName         *string                       `json:"displayName,omitempty"`
	Inputs              *map[string]string            `json:"inputs,omitempty"`
	LinkedVariableGroup *string                       `json:"linkedVariableGroup,omitempty"`
	RetryInterval       *int                          `json:"retryInterval,omitempty"`
}

type TaskCheckDefinitionReference struct {
	Id      *uuid.UUID `json:"id,omitempty"`
	Name    *string    `json:"name,omitempty"`
	Version *string    `json:"version,omitempty"`
}
