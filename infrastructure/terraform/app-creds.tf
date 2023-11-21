resource "random_password" "db" {
  length  = 16
  special = false
}

resource "random_password" "jwt_key" {
  length  = 32
  special = false
}
