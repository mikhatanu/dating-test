# dating-test

## Structure

creates one app.db file (sqlite3) that serves as the database. Database called from auth folder.

### Functionality

#### **Login (/rest/v1/login)**

Logs the user in

#### **Signup (/rest/v1/signup)**

Create user

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

1. Git clone

```sh
git clone https://github.com/mikhatanu/dating-test.git
```

2. Build image using docker or similar tools

```sh
docker build -t <your_repo>/dating-test:1.0.0 .
```

3. Upload to docker image repository

```sh
docker push <your_repo>/dating-test:1.0.0
```

4. Change docker image url in kubernetes deployment to

```yaml
spec:
  containers:
    - name: dating-test
      imagePullPolicy: IfNotPresent
      image: <your_repo>/dating-test:1.0.0 #here
      ports:
        - containerPort: 3000
```

5. Apply the file

```bash
kubectl apply -f kubernetes/deployment.yaml
```

Note: There is no PV in deployement.yaml, so data will be lost on restart.
