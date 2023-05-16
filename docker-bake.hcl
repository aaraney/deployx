variable "GO_VERSION" {
  # default ARG value set in Dockerfile
  default = null
}

variable "VERSION_TAG" {
  default = "0.0.1"
}

variable "DOCKERFILE" {
    default = "./Dockerfile"
}

# Special target: https://github.com/docker/metadata-action#bake-definition
target "meta-helper" {}

target "default" {
  tags = ["aaraney/docker-deployx:latest", "aaraney/docker-deployx:${VERSION_TAG}"]
  args = {
    GO_VERSION = GO_VERSION
  }
  pull = true
  target = "shell"
  dockerfile = "${DOCKERFILE}"
  output = ["type=image"]

  platforms = [
    // "linux/386", // TODO: determine why this is failing
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
  ]
}
