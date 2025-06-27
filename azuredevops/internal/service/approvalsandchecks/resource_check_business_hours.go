package approvalsandchecks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk/pipelineschecksextras"
)

var (
	evaluateBusinessHoursDefVersion = "0.0.1"
	evaluateBusinessHoursDefId      = "445fde2f-6c39-441c-807f-8a59ff2e075f"
)

var evaluateBusinessHoursDef = map[string]interface{}{
	"id":      evaluateBusinessHoursDefId,
	"name":    "evaluateBusinessHours",
	"version": evaluateBusinessHoursDefVersion,
}

var validTimezoneIds = []string{"AUS Central Standard Time", "AUS Eastern Standard Time", "Afghanistan Standard Time", "Alaskan Standard Time", "Aleutian Standard Time", "Altai Standard Time", "Arab Standard Time", "Arabian Standard Time", "Arabic Standard Time", "Argentina Standard Time", "Astrakhan Standard Time", "Atlantic Standard Time", "Aus Central W. Standard Time", "Azerbaijan Standard Time", "Azores Standard Time", "Bahia Standard Time", "Bangladesh Standard Time", "Belarus Standard Time", "Bougainville Standard Time", "Canada Central Standard Time", "Cape Verde Standard Time", "Caucasus Standard Time", "Cen. Australia Standard Time", "Central America Standard Time", "Central Asia Standard Time", "Central Brazilian Standard Time", "Central Europe Standard Time", "Central European Standard Time", "Central Pacific Standard Time", "Central Standard Time (Mexico)", "Central Standard Time", "Chatham Islands Standard Time", "China Standard Time", "Cuba Standard Time", "Dateline Standard Time", "E. Africa Standard Time", "E. Australia Standard Time", "E. Europe Standard Time", "E. South America Standard Time", "Easter Island Standard Time", "Eastern Standard Time (Mexico)", "Eastern Standard Time", "Egypt Standard Time", "Ekaterinburg Standard Time", "FLE Standard Time", "Fiji Standard Time", "GMT Standard Time", "GTB Standard Time", "Georgian Standard Time", "Greenland Standard Time", "Greenwich Standard Time", "Haiti Standard Time", "Hawaiian Standard Time", "India Standard Time", "Iran Standard Time", "Israel Standard Time", "Jordan Standard Time", "Kaliningrad Standard Time", "Kamchatka Standard Time", "Korea Standard Time", "Libya Standard Time", "Line Islands Standard Time", "Lord Howe Standard Time", "Magadan Standard Time", "Magallanes Standard Time", "Marquesas Standard Time", "Mauritius Standard Time", "Mid-Atlantic Standard Time", "Middle East Standard Time", "Montevideo Standard Time", "Morocco Standard Time", "Mountain Standard Time (Mexico)", "Mountain Standard Time", "Myanmar Standard Time", "N. Central Asia Standard Time", "Namibia Standard Time", "Nepal Standard Time", "New Zealand Standard Time", "Newfoundland Standard Time", "Norfolk Standard Time", "North Asia East Standard Time", "North Asia Standard Time", "North Korea Standard Time", "Omsk Standard Time", "Pacific SA Standard Time", "Pacific Standard Time (Mexico)", "Pacific Standard Time", "Pakistan Standard Time", "Paraguay Standard Time", "Qyzylorda Standard Time", "Romance Standard Time", "Russia Time Zone 10", "Russia Time Zone 11", "Russia Time Zone 3", "Russian Standard Time", "SA Eastern Standard Time", "SA Pacific Standard Time", "SA Western Standard Time", "SE Asia Standard Time", "Saint Pierre Standard Time", "Sakhalin Standard Time", "Samoa Standard Time", "Sao Tome Standard Time", "Saratov Standard Time", "Singapore Standard Time", "South Africa Standard Time", "South Sudan Standard Time", "Sri Lanka Standard Time", "Sudan Standard Time", "Syria Standard Time", "Taipei Standard Time", "Tasmania Standard Time", "Tocantins Standard Time", "Tokyo Standard Time", "Tomsk Standard Time", "Tonga Standard Time", "Transbaikal Standard Time", "Turkey Standard Time", "Turks And Caicos Standard Time", "US Eastern Standard Time", "US Mountain Standard Time", "UTC", "UTC+12", "UTC+13", "UTC-02", "UTC-08", "UTC-09", "UTC-11", "Ulaanbaatar Standard Time", "Venezuela Standard Time", "Vladivostok Standard Time", "Volgograd Standard Time", "W. Australia Standard Time", "W. Central Africa Standard Time", "W. Europe Standard Time", "W. Mongolia Standard Time", "West Asia Standard Time", "West Bank Standard Time", "West Pacific Standard Time", "Yakutsk Standard Time", "Yukon Standard Time"}

type DayOfBusinessWeek struct {
	TfName  string
	AdoName string
}

