# Build for arm

  * https://collabnix.com/building-arm-based-docker-images-on-docker-desktop-made-possible-using-buildx/
  * https://www.docker.com/blog/multi-arch-build-and-images-the-simple-way/
  * https://www.docker.com/blog/multi-arch-build-what-about-gitlab-ci/
  * https://github.com/features/packages
  * https://docs.github.com/en/free-pro-team@latest/packages/getting-started-with-github-container-registry/migrating-to-github-container-registry-for-docker-images

```
docker buildx create --name testbuilder
docker buildx use testbuilder
docker buildx ls
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t andrelohmann/k8s-in-a-box .
```
