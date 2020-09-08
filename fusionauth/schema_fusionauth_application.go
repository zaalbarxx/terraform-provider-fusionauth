package fusionauth

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func newApplication() *schema.Resource {
	return &schema.Resource{
		Create: createApplication,
		Read:   readApplication,
		Update: updateApplication,
		Delete: deleteApplication,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id to use for the new Application. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"authentication_token_configuration_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if Users can have Authentication Tokens associated with this Application. This feature may not be enabled for the FusionAuth application.",
			},
			"clean_speak_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_ids": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "An array of UUIDs that map to the CleanSpeak applications for this Application. It is possible that a single Application in FusionAuth might have multiple Applications in CleanSpeak. For example, a FusionAuth Application for a game might have one CleanSpeak Application for usernames and another Application for chat.",
						},
						"username_moderation": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "True if CleanSpeak username moderation is enabled.",
									},
									"application_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Id of the CleanSpeak application that usernames are sent to for moderation.",
									},
								},
							},
						},
					},
				},
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the Application that should be persisted.",
			},
			"jwt_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     newJWTConfiguration(),
				Optional: true,
			},
			"lambda_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_token_populate_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Id of the Lambda that will be invoked when an access token is generated for this application. This will be utilized during OAuth2 and OpenID Connect authentication requests as well as when an access token is generated for the Login API.",
						},
						"id_token_populate_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Id of the Lambda that will be invoked when an Id token is generated for this application during an OpenID Connect authentication request.",
						},
						"samlv2_populate_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Id of the Lambda that will be invoked when a a SAML response is generated during a SAML authentication request.",
						},
					},
				},
			},
			"login_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_token_refresh": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates if a JWT may be refreshed using a Refresh Token for this application. This configuration is separate from issuing new Refresh Tokens which is controlled by the generateRefreshTokens parameter. This configuration indicates specifically if an existing Refresh Token may be used to request a new JWT using the Refresh API.",
						},
						"generate_refresh_tokens": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates if a Refresh Token should be issued from the Login API",
						},
						"require_authentication": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Indicates if the Login API should require an API key. If you set this value to false and your FusionAuth API is on a public network, anyone may attempt to use the Login API.",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Application.",
			},
			"oauth_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     newOAuthConfiguration(),
				Optional: true,
			},
			"registration_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     newRegistrationConfiguration(),
				Optional: true,
			},
			"passwordless_configuration_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if passwordless login is enabled for this application.",
			},
			"registration_delete_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unverified_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates that users without a verified registration for this application will have their registration permanently deleted after application.registrationDeletePolicy.unverified.numberOfDaysToRetain days.",
						},
						"unverified_number_of_days_to_retain": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of days from registration a user’s registration will be retained before being deleted for not completing registration verification. This field is required when application.registrationDeletePolicy.enabled is set to true. Value must be greater than 0.",
						},
					},
				},
			},
			"samlv2_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     newSamlv2Configuration(),
				Optional: true,
			},
			"verification_email_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of the Email Template that is used to send the Registration Verification emails to users. If the verifyRegistration field is true this field is required.",
			},
			"verify_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not registrations to this Application may be verified. When this is set to true the verificationEmailTemplateId parameter is also required.",
			},
			"webhook_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of Webhook Ids. For Webhooks that are not already configured for All Applications, specifying an Id on this request will indicate the associated Webhook should handle events for this application.",
			},
			"email_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_verification_template_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "New configuration per application to specify a more specific template.",
							ValidateFunc: validation.IsUUID,
						},
						"forgot_password_template_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "New configuration per application to specify a more specific template.",
							ValidateFunc: validation.IsUUID,
						},
						"passwordless_email_template_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "New configuration per application to specify a more specific template.",
							ValidateFunc: validation.IsUUID,
						},
						"set_password_email_template_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "New configuration per application to specify a more specific template.",
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
		},
	}
}

