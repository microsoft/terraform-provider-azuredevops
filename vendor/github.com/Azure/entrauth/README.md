# Intro

`entrauth` provides a customizable chained token credential for authenticating to Microsoft Entra ID. Based on this, it contains a sub package `aztfauth` that provides an opinionated chained token credential, which is meant to be used for Azure Terraform providers.

# Credentials

The high level structure of the basic supported credentials are listed below:

```
Auth
  |
  +--> OAuth2 Client Credential
  |      |
  |      +------ client secret ----------------------> "client-secret"
  |      |
  |      +------ client assertion
  |                    |
  |                    +----- plain assertion -------> "assertion-plain" 
  |                    |
  |                    +----- assertion file --------> "assertion-file" 
  |                    |
  |                    +----- client certificate ----> "client-certificate" 
  |                    |       (build assertion)
  |                    |
  |                    +------ request --------------> "assertion-request"
  |                                                    (Github, AzureDevOps)
  +--> Token Provider
         |
         +------ Azure Managed Identity -------------> "managed-identity"
         |
         +------ Azure CLI delegation ---------------> "azure-cli"
         |
         +------ Azure Developer CLI delegation -----> "azure-dev-cli"
```

Based on above, the `aztfauth` provides the following chained token credential:

```
            "assertion-plain"
                   |
                   v
            "assertion-file"
                   |
                   v
           "assertion-request"
                   |
                   v
    ADOServiceConnectionId == "" ?
                  / \
               y /   \ n
                /     \
           Github    AzureDevOps
                \     /
                 \   /
                  \ /
                   v
            "client-secret"
                   |
                   v
          "client-certificate"
                   |
                   v
           "managed-identity"
                   |
                   v
              "Azure CLI"
                   |
                   v
            "Azure Dev CLI"
```

Note that each token credential can be enabled/disabled.
