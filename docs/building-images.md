# ContainerCafe images and how to build them

## openradiant/remoteabac

```
cd openradiant/remoteabac
docker login (with dockerhub creds) 
docker build -t openradiant/remoteabac .
docker push openradiant/remoteabac
```

## containercafe/proxy

(to be written)

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
