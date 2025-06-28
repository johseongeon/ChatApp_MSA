# ChatApp server for Microservice Architecture

---

The MongoDB connection URI(`MONGO_URI`) is defined on line 26 of each Dockerfile.

Unless you change `MONGO_URI`, the connection will default to the MongoDB instance on the host’s port 27017.

You can connect to an external MongoDB instance by modifying the `MONGO_URI` environment variable.

---

## Using Guide

### 1. build image

```cmd
cd {directory}                   // navigate to specific directory
docker build -t {image_name} .   // and build image
```

### 2. run

```cmd
docker run -it -p {external port}:{internal port} {image_name}
```

---

For examgple, to run user_manager server

```cmd
cd user_manager
docker build -t user-manager .
docker run -it -p 8082:8082 user-manager
```
