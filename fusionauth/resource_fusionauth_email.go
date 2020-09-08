package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func newEmail() *schema.Resource {
	return &schema.Resource{
		Create: createEmail,
		Read:   readEmail,
		Update: updateEmail,
		Delete: deleteEmail,
		Schema: map[string]*schema.Schema{
			"default_from_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The default From Name used when sending emails. If not provided, and a localized value cannot be determined, the default value for the tenant will be used. This is the display name part of the email address ( i.e. Jared Dunn <jared@piedpiper.com>).",
			},
			"default_html_template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The default HTML Email Template.",
			},
			"default_subject": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The default Subject used when sending emails.",
			},
			"default_text_template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The default Text Email Template.",
			},
			"from_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email address that this email will be sent from. If not provided, the default value for the tenant will be used. This is the address part email address (i.e. Jared Dunn <jared@piedpiper.com>).",
			},
			"localized_from_names": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The From Name used when sending emails to users who speak other languages. This overrides the default From Name based on the user’s list of preferred languages.",
			},
			"localized_html_templates": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The HTML Email Template used when sending emails to users who speak other languages. This overrides the default HTML Email Template based on the user’s list of preferred languages.",
			},
			"localized_subjects": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The Subject used when sending emails to users who speak other languages. This overrides the default Subject based on the user’s list of preferred languages.",
			},
			"localized_text_templates": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The Text Email Template used when sending emails to users who speak other languages. This overrides the default Text Email Template based on the user’s list of preferred languages.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `A descriptive name for the email template (i.e. "April 2016 Coupon Email")`,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func buildEmail(data *schema.ResourceData) fusionauth.EmailTemplate {
	e := fusionauth.EmailTemplate{
		DefaultFromName:     data.Get("default_from_name").(string),
		DefaultHtmlTemplate: data.Get("default_html_template").(string),
		DefaultSubject:      data.Get("default_subject").(string),
		DefaultTextTemplate: data.Get("default_text_template").(string),
		FromEmail:           data.Get("from_email").(string),
		Name:                data.Get("name").(string),
	}

	if i, ok := data.GetOk("localized_from_names"); ok {
		e.LocalizedFromNames = i.(map[string]string)
	}

	if i, ok := data.GetOk("localized_html_templates"); ok {
		e.LocalizedHtmlTemplates = i.(map[string]string)
	}

	if i, ok := data.GetOk("localized_subjects"); ok {
		e.LocalizedSubjects = i.(map[string]string)
	}

	if i, ok := data.GetOk("localized_text_templates"); ok {
		e.LocalizedTextTemplates = i.(map[string]string)
	}
	return e
}

func createEmail(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	e := buildEmail(data)

	resp, faErrs, err := client.FAClient.CreateEmailTemplate("", fusionauth.EmailTemplateRequest{
		EmailTemplate: e,
	})

	if err != nil {
		return fmt.Errorf("CreateEmailTemplate err: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("CreateEmailTemplate errors: %v", faErrs)
	}

	data.SetId(resp.EmailTemplate.Id)
	return nil
}

func readEmail(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveEmailTemplate(id)
	if err != nil {
		return err
	}

	t := resp.EmailTemplate
	_ = data.Set("default_from_name", t.DefaultFromName)
	_ = data.Set("default_html_template", t.DefaultHtmlTemplate)
	_ = data.Set("default_subject", t.DefaultSubject)
	_ = data.Set("default_text_template", t.DefaultTextTemplate)
	_ = data.Set("from_email", t.FromEmail)
	_ = data.Set("localized_from_names", t.LocalizedFromNames)
	_ = data.Set("localized_html_templates", t.LocalizedHtmlTemplates)
	_ = data.Set("localized_subjects", t.LocalizedSubjects)
	_ = data.Set("localized_text_templates", t.LocalizedTextTemplates)
	_ = data.Set("name", t.Name)

	return nil
}

func updateEmail(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	e := buildEmail(data)

	_, faErrs, err := client.FAClient.UpdateEmailTemplate(data.Id(), fusionauth.EmailTemplateRequest{
		EmailTemplate: e,
	})

	if err != nil {
		return fmt.Errorf("UpdateEmailTemplate err: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("UpdateEmailTemplate errors: %v", faErrs)
	}

	return nil
}

func deleteEmail(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	_, faErrs, err := client.FAClient.DeleteEmailTemplate(id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("DeleteEmailTemplate errors: %v", faErrs)
	}

	return nil
}
