// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package wiki

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

// Defines a wiki repository which encapsulates the git repository backing the wiki.
type Wiki struct {
	// Wiki name.
	Name *string `json:"name,omitempty"`
	// ID of the project in which the wiki is to be created.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// The head commit associated with the git repository backing up the wiki.
	HeadCommit *string `json:"headCommit,omitempty"`
	// The ID of the wiki which is same as the ID of the Git repository that it is backed by.
	Id *uuid.UUID `json:"id,omitempty"`
	// The git repository that backs up the wiki.
	Repository *git.GitRepository `json:"repository,omitempty"`
}

// Defines properties for wiki attachment file.
type WikiAttachment struct {
	// Name of the wiki attachment file.
	Name *string `json:"name,omitempty"`
	// Path of the wiki attachment file.
	Path *string `json:"path,omitempty"`
}

// Response contract for the Wiki Attachments API
type WikiAttachmentResponse struct {
	// Defines properties for wiki attachment file.
	Attachment *WikiAttachment `json:"attachment,omitempty"`
	// Contains the list of ETag values from the response header of the attachments API call. The first item in the list contains the version of the wiki attachment.
	ETag *[]string `json:"eTag,omitempty"`
}

// Base wiki creation parameters.
type WikiCreateBaseParameters struct {
	// Folder path inside repository which is shown as Wiki. Not required for ProjectWiki type.
	MappedPath *string `json:"mappedPath,omitempty"`
	// Wiki name.
	Name *string `json:"name,omitempty"`
	// ID of the project in which the wiki is to be created.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// ID of the git repository that backs up the wiki. Not required for ProjectWiki type.
	RepositoryId *uuid.UUID `json:"repositoryId,omitempty"`
	// Type of the wiki.
	Type *WikiType `json:"type,omitempty"`
}

// Wiki creations parameters.
type WikiCreateParameters struct {
	// Wiki name.
	Name *string `json:"name,omitempty"`
	// ID of the project in which the wiki is to be created.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
}

// Wiki creation parameters.
type WikiCreateParametersV2 struct {
	// Folder path inside repository which is shown as Wiki. Not required for ProjectWiki type.
	MappedPath *string `json:"mappedPath,omitempty"`
	// Wiki name.
	Name *string `json:"name,omitempty"`
	// ID of the project in which the wiki is to be created.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// ID of the git repository that backs up the wiki. Not required for ProjectWiki type.
	RepositoryId *uuid.UUID `json:"repositoryId,omitempty"`
	// Type of the wiki.
	Type *WikiType `json:"type,omitempty"`
	// Version of the wiki. Not required for ProjectWiki type.
	Version *git.GitVersionDescriptor `json:"version,omitempty"`
}

// Defines a page in a wiki.
type WikiPage struct {
	// Content of the wiki page.
	Content *string `json:"content,omitempty"`
	// Path of the git item corresponding to the wiki page stored in the backing Git repository.
	GitItemPath *string `json:"gitItemPath,omitempty"`
	// When present, permanent identifier for the wiki page
	Id *int `json:"id,omitempty"`
	// True if a page is non-conforming, i.e. 1) if the name doesn't match page naming standards. 2) if the page does not have a valid entry in the appropriate order file.
	IsNonConformant *bool `json:"isNonConformant,omitempty"`
	// True if this page has subpages under its path.
	IsParentPage *bool `json:"isParentPage,omitempty"`
	// Order of the wiki page, relative to other pages in the same hierarchy level.
	Order *int `json:"order,omitempty"`
	// Path of the wiki page.
	Path *string `json:"path,omitempty"`
	// Remote web url to the wiki page.
	RemoteUrl *string `json:"remoteUrl,omitempty"`
	// List of subpages of the current page.
	SubPages *[]WikiPage `json:"subPages,omitempty"`
	// REST url for this wiki page.
	Url *string `json:"url,omitempty"`
}

// Contract encapsulating parameters for the page create or update operations.
type WikiPageCreateOrUpdateParameters struct {
	// Content of the wiki page.
	Content *string `json:"content,omitempty"`
}

// Defines a page with its metedata in a wiki.
type WikiPageDetail struct {
	// When present, permanent identifier for the wiki page
	Id *int `json:"id,omitempty"`
	// Path of the wiki page.
	Path *string `json:"path,omitempty"`
	// Path of the wiki page.
	ViewStats *[]WikiPageStat `json:"viewStats,omitempty"`
}

// Request contract for Wiki Page Move.
type WikiPageMove struct {
	// New order of the wiki page.
	NewOrder *int `json:"newOrder,omitempty"`
	// New path of the wiki page.
	NewPath *string `json:"newPath,omitempty"`
	// Current path of the wiki page.
	Path *string `json:"path,omitempty"`
	// Resultant page of this page move operation.
	Page *WikiPage `json:"page,omitempty"`
}

