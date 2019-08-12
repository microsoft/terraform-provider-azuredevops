import os, sys
from dotenv import load_dotenv



def get_env(key):
    value = os.getenv(key)
    if not value:
        print('No value in environment for key {}. Exiting'.format(key))
        sys.exit(-1)
    return value


AZDO_HOST = "dev.azure.com"
AZDO_API_VERSION = "api-version=5.1-preview"


load_dotenv()
AZDO_PAT = get_env('AZDO_PAT')
AZDO_ORGANIZATION = get_env('AZDO_ORGANIZATION')
AZDO_NEW_PROJECT_NAME = get_env('AZDO_NEW_PROJECT_NAME')
AZDO_NEW_PROJECT_DESCRIPTION = get_env('AZDO_NEW_PROJECT_DESCRIPTION')
AZDO_NEW_PIPELINE_NAME = get_env('AZDO_NEW_PIPELINE_NAME')
AZDO_PIPELINE_YML_GIT_REPO = get_env('AZDO_PIPELINE_YML_GIT_REPO')
AZDO_PIPELINE_YML_GIT_REPO_BRANCH = get_env('AZDO_PIPELINE_YML_GIT_REPO_BRANCH')
AZDO_PIPELINE_YML_FILENAME = get_env('AZDO_PIPELINE_YML_FILENAME')
AZDO_GITHUB_SERVICE_CONNECTION_NAME = get_env('AZDO_GITHUB_SERVICE_CONNECTION_NAME')
AZDO_GITHUB_SERVICE_CONNECTION_PAT = get_env('AZDO_GITHUB_SERVICE_CONNECTION_PAT')
