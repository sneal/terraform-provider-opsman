package opsman

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"

	"github.com/gosuri/uilive"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pivotal-cf/om/api"
	omcmd "github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/formcontent"
	"github.com/pivotal-cf/om/progress"
)

const applySleepSeconds = 10

func resourceOpsmanDirector() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsmanDirectorCreate,
		Read:   resourceOpsmanDirectorRead,
		Update: resourceOpsmanDirectorUpdate,
		Delete: resourceOpsmanDirectorDelete,

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
			"decryption_passphrase": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"installation_settings_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_access_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_private_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zones": { // can this just come from all AZs found in networks?
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"director_network": { // this should probably just default to the same network as opsman
				Type:     schema.TypeString,
				Required: true,
			},
			"database": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
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
					},
				},
			},
			"blobstore": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"s3_endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"network": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_subnet_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cidr": {
										Type:     schema.TypeString,
										Required: true,
									},
									"reserved_ip_ranges": {
										Type:     schema.TypeString,
										Required: true,
									},
									"dns": {
										Type:     schema.TypeString,
										Required: true,
									},
									"gateway": {
										Type:     schema.TypeString,
										Required: true,
									},
									"availability_zone": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type networkAndAz struct {
	Network                   network              `json:"network"`
	SingletonAvailabilityZone availabilityZoneName `json:"singleton_availability_zone"`
}

type availabilityZoneName struct {
	Name string `json:"name"`
}

type subnet struct {
	IaasIdentifier    string   `json:"iaas_identifier"`
	Cidr              string   `json:"cidr"`
	ReservedIPRanges  string   `json:"reserved_ip_ranges"`
	DNS               string   `json:"dns"`
	Gateway           string   `json:"gateway"`
	AvailabilityZones []string `json:"availability_zone_names"`
}

type network struct {
	Name    string   `json:"name"`
	Subnets []subnet `json:"subnets"`
}

type networkConfiguration struct {
	IcmpChecksEnabled bool      `json:"icmp_checks_enabled"`
	Networks          []network `json:"networks"`
}

type externalDatabaseOptions struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type s3BlobstoreOptions struct {
	Endpoint         string `json:"endpoint"`
	BucketName       string `json:"bucket_name"`
	AccessKey        string `json:"access_key"`
	SecretKey        string `json:"secret_key"`
	SignatureVersion string `json:"signature_version"`
	Region           string `json:"region"`
}

type directorConfiguration struct {
	NtpServers              string                  `json:"ntp_servers_string"`
	ResurrectorEnabled      bool                    `json:"resurrector_enabled"`
	MaxThreads              uint8                   `json:"max_threads"`
	DatabaseType            string                  `json:"database_type"`
	ExternalDatabaseOptions externalDatabaseOptions `json:"external_database_options"`
	BlobstoreType           string                  `json:"blobstore_type"`
	S3BlobstoreOptions      s3BlobstoreOptions      `json:"s3_blobstore_options"`
}

type awsConfiguration struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	VpcID           string `json:"vpc_id"`
	SecurityGroup   string `json:"security_group"`
	KeyPairName     string `json:"key_pair_name"`
	SSHPrivateKey   string `json:"ssh_private_key"`
	Region          string `json:"region"`
	Encrypted       bool   `json:"encrypted"`
}

func resourceOpsmanDirectorCreate(d *schema.ResourceData, meta interface{}) error {
	_, err := waitForOpsmanUp(d.Get("address").(string))
	if err != nil {
		return err
	}
	err = configureOpsmanInitialAuth(d)
	if err != nil {
		return err
	}
	return configureAndDeployDirector(d)
}

func resourceOpsmanDirectorRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOpsmanDirectorUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("instance_id") {
		settingsFile := d.Get("installation_settings_file").(string)
		if settingsFile == "" {
			return fmt.Errorf("Changing the OpsMan instance without having an " +
				"associated installation_settings_file will cause all OpsMan state to be lost")
		}
		if d.HasChange("decryption_passphrase") {
			return fmt.Errorf("Changing the OpsMan instance and the decryption passphrase at the same " +
				"time will cause the installation settings import to fail")
		}

		// since we're changing instances, we need to wait until the new instance is ready
		_, err := waitForOpsmanUp(d.Get("address").(string))
		if err != nil {
			return err
		}

		log.Printf("[INFO] Importing OpsMan installation settings into new instance %s", d.Get("instance_id").(string))
		return importInstallationSettings(d)
	}

	return nil
}

