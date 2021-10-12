# VPC
resource "aws_vpc" "customer_vpc" {
  cidr_block = var.vpc_cidr
}

# Internet Gateway
resource "aws_internet_gateway" "customer_igw" {
  vpc_id = aws_vpc.customer_vpc.id
}

# Subnets : public
resource "aws_subnet" "public" {
  count                   = length(var.subnets_cidr)
  vpc_id                  = aws_vpc.customer_vpc.id
  cidr_block              = element(var.subnets_cidr, count.index)
  availability_zone       = element(var.azs, count.index)
  map_public_ip_on_launch = true
}

# Route table: attach Internet Gateway 
resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.customer_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.customer_igw.id
  }
}

# Route table association with public subnets
resource "aws_route_table_association" "a" {
  count          = length(var.subnets_cidr)
  subnet_id      = element(aws_subnet.public.*.id, count.index)
  route_table_id = aws_route_table.public_rt.id
}

module "ec2_sensor" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "~> 3.0"

  name = "sensor"

  ami                    = "ami-ebd02392"
  instance_type          = "t2.micro"
  vpc_security_group_ids = ["sg-12345678"]
  subnet_id              = aws_subnet.public.id

  provisioner "file" {
    source      = "../bin/sensor-script.sh"
    destination = "/tmp/script.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/script.sh",
      "/tmp/script.sh",
    ]
  }
}

module "ec2_rev_proxy" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "~> 3.0"

  name = "rev-proxy"

  ami                    = "ami-ebd02392"
  instance_type          = "t2.micro"
  vpc_security_group_ids = ["sg-12345678"]
  subnet_id              = aws_subnet.public.id

  provisioner "file" {
    source      = "../bin/proxy-script.sh"
    destination = "/tmp/script.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/script.sh",
      "/tmp/script.sh",
    ]
  }

resource "aws_eip" "ip" {
  instance = aws_instance.ec2_rev_proxy.id
}

