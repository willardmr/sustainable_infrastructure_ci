# sustainable_infrastructure_ci
Continuous integration plugin to enforce sustainability best practices on software infrastructure.

Chat here: https://softwareforsustainability.zulipchat.com/join/f7vanlepyjfivdy35gfhxl63/

View the [project brief](ProjectBrief.pdf)

## Running tests

### Create a free GCP account
https://cloud.google.com/free

### Create a project in your account
https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project

### Install gcloud
https://cloud.google.com/sdk/docs/install
Then run `gcloud auth login`

### Run the tests
The following commands will run the GCP tests:
* `export GOOGLE_PROJECT=your_gcp_project_name`
* `cd /core/gcp`
* `go test -v -tags gcp .`
