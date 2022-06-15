region = "eu-central-1"

availability_zones = ["eu-central-1a", "eu-central-1b", "eu-central-1c"]

deployment_stage = "mainnet"

forbidden_account_ids = ["909899099608"]

# -----------------------------------------------------------------------------
# Module service-ferp
# -----------------------------------------------------------------------------

service_name = "ferp"

service_desired_count = 1

service_container_port = 9002

task_network_mode = "awsvpc"

task_cpu = 256

task_memory = 512

target_health_path = "/health"

target_health_interval = 120

target_health_timeout = 5

target_health_matcher = "200"

subdomain_name = "ferp"

env_currency_converter_api_key = "de5062430b8774b3b35d"

env_rpc_port = 50000

env_shutdown_timeout = 20