// Contract encapsulating parameters for the page move operation.
type WikiPageMoveParameters struct {
	// New order of the wiki page.
	NewOrder *int `json:"newOrder,omitempty"`
	// New path of the wiki page.
	NewPath *string `json:"newPath,omitempty"`
	// Current path of the wiki page.
	Path *string `json:"path,omitempty"`
}

// Response contract for the Wiki Page Move API.
type WikiPageMoveResponse struct {
	// Contains the list of ETag values from the response header of the page move API call. The first item in the list contains the version of the wiki page subject to page move.
	ETag *[]string `json:"eTag,omitempty"`
	// Defines properties for wiki page move.
	PageMove *WikiPageMove `json:"pageMove,omitempty"`
}

// Response contract for the Wiki Pages PUT, PATCH and DELETE APIs.
type WikiPageResponse struct {
	// Contains the list of ETag values from the response header of the pages API call. The first item in the list contains the version of the wiki page.
	ETag *[]string `json:"eTag,omitempty"`
	// Defines properties for wiki page.
	Page *WikiPage `json:"page,omitempty"`
}

// Contract encapsulating parameters for the pages batch.
type WikiPagesBatchRequest struct {
	// If the list of page data returned is not complete, a continuation token to query next batch of pages is included in the response header as "x-ms-continuationtoken". Omit this parameter to get the first batch of Wiki Page Data.
	ContinuationToken *string `json:"continuationToken,omitempty"`
	// last N days from the current day for which page views is to be returned. It's inclusive of current day.
	PageViewsForDays *int `json:"pageViewsForDays,omitempty"`
	// Total count of pages on a wiki to return.
	Top *int `json:"top,omitempty"`
}

// Defines properties for wiki page stat.
type WikiPageStat struct {
	// the count of the stat for the Day
	Count *int `json:"count,omitempty"`
	// Day of the stat
	Day *azuredevops.Time `json:"day,omitempty"`
}

// Defines properties for wiki page view stats.
type WikiPageViewStats struct {
	// Wiki page view count.
	Count *int `json:"count,omitempty"`
	// Wiki page last viewed time.
	LastViewedTime *azuredevops.Time `json:"lastViewedTime,omitempty"`
	// Wiki page path.
	Path *string `json:"path,omitempty"`
}

// Wiki types.
type WikiType string

type wikiTypeValuesType struct {
	ProjectWiki WikiType
	CodeWiki    WikiType
}

var WikiTypeValues = wikiTypeValuesType{
	// Indicates that the wiki is provisioned for the team project
	ProjectWiki: "projectWiki",
	// Indicates that the wiki is published from a git repository
	CodeWiki: "codeWiki",
}

type WikiUpdatedNotificationMessage struct {
	// Collection host Id for which the wikis are updated.
	CollectionId *uuid.UUID `json:"collectionId,omitempty"`
	// Project Id for which the wikis are updated.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// Repository Id associated with the particular wiki which is added, updated or deleted.
	RepositoryId *uuid.UUID `json:"repositoryId,omitempty"`
}

// Wiki update parameters.
type WikiUpdateParameters struct {
	// Name for wiki.
	Name *string `json:"name,omitempty"`
	// Versions of the wiki.
	Versions *[]git.GitVersionDescriptor `json:"versions,omitempty"`
}

// Defines a wiki resource.
type WikiV2 struct {
	// Folder path inside repository which is shown as Wiki. Not required for ProjectWiki type.
	MappedPath *string `json:"mappedPath,omitempty"`
	// Wiki name.
	Name *string `json:"name,omitempty"`
	// ID of the project in which the wiki is to be created.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`
	// ID of the git repository that backs up the wiki. Not required for ProjectWiki type.
	RepositoryId *uuid.UUID `json:"repositoryId,omitempty"`
	// Type of the wiki.
	Type *WikiType `json:"type,omitempty"`
	// ID of the wiki.
	Id *uuid.UUID `json:"id,omitempty"`
	// Is wiki repository disabled
	IsDisabled *bool `json:"isDisabled,omitempty"`
	// Properties of the wiki.
	Properties *map[string]string `json:"properties,omitempty"`
	// Remote web url to the wiki.
	RemoteUrl *string `json:"remoteUrl,omitempty"`
	// REST url for this wiki.
	Url *string `json:"url,omitempty"`
	// Versions of the wiki.
	Versions *[]git.GitVersionDescriptor `json:"versions,omitempty"`
}
