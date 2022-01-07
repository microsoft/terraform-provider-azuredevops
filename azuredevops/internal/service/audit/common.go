package audit

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

const errMsgTfConfigRead = "Error reading terraform configuration: %+v"

type flatFunc func(d *schema.ResourceData, auditStream *audit.AuditStream, daysToBackfill *int)
type expandFunc func(d *schema.ResourceData) (*audit.AuditStream, *int, error)

// genBaseAuditStreamResource creates a Resource with the common parts
// that all Audit Streams require.
func genBaseAuditStreamResource(f flatFunc, e expandFunc) *schema.Resource {
	return &schema.Resource{
		Create: genAuditStreamCreateFunc(f, e),
		Read:   genAuditStreamReadFunc(f),
		Update: genAuditStreamUpdateFunc(f, e),
		Delete: genAuditStreamDeleteFunc(),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"days_to_backfill": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  "The number of days of previously recorded audit data that will be replayed into the stream",
			},
		},
	}
}

// doBaseExpansion performs the expansion for the 'base' attributes that are defined in the schema, above
func doBaseExpansion(d *schema.ResourceData) (*audit.AuditStream, *int) {
	// an "error" is OK here as it is expected in the case that the ID is not set in the resource data
	var auditStreamId *int
	parsedId, err := strconv.Atoi(d.Id())
	if err == nil {
		auditStreamId = &parsedId
	}

	daysToBackfill := converter.Int(d.Get("days_to_backfill").(int))
	auditStream := &audit.AuditStream{
		Id: auditStreamId,
	}

	return auditStream, daysToBackfill
}

// doBaseFlattening performs the flattening for the 'base' attributes that are defined in the schema, above
func doBaseFlattening(d *schema.ResourceData, auditStream *audit.AuditStream, daysToBackfill *int) {
	d.SetId(strconv.Itoa(*auditStream.Id))
	d.Set("days_to_backfill", daysToBackfill)
}

func genAuditStreamCreateFunc(flatFunc flatFunc, expandFunc expandFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		auditStream, daysToBackfill, err := expandFunc(d)
		if err != nil {
			return fmt.Errorf(errMsgTfConfigRead, err)
		}

		createdAuditStream, err := createAuditStream(clients, auditStream, daysToBackfill)
		if err != nil {
			return fmt.Errorf("Error creating audit stream in Azure DevOps: %+v", err)
		}

		d.SetId(strconv.Itoa(*createdAuditStream.Id))
		return genAuditStreamReadFunc(flatFunc)(d, m)
	}
}

func genAuditStreamReadFunc(flatFunc flatFunc) func(d *schema.ResourceData, m interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		streamId, err := strconv.Atoi(d.Id())
		if err != nil {
			return fmt.Errorf("Error parsing the audit stream ID from the Terraform resource data: %v", err)
		}

		daysToBackfill := d.Get("days_to_backfill").(int)

		auditStream, err := readAuditStream(clients, streamId)
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error looking up audit stream with ID %d. Error: %v", streamId, err)
		}

		if auditStream.Id == nil {
			// e.g. audit stream has been deleted separately without TF
			d.SetId("")
		} else {
			flatFunc(d, auditStream, &daysToBackfill)
		}
		return nil
	}
}

func genAuditStreamUpdateFunc(flatFunc flatFunc, expandFunc expandFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		auditStream, daysToBackfill, err := expandFunc(d)
		if err != nil {
			return fmt.Errorf(errMsgTfConfigRead, err)
		}

		updatedAuditStream, err := updateAuditStream(clients, auditStream)
		if err != nil {
			return fmt.Errorf("Error updating audit stream in Azure DevOps: %+v", err)
		}

		flatFunc(d, updatedAuditStream, daysToBackfill)
		return genAuditStreamReadFunc(flatFunc)(d, m)
	}
}

func genAuditStreamDeleteFunc() schema.DeleteFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		streamId, err := strconv.Atoi(d.Id())
		if err != nil {
			return err
		}

		err = deleteAuditStream(clients, streamId)
		if err != nil {
			return fmt.Errorf("Error deleting audit stream in Azure DevOps: %+v", err)
		}

		return nil
	}
}

func createAuditStream(clients *client.AggregatedClient, stream *audit.AuditStream, daysToBackfill *int) (*audit.AuditStream, error) {
	createdAuditStream, err := clients.AuditClient.CreateStream(
		clients.Ctx,
		audit.CreateStreamArgs{
			Stream:         stream,
			DaysToBackfill: daysToBackfill,
		})

	return createdAuditStream, err
}

func readAuditStream(clients *client.AggregatedClient, streamId int) (*audit.AuditStream, error) {
	auditStream, err := clients.AuditClient.QueryStreamById(
		clients.Ctx,
		audit.QueryStreamByIdArgs{
			StreamId: &streamId,
		})

	return auditStream, err
}

func updateAuditStream(clients *client.AggregatedClient, stream *audit.AuditStream) (*audit.AuditStream, error) {
	updatedAuditStream, err := clients.AuditClient.UpdateStream(
		clients.Ctx,
		audit.UpdateStreamArgs{
			Stream: stream,
		})

	return updatedAuditStream, err
}

func deleteAuditStream(clients *client.AggregatedClient, streamId int) error {
	return clients.AuditClient.DeleteStream(
		clients.Ctx,
		audit.DeleteStreamArgs{
			StreamId: &streamId,
		})
}
