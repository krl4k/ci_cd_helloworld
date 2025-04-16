I installed pg as dependecy like this

Added deps into chart file and execute.

helm dependency update helm/hello-service

then deployed to the cluster


cd helm/hello-service && helm upgrade --install hello-service . --namespace dev --create-namespace