//go:build (all || core || resource_team) && !exclude_resource_team
// +build all core resource_team
// +build !exclude_resource_team

package core

// Tests for the waitForTeamStateChange convergence loop. The service can
// normalize name and description, seed its own ACEs, report duplicate
// identities and delay membership reads, which made team creation hang for
// 30 minutes in issue #1582.

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/dashboard"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type waitTestClients struct {
	clients         *client.AggregatedClient
	coreClient      *azdosdkmocks.MockCoreClient
	identityClient  *azdosdkmocks.MockIdentityClient
	securityClient  *azdosdkmocks.MockSecurityClient
	dashboardClient *azdosdkmocks.MockDashboardClient
}

func newWaitTestClients(t *testing.T, ctrl *gomock.Controller) *waitTestClients {
	t.Helper()
	c := &waitTestClients{
		coreClient:      azdosdkmocks.NewMockCoreClient(ctrl),
		identityClient:  azdosdkmocks.NewMockIdentityClient(ctrl),
		securityClient:  azdosdkmocks.NewMockSecurityClient(ctrl),
		dashboardClient: azdosdkmocks.NewMockDashboardClient(ctrl),
	}
	c.clients = &client.AggregatedClient{
		CoreClient:      c.coreClient,
		IdentityClient:  c.identityClient,
		SecurityClient:  c.securityClient,
		DashboardClient: c.dashboardClient,
		Ctx:             context.Background(),
	}
	return c
}

func newWaitTestResourceData(t *testing.T, projectID string) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, ResourceTeam().Schema, nil)
	d.Set("project_id", projectID)
	return d
}

func mockWaitTestTeam(c *waitTestClients, projectUUID, teamUUID uuid.UUID, storedName, storedDescription *string) {
	c.coreClient.
		EXPECT().
		GetTeam(c.clients.Ctx, gomock.Any()).
		Return(&core.WebApiTeam{
			Id:          &teamUUID,
			Name:        storedName,
			Description: storedDescription,
			ProjectId:   &projectUUID,
		}, nil).
		AnyTimes()
}

func mockWaitTestDashboards(c *waitTestClients, dashboards *[]dashboard.Dashboard) {
	c.dashboardClient.
		EXPECT().
		GetDashboardsByProject(c.clients.Ctx, gomock.Any()).
		Return(dashboards, nil).
		AnyTimes()
}

// Control: when the service returns exactly what was configured, the loop converges.
func TestTeam_WaitForStateChange_Synced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := newWaitTestClients(t, ctrl)
	projectUUID := uuid.New()
	teamUUID := uuid.New()
	teamName := "Squad X"
	teamDescription := "[Terraform Managed] squad"

	d := newWaitTestResourceData(t, projectUUID.String())
	mockWaitTestTeam(c, projectUUID, teamUUID, &teamName, &teamDescription)
	mockWaitTestDashboards(c, &[]dashboard.Dashboard{{}})

	err := waitForTeamStateChange(d, c.clients, projectUUID.String(), teamUUID.String(), &teamName, &teamDescription, nil, nil, nil, 2*time.Minute)
	require.NoError(t, err)
}

// The service trims surrounding whitespace when it stores name and description.
func TestTeam_WaitForStateChange_Converges_WhenServiceNormalizesNameAndDescription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := newWaitTestClients(t, ctrl)
	projectUUID := uuid.New()
	teamUUID := uuid.New()

	sentName := " Squad X "
	storedName := "Squad X"
	sentDescription := "[Terraform Managed] squad "
	storedDescription := "[Terraform Managed] squad"

	d := newWaitTestResourceData(t, projectUUID.String())
	mockWaitTestTeam(c, projectUUID, teamUUID, &storedName, &storedDescription)
	mockWaitTestDashboards(c, &[]dashboard.Dashboard{{}})

	err := waitForTeamStateChange(d, c.clients, projectUUID.String(), teamUUID.String(), &sentName, &sentDescription, nil, nil, nil, 2*time.Minute)
	require.NoError(t, err)
}

