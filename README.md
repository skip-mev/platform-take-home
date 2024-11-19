# (INTERNAL - DO NOT USE) Skip Platform Take Home Challenge

## Overview

The goal of this exercise will be to demonstrate your DevOps / automation skills. In this repo you will find
an application that is a simple REST API that returns a list of items. The application is written in Go and 
uses a Postgres database.

## Exercise

Hypothetical developers have been working on this application for a while now and have noticed that the experience 
is not ideal.

Their current process is building the binary locally and then deploying it to a server. They have to manually
configure the database and the application. They have to manually start the application and make sure it is running.

Currently, their main complaints are:
* the deployment process is manual and error-prone. Whenever a deployment happens, we suffer a tiny bit of downtime due to the server being down.
* there's no standardization for code formatting which leads to inconsistencies
* the proto-gen script is ran locally which leads to developers forgetting to do it before pushing code upstream
* downloading the tooling dependencies is a manual, undocumented process
* running the server locally is a bit of a pain, since it requires manually standing up a Postgres database (to replicate prod)
* there's no easy way to share a feature update with our colleagues, since we only have a local and production environment.

Your task is to actionize these complaints and implement features that will make the developers' lives easier.
You should target about 4-6 hours of work for this exercise. You have the freedom to implement any features you think
will make the developers' lives easier, as long as you can justify them. Additionally, you do not have to alleviate all
of the problems above, as long as you reason why you prioritized certain ones versus others.

Lastly, below you will see the tech stack that Skip uses, but you are not constrained to using it. You can use any
technology you think is best for the job (brownie points for pragmatic creativity).

## Skip's tech stack

At Skip, we use the following technologies:
* Go for backend and on-chain services
* Typescript for frontend
* Postgres for databases
* Docker + Kubernetes for container orchestration
* ArgoCD and Github Workflows for continuous, self-service app deployment

## Deliverables

With the link to the take-home, you also should've received a link for submission. 
Please invite @Zygimantass and @bpiv400 to the Github repository and submit the repository in the given link.
In the README, please provide instructions on how to run your solution (whether locally or in cloud).

A short explainer on your solution is also appreciated, but not required.


