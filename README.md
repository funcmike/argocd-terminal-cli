# ArgoCD Terminal Cli ("atc")
ArgoCD Terminal Cli (short: `atc`) purpose is to allow using of ArgoCD Terminal websocket's connection in normal shell without Web Browser.
Additionally, `atc` is providing ability to list resources from ArgoCD to find out information required to start terminal session. 

## Project stage
Currently, project is in working/usable  alpha stage.
All operations are tested and works correctly.
I'm using it daily without problems.

## Install / Build
```shell
go install ./...
```
## Auth
Standard ArgoCD config file (with current-context) is used for auth (how to perform auth - see `argocd login -h`).
There is a possibility to override this by providing command line arguments (-h for help) or environment variables described below.

## ENV
* ARGOCD_CONFIG_FILE - Path to ArgoCD config file
* ARGOCD_SERVER  -  ArgoCD server host
* ARGOCD_AUTH_TOKEN - ArgoCD auth token

## Using - common commands.
Print help
```shell
atc -h 
```

Start terminal session to pod using 
```shell
atc term --app-name test-app --app-namespace test-argocd --project-name argocd-project --namespace staging --pod test-pod --container test-container
```
 `--project-name` is by default `--app-name`

Get all resources (`--output yaml` or `json`)
```shell
atc get all --app-name test-app --app-namespace staging
```

Get specific resource manifest for ex. POD (`--output yaml` or `json`)
```shell
atc get pod test-pod --app-name test-app --app-namespace test-argocd --namespace staging
```
