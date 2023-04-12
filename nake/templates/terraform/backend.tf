terraform {
  backend "s3" {
    encrypt        = true
    dynamodb_table = "terraform-state-lock"
    region         = "eu-west-2"
  }
}