// A missing description in the response used to panic here.
func TestTeam_WaitForStateChange_Converges_WhenResponseOmitsDescription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := newWaitTestClients(t, ctrl)
	projectUUID := uuid.New()
	teamUUID := uuid.New()
	teamName := "Squad X"
	sentDescription := "[Terraform Managed] squad"

	d := newWaitTestResourceData(t, projectUUID.String())
	mockWaitTestTeam(c, projectUUID, teamUUID, &teamName, nil)
	mockWaitTestDashboards(c, &[]dashboard.Dashboard{{}})

	err := waitForTeamStateChange(d, c.clients, projectUUID.String(), teamUUID.String(), &teamName, &sentDescription, nil, nil, nil, 2*time.Minute)
	require.NoError(t, err)
}

// The service can seed its own ACEs on the team token; a count comparison never matched.
func TestTeam_WaitForStateChange_Converges_WhenACLContainsExtraAdministratorACE(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := newWaitTestClients(t, ctrl)
	projectUUID := uuid.New()
	teamUUID := uuid.New()
	teamName := "Squad X"
	teamDescription := "[Terraform Managed] squad"

	d := newWaitTestResourceData(t, projectUUID.String())
	mockWaitTestTeam(c, projectUUID, teamUUID, &teamName, &teamDescription)
	mockWaitTestDashboards(c, &[]dashboard.Dashboard{{}})

	configuredAdmins := schema.NewSet(schema.HashString, []interface{}{"aad.admin1"})

	nsID := uuid.UUID(securityhelper.SecurityNamespaceIDValues.Identity)
	c.securityClient.
		EXPECT().
		QuerySecurityNamespaces(c.clients.Ctx, security.QuerySecurityNamespacesArgs{
			SecurityNamespaceId: &nsID,
		}).
		Return(&[]security.SecurityNamespaceDescription{
			{
				Actions: &[]security.ActionDefinition{
					{Bit: converter.Int(1), Name: converter.String("Read")},
					{Bit: converter.Int(2), Name: converter.String("Write")},
					{Bit: converter.Int(4), Name: converter.String("Delete")},
					{Bit: converter.Int(8), Name: converter.String("ManageMembership")},
					{Bit: converter.Int(16), Name: converter.String("CreateScope")},
				},
			},
		}, nil).
		AnyTimes()

	token := projectUUID.String() + "\\" + teamUUID.String()
	fullBits := 1 | 2 | 4 | 8 | 16

	// ACL contains the configured admin ACE plus one service-seeded extra ACE,
	// both with the full permission bit set
	c.securityClient.
		EXPECT().
		QueryAccessControlLists(c.clients.Ctx, security.QueryAccessControlListsArgs{
			SecurityNamespaceId: &nsID,
			Token:               &token,
			IncludeExtendedInfo: converter.Bool(true),
		}).
		Return(&[]security.AccessControlList{
			{
				AcesDictionary: &map[string]security.AccessControlEntry{
					"ace-configured-admin": {Allow: converter.Int(fullBits), Descriptor: converter.String("ace-configured-admin")},
					"ace-extra-admin":      {Allow: converter.Int(fullBits), Descriptor: converter.String("ace-extra-admin")},
				},
			},
		}, nil).
		AnyTimes()

	c.identityClient.
		EXPECT().
		ReadIdentities(c.clients.Ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, args identity.ReadIdentitiesArgs) (*[]identity.Identity, error) {
			require.NotNil(t, args.Descriptors)
			require.Contains(t, *args.Descriptors, "ace-configured-admin")
			require.Contains(t, *args.Descriptors, "ace-extra-admin")
			return &[]identity.Identity{
				{SubjectDescriptor: converter.String("aad.admin1")},
				{SubjectDescriptor: converter.String("aad.extra-admin")},
			}, nil
		}).
		AnyTimes()

	err := waitForTeamStateChange(d, c.clients, projectUUID.String(), teamUUID.String(), &teamName, &teamDescription, nil, configuredAdmins, nil, 2*time.Minute)
	require.NoError(t, err)
}

