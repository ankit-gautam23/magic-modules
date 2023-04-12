package google_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	google "internal/terraform-provider-google"
	"testing"
)

func TestAccComputeDiskResourcePolicyAttachment_update(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", RandString(t, 10))
	policyName := fmt.Sprintf("tf-test-policy-%s", RandString(t, 10))
	policyName2 := fmt.Sprintf("tf-test-policy-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName),
			},
			{
				ResourceName: "google_compute_disk_resource_policy_attachment.foobar",
				// ImportStateId:     fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", GetTestProjectFromEnv(), "us-central1", policyName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName2),
			},
			{
				ResourceName: "google_compute_disk_resource_policy_attachment.foobar",
				// ImportStateId:     fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", GetTestProjectFromEnv(), "us-central1", policyName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 1000
  type  = "pd-extreme"
  zone  = "us-central1-c"
  labels = {
    my-label = "my-label-value"
  }
  provisioned_iops = 90000
}

resource "google_compute_resource_policy" "foobar" {
  name = "%s"
  region = "us-central1"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time = "04:00"
      }
    }
  }
}

resource "google_compute_disk_resource_policy_attachment" "foobar" {
  name = google_compute_resource_policy.foobar.name
  disk = google_compute_disk.foobar.name
  zone = "us-central1-c"
}
`, diskName, policyName)
}
