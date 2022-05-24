# Dockerfile for building your infra projects based on cops-hq

Since cops-hq requires quite a lot of dependencies, best practice is to package your infra CLI along with all the 
dependencies in a Dockerfile, which will then be used in CI/CD. For this purpose, you can use the Dockerfile in the
root of this project as a boilerplate. 