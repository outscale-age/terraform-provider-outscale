package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPINicsDataSource(t *testing.T) {
	subregion := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPINicsDataSourceConfig(subregion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_nics.outscale_nics", "nics.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPINicsDataSourceConfig(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}
		
		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = "${outscale_net.outscale_net.net_id}"
		}
		
		resource "outscale_nic" "outscale_nic" {
			subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
		}
		
		data "outscale_nics" "outscale_nics" {
			filter {
				name   = "nic_ids"
				values = ["${outscale_nic.outscale_nic.id}"]
			}
		}
	`, subregion)
}

func testAccCheckOutscaleNicsDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find NICS data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("NICS data source ID not set")
		}
		return nil
	}
}
