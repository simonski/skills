# Required Secrets

This document lists the secrets required to publish releases of `skills`.

Normal development and testing require **no secrets**. Secrets are only needed when publishing a release to GitHub and the Homebrew tap.

## TAP_TOKEN

| Field | Value |
|-------|-------|
| **Purpose** | Allows CI to push the updated Homebrew formula to the `simonski/homebrew-tap` repository |
| **Type** | GitHub Personal Access Token (classic) or fine-grained PAT |
| **Required scope** | `Contents: write` on the `simonski/homebrew-tap` repository |
| **Used by** | `.github/workflows/publish.yml` — Update Homebrew tap step |

### How to create and configure

1. Go to [github.com/settings/tokens](https://github.com/settings/tokens)
2. Generate a new token (classic) with `repo` scope **or** a fine-grained token scoped to `simonski/homebrew-tap` with **Contents: Read and Write**
3. Copy the token value
4. In the `simonski/skills` repository, go to **Settings → Secrets and variables → Actions**
5. Click **New repository secret**
6. Name: `TAP_TOKEN`, Value: *(paste token)*
7. Click **Add secret**

### Verification

After adding the secret, the next push to `main` will trigger the publish workflow. Check the **Actions** tab to confirm the "Update Homebrew tap" step completes successfully.
