# Python Client

### Setup
- The project uses poetry for build, and virutal env management.
- Make sure to install poetry via https://python-poetry.org/docs/#installing-manually
- After installing poetry you can activate the env by `poetry env use`
- Install all dependencies using `poetry install --no-root` (no-root tells that the package is not at the root of the directory)
- For setting up in IDE, make sure to setup the interpreter to use the virtual environment that was created when you activated poetry env.

Note: During development, make sure to open just the python directory, otherwise the IDE misconfigures the imports.