// Team membership can contain members the config does not list.
func TestTeam_WaitForStateChange_Converges_WhenTeamHasExtraMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := newWaitTestClients(t, ctrl)
	projectUUID := uuid.New()
	teamUUID := uuid.New()
	teamName := "Squad X"
	teamDescription := "[Terraform Managed] squad"

	d := newWaitTestResourceData(t, projectUUID.String())
	mockWaitTestTeam(c, projectUUID, teamUUID, &teamName, &teamDescription)
	mockWaitTestDashboards(c, &[]dashboard.Dashboard{{}})

	configuredMembers := schema.NewSet(schema.HashString, []interface{}{"aadgp.group1"})

	c.identityClient.
		EXPECT().
		ReadMembers(c.clients.Ctx, identity.ReadMembersArgs{
			ContainerId: converter.String(teamUUID.String()),
		}).
		Return(&[]string{"aadgp.group1", "aadgp.group-added-via-ui"}, nil).
		AnyTimes()

	c.identityClient.
		EXPECT().
		ReadIdentities(c.clients.Ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, args identity.ReadIdentitiesArgs) (*[]identity.Identity, error) {
			require.NotNil(t, args.Descriptors)
			identities := make([]identity.Identity, 0)
			for _, descriptor := range strings.Split(*args.Descriptors, ",") {
				identities = append(identities, identity.Identity{SubjectDescriptor: converter.String(strings.TrimSpace(descriptor))})
			}
			return &identities, nil
		}).
		AnyTimes()

	err := waitForTeamStateChange(d, c.clients, projectUUID.String(), teamUUID.String(), &teamName, &teamDescription, configuredMembers, nil, nil, 2*time.Minute)
	require.NoError(t, err)
}

// A nil dashboard list used to panic here.
func TestTeam_WaitForStateChange_DoesNotPanic_WhenDashboardsListIsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := newWaitTestClients(t, ctrl)
	projectUUID := uuid.New()
	teamUUID := uuid.New()
	teamName := "Squad X"
	teamDescription := "[Terraform Managed] squad"

	d := newWaitTestResourceData(t, projectUUID.String())
	mockWaitTestTeam(c, projectUUID, teamUUID, &teamName, &teamDescription)
	mockWaitTestDashboards(c, nil)

	err := waitForTeamStateChange(d, c.clients, projectUUID.String(), teamUUID.String(), &teamName, &teamDescription, nil, nil, nil, 2*time.Minute)
	require.NoError(t, err)
}

// A member that never shows up cannot converge; fail instead of hanging for 30 minutes.
func TestTeam_WaitForStateChange_TimesOut_WhenConfiguredMemberNeverAppears(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := newWaitTestClients(t, ctrl)
	projectUUID := uuid.New()
	teamUUID := uuid.New()
	teamName := "Squad X"
	teamDescription := "[Terraform Managed] squad"

	d := newWaitTestResourceData(t, projectUUID.String())
	mockWaitTestTeam(c, projectUUID, teamUUID, &teamName, &teamDescription)
	mockWaitTestDashboards(c, &[]dashboard.Dashboard{{}})

	configuredMembers := schema.NewSet(schema.HashString, []interface{}{"aadgp.group1", "aadgp.group2"})

	// identity service only ever reports group1 as member; group2 was silently dropped
	c.identityClient.
		EXPECT().
		ReadMembers(c.clients.Ctx, identity.ReadMembersArgs{
			ContainerId: converter.String(teamUUID.String()),
		}).
		Return(&[]string{"aadgp.group1"}, nil).
		AnyTimes()

	c.identityClient.
		EXPECT().
		ReadIdentities(c.clients.Ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, args identity.ReadIdentitiesArgs) (*[]identity.Identity, error) {
			require.NotNil(t, args.Descriptors)
			identities := make([]identity.Identity, 0)
			for _, descriptor := range strings.Split(*args.Descriptors, ",") {
				identities = append(identities, identity.Identity{SubjectDescriptor: converter.String(strings.TrimSpace(descriptor))})
			}
			return &identities, nil
		}).
		AnyTimes()

	err := waitForTeamStateChange(d, c.clients, projectUUID.String(), teamUUID.String(), &teamName, &teamDescription, configuredMembers, nil, nil, 12*time.Second)
	require.Error(t, err)
	require.Contains(t, err.Error(), "timeout while waiting for state to become 'Synched' (last state: 'Waiting'")
}
