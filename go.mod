module github.com/gitlabhq/terraform-provider-gitlab

go 1.16

require (
	github.com/hashicorp/go-retryablehttp v0.6.8
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/xanzy/go-gitlab v0.46.0
)

replace github.com/xanzy/go-gitlab => github.com/grv87/go-gitlab v0.48.1-0.20210329170430-df54128faeee
