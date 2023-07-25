package cloud66

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/itsmechlark/cloud66"
)

func resourceCloud66SslCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloud66SslCertificateCreate,
		Read:   resourceCloud66SslCertificateRead,
		Update: resourceCloud66SslCertificateUpdate,
		Delete: resourceCloud66SslCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloud66SslCertificateImport,
		},

		SchemaVersion: 2,
		Schema:        resourceCloud66SslCertificateSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func resourceCloud66SslCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Get("stack_id").(string)
	allowOverwrite := d.Get("allow_overwrite").(bool)

	servernames := []string{}
	servernamesRaw := d.Get("server_names").(*schema.Set)
	for _, h := range servernamesRaw.List() {
		servernames = append(servernames, h.(string))
	}

	newRecord := api.SslCertificate{
		Type:           d.Get("type").(string),
		ServerNames:    strings.Join(servernames, ","),
		SSLTermination: d.Get("ssl_termination").(bool),
	}

	if newRecord.Type == "manual" {
		certificate := d.Get("certificate").(string)
		key := d.Get("key").(string)
		intermediateCertificate := d.Get("intermediate_certificate").(string)
		if certificate != "" && key != "" {
			newRecord.Certificate = &certificate
			newRecord.Key = &key
			newRecord.IntermediateCertificate = &intermediateCertificate
		} else {
			return fmt.Errorf("certificate and key must be set when type is manual")
		}
	}

	record := &api.SslCertificate{}
	err := fmt.Errorf("")

	if allowOverwrite {
		log.Printf("[INFO] Fetching default SSL Certfor stack %s", stackID)
		// check if there is already a certificate
		certs, err := client.ListSslCertificates(stackID)
		if err == nil && len(certs) > 0 {
			cert := &certs[0]
			log.Printf("[INFO] Overwriting SSL Cert %s UUID for stack %s", cert.Uuid, stackID)
			record, err = client.UpdateSslCertificate(stackID, cert.Uuid, &newRecord)

			if record == nil {
				log.Print(fmt.Errorf("error overwriting SSL Certificate %q: %s", stackID, err))
			} else {
				record.Type = d.Get("type").(string)
			}
		}
	}

	if !allowOverwrite || record == nil {
		log.Printf("[INFO] Creating SSL Cert %s for stack %s", newRecord.Type, stackID)
		record, err = client.CreateSslCertificate(stackID, &newRecord)
	}

	if record == nil {
		return fmt.Errorf("error creating SSL Certificate %q: %s", stackID, err)
	}

	setCloud66SslCertificateData(d, record)

	return nil
}

func resourceCloud66SslCertificateRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Id()

	records, err := client.ListSslCertificates(stackID)
	if records != nil {
		for _, record := range records {
			sha256Fingerprint := *record.SHA256Fingerprint
			if sha256Fingerprint == d.Get("sha256_fingerprint").(string) {
				setCloud66SslCertificateData(d, &record)
				break
			}
		}
	} else {
		return fmt.Errorf("error reading SSL Certificate %q: %s", stackID, err)
	}

	return nil
}

func resourceCloud66SslCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Id()
	sslID := d.Get("uuid").(string)

	servernames := []string{}
	servernamesRaw := d.Get("server_names").(*schema.Set)
	for _, h := range servernamesRaw.List() {
		servernames = append(servernames, h.(string))
	}

	newRecord := api.SslCertificate{
		Type:           d.Get("type").(string),
		ServerNames:    strings.Join(servernames, ","),
		SSLTermination: d.Get("ssl_termination").(bool),
	}

	if newRecord.Type == "manual" {
		certificate := d.Get("certificate").(string)
		key := d.Get("key").(string)
		intermediateCertificate := d.Get("intermediate_certificate").(string)
		if certificate != "" && key != "" {
			newRecord.Certificate = &certificate
			newRecord.Key = &key
			newRecord.IntermediateCertificate = &intermediateCertificate
		} else {
			return fmt.Errorf("certificate and key must be set when type is manual")
		}
	}

	log.Printf("[INFO] Updating SSL Cert %s for stack %s", sslID, stackID)

	record, err := client.UpdateSslCertificate(stackID, sslID, &newRecord)

	if record == nil {
		return fmt.Errorf("error updating SSL Certificate %q: %s", stackID, err)
	}

	setCloud66SslCertificateData(d, record)

	return nil
}

func resourceCloud66SslCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Id()
	sslID := d.Get("uuid").(string)

	records, err := client.ListSslCertificates(stackID)
	if records != nil {
		for _, record := range records {
			sha256Fingerprint := *record.SHA256Fingerprint
			if sha256Fingerprint == d.Get("sha256_fingerprint").(string) || (sslID != "" && sslID == record.Uuid) {
				log.Printf("[DEBUG] Deleting SSL Cert %s for stack %s", sslID, stackID)
				record, err := client.DestroySslCertificate(stackID, sslID)

				if record == nil {
					return fmt.Errorf("error deleting SSL Certificate %q: %s", stackID, err)
				}
				break
			}
		}
	} else {
		return fmt.Errorf("error reading SSL Certificate %q: %s", stackID, err)
	}

	return nil
}

func resourceCloud66SslCertificateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// split the id so we can lookup
	idAttr := strings.SplitN(d.Id(), "/", 2)
	if len(idAttr) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"stackID/sha256Fingerprint\"", d.Id())
	}

	stackID, sha256Fingerprint := idAttr[0], idAttr[1]

	log.Printf("[DEBUG] Importing SSL Cert %s for stack %s", sha256Fingerprint, stackID)

	d.Set("sha256_fingerprint", sha256Fingerprint)
	d.SetId(stackID)

	resourceCloud66SslCertificateRead(d, meta)

	return []*schema.ResourceData{d}, nil
}

func setCloud66SslCertificateData(d *schema.ResourceData, ssl *api.SslCertificate) {
	stackID := d.Get("stack_id").(string)

	d.SetId(stackID)
	d.Set("uuid", ssl.Uuid)
	d.Set("name", ssl.Name)
	d.Set("ca_name", ssl.CAName)
	d.Set("type", ssl.Type)
	d.Set("ssl_termination", ssl.SSLTermination)
	d.Set("server_group_id", ssl.ServerGroupID)
	d.Set("has_intermediate_cert", ssl.HasIntermediateCert)
	d.Set("sha256_fingerprint", ssl.SHA256Fingerprint)
	d.Set("status", ssl.Status())

	servernames := schema.NewSet(schema.HashString, []interface{}{})
	for _, h := range strings.Split(ssl.ServerNames, ",") {
		servernames.Add(h)
	}
	d.Set("server_names", servernames)
}
