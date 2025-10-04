# Release Workflow Cheatsheet

This document summarizes the process for creating new releases and fixing failed ones for this project.

## Standard Release Process (One-Step)

Use this method to create a new release after you have committed your code changes.

1.  **Tag the current commit with a new version number:**
    ```bash
    # Example for a new version v0.2.0
    git tag -a v0.2.0 -m "Release v0.2.0"
    ```

2.  **Push the commit and the tag simultaneously:**
    ```bash
    git push origin master --follow-tags
    ```
    This triggers the GitHub Actions workflow that builds the application and creates the GitHub Release.

## How to Fix a Failed Release

If a release fails because of an error in the workflow (`.github/workflows/build.yml`) and you want to re-run it for the **same version number**, follow these steps.

1.  **Fix the workflow file** and commit the change to the `master` branch.
    ```bash
    git add .github/workflows/build.yml
    git commit -m "fix: Correct release workflow"
    git push origin master
    ```

2.  **Delete the failed tag on GitHub:**
    ```bash
    # Example for failed tag v0.1.1
    git push --delete origin v0.1.1
    ```

3.  **Delete the tag on your local machine:**
    ```bash
    git tag -d v0.1.1
    ```

4.  **Re-create the tag on the latest commit** (which now includes your workflow fix):
    ```bash
    git tag -a v0.1.1 -m "Release v0.1.1"
    ```

5.  **Push the new tag to trigger the corrected workflow:**
    ```bash
    git push origin v0.1.1
    ```
