terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.61"
    }
  }
}

provider "aws" {
  region  = "ap-northeast-2"
  profile = "default"
}
