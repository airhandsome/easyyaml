package templates

type K8sTemplate struct {
	Name    string
	Content string
}

func GetK8sDeploymentTemplate() string {
	return `apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: nginx:latest
        ports:
        - containerPort: 80`
}
