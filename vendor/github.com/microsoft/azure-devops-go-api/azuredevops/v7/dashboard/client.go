// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package dashboard

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"net/http"
	"net/url"
)

var ResourceAreaId, _ = uuid.Parse("31c84e0a-3ece-48fd-a29d-100849af99ba")

type Client interface {
	// [Preview API] Create the supplied dashboard.
	CreateDashboard(context.Context, CreateDashboardArgs) (*Dashboard, error)
	// [Preview API] Create a widget on the specified dashboard.
	CreateWidget(context.Context, CreateWidgetArgs) (*Widget, error)
	// [Preview API] Delete a dashboard given its ID. This also deletes the widgets associated with this dashboard.
	DeleteDashboard(context.Context, DeleteDashboardArgs) error
	// [Preview API] Delete the specified widget.
	DeleteWidget(context.Context, DeleteWidgetArgs) (*Dashboard, error)
	// [Preview API] Get a dashboard by its ID.
	GetDashboard(context.Context, GetDashboardArgs) (*Dashboard, error)
	// [Preview API] Get a list of dashboards under a project.
	GetDashboardsByProject(context.Context, GetDashboardsByProjectArgs) (*[]Dashboard, error)
	// [Preview API] Get the current state of the specified widget.
	GetWidget(context.Context, GetWidgetArgs) (*Widget, error)
	// [Preview API] Get the widget metadata satisfying the specified contribution ID.
	GetWidgetMetadata(context.Context, GetWidgetMetadataArgs) (*WidgetMetadataResponse, error)
	// [Preview API] Get widgets contained on the specified dashboard.
	GetWidgets(context.Context, GetWidgetsArgs) (*WidgetsVersionedList, error)
	// [Preview API] Get all available widget metadata in alphabetical order, including widgets marked with isVisibleFromCatalog == false.
	GetWidgetTypes(context.Context, GetWidgetTypesArgs) (*WidgetTypesResponse, error)
	// [Preview API] Replace configuration for the specified dashboard. Replaces Widget list on Dashboard, only if property is supplied.
	ReplaceDashboard(context.Context, ReplaceDashboardArgs) (*Dashboard, error)
	// [Preview API] Update the name and position of dashboards in the supplied group, and remove omitted dashboards. Does not modify dashboard content.
	ReplaceDashboards(context.Context, ReplaceDashboardsArgs) (*DashboardGroup, error)
	// [Preview API] Override the  state of the specified widget.
	ReplaceWidget(context.Context, ReplaceWidgetArgs) (*Widget, error)
	// [Preview API] Replace the widgets on specified dashboard with the supplied widgets.
	ReplaceWidgets(context.Context, ReplaceWidgetsArgs) (*WidgetsVersionedList, error)
	// [Preview API] Perform a partial update of the specified widget.
	UpdateWidget(context.Context, UpdateWidgetArgs) (*Widget, error)
	// [Preview API] Update the supplied widgets on the dashboard using supplied state. State of existing Widgets not passed in the widget list is preserved.
	UpdateWidgets(context.Context, UpdateWidgetsArgs) (*WidgetsVersionedList, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) (Client, error) {
	client, err := connection.GetClientByResourceAreaId(ctx, ResourceAreaId)
	if err != nil {
		return nil, err
	}
	return &ClientImpl{
		Client: *client,
	}, nil
}

// [Preview API] Create the supplied dashboard.
func (client *ClientImpl) CreateDashboard(ctx context.Context, args CreateDashboardArgs) (*Dashboard, error) {
	if args.Dashboard == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Dashboard"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}

	body, marshalErr := json.Marshal(*args.Dashboard)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("454b3e51-2e6e-48d4-ad81-978154089351")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.3", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Dashboard
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateDashboard function
type CreateDashboardArgs struct {
	// (required) The initial state of the dashboard
	Dashboard *Dashboard
	// (required) Project ID or project name
	Project *string
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Create a widget on the specified dashboard.
func (client *ClientImpl) CreateWidget(ctx context.Context, args CreateWidgetArgs) (*Widget, error) {
	if args.Widget == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Widget"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()

	body, marshalErr := json.Marshal(*args.Widget)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("bdcff53a-8355-4172-a00a-40497ea23afc")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.2", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Widget
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateWidget function
type CreateWidgetArgs struct {
	// (required) State of the widget to add
	Widget *Widget
	// (required) Project ID or project name
	Project *string
	// (required) ID of dashboard the widget will be added to.
	DashboardId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Delete a dashboard given its ID. This also deletes the widgets associated with this dashboard.
func (client *ClientImpl) DeleteDashboard(ctx context.Context, args DeleteDashboardArgs) error {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()

	locationId, _ := uuid.Parse("454b3e51-2e6e-48d4-ad81-978154089351")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.3", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteDashboard function
type DeleteDashboardArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) ID of the dashboard to delete.
	DashboardId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Delete the specified widget.
func (client *ClientImpl) DeleteWidget(ctx context.Context, args DeleteWidgetArgs) (*Dashboard, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()
	if args.WidgetId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WidgetId"}
	}
	routeValues["widgetId"] = (*args.WidgetId).String()

	locationId, _ := uuid.Parse("bdcff53a-8355-4172-a00a-40497ea23afc")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.2", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Dashboard
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the DeleteWidget function
type DeleteWidgetArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) ID of the dashboard containing the widget.
	DashboardId *uuid.UUID
	// (required) ID of the widget to update.
	WidgetId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Get a dashboard by its ID.
func (client *ClientImpl) GetDashboard(ctx context.Context, args GetDashboardArgs) (*Dashboard, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()

	locationId, _ := uuid.Parse("454b3e51-2e6e-48d4-ad81-978154089351")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.3", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Dashboard
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetDashboard function
type GetDashboardArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required)
	DashboardId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Get a list of dashboards under a project.
func (client *ClientImpl) GetDashboardsByProject(ctx context.Context, args GetDashboardsByProjectArgs) (*[]Dashboard, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}

	locationId, _ := uuid.Parse("454b3e51-2e6e-48d4-ad81-978154089351")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.3", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Dashboard
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetDashboardsByProject function
type GetDashboardsByProjectArgs struct {
	// (required) Project ID or project name
	Project *string
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Get the current state of the specified widget.
func (client *ClientImpl) GetWidget(ctx context.Context, args GetWidgetArgs) (*Widget, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()
	if args.WidgetId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WidgetId"}
	}
	routeValues["widgetId"] = (*args.WidgetId).String()

	locationId, _ := uuid.Parse("bdcff53a-8355-4172-a00a-40497ea23afc")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.2", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Widget
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWidget function
type GetWidgetArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) ID of the dashboard containing the widget.
	DashboardId *uuid.UUID
	// (required) ID of the widget to read.
	WidgetId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Get the widget metadata satisfying the specified contribution ID.
func (client *ClientImpl) GetWidgetMetadata(ctx context.Context, args GetWidgetMetadataArgs) (*WidgetMetadataResponse, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.ContributionId == nil || *args.ContributionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ContributionId"}
	}
	routeValues["contributionId"] = *args.ContributionId

	locationId, _ := uuid.Parse("6b3628d3-e96f-4fc7-b176-50240b03b515")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WidgetMetadataResponse
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWidgetMetadata function
type GetWidgetMetadataArgs struct {
	// (required) The ID of Contribution for the Widget
	ContributionId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Get widgets contained on the specified dashboard.
func (client *ClientImpl) GetWidgets(ctx context.Context, args GetWidgetsArgs) (*WidgetsVersionedList, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()

	additionalHeaders := make(map[string]string)
	if args.ETag != nil {
		additionalHeaders["ETag"] = *args.ETag
	}
	locationId, _ := uuid.Parse("bdcff53a-8355-4172-a00a-40497ea23afc")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.2", routeValues, nil, nil, "", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseBodyValue []Widget
	err = client.Client.UnmarshalCollectionBody(resp, &responseBodyValue)

	var responseValue *WidgetsVersionedList
	if err == nil {
		responseValue = &WidgetsVersionedList{
			Widgets: &responseBodyValue,
			ETag:    &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the GetWidgets function
type GetWidgetsArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) ID of the dashboard to read.
	DashboardId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
	// (optional) Dashboard Widgets Version
	ETag *string
}

// [Preview API] Get all available widget metadata in alphabetical order, including widgets marked with isVisibleFromCatalog == false.
func (client *ClientImpl) GetWidgetTypes(ctx context.Context, args GetWidgetTypesArgs) (*WidgetTypesResponse, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.Scope == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "scope"}
	}
	queryParams.Add("$scope", string(*args.Scope))
	locationId, _ := uuid.Parse("6b3628d3-e96f-4fc7-b176-50240b03b515")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WidgetTypesResponse
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWidgetTypes function
type GetWidgetTypesArgs struct {
	// (required)
	Scope *WidgetScope
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Replace configuration for the specified dashboard. Replaces Widget list on Dashboard, only if property is supplied.
func (client *ClientImpl) ReplaceDashboard(ctx context.Context, args ReplaceDashboardArgs) (*Dashboard, error) {
	if args.Dashboard == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Dashboard"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()

	body, marshalErr := json.Marshal(*args.Dashboard)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("454b3e51-2e6e-48d4-ad81-978154089351")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.3", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Dashboard
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReplaceDashboard function
type ReplaceDashboardArgs struct {
	// (required) The Configuration of the dashboard to replace.
	Dashboard *Dashboard
	// (required) Project ID or project name
	Project *string
	// (required) ID of the dashboard to replace.
	DashboardId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Update the name and position of dashboards in the supplied group, and remove omitted dashboards. Does not modify dashboard content.
func (client *ClientImpl) ReplaceDashboards(ctx context.Context, args ReplaceDashboardsArgs) (*DashboardGroup, error) {
	if args.Group == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Group"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}

	body, marshalErr := json.Marshal(*args.Group)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("454b3e51-2e6e-48d4-ad81-978154089351")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.3", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue DashboardGroup
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReplaceDashboards function
type ReplaceDashboardsArgs struct {
	// (required)
	Group *DashboardGroup
	// (required) Project ID or project name
	Project *string
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Override the  state of the specified widget.
func (client *ClientImpl) ReplaceWidget(ctx context.Context, args ReplaceWidgetArgs) (*Widget, error) {
	if args.Widget == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Widget"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()
	if args.WidgetId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WidgetId"}
	}
	routeValues["widgetId"] = (*args.WidgetId).String()

	body, marshalErr := json.Marshal(*args.Widget)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("bdcff53a-8355-4172-a00a-40497ea23afc")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.2", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Widget
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReplaceWidget function
type ReplaceWidgetArgs struct {
	// (required) State to be written for the widget.
	Widget *Widget
	// (required) Project ID or project name
	Project *string
	// (required) ID of the dashboard containing the widget.
	DashboardId *uuid.UUID
	// (required) ID of the widget to update.
	WidgetId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Replace the widgets on specified dashboard with the supplied widgets.
func (client *ClientImpl) ReplaceWidgets(ctx context.Context, args ReplaceWidgetsArgs) (*WidgetsVersionedList, error) {
	if args.Widgets == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Widgets"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()

	additionalHeaders := make(map[string]string)
	if args.ETag != nil {
		additionalHeaders["ETag"] = *args.ETag
	}
	body, marshalErr := json.Marshal(*args.Widgets)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("bdcff53a-8355-4172-a00a-40497ea23afc")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.2", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseBodyValue []Widget
	err = client.Client.UnmarshalCollectionBody(resp, &responseBodyValue)

	var responseValue *WidgetsVersionedList
	if err == nil {
		responseValue = &WidgetsVersionedList{
			Widgets: &responseBodyValue,
			ETag:    &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the ReplaceWidgets function
type ReplaceWidgetsArgs struct {
	// (required) Revised state of widgets to store for the dashboard.
	Widgets *[]Widget
	// (required) Project ID or project name
	Project *string
	// (required) ID of the Dashboard to modify.
	DashboardId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
	// (optional) Dashboard Widgets Version
	ETag *string
}

// [Preview API] Perform a partial update of the specified widget.
func (client *ClientImpl) UpdateWidget(ctx context.Context, args UpdateWidgetArgs) (*Widget, error) {
	if args.Widget == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Widget"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()
	if args.WidgetId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WidgetId"}
	}
	routeValues["widgetId"] = (*args.WidgetId).String()

	body, marshalErr := json.Marshal(*args.Widget)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("bdcff53a-8355-4172-a00a-40497ea23afc")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.2", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Widget
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateWidget function
type UpdateWidgetArgs struct {
	// (required) Description of the widget changes to apply. All non-null fields will be replaced.
	Widget *Widget
	// (required) Project ID or project name
	Project *string
	// (required) ID of the dashboard containing the widget.
	DashboardId *uuid.UUID
	// (required) ID of the widget to update.
	WidgetId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
}

// [Preview API] Update the supplied widgets on the dashboard using supplied state. State of existing Widgets not passed in the widget list is preserved.
func (client *ClientImpl) UpdateWidgets(ctx context.Context, args UpdateWidgetsArgs) (*WidgetsVersionedList, error) {
	if args.Widgets == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Widgets"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.DashboardId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DashboardId"}
	}
	routeValues["dashboardId"] = (*args.DashboardId).String()

	additionalHeaders := make(map[string]string)
	if args.ETag != nil {
		additionalHeaders["ETag"] = *args.ETag
	}
	body, marshalErr := json.Marshal(*args.Widgets)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("bdcff53a-8355-4172-a00a-40497ea23afc")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.2", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseBodyValue []Widget
	err = client.Client.UnmarshalCollectionBody(resp, &responseBodyValue)

	var responseValue *WidgetsVersionedList
	if err == nil {
		responseValue = &WidgetsVersionedList{
			Widgets: &responseBodyValue,
			ETag:    &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the UpdateWidgets function
type UpdateWidgetsArgs struct {
	// (required) The set of widget states to update on the dashboard.
	Widgets *[]Widget
	// (required) Project ID or project name
	Project *string
	// (required) ID of the Dashboard to modify.
	DashboardId *uuid.UUID
	// (optional) Team ID or team name
	Team *string
	// (optional) Dashboard Widgets Version
	ETag *string
}
