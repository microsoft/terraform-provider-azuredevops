//go:build (all || resource_check_business_hours) && !exclude_approvalsandchecks

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccCheckBusinessHours_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkName := testutils.GenerateResourceName()
	start_time := "01:20"
	end_time := "03:20"

	resourceType := "azuredevops_check_business_hours"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclCheckBusinessHoursResourceBasic(projectName, checkName, start_time, end_time),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkName),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfCheckNode, "target_resource_id"),
					resource.TestCheckResourceAttrSet(tfCheckNode, "target_resource_type"),
					resource.TestCheckResourceAttr(tfCheckNode, "start_time", start_time),
					resource.TestCheckResourceAttr(tfCheckNode, "end_time", end_time),
				),
			},
		},
	})
}

func TestAccCheckBusinessHours_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkName := testutils.GenerateResourceName()
	start_time := "01:20"
	end_time := "02:20"
	time_zone := "UTC"
	monday := "true"
	tuesday := "true"
	wednesday := "false"
	thursday := "false"
	friday := "false"
	saturday := "false"
	sunday := "false"

	resourceType := "azuredevops_check_business_hours"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclCheckBusinessHoursResourceComplete(projectName, checkName, start_time, end_time, time_zone, monday, tuesday, wednesday, thursday, friday, saturday, sunday),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkName),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfCheckNode, "target_resource_id"),
					resource.TestCheckResourceAttrSet(tfCheckNode, "target_resource_type"),
					resource.TestCheckResourceAttr(tfCheckNode, "start_time", start_time),
					resource.TestCheckResourceAttr(tfCheckNode, "end_time", end_time),
					resource.TestCheckResourceAttr(tfCheckNode, "time_zone", time_zone),
					resource.TestCheckResourceAttr(tfCheckNode, "monday", "true"),
					resource.TestCheckResourceAttr(tfCheckNode, "tuesday", "true"),
					resource.TestCheckResourceAttr(tfCheckNode, "wednesday", "false"),
					resource.TestCheckResourceAttr(tfCheckNode, "thursday", "false"),
					resource.TestCheckResourceAttr(tfCheckNode, "friday", "false"),
					resource.TestCheckResourceAttr(tfCheckNode, "saturday", "false"),
					resource.TestCheckResourceAttr(tfCheckNode, "sunday", "false"),
					resource.TestCheckResourceAttr(tfCheckNode, "timeout", "1440"),
				),
			},
		},
	})
}

