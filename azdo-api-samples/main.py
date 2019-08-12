import requests as r
import base64
import config as c
import backoff
import sys


# This file shows some of the Azure DevOps API primitives that can be used to create the following:
#   (1) New project
#   (2) New service connection with permissions to access GitHub repositorys
#   (3) New build pipeline tied to GitHub repository
#
# This is semi-configurable. You will need to create a file called `.env` and populate it with
# environment variables (see `.env.template`).


# Build headers needed by azure devops to authenticate using a Personal Access Token
def get_auth_headers(extra_headers = None):
    token = base64.b64encode(bytes(':' + c.AZDO_PAT, 'utf-8')).decode('utf-8')
    headers = { 'Authorization': 'Basic ' + token }
    if extra_headers is not None:
        headers.update(extra_headers)
    
    return headers

# Fail if an API response contains a non 200 status code
def fail_on_api_error(api_response):
    if api_response.status_code > 299:
        print('HTTP request failed with code {}. Response was `{}`. Exiting program early...'.format(
            api_response.status_code,
            api_response.text))
        sys.exit(-1)

# Lists the name of instances of a resource within AzureDevOps
def azdo_list_resource_names(url):
    print('HTTP GET: ' + url)
    api_response = r.get(url, headers=get_auth_headers())
    fail_on_api_error(api_response)
    return { resource['name']: resource['id'] for resource in api_response.json()['value'] }

# Create a project in AZDO
def azdo_create_project():
    headers = get_auth_headers({'Content-Type': 'application/json'})
    url = 'https://{0}/{1}/_apis/projects?{2}'.format(c.AZDO_HOST, c.AZDO_ORGANIZATION, c.AZDO_API_VERSION)
    print('HTTP POST: ' + url)
    api_response = r.post(url, headers=headers, json={
        'name': c.AZDO_NEW_PROJECT_NAME,
        'description': c.AZDO_NEW_PROJECT_DESCRIPTION,
        'visibility': 0,
        'capabilities': {
            'versioncontrol': { 'sourceControlType': 'Git' },
            'processTemplate': { 'templateTypeId': 'adcc42ab-9882-485e-a3ed-7678f01f66bc' }
        }
    })

    fail_on_api_error(api_response)

    # This step uses the callback URL from the response of the create operation to wait for
    # the project to finish provisioning
    api_response_json = api_response.json()
    callback_url = api_response_json['url']
    print('Using callback URL to check on project provisioning: ' + callback_url)
    if azdo_is_project_finished_provisioning(callback_url):
        print('Created project with name ' + c.AZDO_NEW_PROJECT_NAME)
    else:
        print('Project was not created. Application will exit early...')
        sys.exit(-1)
    
    return api_response_json['id']
    
@backoff.on_predicate(backoff.fibo, max_time=10)
def azdo_is_project_finished_provisioning(callback_url):
    response = r.get(callback_url, headers=get_auth_headers())
    fail_on_api_error(response)
    response_js = response.json()
    return 'status' in response_js and response_js['status'] == 'succeeded'

# Create a build definition in Azure DevOps
def azdo_create_build_definition(service_connection_id):
    headers = get_auth_headers({'Content-Type': 'application/json'})
    url = 'https://{0}/{1}/{2}/_apis/build/definitions?{3}'.format(c.AZDO_HOST, c.AZDO_ORGANIZATION, c.AZDO_NEW_PROJECT_NAME, c.AZDO_API_VERSION)
    print('HTTP POST: ' + url)
    api_response = r.post(url, headers=headers, json={
        'quality': 'definition',
        'name': c.AZDO_NEW_PIPELINE_NAME,
        'path': '',
        'type': 'build',
        'queueStatus': 'enabled',
        'repository': {
            'url': 'https://github.com/{}.git'.format(c.AZDO_PIPELINE_YML_GIT_REPO),
            'id': c.AZDO_PIPELINE_YML_GIT_REPO,
            'name': c.AZDO_PIPELINE_YML_GIT_REPO,
            'defaultBranch': c.AZDO_PIPELINE_YML_GIT_REPO_BRANCH,
            'type': 'GitHub',
            'clean': True,
            'properties': {
                'connectedServiceId': service_connection_id
            }
        },
        'process': {
            'type': 2,
            'yamlFilename': c.AZDO_PIPELINE_YML_FILENAME
        },
        'queue': {
            'name': 'Hosted Ubuntu 1604',
            'pool': {
                'id': 224,
                'name': 'Hosted Ubuntu 1604'
            }
        }
    })

    fail_on_api_error(api_response)
    print('Created pipeline with name ' + c.AZDO_NEW_PIPELINE_NAME)


