variable "aws_access_key_id" {}
variable "aws_secret_access_key" {}

variable "key_name" {
  description = "AWS key name we're using"
}

variable "ssh_private_key" {
  description = "OpsMan SSH key"
}

variable "opsman_allow_cidr_ranges" {
  type = "list"
  description = "List of CIDRs allowed to access OpsMan, for example your public IP 'x.x.x.x/32'"
}

variable "route53_zone_id" {
  description = "The route 53 zone we can modify"
}

variable "opsman_fqdn" {
  description = "host.fqdn of opsman, i.e. 'opsman.pcf.example.com'"
}

variable "rds_password" {
  description = "Password to set/use for the RDS bosh instance"
}

variable "opsman_password" {
  description = "Password to set/use for the OpsMan instance"
}

variable "opsman_ami" {
  type = "map"
  default = {
 "us-east-1" = "ami-21ee435c"
 "us-east-2" = "ami-ae4575cb"
 "us-west-1" = "ami-b24a5ad2"
 "us-west-2" = "ami-7dbade05"
 "us-gov-west-1" = "ami-025ecb63"
 "ap-south-1" = "ami-d31530bc"
 "eu-west-3" = "ami-cc06b0b1"
 "eu-west-2" = "ami-ed05e48a"
 "eu-west-1" = "ami-50a7f829"
 "ap-northeast-2" = "ami-2d57f843"
 "ap-northeast-1" = "ami-d7c5d2ab"
 "sa-east-1" = "ami-d2590fbe"
 "ca-central-1" = "ami-9e61e7fa"
 "ap-southeast-1" = "ami-0ca5ff70"
 "ap-southeast-2" = "ami-99a36dfb"
 "eu-central-1" = "ami-82bfe069"
}
}

variable "region" {
  description = "AWS region to use"
}

variable "az" {
  description = "AZ to use"
} 

variable "opsman_username" {
  description = "OpsMan instance admin username"
  default = "admin"
}

variable "rds_username" {
  description = "RDS instance username"
  default = "opsmanrdsadmin"
}
