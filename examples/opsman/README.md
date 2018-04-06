# Simple OpsMan Sample Terraform Templates

These Terraform templates create a plain vanilla single AZ OpsMan and director installation, by default in the us-west-2 region.

## Using

Create a terraform.tfvars file in this directory and populate the following reqiured variables:

```
aws_access_key_id = ""
aws_secret_access_key = ""
key_name = ""
ssh_private_key = <<KEY
-----BEGIN RSA PRIVATE KEY example-----
-----END RSA PRIVATE KEY-----
KEY
opsman_allow_cidr_ranges = ["x.x.x.x/32"]
route53_zone_id = ""
opsman_fqdn = ""
rds_password = ""
opsman_password = ""
```

Once those are populated you can directly run terraform or use the `./build_and_run.sh` at the root of the repo.
