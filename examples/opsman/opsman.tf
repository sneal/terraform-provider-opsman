resource "opsman_installation_settings" "opsman" {
  address  = "${var.opsman_fqdn}"
  username = "${var.opsman_username}"
  password = "${var.opsman_password}"
}

resource "aws_instance" "opsman" {
  ami           = "${var.opsman_ami}"
  instance_type = "m4.large"
  key_name      = "${var.key_name}"

  availability_zone      = "${var.az}"
  vpc_security_group_ids = ["${aws_security_group.opsman.id}"]
  subnet_id              = "${aws_subnet.management.id}"
  private_ip             = "10.0.16.7"
  
  root_block_device {
    volume_size = 100
  }
  
  depends_on = ["aws_internet_gateway.pcf", "opsman_installation_settings.opsman"]
}

resource "opsman_director" "opsman" {
  address               = "${aws_route53_record.opsman.fqdn}"
  username              = "${var.opsman_username}"
  password              = "${var.opsman_password}"
  decryption_passphrase = "${var.opsman_password}"

  instance_id                = "${aws_instance.opsman.id}"
  installation_settings_file = "${opsman_installation_settings.opsman.installation_settings_file}"

  access_key_id     = "${var.aws_access_key_id}"
  secret_access_key = "${var.aws_secret_access_key}"
  key_name          = "${var.key_name}"
  ssh_private_key   = "${var.ssh_private_key}"
  
  vpc_id                = "${aws_vpc.pcf.id}"
  vpc_security_group_id = "${aws_security_group.opsman.id}"
  region                = "${var.region}"
  availability_zones    = ["${var.az}"]
  
  availability_zone = "${var.az}"
  director_network  = "management"
  
  database {
    host     = "${aws_db_instance.bosh.address}"
    port     = "${aws_db_instance.bosh.port}"
    username = "${var.rds_username}"
    password = "${var.rds_password}"
  }

  blobstore {
    s3_endpoint = "https://s3.us-west-2.amazonaws.com"
    bucket_name = "${aws_s3_bucket.bosh.id}"
  }

  network {
    name = "management"
    subnet {
      vpc_subnet_id =    "${aws_subnet.management.id}"
      cidr               = "10.0.16.0/28"
      reserved_ip_ranges = "10.0.16.0-10.0.16.5"
      dns                = "169.254.169.253"
      gateway            = "10.0.16.1"
      availability_zone  = "${var.az}"
    }
  }
}

resource "aws_route53_record" "opsman" {
  zone_id = "${var.route53_zone_id}"
  name = "opsman"
  type = "A"
  ttl = "60"
  records = ["${aws_eip.opsman.public_ip}"]
}

resource "aws_eip" "opsman" {
  instance = "${aws_instance.opsman.id}"
  vpc      = true
}

# VPC for OpsMan/PAS install
resource "aws_vpc" "pcf" {
    cidr_block = "10.0.0.0/16"
    enable_dns_hostnames = true
}

# Internet gateway for outbound connections
resource "aws_internet_gateway" "pcf" {
    vpc_id = "${aws_vpc.pcf.id}"
}

resource "aws_route_table" "pcf" {
  vpc_id = "${aws_vpc.pcf.id}"
  route {
      cidr_block = "0.0.0.0/0"
      gateway_id = "${aws_internet_gateway.pcf.id}"
  }
}

resource "aws_route_table_association" "opsman_example_route_table_assoc" {
  subnet_id = "${aws_subnet.management.id}"
  route_table_id = "${aws_route_table.pcf.id}"
}

# Single "public" management subnet for OpsMan and bosh director
resource "aws_subnet" "management" {
    vpc_id = "${aws_vpc.pcf.id}"
    cidr_block = "10.0.16.0/28"
    availability_zone = "${var.az}"
}

# OpsMan security group
resource "aws_security_group" "opsman" {
    name = "opsman-example-opsman-sg"
    description = "Allow incoming connections for Ops Manager."
    vpc_id = "${aws_vpc.pcf.id}"
}

# Allow inbound traffic from other instances associated with this security group
resource "aws_security_group_rule" "opsman_example_allow_ingress_default_sg" {
    type = "ingress"
    from_port = 0
    to_port = 0
    protocol = "all"
    cidr_blocks = ["${aws_subnet.management.cidr_block}"]
    security_group_id = "${aws_security_group.opsman.id}"
}

# Allow inbound SSH from my network
resource "aws_security_group_rule" "opsman_example_allow_ssh_sg" {
    type              = "ingress"
    from_port         = 22
    to_port           = 22
    protocol          = "tcp"
    cidr_blocks       = "${var.opsman_allow_cidr_ranges}"
    security_group_id = "${aws_security_group.opsman.id}"
}

# Allow inbound HTTPS from my network
resource "aws_security_group_rule" "opsman_example_allow_https_sg" {
    type              = "ingress"
    from_port         = 443
    to_port           = 443
    protocol          = "tcp"
    cidr_blocks       = "${var.opsman_allow_cidr_ranges}"
    security_group_id = "${aws_security_group.opsman.id}"
}

# Allow inbound HTTP from my network
resource "aws_security_group_rule" "opsman_example_allow_http_sg" {
    type              = "ingress"
    from_port         = 80
    to_port           = 80
    protocol          = "tcp"
    cidr_blocks       = "${var.opsman_allow_cidr_ranges}"
    security_group_id = "${aws_security_group.opsman.id}"
}

# Allow all VPC instances to egress anywhere in the world
resource "aws_security_group_rule" "opsman_example_allow_egress_default_sg" {
    type = "egress"
    from_port = 0
    to_port = 0
    protocol = "all"
    cidr_blocks = ["0.0.0.0/0"]
    security_group_id = "${aws_security_group.opsman.id}"
}
