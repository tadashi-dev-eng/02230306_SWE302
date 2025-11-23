# PRACTICAL 4-a: Integrating Sonarcloud With GITHUB ACTIONS 

The report documents the successful implementation of Static Application Security Testing (SAST) using SonarCloud integrated with GitHub Actions for the cicd-demo.
## Setup and configuration

![image](./assets/image%20copy%206.png)


## key Achievements:

- SonarCloud account created and configured
- GitHub Actions workflow automated
- Security analysis completed
- Quality gate enforced
- Continuous monitoring implemented

## GitHub Actions Workflow

Key Features:

- Triggered on push to master/main
- Triggered on pull requests
- JDK 17 setup
- Maven caching enabled
- SonarCloud analysis included
- Quality gate wait enabled

Workflow Triggers:

- Push to master/main branches
- Pull requests (opened, synchronize, reopened)

![image](./assets/Screenshot%20From%202025-11-23%2014-56-33.png)

![image](./assets/Screenshot%20From%202025-11-23%2014-57-03.png)


## SONARCLOUD ANALYSIS RESULTS

The sonarcloud shows all the test passed. 

![image](./assets/Screenshot%20From%202025-11-23%2015-01-52.png)


## Security Hotspots Analysis

Hotspot Summary:

- Total Hotspots: 0
- Reviewed: 100%
- Status: No security hotspots requiring review

There are no Security Hotspots to review. Next time you
analyze a piece of code that contains a potential security risk, it will
show up here.

![image](./assets/image%20copy%202.png)


## Security Analysis Findings



Issue 1: Package Naming Convention Violation (Maintainability)
- Package name does not follow Java naming conventions. Java convention
dictates that package names must be all lowercase and follow the regular
expression pattern.

![image](./assets/image.png)


Issue 2: Unnecessary Public Modifier on Test Class
- Test class has an unnecessary "public" modifier. According to Java testing best practices and conventions, test classes should be package-private (no access modifier specified).

![image](./assets/image%20copy.png)

## QUALITY GATE CONFIGURATION

SonarCloud's custom quality gates are a paid feature. On the free plan:

- Default "Sonar way" gate is available
- Provides comprehensive quality checks
- Includes security, reliability, maintainability metrics
- Sufficient for most projects.
- Enforced uisng Maven flag -Dsonar.qualitygate.wait=true

![image](./assets/image%20copy%203.png)

## Quality Gate Conditions

- Conditions on New Code

![image](./assets/image%20copy%205.png)

- Overall Code Conditions:

![image](./assets/image%20copy%204.png)


