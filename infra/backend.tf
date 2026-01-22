terraform {
  backend "gcs" {
    bucket = "cns-tf-state"
    prefix = "cns/prod" # or cns/dev, cns/stage, etc.
  }
}