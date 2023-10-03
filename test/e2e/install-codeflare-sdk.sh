#!/bin/bash

# Install Poetry and configure virtualenvs
pip install poetry
poetry config virtualenvs.create false

# Clone the CodeFlare SDK repository
git clone --branch main https://github.com/project-codeflare/codeflare-sdk.git
cd codeflare-sdk

# Lock dependencies and install them
poetry lock --no-update
poetry install --with test,docs

# Return to the previous directory
cd ..
