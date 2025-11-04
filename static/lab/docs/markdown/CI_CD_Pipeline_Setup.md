## **CI/CD Pipeline Setup**

Use the following questions to gather information about the customer’s current CI/CD usage, if any. If during the technical review call(s) the consultant chooses to ask additional questions regarding a specific topic please add a row to the table and record the information.

| CI/CD | Customer Input |
| ----- | ----- |
| Is there an existing CI/CD solution in use? | *Blue Ocean, Jenkins, JenkinsX, Tekton*  |
| If yes, will the existing CI/CD solution be used to build specific pipelines for this project? |  |
| Will CI/CD perform automated deployments into production? |  |
| Will full deployments be performed in dev/test environments? |  |
| Any specific validation requirements for deployments such as automated testing? |  |
| Are there any existing CI/CD patterns that must be followed? |  |
| Are there limitations to use specific tooling in the CI/CD pipeline? |  |
| Is there an existing solution for storing secrets? |  |

### Customer Current CI/CD Platform  Design

Insert Customer CI/CD Platform Design (if applicable)

## **CI/CD Pipeline Security Considerations**

Use the following questions to gather information about the customer’s current CI/CD pipeline security.

| CI/CD | Customer Input |
| ----- | ----- |
| Describe existing CI/CD security process and procedure, if any?  |   |
| Are commits signed with a GPG key? |  |
| Are dependencies kept up to date in an automated manner? |  |
| Does the repo include a SECURITY.md document that outlines security policies and reporting/escalation procedures for an incident?  |  |
| Does the VCS have a mechanism to prevent secret keys, credentials, etc. from being committed?  | GitGuardian, Bridgecrew.io, etc. |
| Describe use of Docker images. |  |
| Describe use of any linters or language specific security test/tooling. |  |

## **CI/CD Pipeline Input Considerations**

Use the following questions to gather information about the customer’s current CI/CD pipeline security.

| CI/CD | Customer Input |
| ----- | ----- |
| Describe existing CI/CD  |   |
| Dependency versioning, scanning, protection against supply chain attacks. |  |

## **CI/CD Pipeline Output Considerations**

Use the following questions to gather information about the customer’s current CI/CD pipeline security.

| CI/CD | Customer Input |
| ----- | ----- |
| Describe existing CI/CD security process and procedure, if any?  |   |
| Container Registries in use? |  |
| Are containers being used as a means of deployment to Production environments, Kubernetes Clusters, etc.?  |  |
| Are code artifacts being stored or processed in intermediate pipeline steps? | Artifactory  |
|  |  |
