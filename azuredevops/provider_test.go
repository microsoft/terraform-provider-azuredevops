package azuredevops_test

import (
	"fmt"
	"testing"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops"
	"github.com/stretchr/testify/require"
)

func TestProvider_HasChildResources(t *testing.T) {
	expectedResources := []string{
		"azuredevops_agent_pool",
		"azuredevops_agent_queue",
		"azuredevops_area_permissions",
		"azuredevops_branch_policy_auto_reviewers",
		"azuredevops_branch_policy_build_validation",
		"azuredevops_branch_policy_comment_resolution",
		"azuredevops_branch_policy_merge_types",
		"azuredevops_branch_policy_min_reviewers",
		"azuredevops_branch_policy_status_check",
		"azuredevops_branch_policy_work_item_linking",
		"azuredevops_build_definition",
		"azuredevops_build_definition_permissions",
		"azuredevops_build_folder",
		"azuredevops_build_folder_permissions",
		"azuredevops_check_approval",
		"azuredevops_check_branch_control",
		"azuredevops_check_business_hours",
		"azuredevops_check_exclusive_lock",
		"azuredevops_check_required_template",
		"azuredevops_check_rest_api",
		"azuredevops_dashboard",
		"azuredevops_elastic_pool",
		"azuredevops_environment",
		"azuredevops_environment_resource_kubernetes",
		"azuredevops_extension",
		"azuredevops_feed",
		"azuredevops_feed_permission",
		"azuredevops_feed_retention_policy",
		"azuredevops_git_permissions",
		"azuredevops_git_repository",
		"azuredevops_git_repository_branch",
		"azuredevops_git_repository_file",
		"azuredevops_group",
		"azuredevops_group_entitlement",
		"azuredevops_group_membership",
		"azuredevops_iteration_permissions",
		"azuredevops_library_permissions",
		"azuredevops_pipeline_authorization",
		"azuredevops_project",
		"azuredevops_project_features",
		"azuredevops_project_permissions",
		"azuredevops_project_pipeline_settings",
		"azuredevops_project_tags",
		"azuredevops_repository_policy_author_email_pattern",
		"azuredevops_repository_policy_case_enforcement",
		"azuredevops_repository_policy_check_credentials",
		"azuredevops_repository_policy_file_path_pattern",
		"azuredevops_repository_policy_max_file_size",
		"azuredevops_repository_policy_max_path_length",
		"azuredevops_repository_policy_reserved_names",
		"azuredevops_resource_authorization",
		"azuredevops_securityrole_assignment",
		"azuredevops_serviceendpoint_argocd",
		"azuredevops_serviceendpoint_artifactory",
		"azuredevops_serviceendpoint_aws",
		"azuredevops_serviceendpoint_azure_service_bus",
		"azuredevops_serviceendpoint_azurecr",
		"azuredevops_serviceendpoint_azuredevops",
		"azuredevops_serviceendpoint_azurerm",
		"azuredevops_serviceendpoint_bitbucket",
		"azuredevops_serviceendpoint_black_duck",
		"azuredevops_serviceendpoint_checkmarx_one",
		"azuredevops_serviceendpoint_checkmarx_sca",
		"azuredevops_serviceendpoint_checkmarx_sast",
		"azuredevops_serviceendpoint_dockerregistry",
		"azuredevops_serviceendpoint_dynamics_lifecycle_services",
		"azuredevops_serviceendpoint_externaltfs",
		"azuredevops_serviceendpoint_gcp_terraform",
		"azuredevops_serviceendpoint_generic",
		"azuredevops_serviceendpoint_generic_git",
		"azuredevops_serviceendpoint_github",
		"azuredevops_serviceendpoint_github_enterprise",
		"azuredevops_serviceendpoint_gitlab",
		"azuredevops_serviceendpoint_incomingwebhook",
		"azuredevops_serviceendpoint_jenkins",
		"azuredevops_serviceendpoint_jfrog_artifactory_v2",
		"azuredevops_serviceendpoint_jfrog_distribution_v2",
		"azuredevops_serviceendpoint_jfrog_platform_v2",
		"azuredevops_serviceendpoint_jfrog_xray_v2",
		"azuredevops_serviceendpoint_kubernetes",
		"azuredevops_serviceendpoint_maven",
		"azuredevops_serviceendpoint_nexus",
		"azuredevops_serviceendpoint_npm",
		"azuredevops_serviceendpoint_nuget",
		"azuredevops_serviceendpoint_octopusdeploy",
		"azuredevops_serviceendpoint_openshift",
		"azuredevops_serviceendpoint_permissions",
		"azuredevops_serviceendpoint_runpipeline",
		"azuredevops_serviceendpoint_servicefabric",
		"azuredevops_serviceendpoint_snyk",
		"azuredevops_serviceendpoint_sonarcloud",
		"azuredevops_serviceendpoint_sonarqube",
		"azuredevops_serviceendpoint_ssh",
		"azuredevops_serviceendpoint_visualstudiomarketplace",
		"azuredevops_servicehook_permissions",
		"azuredevops_servicehook_storage_queue_pipelines",
		"azuredevops_service_principal_entitlement",
		"azuredevops_tagging_permissions",
		"azuredevops_team",
		"azuredevops_team_administrators",
		"azuredevops_team_members",
		"azuredevops_user_entitlement",
		"azuredevops_variable_group",
		"azuredevops_variable_group_permissions",
		"azuredevops_wiki",
		"azuredevops_wiki_page",
		"azuredevops_workitem",
		"azuredevops_workitemquery",
		"azuredevops_workitemquery_folder",
		"azuredevops_workitemquery_permissions",
	}

	resources := azuredevops.Provider().ResourcesMap

	for _, resource := range expectedResources {
		require.Contains(t, resources, resource, fmt.Sprintf("An expected resource (%s) was not registered", resource))
		require.NotNil(t, resources[resource], "A resource cannot have a nil schema")
	}
	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")
}

