import sys
import os

from time import sleep

from torchx.specs.api import AppState, is_terminal

from codeflare_sdk.cluster.cluster import get_cluster
from codeflare_sdk.job.jobs import DDPJobDefinition

namespace = sys.argv[1]

host = os.getenv('CLUSTER_HOSTNAME')

ingress_options = {}
if host is not None:
    ingress_options = {
        "ingresses": [
            {
                "ingressName": "ray-dashboard",
                "port": 8265,
                "pathType": "Prefix",
                "path": "/",
                "host": host,
            },
        ]
    }



cluster = get_cluster('mnist',namespace)

print(cluster.details())

jobdef = DDPJobDefinition(
    name="mnist",
    script="mnist.py",
    scheduler_args={"requirements": "requirements.txt"},
)
job = jobdef.submit(cluster)

done = False
time = 0
timeout = 300
while not done:
    status = job.status()
    if is_terminal(status.state):
        break
    if not done:
        print(status)
        if timeout and time >= timeout:
            raise TimeoutError(f"job has timed out after waiting {timeout}s")
        sleep(5)
        time += 5

print(f"Job has completed: {status.state}")

print(job.logs())

cluster.down()

if not status.state == AppState.SUCCEEDED:
    exit(1)
else:
    exit(0)
