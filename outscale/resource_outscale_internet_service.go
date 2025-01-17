package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func resourceOutscaleOAPIInternetService() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIInternetServiceCreate,
		Read:   resourceOutscaleOAPIInternetServiceRead,
		Update: resourceOutscaleOAPIInternetServiceUpdate,
		Delete: resourceOutscaleOAPIInternetServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsListOAPISchema(),
		},
	}
}

func resourceOutscaleOAPIInternetServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	resp, _, err := conn.InternetServiceApi.CreateInternetService(context.Background()).CreateInternetServiceRequest(oscgo.CreateInternetServiceRequest{}).Execute()
	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating Internet Service: %s", utils.GetErrorResponse(err))
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), resp.InternetService.GetInternetServiceId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.InternetService.GetInternetServiceId())

	return resourceOutscaleOAPIInternetServiceRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	log.Printf("[DEBUG] Reading Internet Service id (%s)", id)

	req := oscgo.ReadInternetServicesRequest{
		Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{id}},
	}

	var resp oscgo.ReadInternetServicesResponse

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		r, _, err := conn.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(req).Execute()

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", utils.GetErrorResponse(err))

	}
	if !resp.HasInternetServices() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(resp.GetInternetServices()) == 0 {
		d.SetId("")
		return fmt.Errorf("InternetServices not found")
	}
	if err := d.Set("internet_service_id", resp.GetInternetServices()[0].GetInternetServiceId()); err != nil {
		return err
	}

	if err := d.Set("net_id", resp.GetInternetServices()[0].GetNetId()); err != nil {
		return err
	}

	if err := d.Set("state", resp.GetInternetServices()[0].GetState()); err != nil {
		return err
	}

	return d.Set("tags", tagsOSCAPIToMap(resp.GetInternetServices()[0].GetTags()))
}

func resourceOutscaleOAPIInternetServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)
	return resourceOutscaleOAPIInternetServiceRead(d, meta)
}

func resourceOutscaleOAPIInternetServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	internetServiceID := d.Id()
	filterReq := oscgo.ReadInternetServicesRequest{
		Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{internetServiceID}},
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"deleted", "available"},
		Refresh:    LISOAPIStateRefreshFunction(conn, filterReq, "failed"),
		Timeout:    10 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Internet Service (%s) to become deleted: %s", d.Id(), err)
	}

	req := oscgo.DeleteInternetServiceRequest{
		InternetServiceId: internetServiceID,
	}

	_, _, err = conn.InternetServiceApi.DeleteInternetService(context.Background()).DeleteInternetServiceRequest(req).Execute()
	if err != nil {
		return fmt.Errorf("[DEBUG] Error deleting Internet Service id (%s)", err)
	}

	return nil
}
