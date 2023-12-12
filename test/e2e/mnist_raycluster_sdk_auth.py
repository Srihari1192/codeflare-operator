import sys
import os

from time import sleep

from codeflare_sdk.cluster.cluster import Cluster, ClusterConfiguration
# from codeflare_sdk.cluster.auth import TokenAuthentication
from codeflare_sdk.job.jobs import DDPJobDefinition

namespace = sys.argv[1]
ray_image = os.getenv('RAY_IMAGE')

cluster = Cluster(ClusterConfiguration(
    name='mnist',
    namespace=namespace,
    num_workers=1,
    head_cpus='500m',
    head_memory=2,
    min_cpus='500m',
    max_cpus=1,
    min_memory=0.5,
    max_memory=2,
    num_gpus=0,
    instascale=False,
    image=ray_image,
     openshift_oauth=True,
))

# auth = TokenAuthentication(
#     token = "sha256~sL0lb8DDTxx3yPIJCbD6voVqi02Ue0QIp7ln1Aml7aM",
#     server = "https://api.svenkata.dev.datahub.redhat.com:6443",
#     skip_tls= True
# )
# auth.login()

cluster.up()

cluster.status()

cluster.wait_ready()

cluster.status()

cluster.details()

jobdef = DDPJobDefinition(
    name="mnist",
    script="mnist.py",
    scheduler_args={"requirements": "requirements.txt"},
)
job = jobdef.submit(cluster)

