# RDS instance
resource "aws_db_instance" "bosh" {
    identifier              = "opsman-example-rds-bosh"
    allocated_storage       = 100
    engine                  = "mariadb"
    engine_version          = "10.1.19"
    iops                    = 1000
    instance_class          = "db.m4.large"
    name                    = "bosh"
    username                = "${var.rds_username}"
    password                = "${var.rds_password}"
    db_subnet_group_name    = "${aws_db_subnet_group.bosh_rds.name}"
    parameter_group_name    = "default.mariadb10.1"
    vpc_security_group_ids  = ["${aws_security_group.bosh_rds.id}"]
    multi_az                = false
    backup_retention_period = 1
    apply_immediately       = true
    skip_final_snapshot     = true
}

resource "aws_db_subnet_group" "bosh_rds" {
    name = "opsman_example_rds_subnet_group"
    subnet_ids = [
      "${aws_subnet.bosh_rds_az1.id}",
      "${aws_subnet.bosh_rds_az2.id}"
    ]
}

resource "aws_security_group" "bosh_rds" {
    name = "bosh_rds"
    description = "Allow incoming connections for RDS."
    vpc_id = "${aws_vpc.pcf.id}"
    ingress {
        from_port = 3306
        to_port = 3306
        protocol = "tcp"
        cidr_blocks = ["${aws_subnet.management.cidr_block}"]
    }
    egress {
        from_port = 0
        to_port = 0
        protocol = -1
        cidr_blocks = ["0.0.0.0/0"]
    }
}

resource "aws_subnet" "bosh_rds_az1" {
    vpc_id = "${aws_vpc.pcf.id}"
    cidr_block = "10.0.12.0/24"
    availability_zone = "${var.region}a"
}

resource "aws_subnet" "bosh_rds_az2" {
    vpc_id = "${aws_vpc.pcf.id}"
    cidr_block = "10.0.13.0/24"
    availability_zone = "${var.region}b"
}