var daysOfBusinessWeek = []DayOfBusinessWeek{
	{
		TfName:  "monday",
		AdoName: "Monday",
	},
	{
		TfName:  "tuesday",
		AdoName: "Tuesday",
	},
	{
		TfName:  "wednesday",
		AdoName: "Wednesday",
	},
	{
		TfName:  "thursday",
		AdoName: "Thursday",
	},
	{
		TfName:  "friday",
		AdoName: "Friday",
	},
	{
		TfName:  "saturday",
		AdoName: "Saturday",
	},
	{
		TfName:  "sunday",
		AdoName: "Sunday",
	},
}

// ResourceCheckBusinessHours schema and implementation for build definition resource
func ResourceCheckBusinessHours() *schema.Resource {
	r := genBaseCheckResource(flattenBusinessHours, expandBusinessHours)
	for _, day := range daysOfBusinessWeek {
		r.Schema[day.TfName] = &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		}
	}

	r.Schema["time_zone"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(validTimezoneIds, false),
	}

	timeRegExp := regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`)

	r.Schema["start_time"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringMatch(timeRegExp, "Must be a 24 hour time with leading zeros"),
	}
	r.Schema["end_time"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringMatch(timeRegExp, "Must be a 24 hour time with leading zeros"),
	}

	r.Schema["display_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "Managed by Terraform",
		ValidateFunc: validation.StringIsNotEmpty,
	}

	r.Schema["timeout"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1440,
		ValidateFunc: validation.IntBetween(1, 2147483647),
	}

	return r
}

func flattenBusinessHours(d *schema.ResourceData, businessHoursCheck *pipelineschecksextras.CheckConfiguration, projectID string) error {
	err := doBaseFlattening(d, businessHoursCheck, projectID)
	if err != nil {
		return err
	}

	businessHoursCheck.Type.Id = converter.UUID("fe1de3ee-a436-41b4-bb20-f6eb4cb879a7")

	if businessHoursCheck.Settings == nil {
		return fmt.Errorf("Settings nil")
	}

	if displayName, found := businessHoursCheck.Settings.(map[string]interface{})["displayName"]; found {
		d.Set("display_name", displayName.(string))
	} else {
		return fmt.Errorf("displayName setting not found")
	}

	if definitionRefMap, found := businessHoursCheck.Settings.(map[string]interface{})["definitionRef"]; found {
		definitionRef := definitionRefMap.(map[string]interface{})
		if id, found := definitionRef["id"]; found {
			if !strings.EqualFold(id.(string), evaluateBusinessHoursDefId) {
				return fmt.Errorf("invalid definitionRef id")
			}
		} else {
			return fmt.Errorf("definitionRef ID not found. Expect ID: %s", evaluateBusinessHoursDefId)
		}
		if version, found := definitionRef["version"]; found {
			if version != evaluateBusinessHoursDefVersion {
				return fmt.Errorf("unsupported definitionRef version. Expect version: %s", evaluateBusinessHoursDefVersion)
			}
		} else {
			return fmt.Errorf("unsupported definitionRef version")
		}
	} else {
		return fmt.Errorf("definitionRef not found")
	}

	if inputMap, found := businessHoursCheck.Settings.(map[string]interface{})["inputs"]; found {
		inputs := inputMap.(map[string]interface{})
		if businessDays, found := inputs["businessDays"]; found {
			for _, day := range daysOfBusinessWeek {
				d.Set(day.TfName, strings.Contains(businessDays.(string), day.AdoName))
			}
		} else {
			return fmt.Errorf("businessDays input not found")
		}
		if timeZone, found := inputs["timeZone"]; found {
			d.Set("time_zone", timeZone)
		} else {
			return fmt.Errorf("timeZone input not found")
		}
		if startTime, found := inputs["startTime"]; found {
			d.Set("start_time", startTime)
		} else {
			return fmt.Errorf("startTime input not found")
		}
		if endTime, found := inputs["endTime"]; found {
			d.Set("end_time", endTime)
		} else {
			return fmt.Errorf("endTime input not found")
		}
	} else {
		return fmt.Errorf("inputs not found")
	}

	if businessHoursCheck.Timeout != nil {
		d.Set("timeout", *businessHoursCheck.Timeout)
	}

	return nil
}

func expandBusinessHours(d *schema.ResourceData) (*pipelineschecksextras.CheckConfiguration, string, error) {
	var days []string
	for _, day := range daysOfBusinessWeek {
		if d.Get(day.TfName).(bool) {
			days = append(days, day.AdoName)
		}
	}

	inputs := map[string]interface{}{
		"businessDays": strings.Join(days, ", "),
		"startTime":    d.Get("start_time").(string),
		"endTime":      d.Get("end_time").(string),
		"timeZone":     d.Get("time_zone").(string),
	}

	settings := map[string]interface{}{}
	settings["inputs"] = inputs
	settings["definitionRef"] = evaluateBusinessHoursDef
	settings["displayName"] = d.Get("display_name").(string)

	return doBaseExpansion(d, approvalAndCheckType.BusinessHours, settings, converter.ToPtr(d.Get("timeout").(int)))
}
