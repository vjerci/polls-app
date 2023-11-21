# defines ecr to push docker images to
resource "aws_ecr_repository" "polls_api" {
  name                 = "polls-api"
  image_tag_mutability = "MUTABLE"

  force_delete = true

  image_scanning_configuration {
    scan_on_push = true
  }
}


resource "aws_ecr_repository" "polls_front" {
  name                 = "polls-front"
  image_tag_mutability = "MUTABLE"

  force_delete = true

  image_scanning_configuration {
    scan_on_push = true
  }
}
