# Python Client

### Setup
- The project uses poetry for build, and virutal env management.
- Make sure to install poetry via https://python-poetry.org/docs/#installing-manually
- After installing poetry you can activate the env by `poetry env use`
- Install all dependencies using `poetry install --no-root --with=dev` (no-root tells that the package is not at the root of the directory)
- For setting up in IDE, make sure to setup the interpreter to use the virtual environment that was created when you activated poetry env.

### Lint and Formatting
- We use black for formatting of python files and pylint, ruff for linting the python files.
- You can check the command for running lint and formating by referring to `test-python-client.yml` workflow.

### Usage
- You can use the raccoon by installing it from PyPi by the following command
  - From Pypi
  ```pip install raccoon_client``` 
  - From Github 
  ```pip install raccoon_client@git+https://github.com/raystack/raccoon@$VERSION#subdirectory=clients/python```
    where $VERSION is a git tag.
- An example on how to use the client is under the [examples](examples) package.

### Confiugration
The client supports the following configuration:

| Name    | Description                                                                       | Type                              | Default |
|---------|-----------------------------------------------------------------------------------|-----------------------------------|---------|
| url     | The remote server url to connect to                                               | string                            | ""      |
| retries | The max number of retries to be attempted before an event is considered a failure | int (<10)                         | 0       |
| timeout | The number of seconds to wait before timing out the request                       | float                             | 1.0     |
| serialiser | The format to which event field of client.Event serialises it's data to           | Serialiser Enum(JSON or PROTOBUF) | JSON    |
|wire_type | The format in which the request payload should be sent to server                  | Wire Type Enum(JSON or PROTOBUF)  | JSON | 
| headers | HTTP header key value pair to be sent along with each request | dict                              | {} |


Note: 
- During development, make sure to open just the python directory, otherwise the IDE misconfigures the imports.
- The protos package contain generated code and should not be edited manually.

