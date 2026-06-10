package teamsettings

type TeamFieldValues struct {
	DefaultValue *string                    `json:"defaultValue"`
	Values       *[]TeamFieldValueReference `json:"values,omitempty"`
}

type TeamFieldValueReference struct {
	Value           *string `json:"value,omitempty"`
	IncludeChildren *bool   `json:"includeChildren,omitempty"`
}

type GetTeamFieldValuesArgs struct {
	Project *string
	Team    *string
}

type UpdateTeamFieldValuesArgs struct {
	Project         *string
	Team            *string
	TeamFieldValues *TeamFieldValues
}
