resource "appclacks_healthcheck_command" "test_command" {
  name = "test-command"
  interval = "30s"
  timeout = "5s"
  description = "example command health check"
  labels = {
    "env": "prod"
  }
  command = "ls"
  arguments = ["/"]
}
