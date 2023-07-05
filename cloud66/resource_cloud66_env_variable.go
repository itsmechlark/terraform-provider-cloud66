package cloud66

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/itsmechlark/cloud66"
)

func resourceCloud66EnvVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloud66EnvVariableCreate,
		Read:   resourceCloud66EnvVariableRead,
		Update: resourceCloud66EnvVariableUpdate,
		Delete: resourceCloud66EnvVariableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloud66EnvVariableImport,
		},

		SchemaVersion: 2,
		Schema:        resourceCloud66EnvVariableSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func resourceCloud66EnvVariableCreate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Get("stack_id").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	applyStrategy := d.Get("apply_strategy").(string)

	log.Printf("[INFO] Creating %s Env Variable for stack %s", key, stackID)

	record, err := client.StackEnvVarNew(stackID, key, value, applyStrategy)

	if record == nil {
		return fmt.Errorf("error creating Env Variable %q: %s", stackID, err)
	}

	envVar := api.StackEnvVar{
		Key:      key,
		Value:    value,
		Readonly: d.Get("readonly").(bool),
	}

	setCloud66EnvVariableData(d, &envVar)

	return nil
}

func resourceCloud66EnvVariableRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Get("stack_id").(string)
	key := d.Id()

	records, err := client.StackEnvVars(stackID)
	if records != nil {
		for _, record := range records {
			if record.Key == key {
				setCloud66EnvVariableData(d, &record)
				break
			}
		}
	} else {
		return fmt.Errorf("error reading Env Variable %q: %s", stackID, err)
	}

	return nil
}

func resourceCloud66EnvVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Get("stack_id").(string)
	key := d.Id()
	value := d.Get("value").(string)
	applyStrategy := d.Get("apply_strategy").(string)

	log.Printf("[INFO] Updating %s Env Variable for stack %s", key, stackID)

	record, err := client.StackEnvVarSet(stackID, key, value, applyStrategy)

	if record == nil {
		return fmt.Errorf("error updating Env Variable %q: %s", stackID, err)
	}

	envVar := api.StackEnvVar{
		Key:      key,
		Value:    value,
		Readonly: d.Get("readonly").(bool),
	}

	setCloud66EnvVariableData(d, &envVar)

	return nil
}

func resourceCloud66EnvVariableDelete(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Get("stack_id").(string)
	key := d.Id()

	req, err := client.NewRequest("DELETE", "/stacks/"+stackID+"/environments/"+key+".json", nil, nil)
	if req == nil {
		return fmt.Errorf("error deleting Env Variable %q: %s", stackID, err)
	}

	return nil
}

func resourceCloud66EnvVariableImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// split the id so we can lookup
	idAttr := strings.SplitN(d.Id(), "/", 2)
	if len(idAttr) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"stackID/key\"", d.Id())
	}

	stackID, key := idAttr[0], idAttr[1]

	log.Printf("[DEBUG] Importing %s Env Variable for stack %s", key, stackID)

	d.Set("stack_id", stackID)
	d.SetId(key)

	resourceCloud66EnvVariableRead(d, meta)

	return []*schema.ResourceData{d}, nil
}

func setCloud66EnvVariableData(d *schema.ResourceData, envVar *api.StackEnvVar) {
	stackID := d.Get("stack_id").(string)

	d.SetId(envVar.Key)
	d.Set("stack_id", stackID)
	d.Set("key", envVar.Key)
	d.Set("value", envVar.Value)
	d.Set("readonly", envVar.Readonly)
}
