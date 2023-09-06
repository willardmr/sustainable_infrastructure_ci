# sustainable_infrastructure_ci
Continuous integration plugin to enforce sustainability best practices on software infrastructure.

[![project chat](https://img.shields.io/badge/zulip-join_chat-brightgreen.svg)](https://softwareforsustainability.zulipchat.com/join/f7vanlepyjfivdy35gfhxl63/)

View the [project brief](ProjectBrief.pdf)

## Getting started

### Dependencies

1. Terraform: https://developer.hashicorp.com/terraform/downloads
2. Google Cloud: https://cloud.google.com/sdk/docs/install or https://snapcraft.io/google-cloud-cli

Make sure the directories containing terraform and gcloud can be found in your %PATH%

### Installation

1. Create a free GCP account: https://cloud.google.com/free
2. Create a project in your account: https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project
3. Clone & cd:
```sh
git clone https://github.com/AntoineSebert/sustainable_infrastructure_ci.git
cd sustainable_infrastructure_ci
```

### Running the tests

1. Login to Google Cloud:
```sh
gcloud auth login
# if you encounter an authentication error at a later stages you can try instead:
# gcloud auth application-default login
```
2. Set your project name in the environment: \
For linux:
```sh
export GOOGLE_PROJECT=your_gcp_project_name
```
&emsp; &emsp; For Windows: 
```sh
set GOOGLE_PROJECT=your_gcp_project_name
```
3. Run the test:
```sh
cd core/gcp
go test -v -tags gcp .
```
