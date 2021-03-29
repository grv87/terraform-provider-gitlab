# gitlab\_service\_jenkins\_ci

This resource allows you to manage Jenkins CI integration.

## Example Usage

```hcl
resource "gitlab_project" "awesome_project" {
  name = "awesome_project"
  description = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_jenkins_ci" "jenkins_ci" {
  project      = gitlab_project.awesome_project.id
  url          = "https://jenkinsci.example.com"
  project_name = "awesome_project"
  username     = "user"
  password     = "mypass"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) ID of the project you want to activate integration on.

* `url` - (Required) The URL to the Jenkins instance. For example, https://jenkinsci.example.com.

* `project_name` - (Required) The URL-friendly project name, e.g., `awesome_project`

* `username` - (Required) The username of the user created to be used with GitLab/Jenkins CI.

* `password` - (Required) The password of the user created to be used with GitLab/Jenkins CI.

* `push_events` - (Optional) Trigger event when a new commit is pushed

* `merge_requests_events` - (Optional) Trigger event when a merge request is created/updated/merged

* `tag_push_events` - (Optional) Trigger event when a new tag is pushed

## Importing Jenkins CI service

 You can import a service_jenkins_ci state using `terraform import <resource> <project_id>`:

```bash
$ terraform import gitlab_service_jenkins_ci.jenkins_ci 1
```
