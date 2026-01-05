package core

type TeamProjectCapabilities struct {
	Versioncontrol  *TeamProjectCapabilitiesVersionControl  `json:"versioncontrol,omitempty"`
	ProcessTemplate *TeamProjectCapabilitiesProcessTemplate `json:"processTemplate,omitempty"`
}

type TeamProjectCapabilitiesVersionControl struct {
	SourceControlType *SourceControlTypes `json:"sourceControlType,omitempty"`
}

type TeamProjectCapabilitiesProcessTemplate struct {
	TemplateId *string `json:"templateTypeId,omitempty"`
}
