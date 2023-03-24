package git

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/git"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestGitRepositoryBranch_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *schema.ResourceData
		m   interface{}
	}
	tests := []struct {
		name string
		args func(g *azdosdkmocks.MockGitClient) args
		want diag.Diagnostics
	}{
		{
			"When ref is given, refs update does not swallow error",
			func(g *azdosdkmocks.MockGitClient) args {
				clients := &client.AggregatedClient{
					GitReposClient: g,
					Ctx:            context.Background(),
				}
				d := schema.TestResourceDataRaw(t, ResourceGitRepositoryBranch().Schema, nil)
				ref := "refs/heads/another-branch"
				fakeCommitId := "a-commit"
				branchName := "a-branch"
				repoId := "a-repo"
				d.Set("ref_branch", ref)
				d.Set("name", branchName)
				d.Set("repository_id", repoId)

				g.EXPECT().
					GetRefs(clients.Ctx, git.GetRefsArgs{
						RepositoryId: &repoId,
						Filter:       converter.String(strings.TrimPrefix(ref, "refs/")),
						Top:          converter.Int(1),
						PeelTags:     converter.Bool(true),
					}).
					Return(&git.GetRefsResponseValue{
						Value: []git.GitRef{{
							Name:     &ref,
							ObjectId: &fakeCommitId,
						}},
					}, nil)

				g.EXPECT().
					UpdateRefs(clients.Ctx, git.UpdateRefsArgs{
						RefUpdates: &[]git.GitRefUpdate{{
							Name:        converter.String(withPrefix("refs/heads/", "a-branch")),
							NewObjectId: &fakeCommitId,
							OldObjectId: converter.String("0000000000000000000000000000000000000000"),
						}},
						RepositoryId: converter.String("a-repo"),
					}).
					Return(nil, fmt.Errorf("an-error"))
				return args{
					context.Background(),
					d,
					clients,
				}
			},
			diag.FromErr(fmt.Errorf("Error creating branch \"a-branch\": an-error")),
		},
		{
			"When more than one of ref_commit_id, ref_tag, or ref_branch are given, the first one from left to right wins",
			func(g *azdosdkmocks.MockGitClient) args {
				clients := &client.AggregatedClient{
					GitReposClient: g,
					Ctx:            context.Background(),
				}
				d := schema.TestResourceDataRaw(t, ResourceGitRepositoryBranch().Schema, nil)
				tag := "refs/tags/v1.0.0"
				fakeCommitId := "a-commit"
				branchName := "a-branch"
				repoId := "a-repo"
				d.Set("ref_commit_id", fakeCommitId)
				d.Set("ref_tag", tag)
				d.Set("name", branchName)
				d.Set("repository_id", repoId)

				g.EXPECT().
					UpdateRefs(clients.Ctx, git.UpdateRefsArgs{
						RefUpdates: &[]git.GitRefUpdate{{
							Name:        converter.String(withPrefix("refs/heads/", "a-branch")),
							NewObjectId: &fakeCommitId,
							OldObjectId: converter.String("0000000000000000000000000000000000000000"),
						}},
						RepositoryId: converter.String("a-repo"),
					}).
					Return(nil, fmt.Errorf("an-error"))
				return args{
					context.Background(),
					d,
					clients,
				}
			},
			diag.FromErr(fmt.Errorf("Error creating branch \"a-branch\": an-error")),
		},
		{
			"When invalid RefUpdate UpdateStatus, throw error",
			func(g *azdosdkmocks.MockGitClient) args {
				clients := &client.AggregatedClient{
					GitReposClient: g,
					Ctx:            context.Background(),
				}
				d := schema.TestResourceDataRaw(t, ResourceGitRepositoryBranch().Schema, nil)
				ref := "refs/heads/another-branch"
				fakeCommitId := "a-commit"
				branchName := "a-branch"
				repoId := "a-repo"
				d.Set("ref_branch", ref)
				d.Set("name", branchName)
				d.Set("repository_id", repoId)

				g.EXPECT().
					GetRefs(clients.Ctx, git.GetRefsArgs{
						RepositoryId: &repoId,
						Filter:       converter.String(strings.TrimPrefix(ref, "refs/")),
						Top:          converter.Int(1),
						PeelTags:     converter.Bool(true),
					}).
					Return(&git.GetRefsResponseValue{
						Value: []git.GitRef{{
							Name:     &ref,
							ObjectId: &fakeCommitId,
						}},
					}, nil)

				g.EXPECT().
					UpdateRefs(clients.Ctx, git.UpdateRefsArgs{
						RefUpdates: &[]git.GitRefUpdate{{
							Name:        converter.String(withPrefix("refs/heads/", "a-branch")),
							NewObjectId: &fakeCommitId,
							OldObjectId: converter.String("0000000000000000000000000000000000000000"),
						}},
						RepositoryId: converter.String("a-repo"),
					}).
					Return(&[]git.GitRefUpdateResult{{
						Success:      converter.Bool(false),
						UpdateStatus: &git.GitRefUpdateStatusValues.InvalidRefName,
					}}, nil)
				return args{
					context.Background(),
					d,
					clients,
				}
			},
			diag.FromErr(fmt.Errorf("Error creating branch \"a-branch\": Error got invalid GitRefUpdate.UpdateStatus: invalidRefName")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gitClient := azdosdkmocks.NewMockGitClient(ctrl)
			testArgs := tt.args(gitClient)

			if got := resourceGitRepositoryBranchCreate(testArgs.ctx, testArgs.d, testArgs.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceGitRepositoryBranchCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepositoryBranch_Read(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *schema.ResourceData
		m   interface{}
	}
	tests := []struct {
		name string
		args func(g *azdosdkmocks.MockGitClient) args
		want diag.Diagnostics
	}{
		{
			"Read does not swallow error.",
			func(g *azdosdkmocks.MockGitClient) args {
				clients := &client.AggregatedClient{
					GitReposClient: g,
					Ctx:            context.Background(),
				}

				d := schema.TestResourceDataRaw(t, ResourceGitRepositoryBranch().Schema, nil)
				d.Set("ref_branch", "another-branch")
				d.Set("name", "a-branch")
				d.Set("repository_id", "a-repo")
				d.SetId("a-repo:a-branch")

				g.EXPECT().
					GetBranch(clients.Ctx, git.GetBranchArgs{
						RepositoryId: converter.String("a-repo"),
						Name:         converter.String("a-branch"),
					}).
					Return(nil, fmt.Errorf("an-error"))

				return args{
					ctx: context.Background(),
					d:   d,
					m:   clients,
				}
			},
			diag.FromErr(fmt.Errorf("Error reading branch \"a-branch\": an-error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gitClient := azdosdkmocks.NewMockGitClient(ctrl)
			testArgs := tt.args(gitClient)

			if got := resourceGitRepositoryBranchRead(testArgs.ctx, testArgs.d, testArgs.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceGitRepositoryBranchCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepositoryBranch_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *schema.ResourceData
		m   interface{}
	}
	tests := []struct {
		name string
		args func(g *azdosdkmocks.MockGitClient) args
		want diag.Diagnostics
	}{
		{
			"Delete based on repositoryId:branchName does not swallow error.",
			func(g *azdosdkmocks.MockGitClient) args {
				clients := &client.AggregatedClient{
					GitReposClient: g,
					Ctx:            context.Background(),
				}

				d := schema.TestResourceDataRaw(t, ResourceGitRepositoryBranch().Schema, nil)
				d.Set("ref_branch", "another-branch")
				d.Set("name", "a-branch")
				d.Set("repository_id", "a-repo")
				d.SetId("a-repo:a-branch")

				g.EXPECT().
					GetBranch(clients.Ctx, git.GetBranchArgs{
						RepositoryId: converter.String("a-repo"),
						Name:         converter.String("a-branch"),
					}).
					Return(&git.GitBranchStats{
						Commit: &git.GitCommitRef{
							CommitId: converter.String("a-commit"),
						},
					}, nil)

				g.EXPECT().
					UpdateRefs(clients.Ctx, git.UpdateRefsArgs{
						RefUpdates: &[]git.GitRefUpdate{{
							Name:        converter.String(withPrefix("refs/heads/", "a-branch")),
							OldObjectId: converter.String("a-commit"),
							NewObjectId: converter.String("0000000000000000000000000000000000000000"),
						}},
						RepositoryId: converter.String("a-repo"),
					}).
					Return(nil, fmt.Errorf("an-error"))

				return args{
					ctx: clients.Ctx,
					d:   d,
					m:   clients,
				}
			},
			diag.FromErr(fmt.Errorf("Error deleting branch \"a-branch\": an-error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gitClient := azdosdkmocks.NewMockGitClient(ctrl)
			testArgs := tt.args(gitClient)

			if got := resourceGitRepositoryBranchDelete(testArgs.ctx, testArgs.d, testArgs.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceGitRepositoryBranchDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}