func newSamlv2Configuration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"audience": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The audience for the SAML response sent to back to the service provider from FusionAuth. Some service providers require different audience values than the issuer and this configuration option lets you change the audience in the response.",
			},
			"callback_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the callback (sometimes called the Assertion Consumer Service or ACS). This is where FusionAuth sends the browser after the user logs in via SAML.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not FusionAuth will log SAML debug messages to the event log. This is useful for debugging purposes.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the SAML IdP for this Application is enabled or not.",
			},
			"issuer": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The issuer that identifies the service provider and allows FusionAuth to load the correct Application and SAML configuration. If you don’t know the issuer, you can often times put in anything here and FusionAuth will display an error message with the issuer from the service provider when you test the SAML login.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the Key used to sign the SAML response. If you do not specify this property, FusionAuth will create a new key and associate it with this Application.",
			},
			"logout_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL that the browser is taken to after the user logs out of the SAML service provider. Often service providers need this URL in order to correctly hook up single-logout. Note that FusionAuth does not support the SAML single-logout profile because most service providers to not support it properly.",
			},
			"xml_signature_canonicalization_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "exclusive_with_comments",
				ValidateFunc: validation.StringInSlice([]string{
					"exclusive",
					"exclusive_with_comments",
					"inclusive",
					"inclusive_with_comments",
				}, false),
				Description: "The XML signature canonicalization method used when digesting and signing the SAML response. Unfortunately, many service providers do not correctly implement the XML signature specifications and force a specific canonicalization method. This setting allows you to change the canonicalization method to match the service provider. Often, service providers don’t even document their required method. You might need to contact enterprise support at the service provider to figure out what method they use.",
			},
		},
	}
}

func newOAuthConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"authorized_origin_urls": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of URLs that are the authorized origins for FusionAuth OAuth.",
			},
			"authorized_redirect_urls": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of URLs that are the authorized redirect URLs for FusionAuth OAuth.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OAuth 2.0 client secret. If you leave this blank during a POST, a secure secret will be generated for you. If you leave this blank during PUT, the previous value will be maintained. For both POST and PUT you can provide a value and it will be stored.",
			},
			"device_verification_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The device verification URL to be used with the Device Code grant type, this field is required when device_code is enabled.",
			},
			"enabled_grants": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The enabled grants for this application. In order to utilize a particular grant with the OAuth 2.0 endpoints you must have enabled the grant.",
			},
			"generate_refresh_tokens": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if the OAuth 2.0 Token endpoint will generate a refresh token when the offline_access scope is requested.",
			},
			"logout_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RedirectOnly",
					"AllApplications",
				}, false),
				Default:     "AllApplications",
				Description: "Behavior when /oauth2/logout is called.",
			},
			"logout_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The logout URL for the Application. FusionAuth will redirect to this URL after the user logs out of OAuth.",
			},
			"require_client_authentication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Determines if the OAuth 2.0 Token endpoint requires client authentication. If this is enabled, the client must provide client credentials when using the Token endpoint. The client_id and client_secret may be provided using a Basic Authorization HTTP header, or by sending these parameters in the request body using POST data.",
			},
		},
	}
}

func newJWTConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_token_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of the signing key used to sign the access token.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if this application is using the JWT configuration defined here or the global JWT configuration defined by the System Configuration. If this is false the signing algorithm configured in the System Configuration will be used. If true the signing algorithm defined in this application will be used.",
			},
			"id_token_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of the signing key used to sign the Id token.",
			},
			"refresh_token_ttl_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The length of time in minutes the JWT refresh token will live before it is expired and is not able to be exchanged for a JWT.",
			},
			"ttl_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The length of time in seconds the JWT will live before it is expired and no longer valid.",
			},
		},
	}
}

func newRegistrationConfiguration() *schema.Resource {
	requireable := func() *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"required": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		}
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"birth_date": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     requireable(),
				Optional: true,
			},
			"confirm_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if the password should be confirmed during self service registration, this means that the user will be required to type the password twice.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if self service registration is enabled for this application. When this value is false, you may still use the Registration API, this only affects if the self service option is available during the OAuth 2.0 login.",
			},
			"first_name": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     requireable(),
				Optional: true,
			},
			"full_name": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     requireable(),
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     requireable(),
				Optional: true,
			},
			"login_id_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"email",
					"username",
				}, false),
				Default:     "email",
				Description: "The unique login Id that will be collected during registration, this value can be email or username. Leaving the default value of email is preferred because an email address is globally unique.",
			},
			"middle_name": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     requireable(),
				Optional: true,
			},
			"mobile_phone": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     requireable(),
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"basic",
					"advanced",
				}, false),
				Default:     "basic",
				Description: "The type of registration flow.",
			},
			"form_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of an associated Form when using advanced registration configuration type. This field is required when application.registrationConfiguration.type is set to advanced.",
			},
		},
	}
}
