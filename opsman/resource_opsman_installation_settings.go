package opsman

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gosuri/uilive"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pivotal-cf/om/api"
	omcmd "github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/progress"
)

func resourceOpsmanInstallationSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsmanInstallationSettingsCreate,
		Read:   resourceOpsmanInstallationSettingsRead,
		Update: resourceOpsmanInstallationSettingsUpdate,
		Delete: resourceOpsmanInstallationSettingsDelete,

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"installation_settings_file": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOpsmanInstallationSettingsCreate(d *schema.ResourceData, meta interface{}) error {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("installation-settings-%s", id))
	return resourceOpsmanInstallationSettingsRead(d, meta)
}

func resourceOpsmanInstallationSettingsRead(d *schema.ResourceData, meta interface{}) error {
	opsmanExportSettingsFile := filepath.Join(os.TempDir(), "installation.zip")
	log.Printf("[INFO] Exporting the OpsMan installation settings to %s", opsmanExportSettingsFile)

	// set a default so at least TF won't error out about a missing attribute
	d.Set("installation_settings_file", "")

	authedClient, err := createAuthedClient(d)
	if err != nil {
		// if we can't export the settings, that's OK - happens before OpsMan exists
		return nil
	}

	exportInstallationService := api.NewInstallationAssetService(authedClient, progress.NewBar(), uilive.New())
	exportInstallationCmd := omcmd.NewExportInstallation(exportInstallationService, log.New(omWriter{}, "", 0))

	err = exportInstallationCmd.Execute([]string{
		"--output-file", opsmanExportSettingsFile,
	})
	if err != nil {
		// if we can't export the settings, that's OK - happens before OpsMan exists
		log.Printf("Failed to export OpsMan installation settings: %s", err.Error())
		return nil
	}

	// we exported the settings, set the file path
	d.Set("installation_settings_file", opsmanExportSettingsFile)

	return nil
}

func resourceOpsmanInstallationSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOpsmanInstallationSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
