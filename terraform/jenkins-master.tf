# resource "aws_vpc" "jenkins" {
#   cidr_block = "10.0.0.0/16"
#   tags = {
#     Name : "jenkins-vpc"
#   }
# }

# resource "aws_subnet" "jenkins" {
#   vpc_id = aws_vpc.jenkins.id
#   tags = {
#     Name : "blog_subnet"
#   }
#   map_public_ip_on_launch = true
#   cidr_block              = "10.0.1.0/24"
# }

# resource "aws_internet_gateway" "jenkins" {
#   vpc_id = aws_vpc.jenkins.id
# }

# resource "aws_route_table" "jenkins" {
#   vpc_id = aws_vpc.jenkins.id
#   route {
#     cidr_block = "0.0.0.0/0"
#     gateway_id = aws_internet_gateway.jenkins.id
#   }
# }

# resource "aws_route_table_association" "jenkins" {
#   subnet_id      = aws_subnet.jenkins.id
#   route_table_id = aws_route_table.jenkins.id
# }

# resource "aws_security_group" "jenkins" {
#   vpc_id = aws_vpc.jenkins.id
#   name   = "jenkins_sg"

#   ingress {
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = ["0.0.0.0/0"]
#   }

#   egress {
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = ["0.0.0.0/0"]
#   }
# }

# resource "null_resource" "pemkey" {
#   connection {
#     host     = "${aws_instance.jenkins_master.public_ip}"
#     type        = "ssh"
#     user        = "ubuntu"
#     private_key = file("~/pemkey/blog.pem")
#   }

#   provisioner "file" {
#     source      = "~/pemkey/blog.pem"    
#     destination = "blog.pem"
#   }
# }

# resource "aws_instance" "jenkins_master" {
#   ami                    = "ami-0c9c942bd7bf113a2"
#   instance_type          = "t3.micro"
#   vpc_security_group_ids = [aws_security_group.jenkins.id]
#   subnet_id              = aws_subnet.jenkins.id
#   key_name               = "blog"
#   # iam_instance_profile   = aws_iam_instance_profile.node.name

#   tags = {
#     Name = "master_node"
#   }

#   connection {
#     type        = "ssh"
#     user        = "ubuntu"
#     private_key = file("../../pemkey/blog.pem")
#     host        = self.public_ip
#   }

#   root_block_device {
#     volume_size = 8
#     volume_type = "gp2"
#   }  
  
#   provisioner "remote-exec" {

#     inline = [
#       "sudo apt-get update",
#       "sudo apt-get install curl",
#       "curl -fsSL https://get.docker.com -o get-docker.sh",
#       "sudo sh get-docker.sh",
#       "docker pull achoistic98/blog_jenkins:latest",
#       "docker run -it -p 8080:8080 achoistic98/blog_jenkins:latest",
#     ]
#   }
# }

# output "jenkins_master-ip" {
#   value = aws_instance.jenkins_master.public_ip
# }