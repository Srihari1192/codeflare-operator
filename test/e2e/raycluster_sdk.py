import sys

from time import sleep

from codeflare_sdk.cluster.cluster import Cluster, ClusterConfiguration

namespace = sys.argv[1]

cluster = Cluster(ClusterConfiguration(
    name='mnist',
    namespace=namespace,
    num_workers=1,
    min_cpus='500m',
    max_cpus=1,
    min_memory=0.5,
    max_memory=1,
    num_gpus=0,
    image="quay.io/project-codeflare/ray:2.5.0-py38-cu116", #current default
    instascale=False,
))

cluster.up()

sleep(120)

cluster.status()

cluster.wait_ready()

cluster.status()

cluster.details()
