terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.11, < 4.0"
    }
  }
}

provider "aws" {
  region  = "ap-northeast-2"
  profile = "default"
}
