# Leanix-Exporter

Extract Kubernetes Cluster data for different source and expose it is a simple format for Leanix Integration

## What LeanIX should store

State Data, e.g. what is running (meta data)
 * Taking Kubernetes Services, should give us a general overview of what is running in a cluster that is published
 * Then possibly take other more internal things like deamonsets for internal things
 
Time Series Data
 * When did a deployment takes place and how many deployments are taking place
 * Resource Usage of Containers (at the time of deployment)
 * Total Resource Usage and Availability of Cluster (check Giant Swarm API)

## When the LeanIX Exporter is called

At best, we run on changes. otherwise we have two options.

 * Run it every hour, minute, ... (resource usage needs to be checked as the exporter is currently querying the entire cluster)
 * Allow the customer to tell the exporter to run based on e.g. deployments (check if we can listen to certain events)

## Result example

```json
{
  "Namespaces": [
    {
      "Name": "default",
      "Pods": [
        {
          "Name": "demo-2476732077-wr8nb",
          "Status": "Pending",
          "ContainerStatuses": [
            {
              "name": "demo",
              "state": {
                "waiting": {
                  "reason": "ErrImageNeverPull",
                  "message": "Container image \"giantswarm/leanix-sexporter\" is not present with pull policy of Never"
                }
              },
              "lastState": {},
              "ready": false,
              "restartCount": 0,
              "image": "giantswarm/leanix-sexporter",
              "imageID": ""
            }
          ]
        }
      ]
    },
    {
      "Name": "giantswarm",
      "Pods": [
        {
          "Name": "leanix-exporter-4242224814-6lzll",
          "Status": "Running",
          "ContainerStatuses": [
            {
              "name": "leanix-exporter",
              "state": {
                "running": {
                  "startedAt": "2017-07-18T12:32:01Z"
                }
              },
              "lastState": {},
              "ready": true,
              "restartCount": 0,
              "image": "giantswarm/leanix-exporter:latest",
              "imageID": "docker://sha256:012b66025e29c5f50abdfa08004d124dbffe90cfc6c216753980ff148781ae3d",
              "containerID": "docker://ca516a6ee65e870883b80e0fa57318600d3ed177ee88e56e9633f4d8e1a92754"
            }
          ]
        }
      ]
    }
  ],
  "LastUpdate": "2017-07-18T12:32:16.298303791Z"
```

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/leanix-exporter/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the
contribution workflow as well as reporting bugs.

## License

leanix-exporter is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.

## Credit
- https://golang.org
- https://github.com/giantswarm/microkit
- https://github.com/kubernetes/client-go

