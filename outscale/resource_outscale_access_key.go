package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
)

func resourceOutscaleIamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleIamAccessKeyCreate,
		Read:   resourceOutscaleIamAccessKeyRead,
		Delete: resourceOutscaleIamAccessKeyDelete,

		Schema: map[string]*schema.Schema{
			"access_key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tag":     tagsSchema(),
			"tag_set": tagsSchemaComputed(),
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleIamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	icu_client := meta.(*OutscaleClient).ICU
	fcu_client := meta.(*OutscaleClient).FCU

	request := &icu.CreateAccessKeyInput{}
	if v, ok := d.GetOk("access_key_id"); ok {
		request.AccessKeyId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("secret_access_key"); ok {
		request.SecretAccessKey = aws.String(v.(string))
	}
	if d.IsNewResource() {
		if err := setTags(fcu_client, d); err != nil {
			return err
		}
		d.SetPartial("tag_set")
	}
	var err error
	var createResp *icu.CreateAccessKeyOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		createResp, err = icu_client.ICU.CreateAccessKey(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(*createResp.AccessKey.AccessKeyId)

	if createResp.AccessKey == nil || createResp.AccessKey.SecretAccessKey == nil {
		return fmt.Errorf("[ERR] CreateAccessKey response did not contain a Secret Access Key as expected")
	}

	return resourceOutscaleIamAccessKeyReadResult(d, &icu.AccessKeyMetadata{
		AccessKeyId: createResp.AccessKey.AccessKeyId,
		Status:      createResp.AccessKey.Status,
	})
}

func resourceOutscaleIamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	icu_client := meta.(*OutscaleClient).ICU

	request := &icu.ListAccessKeysInput{}

	var err error
	var getResp *icu.ListAccessKeysOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		getResp, err = icu_client.ICU.ListAccessKeys(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") { // XXX TEST ME
			// the user does not exist, so the key can't exist.
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading IAM acces key: %s", err)
	}

	if getResp.AccessKeyMetadata[0].AccessKeyId != nil {
		d.SetId(*getResp.AccessKeyMetadata[0].AccessKeyId)
	}
	if getResp.AccessKeyMetadata[0].Status != nil {
		d.Set("status", getResp.AccessKeyMetadata[0].Status)
	}
	if getResp.AccessKeyMetadata[0].OwnerId != nil {
		d.Set("owner_id", getResp.AccessKeyMetadata[0].OwnerId)
	}
	if getResp.AccessKeyMetadata[0].SecretAccessKey != nil {
		d.Set("secret_access_key", getResp.AccessKeyMetadata[0].SecretAccessKey)
	}
	if getResp.AccessKeyMetadata[0].Tags != nil {
		d.Set("tag_set", tagsToMapss(getResp.AccessKeyMetadata[0].Tags))
	}
	return nil
}

func resourceOutscaleIamAccessKeyReadResult(d *schema.ResourceData, key *icu.AccessKeyMetadata) error {
	d.SetId(*key.AccessKeyId)
	if err := d.Set("status", key.Status); err != nil {
		return err
	}
	return nil
}

func resourceOutscaleIamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	icu_client := meta.(*OutscaleClient).ICU

	request := &icu.DeleteAccessKeyInput{
		AccessKeyId: aws.String(d.Id()),
	}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err = icu_client.ICU.DeleteAccessKey(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting access key %s: %s", d.Id(), err)
	}

	return nil
}
func resourceOutscaleIamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	icu_client := meta.(*OutscaleClient).ICU

	request := &icu.UpdateAccessKeyInput{}
	if v, ok := d.GetOk("access_key_id"); ok {
		request.AccessKeyId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		request.Status = aws.String(v.(string))
	}

	var err error
	var createResp *icu.UpdateAccessKeyOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		createResp, err = icu_client.ICU.UpdateAccessKey(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceOutscaleIamAccessKeyRead(d, meta)
}
