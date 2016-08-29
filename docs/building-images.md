# ContainerCafe images and how to build them

## openradiant/remoteabac

```
cd openradiant/remoteabac
docker login (with dockerhub creds)
docker build -t openradiant/remoteabac .
docker push openradiant/remoteabac
```

## containercafe/proxy

```bash
docker login (with dockerhub creds)
cd proxy
./builddocker.sh
# see you new image:
docker images
docker tag api-proxy containercafe/api-proxy
docker push containercafe/api-proxy
```

## openradiant/km

Building is straightforward:

```
cd openradiant/misc/dockerfiles/km
docker login (with dockerhub creds)
docker build -t openradiant/km .
docker push openradiant/km
```

The selection of *what* to build is more complicated.  The choice is
in the Dockerfile.  Documentation needed for the why behind that
choice.
