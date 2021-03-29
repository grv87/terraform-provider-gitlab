package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabServiceJenkinsCI_basic(t *testing.T) {
	var jenkinsCIService gitlab.JenkinsCIService
	rInt := acctest.RandInt()
	jenkinsCIResourceName := "gitlab_service_jenkins_ci.jenkins_ci"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabServiceJenkinsCIDestroy,
		Steps: []resource.TestStep{
			// Create a project and a jenkins ci service
			{
				Config: testAccGitlabServiceJenkinsCIConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceJenkinsCIExists(jenkinsCIResourceName, &jenkinsCIService),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "project_name", "foo"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "username", "user1"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "password", "mypass"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "push_events", "true"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "merge_requests_events", "false"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "tag_push_events", "true"),
				),
			},
			// Update the jenkins ci service
			{
				Config: testAccGitlabServiceJenkinsCIUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceJenkinsCIExists(jenkinsCIResourceName, &jenkinsCIService),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "url", "https://testurl.com"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "project_name", "bar"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "username", "user2"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "password", "mypass_update"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "push_events", "false"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "merge_requests_events", "true"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "tag_push_events", "false"),
				),
			},
			// Update the jenkins ci service to get back to previous settings
			{
				Config: testAccGitlabServiceJenkinsCIConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceJenkinsCIExists(jenkinsCIResourceName, &jenkinsCIService),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "project_name", "foo"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "username", "user1"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "password", "mypass"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "push_events", "true"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "merge_requests_events", "false"),
					resource.TestCheckResourceAttr(jenkinsCIResourceName, "tag_push_events", "true"),
				),
			},
		},
	})
}

func TestAccGitlabServiceJenkinsCI_import(t *testing.T) {
	jenkinsCIResourceName := "gitlab_service_jenkins_ci.jenkins_ci"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabServiceJenkinsCIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabServiceJenkinsCIConfig(rInt),
			},
			{
				ResourceName:      jenkinsCIResourceName,
				ImportStateIdFunc: getJenkinsCIProjectID(jenkinsCIResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testAccCheckGitlabServiceJenkinsCIExists(n string, service *gitlab.JenkinsCIService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		jenkinsCIService, _, err := conn.Services.GetJenkinsCIService(project)
		if err != nil {
			return fmt.Errorf("Jenkins CI service does not exist in project %s: %v", project, err)
		}
		*service = *jenkinsCIService

		return nil
	}
}

func testAccCheckGitlabServiceJenkinsCIDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists")
				}
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func getJenkinsCIProjectID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return "", fmt.Errorf("No project ID is set")
		}

		return project, nil
	}
}

func testAccGitlabServiceJenkinsCIConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"
  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_jenkins_ci" "jenkins_ci" {
  project      = "${gitlab_project.foo.id}"
  url          = "https://test.com"
  project_name = "foo"
  username     = "user1"
  password     = "mypass"
  push_events           = true
  merge_requests_events = false
  tag_push_events       = true
}
`, rInt)
}

func testAccGitlabServiceJenkinsCIUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"
  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_jenkins_ci" "jenkins_ci" {
  project      = "${gitlab_project.foo.id}"
  url          = "https://testurl.com"
  project_name = "bar"
  username     = "user2"
  password     = "mypass_update"
  push_events           = false
  merge_requests_events = true
  tag_push_events       = false
}
`, rInt)
}
