package azurerm

import (
	"fmt"
	"log"
	//	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/frontdoor/mgmt/2019-04-01/frontdoor"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	//	"github.com/hashicorp/terraform/helper/validation"
	helpers "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFrontDoor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFrontDoorCreateUpdate,
		Read:   resourceArmFrontDoorRead,
		Update: resourceArmFrontDoorCreateUpdate,
		Delete: resourceArmFrontDoorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: helpers.ValidateFrontDoorName,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmFrontDoorCreate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Not implemented")
}
func resourceArmFrontDoorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontDoorsClient
	managementClient := meta.(*ArmClient).frontDoorManagementClient

	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Front Door creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Azure Front Door %s (resource group %s) ID", name, resGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_frontdoor", *existing.ID)
		}
	}

	// Check that the front door name is available
	nameCheckInput := frontdoor.CheckNameAvailabilityInput{
		Name: &name,
		Type: frontdoor.MicrosoftNetworkfrontDoors,
	}
	nameCheckResult, err := managementClient.CheckFrontDoorNameAvailability(ctx, nameCheckInput)
	if err != nil {
		return fmt.Errorf("Error checking the name availability for Azure Front Door %s (resource group %s)", name, resGroup)
	}
	if nameCheckResult.NameAvailability == "Unavailable" {
		return fmt.Errorf("The Azure Front Door name you selected (%s) is unavailable: %s.", &nameCheckResult.Reason)
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	enabled := d.Get("enabled").(bool)
	tags := d.Get("tags").(map[string]interface{})
	friendlyName := d.Get("friendly_name").(string)

	state := frontdoor.EnabledStateEnabled
	if !enabled {
		state = frontdoor.EnabledStateDisabled
	}

	frontDoorProperties := frontdoor.FrontDoor{
		Name:     &name,
		Type:     utils.String(string(frontdoor.MicrosoftNetworkfrontDoors)),
		Location: &location,
		Properties: &frontdoor.Properties{
			EnabledState: state,
			FriendlyName: &friendlyName,
		},
		Tags: expandTags(tags),
	}

	createFuture, err := client.CreateOrUpdate(ctx, resGroup, name, frontDoorProperties)
	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceRead(d, meta)
}

func resourceArmFrontDoorRead(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Not implemented")
}

func resourceArmFrontDoorDelete(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Not implemented")
}
