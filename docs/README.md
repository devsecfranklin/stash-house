# Stash

## Usage

To import a multiline secret into an env var from `pass` database:

```sh
export GOOGLE_APPLICATION_CREDENTIALS=(pass show gcp-gcs-pso-f9322c73b9bf.json| string split0)
```
