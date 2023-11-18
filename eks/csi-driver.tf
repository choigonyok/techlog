# data "aws_iam_policy" "efs_csi_policy" {
#   arn = "arn:aws:iam::aws:policy/service-role/AmazonEFSCSIDriverPolicy"
# }

# module "irsa-efs-csi" {
#   source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
#   version = "4.7.0"

#   create_role                   = true
#   role_name                     = "AmazonEKSTFEFSCSIRole-${module.eks.cluster_name}"
#   provider_url                  = module.eks.oidc_provider
#   role_policy_arns              = [data.aws_iam_policy.efs_csi_policy.arn]
#   oidc_fully_qualified_subjects = ["system:serviceaccount:kube-system:efs-csi-controller-sa"]
# }

# resource "aws_eks_addon" "efs-csi" {
#   cluster_name             = module.eks.cluster_name
#   addon_name               = "aws-efs-csi-driver"
#   addon_version            = "v1.7.0-eksbuild.1"
#   # aws eks describe-addon-versions --kubernetes-version 1.27 --addon-name aws-efs-csi-driver
#   # 현재 클러스터 버전에서 사용가능한 addon 버전목록 확인 가능

#   service_account_role_arn = module.irsa-efs-csi.iam_role_arn
#   tags = {
#     "eks_addon" = "efs-csi"
#     "terraform" = "true"
#   }
# }

# resource "aws_efs_file_system" "example" {
#   creation_token = "${local.cluster_name}-example"

#   tags = {
#     Name = "${local.cluster_name}-example"
#   }
# }

# resource "aws_efs_mount_target" "example" {
#   count          = 3
#   file_system_id = aws_efs_file_system.example.id
#   #   subnet_id = aws_subnet.this.*.id[count.index]
#   subnet_id = module.vpc.private_subnets[count.index]
#   #   security_groups = [aws_security_group.eks-cluster.id]
#   security_groups = [module.eks.node_security_group_id]
# }