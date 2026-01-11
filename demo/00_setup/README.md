# 00 Setup

* In this step we run the setup tool to prepare the local environment.

##  `config.yml`

This YAML is meant to ease the setup concerns on the end user.

## Backends

These are some options for storing your credentials once they have been encrypted.

### Github

You can use GH or other git related revision control system.

### Database

### GCloud Storage Bucket

```sh
sudo apt install npm
npm install -g @google/gemini-cli


# Creates a simulated folder by uploading a 0-byte object
gcloud storage cp /dev/null gs://lab-franklin/stash-house/
# Create a new bucket with HNS
gcloud storage buckets create gs://new-hierarchical-bucket --enable-hierarchical-namespace

# Define the repository for Bullseye
export GCSFUSE_REPO=gcsfuse-bullseye

# Add the repository to your apt sources
echo "deb https://packages.cloud.google.com/apt $GCSFUSE_REPO main" | sudo tee /etc/apt/sources.list.d/gcsfuse.list

# Import the Google Cloud public GPG key
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -

sudo apt-get update
sudo apt-get install fuse gcsfuse

Once installed, you must provide the instance with credentials to access your bucket: 
Authorize: gcloud auth application-default login
Create Mount Point: mkdir ~/my_bucket_mount
Mount: gcsfuse YOUR_BUCKET_NAME ~/my_bucket_mount
```
