provider "aws" {
  region = "ap-northeast-2"
}

resource "aws_vpc" "blog_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name : "blog-vpc"
  }
}

resource "aws_subnet" "blog_subnet" {
  vpc_id = aws_vpc.blog_vpc.id
  tags = {
    Name : "blog_subnet"
  }
  map_public_ip_on_launch = true
  cidr_block              = "10.0.1.0/24"
}

resource "aws_internet_gateway" "blog_ingernet_gateway" {
  vpc_id = aws_vpc.blog_vpc.id
}

resource "aws_route_table" "blog_route_table" {
  vpc_id = aws_vpc.blog_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.blog_ingernet_gateway.id
  }
}

resource "aws_route_table_association" "blog_association" {
  subnet_id      = aws_subnet.blog_subnet.id
  route_table_id = aws_route_table.blog_route_table.id
}

resource "aws_security_group" "blog_sg" {
  vpc_id = aws_vpc.blog_vpc.id
  name   = "blog_sg"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "blog_nlb" {
  name               = "blog-nlb"
  internal           = false
  load_balancer_type = "network"

  subnet_mapping {
    subnet_id = aws_subnet.blog_subnet.id
  }

  enable_deletion_protection = false

  tags = {
    Environment = "production"
  }
}

resource "aws_lb_target_group" "blog_target_group" {
  name        = "blog-target-group"
  port        = 80
  protocol    = "TCP"
  target_type = "ip"
  vpc_id      = aws_vpc.blog_vpc.id
}

resource "aws_lb_target_group_attachment" "blog_target_group_ip" {
  target_group_arn = aws_lb_target_group.blog_target_group.arn
  target_id        = aws_instance.blog_worker1.private_ip
  port             = 32665
}


resource "aws_lb_listener" "blog_nlb_listner" {
  load_balancer_arn = aws_lb.blog_nlb.arn
  port              = "80"
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.blog_target_group.arn
  }
}

resource "aws_ebs_volume" "jenkins_volume" {
  availability_zone = "ap-northeast-2b"
  size              = 10

  tags = {
    Name = "jenkins-volume"
  }
}

resource "aws_instance" "blog_master" {
  ami                    = "ami-0c9c942bd7bf113a2"
  instance_type          = "t3.small"
  vpc_security_group_ids = [aws_security_group.blog_sg.id]
  subnet_id              = aws_subnet.blog_subnet.id
  key_name               = "blog"

  tags = {
    Name = "master_node"
  }

  connection {
    type        = "ssh"
    user        = "ubuntu"
    private_key = file("../../pemkey/blog.pem")
    host        = self.public_ip
  }

  root_block_device {
    volume_size = 8
    volume_type = "gp2"
  }

  provisioner "remote-exec" {

    inline = [
      "cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf",
      "overlay",
      "br_netfilter",
      "EOF",
      "sudo modprobe overlay",
      "sudo modprobe br_netfilter",
      "cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf",
      "net.bridge.bridge-nf-call-iptables  = 1",
      "net.bridge.bridge-nf-call-ip6tables = 1",
      "net.ipv4.ip_forward = 1",
      "EOF",
      "sudo sysctl --system",
      "sudo swapoff -a",
      "(crontab -l 2>/dev/null; echo '@reboot /sbin/swapoff -a') | crontab - || true",
      "wget https://github.com/containerd/containerd/releases/download/v1.7.3/containerd-1.7.3-linux-amd64.tar.gz",
      "sudo tar Czxvf /usr/local containerd-1.7.3-linux-amd64.tar.gz",
      "wget https://raw.githubusercontent.com/containerd/containerd/main/containerd.service",
      "sudo mv containerd.service /usr/lib/systemd/system/",
      "sudo systemctl daemon-reload",
      "sudo systemctl enable --now containerd",
      "wget https://github.com/opencontainers/runc/releases/download/v1.1.8/runc.amd64",
      "sudo install -m 755 runc.amd64 /usr/local/sbin/runc",
      "sudo mkdir -p /etc/containerd/",
      "containerd config default | sudo tee /etc/containerd/config.toml",
      "sudo sed -i 's/SystemdCgroup \\= false/SystemdCgroup \\= true/g' /etc/containerd/config.toml",
      "sudo systemctl restart containerd",
      "sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates curl",
      "curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -",
      "cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list",
      "deb https://apt.kubernetes.io/ kubernetes-xenial main",
      "EOF",
      "sudo apt-get update",
      "sudo apt-get install -y kubelet kubeadm kubectl",
      "sudo apt-mark hold kubelet kubeadm kubectl containerd",
      "IPADDR=$(curl ifconfig.me && echo \"\")",
      "NODENAME=$(hostname -s)",
      "POD_CIDR=\"192.168.0.0/16\"",
      "touch text.txt",
      "sudo kubeadm init --control-plane-endpoint=$IPADDR  --apiserver-cert-extra-sans=$IPADDR  --pod-network-cidr=$POD_CIDR --node-name $NODENAME --ignore-preflight-errors Swap > text.txt",
      "sudo mkdir -p $HOME/.kube",
      "sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config",
      "sudo chown $(id -u):$(id -g) $HOME/.kube/config",
      "kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.0/manifests/calico.yaml",
      "kubectl apply -f https://raw.githubusercontent.com/techiescamp/kubeadm-scripts/main/manifests/metrics-server.yaml",
    ]
  }
}