func resourceOpsmanDirectorDelete(d *schema.ResourceData, meta interface{}) error {
	// Since there's no dependency between the director and network resources, those
	// resources are getting destroyed before we can make/finish the destroy director
	// API call to OpsMan.

	// For now just orphan the bosh director and anything it has deployed...
	return nil

	/*
		log.Print("[INFO] Deleting opsman bosh director")
		authedClient, err := createAuthedClient(d)
		if err != nil {
			return err
		}
		deleteInstallationService := api.NewInstallationAssetService(authedClient, nil, nil)
		installationsService := api.NewInstallationsService(authedClient)
		logWriter := omcmd.NewLogWriter(omWriter{})
		logger := log.New(omWriter{}, "", 0)
		delInstallationCmd := omcmd.NewDeleteInstallation(deleteInstallationService, installationsService, logWriter, logger, applySleepSeconds)

		return delInstallationCmd.Execute([]string{})
	*/
}

func configureOpsmanInitialAuth(d *schema.ResourceData) error {
	log.Printf("[INFO] Configuring OpsMan authentication system")
	setupService := api.NewSetupService(createUnauthedClient(d))
	configAuthCmd := omcmd.NewConfigureAuthentication(setupService, log.New(os.Stdout, "", 0))
	err := configAuthCmd.Execute([]string{
		"--username", d.Get("username").(string),
		"--password", d.Get("password").(string),
		"--decryption-passphrase", d.Get("decryption_passphrase").(string),
	})
	return err
}

func configureAndDeployDirector(d *schema.ResourceData) error {
	log.Printf("[INFO] Configuring OpsMan director tile")

	// map director network assignment
	networkAndAz := &networkAndAz{
		Network: network{
			Name: d.Get("director_network").(string),
		},
		SingletonAvailabilityZone: availabilityZoneName{
			Name: d.Get("availability_zone").(string),
		},
	}

	networkAndAzJSON, err := json.Marshal(networkAndAz)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Network assignment JSON: %s", networkAndAzJSON)

	// map AZs
	availabilityZones := []availabilityZoneName{}

	if v, ok := d.GetOk("availability_zones"); ok {
		for _, n := range v.([]interface{}) {
			az := availabilityZoneName{
				Name: n.(string),
			}
			availabilityZones = append(availabilityZones, az)
		}
	}

	availabilityZonesJSON, err := json.Marshal(availabilityZones)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] AZs JSON: %s", availabilityZonesJSON)

	// map bosh networks
	networkConfiguration := &networkConfiguration{
		IcmpChecksEnabled: false,
		Networks:          []network{},
	}
	for _, rawNetwork := range d.Get("network").(*schema.Set).List() {
		netMap := rawNetwork.(map[string]interface{})
		network := network{Name: netMap["name"].(string), Subnets: []subnet{}}
		for _, rawSubnet := range netMap["subnet"].(*schema.Set).List() {
			subnetMap := rawSubnet.(map[string]interface{})
			subnet := subnet{
				IaasIdentifier:    subnetMap["vpc_subnet_id"].(string),
				Cidr:              subnetMap["cidr"].(string),
				ReservedIPRanges:  subnetMap["reserved_ip_ranges"].(string),
				DNS:               subnetMap["dns"].(string),
				Gateway:           subnetMap["gateway"].(string),
				AvailabilityZones: []string{subnetMap["availability_zone"].(string)},
			}
			network.Subnets = append(network.Subnets, subnet)
		}
		networkConfiguration.Networks = append(networkConfiguration.Networks, network)
	}
	networkJSON, err := json.Marshal(networkConfiguration)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Bosh network JSON: %s", networkJSON)

	// map director configuration
	rawDatabase := d.Get("database").(*schema.Set).List()[0]
	dbMap := rawDatabase.(map[string]interface{})

	rawBlobstore := d.Get("blobstore").(*schema.Set).List()[0]
	blobstoreMap := rawBlobstore.(map[string]interface{})

	directorConfiguration := &directorConfiguration{
		NtpServers:         "0.amazon.pool.ntp.org,1.amazon.pool.ntp.org,2.amazon.pool.ntp.org,3.amazon.pool.ntp.org",
		ResurrectorEnabled: true,
		MaxThreads:         30,
		DatabaseType:       "external",

		ExternalDatabaseOptions: externalDatabaseOptions{
			Host:     dbMap["host"].(string),
			Port:     uint16(dbMap["port"].(int)),
			User:     dbMap["username"].(string),
			Password: dbMap["password"].(string),
			Database: "bosh", // db name
		},
		BlobstoreType: "s3",
		S3BlobstoreOptions: s3BlobstoreOptions{
			Endpoint:         blobstoreMap["s3_endpoint"].(string),
			BucketName:       blobstoreMap["bucket_name"].(string),
			AccessKey:        d.Get("access_key_id").(string),
			SecretKey:        d.Get("secret_access_key").(string),
			SignatureVersion: "4",
			Region:           d.Get("region").(string),
		},
	}

	directorJSON, err := json.Marshal(directorConfiguration)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Director configuration JSON: %s", directorJSON)

	// AWS specific IaaS configuration
	awsConfiguration := &awsConfiguration{
		AccessKeyID:     d.Get("access_key_id").(string),
		SecretAccessKey: d.Get("secret_access_key").(string),
		VpcID:           d.Get("vpc_id").(string),
		SecurityGroup:   d.Get("vpc_security_group_id").(string),
		KeyPairName:     d.Get("key_name").(string),
		SSHPrivateKey:   d.Get("ssh_private_key").(string),
		Region:          d.Get("region").(string),
		Encrypted:       false,
	}

	awsJSON, err := json.Marshal(awsConfiguration)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] AWS IaaS configuration JSON: %s", awsJSON)

	// create om services
	authedClient, err := createAuthedClient(d)
	if err != nil {
		return err
	}

	logger := log.New(omWriter{}, "", 0)
	jobsService := api.NewJobsService(authedClient)
	directorService := api.NewDirectorService(authedClient)
	stagedProductsService := api.NewStagedProductsService(authedClient)
	configDirectorCmd := omcmd.NewConfigureDirector(directorService, jobsService, stagedProductsService, logger)

	// configure the director tile
	err = configDirectorCmd.Execute([]string{
		"--network-assignment", string(networkAndAzJSON),
		"--az-configuration", string(availabilityZonesJSON),
		"--networks-configuration", string(networkJSON),
		"--director-configuration", string(directorJSON),
		"--iaas-configuration", string(awsJSON),
		"--security-configuration", `{"trusted_certificates": "", "vm_password_type": "generate"}`,
		//"--syslog-configuration", `{"some-syslog-assignment": "syslog"}`,
		//"--resource-configuration", `{"resource": {"instance_type": {"id": "some-type"}}}`,
	})

	if err != nil {
		return err
	}

	// Apply Changes
	installationsService := api.NewInstallationsService(authedClient)
	logWriter := omcmd.NewLogWriter(omWriter{})
	applyChangesCmd := omcmd.NewApplyChanges(installationsService, logWriter, logger, applySleepSeconds)
	err = applyChangesCmd.Execute([]string{"--skip-deploy-products"})
	if err != nil {
		return err
	}

	// success, save state
	c := hashcode.String(d.Get("address").(string))
	d.SetId(fmt.Sprintf("opsman-%d", c))

	return nil
}

