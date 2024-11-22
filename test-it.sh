LB_IP=$(kubectl get svc/foo-service -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
for _ in {1..10}; do
  echo " -- "
  curl ${LB_IP}:5678
  echo " -- "
done
