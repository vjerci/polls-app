# defines ecr to push docker images to
resource "aws_ecr_repository" "polls_images" {
  name                 = "polls-app"
  image_tag_mutability = "MUTABLE"

  force_delete = true

  image_scanning_configuration {
    scan_on_push = true
  }
}
