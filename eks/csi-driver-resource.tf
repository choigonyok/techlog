data "aws_iam_policy" "efs_csi_policy" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonEFSCSIDriverPolicy"
}

resource "aws_iam_role_policy_attachment" "policy_attachment" {
  #   policy_arn = aws_iam_policy.policy.arn
  policy_arn = data.aws_iam_policy.efs_csi_policy.arn
  role       = aws_iam_role.iamrole.name
}

##########################################

# import {
#     to = aws_iam_policy.policy
#     id = "arn:aws:iam::aws:policy/service-role/AmazonEFSCSIDriverPolicy"
# }

# resource "aws_iam_policy" "policy" {
#     # policy = jsonencode({}) 
# }

##########################################

# resource "aws_iam_role_policy" "node" {
#   name = "node"
#   role = "${aws_iam_role.iamrole.name}"
# #   policy =  jsonencode({
# #     "Version": "2012-10-17",
# #     "Statement": [
# #       {
# #         "Action": [
# #           "sts:AssumeRole"
# #         ],
# #         "Effect": "Allow",
# #         "Resource": "${aws_iam_role.node.arn}"
# #       },
# #       {
# #         "Sid": "",
# #         "Effect": "Allow",
# #         "Action": [
# #             "*"
# #         ],
# #         "Resource": "*"
# #       }
# #     ]
# #   })
#     policy = jsonencode({})
# }

##########################################

resource "aws_efs_file_system" "example" {
  creation_token = "${local.cluster_name}-example"

  tags = {
    Name = "${local.cluster_name}-example"
  }
}

resource "aws_efs_mount_target" "example" {
  count          = 3
  file_system_id = aws_efs_file_system.example.id
  #   subnet_id = aws_subnet.this.*.id[count.index]
  subnet_id = module.vpc.private_subnets[count.index]
  #   security_groups = [aws_security_group.eks-cluster.id]
  security_groups = [module.eks.node_security_group_id]
}

resource "aws_eks_addon" "efs-csi" {
  cluster_name  = module.eks.cluster_name
  addon_name    = "aws-efs-csi-driver"
  addon_version = "v1.7.0-eksbuild.1"
  # aws eks describe-addon-versions --kubernetes-version 1.27 --addon-name aws-efs-csi-driver
  # 현재 클러스터 버전에서 사용가능한 addon 버전목록 확인 가능
  service_account_role_arn = aws_iam_role.iamrole.arn
  tags = {
    "eks_addon" = "efs-csi"
    "terraform" = "true"
  }
}

resource "aws_iam_role" "iamrole" {
  name = "efs"

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

##########################################

# data "aws_iam_policy_document" "assume_role_policy" {
#   statement {
#     actions = ["sts:AssumeRoleWithWebIdentity"]

#     principals {
#       type = "federated"
#       #   identifiers = [aws_iam_openid_connect_provider.my_eks_oidc.arn]
#       identifiers = [module.eks.oidc_provider_arn]
#     }

#     condition {
#       test     = "StringEquals"
#       variable = "${replace(module.eks.cluster_oidc_issuer_url, "https://", "")}:sub"
#       values   = ["system:serviceaccount:default:debug-sa"]
#     }
#   }
# }

