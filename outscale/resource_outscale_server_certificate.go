package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/openlyinc/pointy"
	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func resourceOutscaleOAPIServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIServerCertificateCreate,
		Read:   resourceOutscaleOAPIServerCertificateRead,
		Update: resourceOutscaleOAPIServerCertificateUpdate,
		Delete: resourceOutscaleOAPIServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"body": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"chain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upload_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateServerCertificateRequest{
		Body:       d.Get("body").(string),
		Name:       d.Get("name").(string),
		PrivateKey: d.Get("private_key").(string),
	}

	if v, ok := d.GetOk("chain"); ok {
		req.Chain = pointy.String(v.(string))
	}
	if v, ok := d.GetOk("dry_run"); ok {
		req.DryRun = pointy.Bool(v.(bool))
	}
	if v, ok := d.GetOk("path"); ok {
		req.Path = pointy.String(v.(string))
	}

	resp, _, err := conn.ServerCertificateApi.CreateServerCertificate(context.Background()).CreateServerCertificateRequest(req).Execute()
	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating Server Certificate: %s", utils.GetErrorResponse(err))
	}

	d.SetId(cast.ToString(resp.ServerCertificate.Id))

	return resourceOutscaleOAPIServerCertificateRead(d, meta)
}

func resourceOutscaleOAPIServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	log.Printf("[DEBUG] Reading Server Certificate id (%s)", id)

	var resp oscgo.ReadServerCertificatesResponse

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		r, _, err := conn.ServerCertificateApi.ReadServerCertificates(context.Background()).ReadServerCertificatesRequest(oscgo.ReadServerCertificatesRequest{}).Execute()

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
		return fmt.Errorf("[DEBUG] Error reading Server Certificate id (%s)", utils.GetErrorResponse(err))

	}
	if !resp.HasServerCertificates() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(resp.GetServerCertificates()) == 0 {
		d.SetId("")
		return fmt.Errorf("ServerCertificates not found")
	}

	var server oscgo.ServerCertificate

	for _, serv := range resp.GetServerCertificates() {
		if serv.GetId() == d.Id() {
			server = serv
		}
	}

	d.Set("expiration_date", server.ExpirationDate)
	d.Set("name", server.Name)
	d.Set("path", server.Path)
	d.Set("upload_date", server.UploadDate)

	return nil
}

func resourceOutscaleOAPIServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	oldName, newName := d.GetChange("name")
	req := oscgo.UpdateServerCertificateRequest{
		Name: oldName.(string),
	}

	if d.HasChange("name") {
		req.NewName = pointy.String(newName.(string))
	}
	if d.HasChange("path") {
		req.NewPath = pointy.String(d.Get("path").(string))
	}

	_, _, err := conn.ServerCertificateApi.UpdateServerCertificate(context.Background()).UpdateServerCertificateRequest(req).Execute()
	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating Server Certificate: %s", utils.GetErrorResponse(err))
	}
	return resourceOutscaleOAPIServerCertificateRead(d, meta)
}

func resourceOutscaleOAPIServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteServerCertificateRequest{
		Name: d.Get("name").(string),
	}

	_, _, err := conn.ServerCertificateApi.DeleteServerCertificate(context.Background()).DeleteServerCertificateRequest(req).Execute()
	if err != nil {
		return fmt.Errorf("[DEBUG] Error deleting Server Certificate id (%s)", err)
	}

	return nil
}
