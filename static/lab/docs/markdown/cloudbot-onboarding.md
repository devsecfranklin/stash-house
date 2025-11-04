# On-boarding

* Invite bot user "n0ctilucent" as a collaborator on the repo.

## Add the GCP Cloud function

* Add [the Github action](../actions/cloudbot-call.yml) to folder `.github/workflows`
* Under "Settings" -> "Secrets" add a new secret `GOOGLE_APPLICATION_CREDENTIALS` with GCP JSON.
* Under "Settings" -> "Secrets" add a new secret `PROJECT_ID` with value `gcp-gcs-pso`

## Add the Tekton CI webhook

* Click on "Settings" -> "Webhooks". Click the "Add webhook" button.
* Payload URL `http://34.132.46.121/hooks`
* Content type `application/json`
* Secret (get k8s secret, then base64 decode it):

```sh
k get secret github-interceptor-secret -n tekton-pipelines -o jsonpath='{.data}'
echo 'cG9vcCBzYW5kd2ljaAo=' | base64 --decode
```

* "Let me select individual events"
  * tick the "Pull requests" box
  * untick the "Pushes" box
