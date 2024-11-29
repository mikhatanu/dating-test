# dating-test

## Structure

main.go -> http handler -> db cal -> sqlite

Uses go serve mux as the http handler. Main file will use handlerfunction from another folder. Handler calls function that does the actual database query

## How to run
### Localy
```sh
git clone https://github.com/mikhatanu/dating-test.git
go build

#execute the binary, example in windows
dating-test.exe
```

### Kubernetes
#### Prerequisites
1. Kubernetes cluster installed with ingress

#### Steps
1. Build image using docker or similar tools
```sh 
docker build -t <your_repo>/dating-test:1.0.0 .
```
2. Upload to docker image repository
3. Change docker image to the one you upload
4. Apply the file
```bash 
kubectl apply -f kubernetes/deployment.yaml
``` 