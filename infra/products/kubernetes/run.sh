kubectl create configmap products-config --from-file=config.development.json
kubectl apply -f ./app.yaml 
kubectl apply -f ./config.yaml 
kubectl apply -f ./secret.yaml 


kubectl delete -f ./app.yaml 
kubectl delete -f ./config.yaml 
kubectl delete -f ./secret.yaml 