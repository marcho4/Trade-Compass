import subprocess
import time
from pathlib import Path

import pytest
import requests
from testcontainers.core.container import DockerContainer
from testcontainers.core.network import Network
from testcontainers.core.waiting_utils import wait_for_logs
from testcontainers.postgres import PostgresContainer

ADMIN_API_KEY = "test-admin-key-integration"
PARSER_DIR = Path(__file__).parent.parent.parent


@pytest.fixture(scope="session", autouse=True)
def build_parser_image():
    already_built = subprocess.run(
        ["docker", "image", "inspect", "trade-compass-parser-test:latest"],
        capture_output=True,
    ).returncode == 0

    if not already_built:
        subprocess.run(
            ["docker", "build", "-t", "trade-compass-parser-test", "."],
            cwd=str(PARSER_DIR),
            check=True,
        )


@pytest.fixture(scope="session")
def docker_network():
    with Network() as network:
        yield network


@pytest.fixture(scope="session")
def postgres_container(docker_network):
    container = (
        PostgresContainer(image="postgres:16-alpine", username="tc_user", password="tc_pass", dbname="tc_db")
        .with_network(docker_network)
        .with_network_aliases("tc-postgres")
    )
    with container:
        yield container


@pytest.fixture(scope="session")
def kafka_container(docker_network):
    container = (
        DockerContainer(image="confluentinc/cp-kafka:7.6.0")
        .with_network(docker_network)
        .with_network_aliases("tc-kafka")
        .with_env("CLUSTER_ID", "MkU3OEVBNTcwNTJENDM2Qk")
        .with_env("KAFKA_NODE_ID", "1")
        .with_env("KAFKA_PROCESS_ROLES", "broker,controller")
        .with_env("KAFKA_LISTENERS", "PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093")
        .with_env("KAFKA_ADVERTISED_LISTENERS", "PLAINTEXT://tc-kafka:9092")
        .with_env(
            "KAFKA_LISTENER_SECURITY_PROTOCOL_MAP",
            "PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
        )
        .with_env("KAFKA_INTER_BROKER_LISTENER_NAME", "PLAINTEXT")
        .with_env("KAFKA_CONTROLLER_LISTENER_NAMES", "CONTROLLER")
        .with_env("KAFKA_CONTROLLER_QUORUM_VOTERS", "1@localhost:9093")
        .with_env("KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR", "1")
        .with_env("KAFKA_TRANSACTION_STATE_LOG_MIN_ISR", "1")
        .with_env("KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR", "1")
        .with_env("KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS", "0")
    )
    with container:
        wait_for_logs(container, "started", timeout=60)
        yield container


@pytest.fixture(scope="session")
def qdrant_container(docker_network):
    container = (
        DockerContainer(image="qdrant/qdrant:latest")
        .with_network(docker_network)
        .with_network_aliases("tc-qdrant")
        .with_exposed_ports(6333)
    )
    with container:
        wait_for_logs(container, "Qdrant HTTP listening on", timeout=30)
        yield container


@pytest.fixture(scope="session")
def parser_url(build_parser_image, postgres_container, kafka_container, qdrant_container, docker_network):
    container = (
        DockerContainer(image="trade-compass-parser-test")
        .with_network(docker_network)
        .with_exposed_ports(8081)
        .with_env("POSTGRES_HOST", "tc-postgres")
        .with_env("POSTGRES_PORT", "5432")
        .with_env("POSTGRES_DB", "tc_db")
        .with_env("POSTGRES_USER", "tc_user")
        .with_env("POSTGRES_PASSWORD", "tc_pass")
        .with_env("KAFKA_BOOTSTRAP_SERVERS", "tc-kafka:9092")
        .with_env("ADMIN_API_KEY", ADMIN_API_KEY)
        .with_env("YANDEX_CLOUD_S3_ACCESS_KEY_ID", "dummy")
        .with_env("YANDEX_CLOUD_S3_SECRET_ACCESS_KEY", "dummy")
        .with_env("BUCKET_NAME", "dummy-bucket")
        .with_env("GEMINI_API_KEY", "dummy")
        .with_env("QDRANT_HOST", "tc-qdrant")
        .with_env("QDRANT_PORT", "6333")
    )
    with container:
        host = container.get_container_host_ip()
        port = container.get_exposed_port(8081)
        base_url = f"http://{host}:{port}"

        deadline = time.time() + 15
        while time.time() < deadline:
            try:
                response = requests.get(f"{base_url}/health", timeout=2)
                if response.status_code == 200:
                    break
            except requests.RequestException:
                pass
            time.sleep(1)
        else:
            logs = container.get_logs()
            raise RuntimeError(
                f"Parser service did not become healthy within 15 seconds.\nLogs:\n{logs}"
            )

        yield base_url
