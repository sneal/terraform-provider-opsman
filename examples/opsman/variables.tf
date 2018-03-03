variable "aws_access_key_id" {}
variable "aws_secret_access_key" {}
variable "key_name" {}
variable "ssh_private_key" {}
variable "opsman_allow_cidr_ranges" {
  type = "list"
  description = "List of CIDRs allowed to access OpsMan"
}
variable "route53_zone_id" {}
variable "opsman_fqdn" {}

variable "opsman_ami" {
  description = "OpsMan AMI"
  default = "ami-5cd04f24"
}

variable "opsman_username" {
  description = "OpsMan instance admin username"
  default = "admin"
}

variable "opsman_password" {
  description = "OpsMan instance admin password"
  default = "opsmanpassw0rd"
}

variable "region" {
  description = "AWS region to use"
  default = "us-west-2"
}

variable "az" {
  description = "AZ to use"
  default = "us-west-2a"
}

variable "rds_username" {
  description = "RDS instance username"
  default = "opsmanrdsadmin"
}

variable "rds_password" {
  description = "RDS instance password"
  default = "opsmanrdspassw0rd"
}