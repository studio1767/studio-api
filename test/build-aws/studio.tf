## ssh key for servers

resource "aws_key_pair" "ssh_key" {
  key_name   = var.studio_code
  public_key = tls_private_key.ssh_key.public_key_openssh
}

## site vpc

resource "aws_vpc" "site" {
  cidr_block = var.cidr_block
  enable_dns_support = true
  enable_dns_hostnames = false
  tags = {
    Name = var.studio_code
  }
}

## tag default route table and security group
##   ... and neuter the default security group

resource "aws_default_route_table" "site" {
  default_route_table_id = aws_vpc.site.default_route_table_id
  tags = {
    Name = "${var.studio_code}-default"
  }
}

resource "aws_default_security_group" "site" {
  vpc_id = aws_vpc.site.id
  tags = {
    Name = "${var.studio_code}-default"
  }
}

## internet gateway

resource "aws_internet_gateway" "site" {
  vpc_id = aws_vpc.site.id
  tags = {
    Name = var.studio_code
  }
}

## subnet

resource "aws_subnet" "subnet" {
  vpc_id     = aws_vpc.site.id
  cidr_block = var.cidr_block
  availability_zone = local.aws_availability_zone

  tags = {
    Name = var.studio_code
  }
}

resource "aws_route_table" "subnet" {
  vpc_id = aws_vpc.site.id
  tags = {
    Name = var.studio_code
  }
}

resource "aws_route_table_association" "subnet" {
  route_table_id = aws_route_table.subnet.id
  subnet_id      = aws_subnet.subnet.id
}

resource "aws_route" "gateway_default" {
  route_table_id         = aws_route_table.subnet.id
  gateway_id             = aws_internet_gateway.site.id
  destination_cidr_block = "0.0.0.0/0"
}

resource "aws_security_group" "subnet" {
  vpc_id = aws_vpc.site.id
  tags = {
    Name = var.studio_code
  }
}

resource "aws_security_group_rule" "all_out" {
  security_group_id = aws_security_group.subnet.id
  type        = "egress"
  protocol    = -1
  from_port   = 0
  to_port     = 0
  cidr_blocks = [ "0.0.0.0/0" ]
}

resource "aws_security_group_rule" "ssh" {
  security_group_id = aws_security_group.subnet.id
  type        = "ingress"
  protocol    = "tcp"
  from_port   = 22
  to_port     = 22
  cidr_blocks = [ local.management_net ]
}

resource "aws_security_group_rule" "ldaps" {
  security_group_id = aws_security_group.subnet.id
  type        = "ingress"
  protocol    = "tcp"
  from_port   = 636
  to_port     = 636
  cidr_blocks = [ local.management_net ]
}

resource "aws_security_group_rule" "db" {
  security_group_id = aws_security_group.subnet.id
  type        = "ingress"
  protocol    = "tcp"
  from_port   = 3306
  to_port     = 3306
  cidr_blocks = [ local.management_net ]
}