func TestAccCheckBusinessHours_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkNameFirst := testutils.GenerateResourceName()
	start_time_first := "01:20"
	end_time_first := "02:20"
	time_zone_first := "UTC"
	monday_first := "true"
	tuesday_first := "true"
	wednesday_first := "false"
	thursday_first := "false"
	friday_first := "false"
	saturday_first := "false"
	sunday_first := "false"

	checkNameSecond := testutils.GenerateResourceName()
	start_time_second := "03:20"
	end_time_second := "04:20"
	time_zone_second := "AUS Central Standard Time"
	monday_second := "false"
	tuesday_second := "false"
	wednesday_second := "true"
	thursday_second := "true"
	friday_second := "true"
	saturday_second := "true"
	sunday_second := "true"

	resourceType := "azuredevops_check_business_hours"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclCheckBusinessHoursResourceComplete(projectName, checkNameFirst, start_time_first, end_time_first, time_zone_first, monday_first, tuesday_first, wednesday_first, thursday_first, friday_first, saturday_first, sunday_first),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkNameFirst),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "start_time", start_time_first),
					resource.TestCheckResourceAttr(tfCheckNode, "end_time", end_time_first),
					resource.TestCheckResourceAttr(tfCheckNode, "time_zone", time_zone_first),
					resource.TestCheckResourceAttr(tfCheckNode, "monday", monday_first),
					resource.TestCheckResourceAttr(tfCheckNode, "tuesday", tuesday_first),
					resource.TestCheckResourceAttr(tfCheckNode, "wednesday", wednesday_first),
					resource.TestCheckResourceAttr(tfCheckNode, "thursday", thursday_first),
					resource.TestCheckResourceAttr(tfCheckNode, "friday", friday_first),
					resource.TestCheckResourceAttr(tfCheckNode, "saturday", saturday_first),
					resource.TestCheckResourceAttr(tfCheckNode, "sunday", sunday_first),
				),
			},
			{
				Config: hclCheckBusinessHoursResourceUpdate(projectName, checkNameSecond, start_time_second, end_time_second, time_zone_second, monday_second, tuesday_second, wednesday_second, thursday_second, friday_second, saturday_second, sunday_second),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkNameSecond),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "start_time", start_time_second),
					resource.TestCheckResourceAttr(tfCheckNode, "end_time", end_time_second),
					resource.TestCheckResourceAttr(tfCheckNode, "time_zone", time_zone_second),
					resource.TestCheckResourceAttr(tfCheckNode, "monday", monday_second),
					resource.TestCheckResourceAttr(tfCheckNode, "tuesday", tuesday_second),
					resource.TestCheckResourceAttr(tfCheckNode, "wednesday", wednesday_second),
					resource.TestCheckResourceAttr(tfCheckNode, "thursday", thursday_second),
					resource.TestCheckResourceAttr(tfCheckNode, "friday", friday_second),
					resource.TestCheckResourceAttr(tfCheckNode, "saturday", saturday_second),
					resource.TestCheckResourceAttr(tfCheckNode, "sunday", sunday_second),
					resource.TestCheckResourceAttr(tfCheckNode, "version", "2"),
				),
			},
		},
	})
}

func hclCheckBusinessHoursResourceBasic(projectName string, checkName string, start_time string, end_time string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_business_hours" "test" {
  project_id           = azuredevops_project.project.id
  display_name         = "%s"
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  target_resource_type = "endpoint"
  time_zone            = "UTC"
  start_time           = "%s"
  end_time             = "%s"
  monday               = true
}`, checkName, start_time, end_time)

	genericServiceEndpointResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericServiceEndpointResource, checkResource)
}

func hclCheckBusinessHoursResourceComplete(projectName string, checkName string, start_time string, end_time string, time_zone string,
	monday string, tuesday string, wednesday string, thursday string, friday string, saturday string, sunday string,
) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_business_hours" "test" {
  project_id           = azuredevops_project.project.id
  display_name         = "%s"
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  target_resource_type = "endpoint"
  start_time           = "%s"
  end_time             = "%s"
  time_zone            = "%s"
  monday               = "%s"
  tuesday              = "%s"
  wednesday            = "%s"
  thursday             = "%s"
  friday               = "%s"
  saturday             = "%s"
  sunday               = "%s"
}`, checkName, start_time, end_time, time_zone, monday, tuesday, wednesday, thursday, friday, saturday, sunday,
	)

	genericServiceEndpointResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericServiceEndpointResource, checkResource)
}

func hclCheckBusinessHoursResourceUpdate(projectName string, checkName string, start_time string, end_time string, time_zone string,
	monday string, tuesday string, wednesday string, thursday string, friday string, saturday string, sunday string,
) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_business_hours" "test" {
  project_id           = azuredevops_project.project.id
  display_name         = "%s"
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  target_resource_type = "endpoint"
  start_time           = "%s"
  end_time             = "%s"
  time_zone            = "%s"
  monday               = "%s"
  tuesday              = "%s"
  wednesday            = "%s"
  thursday             = "%s"
  friday               = "%s"
  saturday             = "%s"
  sunday               = "%s"
  timeout              = 50000
}`, checkName, start_time, end_time, time_zone, monday, tuesday, wednesday, thursday, friday, saturday, sunday,
	)

	genericServiceEndpointResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericServiceEndpointResource, checkResource)
}
