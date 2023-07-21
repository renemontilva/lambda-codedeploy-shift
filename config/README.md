
# Dockerfile for Go Lambda Function


This Dockerfile provides a template for building and running a Go Lambda function using Docker. It includes separate stages for building and development (dev) and production (prod) environments.




## Usage/Examples

To build and run the Go Lambda function locally, follow these steps:

Install Docker on your machine. Refer to the Docker documentation for instructions specific to your operating system.

Open a terminal and navigate to the project directory containing the Dockerfile.

Build the Docker dev image using the following command:

```bash
docker image build -t codedeployshift:dev --target dev -f config/Dockerfile .
```

Build the Docker prod image using the following command:

```bash
docker image build -t codedeployshift:prod --target prod -f config/Dockerfile .

```

Once the Docker image is built, you can run it in either development or production mode:

### Development Mode:
Use the following command to start a development container:
```bash
docker container run -p 9000:8080  codedeployshift:dev 
```
From a new terminal, post an event to the following endpoint using a curl command:

```bash
curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
```


## Additional Notes

* The Dockerfile uses the public.ecr.aws registry to pull the base images for the builder, development, and production stages. If you encounter any issues related to image availability, ensure that you have the necessary permissions and access to the specified registry.

* The builder stage is responsible for building the Go application. It copies the project files, downloads the dependencies, and builds the main executable.

* The development (dev) and production (prod) stages copy the main executable from the builder stage into the respective containers. The CMD instruction runs the main executable when the container starts.

* Feel free to modify the Dockerfile as needed to incorporate additional dependencies, environment variables, or build steps specific to your Go Lambda function.

* Remember to handle any security or authentication requirements when deploying your Go Lambda function to AWS Lambda, as this Dockerfile focuses primarily on local development and testing.