## Executive summary

The repository was configured to run SonarCloud analysis via GitHub Actions using Maven. The CI pipeline performs build, test, and Sonar analysis on pushes to `main` and on pull requests. The deliverables required for submission are screenshots of the SonarCloud project dashboard, a successful GitHub Actions workflow run, and at least two security issues (vulnerabilities or hotspots) shown in SonarCloud.


 ## configuration snippets

### sonar-project.properties

```properties
# sonar-project.properties - minimal example for a Maven project
sonar.projectKey=your-org_key:your-project-key
sonar.organization=your-org_key
sonar.host.url=https://sonarcloud.io

# Source locations
sonar.sources=src/main/java
sonar.tests=src/test/java
sonar.java.binaries=target/classes

# Exclusions (optional)
sonar.exclusions=**/target/**,**/generated-sources/**

# Encoding
sonar.sourceEncoding=UTF-8
```

Replace `your-org_key` and `your-project-key` with values from SonarCloud.

### .github/workflows/sonarcloud.yml

```yaml
name: CI - Maven Build and SonarCloud

on:
  push:
    branches: [ main ]
  pull_request:
    types: [opened, synchronize, reopened, closed]

jobs:
  build:
    name: Build, test, and analyze on SonarCloud
    runs-on: ubuntu-latest

    permissions:
      contents: read
      packages: read
      issues: write
      pull-requests: write

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up JDK 17
        uses: actions/setup-java@v4
        with:
          distribution: temurin
          java-version: '17'

      - name: Cache Maven packages
        uses: actions/cache@v4
        with:
          path: ~/.m2/repository
          key: ${{ runner.os }}-maven-${{ hashFiles('**/pom.xml') }}
          restore-keys: |
            ${{ runner.os }}-maven-

      - name: Build and run tests
        run: mvn -B -DskipTests=false verify

      - name: SonarCloud Scan
        env:
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
        run: mvn -B sonar:sonar \
              -Dsonar.host.url=https://sonarcloud.io \
              -Dsonar.organization=your-org_key \
              -Dsonar.login=${{ secrets.SONAR_TOKEN }}
```

Notes: add a repository secret `SONAR_TOKEN` with your SonarCloud token. Replace `your-org_key`.

### pom.xml snippet (optional)

```xml
<!-- in <properties> -->
<properties>
  <!-- ...existing properties... -->
  <sonar.host.url>https://sonarcloud.io</sonar.host.url>
</properties>

<!-- in <build><plugins> -->
<plugin>
  <groupId>org.sonarsource.scanner.maven</groupId>
  <artifactId>sonar-maven-plugin</artifactId>
  <version>3.11.0.2602</version>
</plugin>
```

The plugin is optional; `mvn sonar:sonar` will download it when needed.


## Test proves

1. SonarCloud Project Dashboard
   
   ![alt text ](./Assets/sonarcloud.png)

   ![alt text ](./Assets/sonarcloud1.png)

The sonarcloud dashboard shows the repository has no issue when analysing but the github cicd shows many security issues.

## Finding 1
   
   ![alt text ](./Assets/scyml.png)

   updated the build.yml the issue was sloved.

   ![alt text ](./Assets/vscode.png)

  ![alt text ](./Assets/solved.png)


## Finding 2
   ![alt text ](./Assets/docker.png)

The Snyk container scan previously failed with `SNYK-CLI-0000: Could not detect package manager for file: dockerfile`. I adjusted the container-scan steps in `.github/workflows/enhanced-security.yml` to:

- Print the repository `Dockerfile` (debug) so you can confirm the file name and content seen by the runner.
- Remove `--file` and the `--exclude-base-image-vulns` flag, allowing Snyk to scan the built image for base-image vulnerabilities and application deps when present.

This change prevents Snyk from failing when the Dockerfile contains no application package manager layers (common for JAR-based images). To re-run the container scan and validate:

