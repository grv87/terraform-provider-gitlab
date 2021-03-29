package gitlab

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServiceJenkinsCI() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabServiceJenkinsCICreate,
		Read:   resourceGitlabServiceJenkinsCIRead,
		Update: resourceGitlabServiceJenkinsCIUpdate,
		Delete: resourceGitlabServiceJenkinsCIDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabServiceJenkinsCIImportState,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateURLFunc,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"push_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"merge_requests_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tag_push_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceGitlabServiceJenkinsCICreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)

	jenkinsCIOptions, err := expandJenkinsCIOptions(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Create Gitlab Jenkins CI service")

	if _, err := client.Services.SetJenkinsCIService(project, jenkinsCIOptions); err != nil {
		return fmt.Errorf("couldn't create Gitlab Jenkins CI service: %w", err)
	}

	d.SetId(project)

	return resourceGitlabServiceJenkinsCIRead(d, meta)
}

func resourceGitlabServiceJenkinsCIRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	p, resp, err := client.Projects.GetProject(project, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Removing Gitlab Jenkins CI service %s because project %s not found", d.Id(), p.Name)
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[DEBUG] Read Gitlab Jenkins CI service %s", d.Id())

	jenkinsCIService, _, err := client.Services.GetJenkinsCIService(project)
	if err != nil {
		return err
	}

	if v := jenkinsCIService.Properties.URL; v != "" {
		d.Set("url", v)
	}
	if v := jenkinsCIService.Properties.ProjectName; v != "" {
		d.Set("project_name", v)
	}
	if v := jenkinsCIService.Properties.Username; v != "" {
		d.Set("username", v)
	}

	d.Set("title", jenkinsCIService.Title)
	d.Set("created_at", jenkinsCIService.CreatedAt.String())
	d.Set("updated_at", jenkinsCIService.UpdatedAt.String())
	d.Set("active", jenkinsCIService.Active)
	d.Set("push_events", jenkinsCIService.PushEvents)
	d.Set("issues_events", jenkinsCIService.IssuesEvents)
	d.Set("commit_events", jenkinsCIService.CommitEvents)
	d.Set("merge_requests_events", jenkinsCIService.MergeRequestsEvents)
	d.Set("tag_push_events", jenkinsCIService.TagPushEvents)
	d.Set("note_events", jenkinsCIService.NoteEvents)
	d.Set("pipeline_events", jenkinsCIService.PipelineEvents)
	d.Set("job_events", jenkinsCIService.JobEvents)

	return nil
}

func resourceGitlabServiceJenkinsCIUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceGitlabServiceJenkinsCICreate(d, meta)
}

func resourceGitlabServiceJenkinsCIDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)

	log.Printf("[DEBUG] Delete Gitlab Jenkins CI service %s", d.Id())

	_, err := client.Services.DeleteJenkinsCIService(project)

	return err
}

func expandJenkinsCIOptions(d *schema.ResourceData) (*gitlab.SetJenkinsCIServiceOptions, error) {
	setJenkinsCIServiceOptions := gitlab.SetJenkinsCIServiceOptions{}

	// Set required properties
	setJenkinsCIServiceOptions.URL = gitlab.String(d.Get("url").(string))
	setJenkinsCIServiceOptions.ProjectName = gitlab.String(d.Get("project_name").(string))
	setJenkinsCIServiceOptions.Username = gitlab.String(d.Get("username").(string))
	setJenkinsCIServiceOptions.Password = gitlab.String(d.Get("password").(string))
	setJenkinsCIServiceOptions.PushEvents = gitlab.Bool(d.Get("push_events").(bool))
	setJenkinsCIServiceOptions.MergeRequestsEvents = gitlab.Bool(d.Get("merge_requests_events").(bool))
	setJenkinsCIServiceOptions.TagPushEvents = gitlab.Bool(d.Get("tag_push_events").(bool))

	return &setJenkinsCIServiceOptions, nil
}

func resourceGitlabServiceJenkinsCIImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("project", d.Id())

	return []*schema.ResourceData{d}, nil
}
