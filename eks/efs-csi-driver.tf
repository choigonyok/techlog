data "aws_iam_policy" "efs_csi_driver" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonEFSCSIDriverPolicy"
}

resource "aws_iam_role_policy_attachment" "efs_csi_driver" {
  policy_arn = data.aws_iam_policy.efs_csi_driver.arn
  role       = aws_iam_role.efs_csi_driver.name
}

resource "aws_efs_file_system" "jenkins_volume" {
  creation_token = "${local.cluster_name}-jenkins_volume"
  tags = {
    Name = "${local.cluster_name}-jenkins_volume"
  }
}

resource "aws_eks_addon" "efs-csi" {
  cluster_name             = module.eks.cluster_name
  addon_name               = "aws-efs-csi-driver"
  addon_version            = "v1.7.0-eksbuild.1"
  service_account_role_arn = aws_iam_role.efs_csi_driver.arn
  tags = {
    "eks_addon" = "efs-csi"
    "terraform" = "true"
  }
}

resource "aws_efs_mount_target" "private_subnet" {
  count           = 3
  file_system_id  = aws_efs_file_system.jenkins_volume.id
  subnet_id       = module.vpc.private_subnets[count.index]
  security_groups = [module.eks.node_security_group_id]
}

resource "aws_iam_role" "efs_csi_driver" {
  name = "efs_csi_driver"
  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Principal" : {
          "Federated" : "${module.eks.oidc_provider_arn}"
        },
        "Action" : "sts:AssumeRoleWithWebIdentity",
        "Condition" : {
          "StringLike" : {
            "${module.eks.oidc_provider}:sub" : "system:serviceaccount:kube-system:efs-csi*",
            "${module.eks.oidc_provider}:aud" : "sts.amazonaws.com"
          }
        }
      }
    ]
    }
  )
} 