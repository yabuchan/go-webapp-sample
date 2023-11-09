#!/bin/bash
LOG_LEVEL=$1
pod=$(kubectl get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | grep app-server)
echo "$pod"
kubectl exec ${pod} /bin/bash -c  'echo $(level) >  /logLevel';