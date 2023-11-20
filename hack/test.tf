locals {
  cluster_name = "testcluster"
}

data "external" "example" {
  program = ["bash", "-c", "./test.sh"]
  query = {
    name = "bsaba"
  }
}

output "example_output" {
  value = data.external.example.result.output_value
}

// import 블록
// document 리소스