func TestProvider_HasChildDataSources(t *testing.T) {
	expectedDataSources := []string{
		"azuredevops_agent_pool",
		"azuredevops_agent_pools",
		"azuredevops_agent_queue",
		"azuredevops_area",
		"azuredevops_build_definition",
		"azuredevops_client_config",
		"azuredevops_descriptor",
		"azuredevops_environment",
		"azuredevops_feed",
		"azuredevops_git_repositories",
		"azuredevops_git_repository",
		"azuredevops_git_repository_file",
		"azuredevops_group",
		"azuredevops_group_membership",
		"azuredevops_groups",
		"azuredevops_identity_group",
		"azuredevops_identity_groups",
		"azuredevops_identity_user",
		"azuredevops_iteration",
		"azuredevops_project",
		"azuredevops_projects",
		"azuredevops_securityrole_definitions",
		"azuredevops_serviceendpoint_azurecr",
		"azuredevops_serviceendpoint_azurerm",
		"azuredevops_serviceendpoint_bitbucket",
		"azuredevops_serviceendpoint_dockerregistry",
		"azuredevops_serviceendpoint_generic",
		"azuredevops_serviceendpoint_github",
		"azuredevops_serviceendpoint_npm",
		"azuredevops_serviceendpoint_sonarcloud",
		"azuredevops_storage_key",
		"azuredevops_service_principal",
		"azuredevops_team",
		"azuredevops_teams",
		"azuredevops_user",
		"azuredevops_users",
		"azuredevops_variable_group",
	}

	dataSources := azuredevops.Provider().DataSourcesMap

	for _, resource := range expectedDataSources {
		require.Contains(t, dataSources, resource, "An expected data source was not registered")
		require.NotNil(t, dataSources[resource], "A data source cannot have a nil schema")
	}
	require.Equal(t, len(expectedDataSources), len(dataSources), "There are an unexpected number of registered data sources")
}

func TestProvider_SchemaIsValid(t *testing.T) {
	type testParams struct {
		name      string
		required  bool
		sensitive bool
	}

	tests := []testParams{
		{"org_service_url", false, false},
		{"personal_access_token", false, true},

		{"client_id", false, false},
		{"client_id_file_path", false, false},
		{"tenant_id", false, false},
		{"auxiliary_tenant_ids", false, false},
		{"client_certificate_path", false, false},
		{"client_certificate", false, true},
		{"client_certificate_password", false, true},
		{"client_secret", false, true},
		{"client_secret_path", false, false},
		{"oidc_request_token", false, false},
		{"oidc_request_url", false, false},
		{"oidc_token", false, true},
		{"oidc_token_file_path", false, false},
		{"oidc_azure_service_connection_id", false, false},
		{"use_oidc", false, false},
		{"use_msi", false, false},
		{"use_cli", false, false},
	}

	schema := azuredevops.Provider().Schema
	require.Equal(t, len(tests), len(schema), "There are an unexpected number of properties in the schema")

	for _, test := range tests {
		require.Contains(t, schema, test.name, "An expected property was not found in the schema")
		require.NotNil(t, schema[test.name], "A property in the schema cannot have a nil value")
		require.Equal(t, test.sensitive, schema[test.name].Sensitive, "A property in the schema has an incorrect sensitivity value")
		require.Equal(t, test.required, schema[test.name].Required, "A property in the schema has an incorrect required value")
	}
}