resource "aws_instance" "blog_worker2" {
  ami                    = "ami-0c9c942bd7bf113a2"
  instance_type          = "t3.micro"
  vpc_security_group_ids = [aws_security_group.blog_sg.id]
  subnet_id              = aws_subnet.blog_subnet.id
  key_name               = "blog"

  tags = {
    Name = "worker_node2"
  }

  connection {
    type        = "ssh"
    user        = "ubuntu"
    private_key = file("../../pemkey/blog.pem")
    host        = self.public_ip
  }

  root_block_device {
    volume_size = 16
    volume_type = "gp2"
  }

  provisioner "remote-exec" {

    inline = [
      "cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf",
      "overlay",
      "br_netfilter",
      "EOF",
      "sudo modprobe overlay",
      "sudo modprobe br_netfilter",
      "cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf",
      "net.bridge.bridge-nf-call-iptables  = 1",
      "net.bridge.bridge-nf-call-ip6tables = 1",
      "net.ipv4.ip_forward = 1",
      "EOF",
      "sudo sysctl --system",
      "sudo swapoff -a",
      "(crontab -l 2>/dev/null; echo '@reboot /sbin/swapoff -a') | crontab - || true",
      "wget https://github.com/containerd/containerd/releases/download/v1.7.3/containerd-1.7.3-linux-amd64.tar.gz",
      "sudo tar Czxvf /usr/local containerd-1.7.3-linux-amd64.tar.gz",
      "wget https://raw.githubusercontent.com/containerd/containerd/main/containerd.service",
      "sudo mv containerd.service /usr/lib/systemd/system/",
      "sudo systemctl daemon-reload",
      "sudo systemctl enable --now containerd",
      "wget https://github.com/opencontainers/runc/releases/download/v1.1.8/runc.amd64",
      "sudo install -m 755 runc.amd64 /usr/local/sbin/runc",
      "sudo mkdir -p /etc/containerd/",
      "containerd config default | sudo tee /etc/containerd/config.toml",
      "sudo sed -i 's/SystemdCgroup \\= false/SystemdCgroup \\= true/g' /etc/containerd/config.toml",
      "sudo systemctl restart containerd",
      "sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates curl",
      "curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -",
      "cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list",
      "deb https://apt.kubernetes.io/ kubernetes-xenial main",
      "EOF",
      "sudo apt-get update",
      "sudo apt-get install -y kubelet kubeadm kubectl",
      "sudo apt-mark hold kubelet kubeadm kubectl containerd",
    ]
  }
}

resource "aws_instance" "blog_worker1" {
  ami                    = "ami-0c9c942bd7bf113a2"
  instance_type          = "t3.small"
  vpc_security_group_ids = [aws_security_group.blog_sg.id]
  subnet_id              = aws_subnet.blog_subnet.id
  key_name               = "blog"

  tags = {
    Name = "worker-node1"
  }

  connection {
    type        = "ssh"
    user        = "ubuntu"
    private_key = file("../../pemkey/blog.pem")
    host        = self.public_ip
  }

  root_block_device {
    volume_size = 8
    volume_type = "gp2"
  }

  provisioner "remote-exec" {

    inline = [
      "cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf",
      "overlay",
      "br_netfilter",
      "EOF",
      "sudo modprobe overlay",
      "sudo modprobe br_netfilter",
      "cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf",
      "net.bridge.bridge-nf-call-iptables  = 1",
      "net.bridge.bridge-nf-call-ip6tables = 1",
      "net.ipv4.ip_forward = 1",
      "EOF",
      "sudo sysctl --system",
      "sudo swapoff -a",
      "(crontab -l 2>/dev/null; echo '@reboot /sbin/swapoff -a') | crontab - || true",
      "wget https://github.com/containerd/containerd/releases/download/v1.7.3/containerd-1.7.3-linux-amd64.tar.gz",
      "sudo tar Czxvf /usr/local containerd-1.7.3-linux-amd64.tar.gz",
      "wget https://raw.githubusercontent.com/containerd/containerd/main/containerd.service",
      "sudo mv containerd.service /usr/lib/systemd/system/",
      "sudo systemctl daemon-reload",
      "sudo systemctl enable --now containerd",
      "wget https://github.com/opencontainers/runc/releases/download/v1.1.8/runc.amd64",
      "sudo install -m 755 runc.amd64 /usr/local/sbin/runc",
      "sudo mkdir -p /etc/containerd/",
      "containerd config default | sudo tee /etc/containerd/config.toml",
      "sudo sed -i 's/SystemdCgroup \\= false/SystemdCgroup \\= true/g' /etc/containerd/config.toml",
      "sudo systemctl restart containerd",
      "sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates curl",
      "curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -",
      "cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list",
      "deb https://apt.kubernetes.io/ kubernetes-xenial main",
      "EOF",
      "sudo apt-get update",
      "sudo apt-get install -y kubelet kubeadm kubectl",
      "sudo apt-mark hold kubelet kubeadm kubectl containerd",
    ]
  }
}

resource "aws_s3_bucket" "blog_bucket" {
  bucket = "blog-bucket-20231106"

  tags = {
    Name        = "blog-bucket-20231106"
    Environment = "Production"
  }
}

output "master-ip" {
  value = aws_instance.blog_master.public_ip
}

output "worker1-ip" {
  value = aws_instance.blog_worker1.public_ip
}

output "worker2-ip" {
  value = aws_instance.blog_worker2.public_ip
}

output "lb-dns-hostname" {
  value = aws_lb.blog_nlb.dns_name
}

output "jenkins-ebs-id" {
  value = aws_ebs_volume.jenkins_volume.id
}
