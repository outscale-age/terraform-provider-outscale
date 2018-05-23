package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleImageExportTasks() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageExportTasksCreate,
		Read:   resourceImageExportTasksRead,
		Delete: resourceImageExportTasksDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"export_to_osu": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_image_format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"osu_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"osu_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"osu_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"aksk": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"secret_key": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"completion": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_export": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"snapshot_export_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImageExportTasksCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	eto, etoOk := d.GetOk("export_to_osu")
	v, ok := d.GetOk("snapshot_id")
	request := &fcu.CreateSnapshotExportTaskInput{}

	if !etoOk && !ok {
		return fmt.Errorf("Please provide the required attributes export_to_osu and image_id")
	}

	request.SnapshotId = aws.String(v.(string))

	if etoOk {
		exp := eto.([]interface{})
		e := exp[0].(map[string]interface{})

		et := &fcu.ExportToOsuTaskSpecification{}

		if v, ok := e["disk_image_format"]; ok {
			et.DiskImageFormat = aws.String(v.(string))
		}
		if v, ok := e["osu_key"]; ok {
			et.OsuKey = aws.String(v.(string))
		}
		if v, ok := e["osu_bucket"]; ok {
			et.OsuBucket = aws.String(v.(string))
		}
		if v, ok := e["osu_prefix"]; ok {
			et.OsuPrefix = aws.String(v.(string))
		}
		if v, ok := e["aksk"]; ok {
			a := v.([]interface{})
			if len(a) > 0 {
				w := a[0].(map[string]interface{})
				et.AkSk = &fcu.ExportToOsuAccessKeySpecification{
					AccessKey: aws.String(w["access_key"].(string)),
					SecretKey: aws.String(w["secret_key"].(string)),
				}
			}
		}
		request.ExportToOsu = et
	}

	var resp *fcu.CreateSnapshotExportTaskOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateSnapshotExportTask(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error image task %s", err)
	}

	id := *resp.SnapshotExportTask.SnapshotExportTaskId
	d.SetId(id)

	_, err = resourceOutscaleSnapshotTaskWaitForAvailable(id, conn, 1)
	if err != nil {
		return err
	}

	return resourceImageExportTasksRead(d, meta)
}

func resourceImageExportTasksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeSnapshotExportTasksOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeSnapshotExportTasks(&fcu.DescribeSnapshotExportTasksInput{
			SnapshotExportTaskId: []*string{aws.String(d.Id())},
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading task image %s", err)
	}

	v := resp.SnapshotExportTask[0]

	d.Set("completion", v.Completion)
	d.Set("snapshot_export_task_id", v.SnapshotExportTaskId)
	d.Set("state", v.State)
	d.Set("completion", v.Completion)
	if v.StatusMessage != nil {
		d.Set("status_message", v.StatusMessage)
	} else {
		d.Set("status_message", "")
	}

	exp := make([]map[string]interface{}, 1)
	exportToOsu := make(map[string]interface{})
	exportToOsu["disk_image_format"] = *v.ExportToOsu.DiskImageFormat
	exportToOsu["osu_bucket"] = *v.ExportToOsu.OsuBucket
	exportToOsu["osu_key"] = *v.ExportToOsu.OsuKey
	if v.ExportToOsu.OsuPrefix != nil {
		exportToOsu["osu_prefix"] = *v.ExportToOsu.OsuPrefix
	} else {
		exportToOsu["osu_prefix"] = ""
	}

	aksk := make([]map[string]interface{}, 1)
	osuAkSk := make(map[string]interface{})
	if v.ExportToOsu.AkSk != nil {
		osuAkSk["access_key"] = *v.ExportToOsu.AkSk.AccessKey
		osuAkSk["secret_key"] = *v.ExportToOsu.AkSk.SecretKey
	} else {
		osuAkSk["access_key"] = ""
		osuAkSk["secret_key"] = ""
	}
	aksk[0] = osuAkSk
	exportToOsu["aksk"] = aksk

	snapExp := make(map[string]interface{})
	snapExp["snapshot_id"] = *v.SnapshotExport.SnapshotId

	d.Set("snapshot_export", snapExp)
	exp[0] = exportToOsu
	if err := d.Set("export_to_osu", exp); err != nil {
		return err
	}
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceImageExportTasksDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")
	d.Set("snapshot_export", nil)
	d.Set("export_to_osu", nil)
	d.Set("request_id", nil)

	return nil
}

func resourceOutscaleSnapshotTaskWaitForAvailable(id string, client *fcu.Client, i int) (*fcu.SnapshotExportTask, error) {
	log.Printf("Waiting for Image Task %s to become available...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued"},
		Target:     []string{"active"},
		Refresh:    SnapshotTaskStateRefreshFunc(client, id),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for OMI (%s) to be ready: %v", id, err)
	}
	return info.(*fcu.SnapshotExportTask), nil
}

// SnapshotTaskStateRefreshFunc ...
func SnapshotTaskStateRefreshFunc(client *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &fcu.DescribeSnapshotExportTasksOutput{}

		var resp *fcu.DescribeSnapshotExportTasksOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = client.VM.DescribeSnapshotExportTasks(&fcu.DescribeSnapshotExportTasksInput{
				SnapshotExportTaskId: []*string{aws.String(id)},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAMIID.NotFound") {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil

			} else if resp != nil && len(resp.SnapshotExportTask) == 0 {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp == nil || resp.SnapshotExportTask == nil || len(resp.SnapshotExportTask) == 0 {
			return emptyResp, "destroyed", nil
		}

		// OMI is valid, so return it's state
		return resp.SnapshotExportTask[0], *resp.SnapshotExportTask[0].State, nil
	}
}
