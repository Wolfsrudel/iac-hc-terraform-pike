resource "aws_iam_policy" "basic" {
  name = "basic"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:CreateTable",
          "dynamodb:TagResource",
          "dynamodb:DescribeTable",
          "dynamodb:DescribeContinuousBackups",
          "dynamodb:DescribeTimeToLive",
          "dynamodb:ListTagsOfResource",
          "dynamodb:DeleteTable",
          "dynamodb:UntagResource",
          "dynamodb:UpdateTable",
          "dynamodb:UpdateTimeToLive",
          "dynamodb:CreateTableReplica",
          "dynamodb:Scan",
          "dynamodb:Query",
          "dynamodb:UpdateItem",
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:DeleteItem",
          "dynamodb:BatchWriteItem",
          "iam:CreateServiceLinkedRole",
          "dynamodb:UpdateContinuousBackups",
          "dynamodb:DeleteTableReplica"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
  tags = {
    createdby = "JamesWoolfenden"
  }
}

resource "aws_iam_role_policy_attachment" "basic" {
  role       = aws_iam_role.basic.name
  policy_arn = aws_iam_policy.basic.arn
}

resource "aws_iam_user_policy_attachment" "basic" {
  # checkov:skip=CKV_AWS_40: By design
  user       = "basic"
  policy_arn = aws_iam_policy.basic.arn
}
