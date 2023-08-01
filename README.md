# sustainable_infrastructure_ci
Continuous integration plugin to enforce sustainability best practices on software infrastructure.

## Running tests

### Create a free GCP account

### Create a project in your account

### Install gcloud https://cloud.google.com/sdk/docs/install
Then run `gcloud auth login`

### Run the tests
The following commands will run the GCP tests:
* `export GOOGLE_PROJECT=your_gcp_project_name`
* `cd /core/gcp`
* `go test -v -tags gcp .`