# Create a service connection in Azure DevOps
def azdo_create_service_connection(project_id):
    headers = get_auth_headers({'Content-Type': 'application/json'})
    url = 'https://{0}/{1}/{2}/_apis/serviceendpoint/endpoints?{3}'.format(c.AZDO_HOST, c.AZDO_ORGANIZATION, c.AZDO_NEW_PROJECT_NAME, c.AZDO_API_VERSION)
    print('HTTP POST: ' + url)

    print('project id', project_id)

    api_response = r.post(url, headers=headers, json={
        'name': c.AZDO_GITHUB_SERVICE_CONNECTION_NAME,
        'type': 'github',
        'url': 'http://github.com',
        'owner': 'library',
        'authorization': {
            'parameters': {
                'accessToken': c.AZDO_GITHUB_SERVICE_CONNECTION_PAT
            },
            'scheme': 'PersonalAccessToken'
        },
        'serviceEndpointProjectReferences': [
            { 
                'projectReference': {
                    'name': project_id
                },
                'name': c.AZDO_GITHUB_SERVICE_CONNECTION_NAME
            }
        ]
    })

    fail_on_api_error(api_response)
    print('Created Service Connection with name ' + c.AZDO_GITHUB_SERVICE_CONNECTION_NAME)
    return api_response.json()['id']




def azdo_example():

    ### Create an Azure DevOps project (if one does not already exist)
    projects = azdo_list_resource_names('https://{0}/{1}/_apis/projects?{2}'.format(
        c.AZDO_HOST,
        c.AZDO_ORGANIZATION,
        c.AZDO_API_VERSION))
    if c.AZDO_NEW_PROJECT_NAME in projects.keys():
        print('Project with name {} already exists. Skipping creation'.format(c.AZDO_NEW_PROJECT_NAME))
        project_id = projects[c.AZDO_NEW_PROJECT_NAME]
    else:
        project_id = azdo_create_project()
    

    ### Create a Service Connection to GitHub (if one does not already exist)
    service_connections = azdo_list_resource_names('https://{0}/{1}/{2}/_apis/serviceendpoint/endpoints?{3}'.format(
        c.AZDO_HOST,
        c.AZDO_ORGANIZATION,
        c.AZDO_NEW_PROJECT_NAME,
        c.AZDO_API_VERSION))
    if c.AZDO_GITHUB_SERVICE_CONNECTION_NAME in service_connections:
        print('Service Connection with name {} already exists. Skipping creation.'.format(c.AZDO_GITHUB_SERVICE_CONNECTION_NAME))
        service_connection_id = service_connections[c.AZDO_GITHUB_SERVICE_CONNECTION_NAME]
    else:
        service_connection_id = azdo_create_service_connection(project_id)

    ### Create a build definition (if one does not already exist)
    build_definitions = azdo_list_resource_names('https://{0}/{1}/{2}/_apis/build/definitions?{3}'.format(
        c.AZDO_HOST,
        c.AZDO_ORGANIZATION,
        c.AZDO_NEW_PROJECT_NAME,
        c.AZDO_API_VERSION))
    if c.AZDO_NEW_PIPELINE_NAME in build_definitions.keys():
        print('Build definition with name {} already exists. Skipping creation'.format(c.AZDO_NEW_PIPELINE_NAME))
    else:
        azdo_create_build_definition(service_connection_id)


if __name__ == '__main__':
    azdo_example()