func importInstallationSettings(d *schema.ResourceData) error {
	logger := log.New(omWriter{}, "", 0)
	unauthedClient := createUnauthedClient(d)
	importInstallationService := api.NewInstallationAssetService(unauthedClient, progress.NewBar(), uilive.New())
	setupService := api.NewSetupService(createUnauthedClient(d))
	form, err := formcontent.NewForm()
	if err != nil {
		return err
	}

	importInstallationCmd := omcmd.NewImportInstallation(form, importInstallationService, setupService, logger)
	return importInstallationCmd.Execute([]string{
		"--installation", d.Get("installation_settings_file").(string),
		"--decryption-passphrase", d.Get("decryption_passphrase").(string),
	})
}

func waitForOpsmanUp(address string) (int, error) {
	log.Printf("[DEBUG] Waiting for (%s) to return positive HTTP status", address)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"running"},
		Refresh:    opsmanStateRefreshFunc(address),
		Timeout:    time.Duration(5) * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	statusCode, err := stateConf.WaitForState()
	return statusCode.(int), err
}

func opsmanStateRefreshFunc(address string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		rs, err := client.Get(fmt.Sprintf("https://%s", address))
		if err != nil {
			log.Printf("[DEBUG] Checking OpsMan availability err: %s", err.Error())
			return nil, "", nil
		}
		defer rs.Body.Close()
		if rs.StatusCode > 299 {
			log.Printf("[DEBUG] Checking OpsMan availability non-success status code: %d", rs.StatusCode)
			return nil, "", nil
		}
		log.Print("[DEBUG] Checking OpsMan availability complete")
		return rs.StatusCode, "running", nil
	}
}
