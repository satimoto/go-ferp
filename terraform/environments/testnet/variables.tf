variable "region" {
  description = "The AWS region"
  default     = "eu-central-1"
}

variable "availability_zones" {
  description = "A list of Availability Zones where subnets and DB instances can be created"
}

variable "deployment_stage" {
  description = "The deployment stage"
  default     = "testnet"
}

variable "forbidden_account_ids" {
  description = "The forbidden account IDs"
  type        = list(string)
  default     = []
}

# -----------------------------------------------------------------------------
# Module service-ferp
# -----------------------------------------------------------------------------

variable "service_name" {
  description = "The name of the service"
}

variable "service_desired_count" {
  description = "The number of instances of the task definition to place and keep running"
}

variable "service_container_port" {
  description = "The port on the container to associate with the load balancer"
}

variable "service_metric_port" {
  description = "The port to associate with metric collection"
}

variable "task_network_mode" {
  description = "The Docker networking mode to use for the containers in the task"
}

variable "task_cpu" {
  description = "The number of cpu units used by the task"
}

variable "task_memory" {
  description = "The amount (in MiB) of memory used by the task"
}

variable "target_health_path" {
  description = "The path to check the target's health"
}

variable "target_health_interval" {
  description = "The approximate amount of time, in seconds, between health checks of an individual target"
}

variable "target_health_timeout" {
  description = "The amount of time, in seconds, during which no response means a failed health check"
}

variable "target_health_matcher" {
  description = "The HTTP codes to use when checking for a successful response from a target"
}

variable "subdomain_name" {
  description = "The subdomain name of the service"
}

variable "env_currency_converter_api_key" {
  description = "The environment variable to set the currency converter api key"
}

variable "env_open_exchange_rate_api_key" {
  description = "The environment variable to set the open exchange rate api key"
}

variable "env_rpc_port" {
  description = "The environment variable to set the RPC port"
}

variable "env_shutdown_timeout" {
  description = "The environment variable to set the shutdown timeout"
}
