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
  description = "OpsMan AMI - defaulted to OpsMan 2.1.0 in us-west-2 region"
  default = "ami-5cd04f24"
}

variable "opsman_username" {
  description = "OpsMan instance admin username"
  default = "admin"
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