1. Ensure the workflow builds the image `cicd-demo:latest` (it does via `docker build -t cicd-demo:${{ github.sha }}` then tags to `latest`).
2. Re-run the workflow; Snyk will scan the image and generate `snyk-container.sarif`.
3. If you prefer to scan application dependencies inside the image, change the `args` to include `--file=Dockerfile` only if the Dockerfile uses a supported package manager and leaves metadata in the image layers.

Example of robust args to scan the image:

```yaml
args: >
  --severity-threshold=high
  --sarif-file-output=snyk-container.sarif
```






2. GitHub Actions Workflow Run
   - Go to your GitHub repo → Actions and open the run triggered by your push/PR.
   - Copy the run URL and screenshot the run summary showing successful steps.
   - Filename suggestion: `github_actions_run.png`


Tips:

- Include timestamps and project key visible in screenshots where possible.
- If no vulnerabilities appear, you may create a temporary demonstration change that triggers a known rule (remove it after capturing evidence).

## Quality Gate configuration and testing

- SonarCloud uses a default Quality Gate. To create or modify a gate: SonarCloud → Administration → Quality Gates.
- Attach the gate to the project and test by creating a PR that introduces an issue on new code. The PR should show decoration or a failed status when the gate fails.

## Submission checklist

- [ ] SonarCloud account created and project configured
- [ ] `sonar-project.properties` added (or equivalent config in workflow)
- [ ] `.github/workflows/sonarcloud.yml` added and committed
- [ ] `pom.xml` Sonar plugin/properties updated (if applicable)
- [ ] `SONAR_TOKEN` secret added to GitHub
- [ ] GitHub Action executed successfully (provide run link)
- [ ] SonarCloud analysis completed (screenshot)
- [ ] At least 2 security issues identified and captured (screenshots)
- [ ] Quality gate configured and tested (screenshot)
- [ ] Continuous monitoring (PR decoration/webhook) verified
- [ ] All screenshots and documentation prepared and submitted

## Suggested README snippet for the submission

```text
Evidence:
- sonarcloud_dashboard.png — SonarCloud project dashboard (Quality Gate: PASSED)
- github_actions_run.png — GitHub Actions workflow run (successful) https://github.com/<owner>/<repo>/actions/runs/<id>
- sonar_vulnerability_1.png — Vulnerability: <rule> in src/...
- sonar_hotspot_1.png — Security Hotspot: <rule> in src/...
```

## Next steps and remediation ideas

- Triage and fix critical/major security issues first.
- Add tests for code paths flagged by Sonar and increase coverage.
- Configure branch protection rules to block merges on failing Quality Gates.
- Set up webhooks or integrations for notifications (Slack/Teams/email).

### Note about Quality Gate failures and CI behaviour

During CI runs we observed `QUALITY GATE STATUS: FAILED` which causes `mvn sonar:sonar` to return a non-zero exit code and fail the build. To make failures clearer and to include the failing conditions in CI logs, the workflow was updated to:

- Run Sonar analysis without blocking (set `sonar.qualitygate.wait=false`), then
- Query the SonarCloud API (`/api/qualitygates/project_status`) to fetch the Quality Gate result and failing conditions, and
- Fail the job with a clear message and conditions list if the gate failed.

This change provides clearer error messages in CI logs and makes it easier to include the exact failing metrics and messages in your submission screenshots. The updated workflow is in `.github/workflows/sonarcloud.yml`.

Remediation steps when the Quality Gate fails:

1. Open the SonarCloud dashboard for the project and inspect the Quality Gate widget.
2. In SonarCloud, go to Issues and filter by New Code or by the metric listed in the failed condition (e.g., coverage, vulnerabilities).
3. Fix or triage issues (mark false-positives where appropriate), add tests to improve coverage, and push a new commit to re-run analysis.
4. Re-run the GitHub Action and capture new screenshots showing the gate passing.

## Final notes

If you want, I can add the example `sonar-project.properties` and the workflow to this repository and create a `SUBMISSION.md` commit (this file), or I can also prepare a PR that demonstrates creating/triggering security issues for capture (temporary demonstration only).

---

Prepared for submission on 2025-10-06.
