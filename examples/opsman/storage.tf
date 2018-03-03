# S3 buckets
resource "random_id" "suffix" {
  byte_length = 4
}

resource "aws_s3_bucket" "bosh" {
    bucket = "opsman-example-bosh-${random_id.suffix.hex}"
    acl = "private"
    force_destroy = true
    tags {
        name = "opsman-example-bosh"
        environment = "opsman-example"
    }
}

resource "aws_s3_bucket" "buildpacks" {
    bucket = "opsman-example-buildpacks-${random_id.suffix.hex}"
    acl = "private"
    force_destroy = true
    tags {
        name = "opsman-example-buildpacks"
        environment = "opsman-example"
    }
}

resource "aws_s3_bucket" "droplets" {
    bucket = "opsman-example-droplets-${random_id.suffix.hex}"
    acl = "private"
    force_destroy = true
    tags {
        name = "opsman-example-droplets"
        environment = "opsman-example"
    }
}

resource "aws_s3_bucket" "packages" {
    bucket = "opsman-example-packages-${random_id.suffix.hex}"
    acl = "private"
    force_destroy = true
    tags {
        name = "opsman-example-packages"
        environment = "opsman-example"
    }
}

resource "aws_s3_bucket" "resources" {
    bucket = "opsman-example-resources-${random_id.suffix.hex}"
    acl = "private"
    force_destroy = true
    tags {
        name = "opsman-example-resources"
        environment = "opsman-example"
    }
}