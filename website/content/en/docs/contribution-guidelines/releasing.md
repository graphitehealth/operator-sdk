---
title: Release Guide
linkTitle: Releasing
weight: 4
---

These steps describe how to conduct a release of the operator-sdk repo using example versions.
Replace these versions with the current and new version you are releasing, respectively.

## Table of Contents:

- [Prerequisites](#prerequisites)
- [Major and minor releases](#major-and-minor-releases)
- [Patch releases](#patch-releases)
- [`scorecard-test-kuttl` image releases](#scorecard-test-kuttl-image-releases)
- [Release tips](#helpful-tips-and-information)

## Prerequisites

The following tools and permissions are needed to conduct a release of the operator-sdk repo.

### Tools

- [`git`](https://git-scm.com/downloads): version 2.2+
- [`make`](https://www.gnu.org/software/make/): version 4.2+
- [`sed`](https://www.gnu.org/software/sed/): version 4.3+

### Permissions

- Must be a [Netlify admin][doc-owners]
- Must be an admin on the [operator-sdk repo](https://github.com/graphitehealth/operator-sdk/settings/access)

### Setting Up Tools for MacOS Users

To install the prerequisite tools on MacOS, complete the following steps:

1. Install GNU `sed` and `make`, which may not be installed by default: 

   - ```sh
      brew install gnu-sed make
     ```

1. Verify that the version of `make` is higher than 4.2 using command `make --version`.

1. Add the gnubin directory for `make` to your PATH from your `~/.bashrc`
to allow you to use `gmake` as `make`:

   - ```sh
     echo 'export PATH="/usr/local/opt/make/libexec/gnubin:$PATH"' >> ~/.bashrc
     ```

1. Verify that the version of `sed` is higher than 4.3 using command `gnu-sed --version`.

1. Add the gnubin directory for `gnu-sed` to your PATH from your `~/.bashrc`
to allow you to use `gnu-sed` as `sed`:

   - ```sh
     echo 'export PATH="/usr/local/opt/gnu-sed/libexec/gnubin:$PATH"' >> ~/.bashrc
     ```

## Major and Minor Releases

We will use the `v1.3.0` release version in this example.

**Be sure to substitute
the version you are releasing into the provided commands.**

To perform a major or minor release, you must perform the following actions:

- Ensure a new Netlify branch is created
- Create a release branch and lock down the master branch
- Create and merge a PR for the release branch
- Unlock the master branch and push a release tag to it
- Perform some clean up actions and announce the new release to the community

### Procedure

1. **Before creating a new release branch**, it is imperative to perform the following initial setup steps:
   1. In the [Branches and deploy contexts](https://app.netlify.com/sites/operator-sdk/settings/deploys#branches)
   pane in Netlify, click into the Additional branches list section and add `v1.13.x`.
      - This will watch the branch when there are changes on Github (creating the branch, or adding a commit).
      - NOTE: You must be a [Netlify admin][doc-owners] in order to edit the branches list.
1. Create a release branch by running the following, assuming the upstream SDK repo is the `upstream` remote on your machine:

   - ```sh
     git checkout master
     git fetch upstream master
     git pull master
     git checkout -b v1.3.x
     git push upstream v1.3.x
     ```

1. Make sure that the list of supported OLM versions is up to date:
   1. Identify if a new version of OLM needs to be officially supported by ensuring that the latest three releases listed on the [OLM release page](https://github.com/operator-framework/operator-lifecycle-manager/releases) are all listed as supported in the [Overview][overview] section of the SDK docs.
   1. If a new version of OLM needs to be added and an old version removed, follow the steps in the [updating OLM bindata](#updating-olm-bindata) section before moving onto the next step.

1. Lock down the `master` branch to prevent further commits before the release completes:
   1. Go to `Settings -> Branches` in the SDK repo.
   1. Under `Branch protection rules`, click `Edit` on the `master` branch rule.
   1. In section `Protect matching branches` of the `Rule settings` box, increase the number of required approving reviewers to 6.
   1. Scroll down to save your changes to protect the `master` branch.

1. Create and push a release commit
   1. Create a new branch to push the release commit:

      - ```sh
        export RELEASE_VERSION=v1.3.0
        git checkout master
        git pull master
        git checkout -b release-$RELEASE_VERSION
        ```

   1. Update the top-level [Makefile] variable `IMAGE_VERSION`
to the upcoming release tag `v1.3.0`. This variable ensures sample projects have been tagged
correctly prior to the release commit.

      - ```sh
        sed -i -E 's/(IMAGE_VERSION = ).+/\1v1\.3\.0/g' Makefile
        ```

        If this command fails on MacOS with a warning "sed is not found", follow the step 5 in the [Setting Up Tools for MacOS Users](#setting-up-tools-for-macos-users) section to map `gsed` to `sed`. 
   1. Run the pre-release `make` target:

      - ```sh
        make prerelease
        ```

      The following changes should be present:
      - `Makefile`: IMAGE_VERSION should be modified to the upcoming release tag. (This variable ensures sampleprojects have been tagged correctly prior to the release commit.)
      - `changelog/generated/v1.3.0.md`: commit changes (created by changelog generation).
      - `changelog/fragments/*`: commit deleted fragment files (deleted by changelog generation).
      - `website/content/en/docs/upgrading-sdk-version/v1.3.0.md`: commit changes (created by changelog generation).
      - `website/config.toml`: commit changes (modified by release script).
      - `testdata/*`: Generated sample code.
   1. Commit these changes and push to your remote (assuming your remote is named `origin`):

      - ```sh
        git add Makefile changelog website testdata
        git commit -sm "Release $RELEASE_VERSION"
        git push origin release-$RELEASE_VERSION
        ```

1. Create and merge a new PR for the release-v1.3.0 branch created in step 5.4.
   - You can force-merge your PR to the locked-down `master`
if you have admin access to the operator-sdk repo, or ask an administrator to do so.
   - Note that the docs PR check will fail because the site isn't published yet; the PR can be merged anyways.

1. Unlock the `master` branch
   1. Go to `Settings -> Branches` in the SDK repo.
   1. Under `Branch protection rules`, click `Edit` on the `master` branch rule.
   1. In section `Protect matching branches` of the `Rule settings` box, reduce the number of required approving reviewers back to 1.

1. Create and push a release tag on `master`
   1. Refresh your local `master` branch, tag the release PR commit, and push to the main operator-sdk repo (assumes the remote's name is `upstream`):

      - ```sh
        git checkout master
        git pull master
        make tag
        git push upstream refs/tags/$RELEASE_VERSION
        ```

1. Fast-forward the `latest` and release branches
   1. The `latest` branch points to the latest release tag to keep the main website subdomain up-to-date.
   Run the following commands to do so:

      - ```sh
        git checkout latest
        git reset --hard refs/tags/$RELEASE_VERSION
        git push -f upstream latest
        ```

   1. Similarly, to update the release branch, run:

      - ```sh
        git checkout v1.3.x
        git reset --hard refs/tags/$RELEASE_VERSION
        git push -f upstream v1.3.x
        ```

1. Post release steps
   1. Publish the new Netlify subdomain for version-specific docs. 
      1. Assuming that the Netlify prestep was done before the new branch was created, a new [branch option](https://app.netlify.com/sites/operator-sdk/settings/domain#branch-subdomains)
      should be visible to Netlify Admins under Domain management > Branch subdomains and can be mapped to a subdomain. (Note: you may have to scroll down to the bottom of the Branch subdomains section to find the branch that is ready to be mapped.)
      1. Please test that this subdomain works by going to the link in a browser. You can use the link in the second column to jump to the docs page for this release.
   1. Make an [operator-framework Google Group][of-ggroup] post.
      - You can use [this post](https://groups.google.com/g/operator-framework/c/2fBHHLQOKs8/m/VAd_zd_IAwAJ) as an example.
   1. Post to Kubernetes slack in #kubernetes-operators and #operator-sdk-dev.
      - You can use [this post](https://kubernetes.slack.com/archives/C017UU45SHL/p1679082546359389) as an example. 
   1. Clean up the GitHub milestone
      1. In the [GitHub milestone][gh-milestones], bump any open issues to the following release.
      1. Close out the milestone.
   1. Update the newly unsupported branch documentation (1.1.x in this example)to mark it as archived. (Note that this step does not need to be merged before the release is complete.)
      1. Checkout the newly unsupported release branch:

         - ```sh
           git checkout v1.1.x
           ```

      1. Modify the `website/config.toml` file on lines 88-90 to be the following:

        - ```toml
          version = "v1.1"
          archived_version = true
          url_latest_version = "https://sdk.operatorframework.io"
          ```

## Patch releases

We will use the `v1.3.1` release version in this example.

### 0. Lock down release branches on GitHub

1. Lock down the `v1.3.x` branch to prevent further commits before the release completes:
   1. Go to `Settings -> Branches` in the SDK repo.
   1. Under `Branch protection rules`, click `Edit` on the `v*.` branch rule.
   1. In section `Protect matching branches` of the `Rule settings` box, increase the number of required approving reviewers to `6`.

### 1. Branch

Create a new branch from the release branch (v1.3.x in this example). This branch should already exist prior to cutting a patch release.

```sh
export RELEASE_VERSION=v1.3.1
git checkout v1.3.x
git pull
git checkout -b release-$RELEASE_VERSION
```

### 2. Prepare the release commit

Using the version for your release as the IMAGE_VERSION, execute the
following commands from the root of the project.

```sh
# Update the IMAGE_VERSION in the Makefile
sed -i -E 's/(IMAGE_VERSION = ).+/\1v1\.3\.1/g' Makefile
#  Run the pre-release `make` target:
make prerelease
```

All of the following changes should be present (and no others).

- Makefile: IMAGE_VERSION should be modified to the upcoming release tag. (This variable ensures sampleprojects have been tagged correctpy priror to the release commit.)
- changelog/: all fragments should be deleted and consolidated into the new file `changelog/generated/v1.3.1.md`
- docs: If there are migration steps, a new migration doc will be created. The installation docs should also contain a link update.
- testdata/: Generated samples and tests should have version bumps

Commit these changes and push these changes **to your fork**:

```sh
git add Makefile changelog website testdata
git commit -sm "Release $RELEASE_VERSION"
git push -u origin release-$RELEASE_VERSION
```

### 3. Create and merge Pull Request

- Create a pull request against the `v1.3.x` branch.
- Once approving review is given, merge. You may have to unlock the branch by setting
"required approving reviewers" to back to `1`. (See step 0).

### 4. Create a release tag

Pull down `v1.3.x` and tag it.

```sh
git checkout v1.3.x
git pull upstream v1.3.x
make tag
git push upstream refs/tags/$RELEASE_VERSION
```

### 5. Fast-forward the `latest` branch

If the patch release is on the latest y-stream (in the example you would
not ff latest if there was a y-stream for v1.4.x), you will need to
fast-forward the `latest` git branch.

(The `latest` branch points to the latest release tag to keep the main website subdomain up-to-date.)

```sh
git checkout latest
git reset --hard tags/$RELEASE_VERSION
git push -f upstream latest
```

### 6. Post release steps

- Make an [operator-framework Google Group][of-ggroup] post.
- Post to Kubernetes slack in #kubernetes-operators and #operator-sdk-dev.
- In the [GitHub milestone][gh-milestones], bump any open issues to the following release.

**Note**
In case there are non-transient errors while building the release job, you must:

1. Revert the release PR. To do so, create a PR which reverts the patch release PR created in step [3](#3-create-and-merge-pull-request).
2. Fix what broke in the release branch.
3. Re-run the release with an incremented minor version to avoid Go module errors (ex. if v1.3.1 broke, then re-run the release as v1.3.2). Patch versions are cheap so this is not a big deal.

## `scorecard-test-kuttl` image releases

The `quay.io/operator-framework/scorecard-test-kuttl` image is released separately from other images because it
contains the [`kudobuilder/kuttl`](https://hub.docker.com/r/kudobuilder/kuttl/tags) image, which is subject to breaking changes.

Release tags of this image are of the form: `scorecard-kuttl/vX.Y.Z`, where `X.Y.Z` is _not_ the current operator-sdk version.
For the latest version, query the [operator-sdk repo tags](https://github.com/graphitehealth/operator-sdk/tags) for `scorecard-kuttl/v`.

The only step required is to create and push a tag.
This example uses version `v2.0.0`, the first independent release version of this image:

```sh
export RELEASE_VERSION=scorecard-kuttl/v2.0.0
make tag
git push upstream refs/tags/$RELEASE_VERSION
```

The [`deploy/image-scorecard-test-kuttl`](https://github.com/graphitehealth/operator-sdk/actions/workflows/deploy.yml)
Action workflow will build and push this image.

## Helpful Tips and Information

### Binaries and Signatures

Binaries will be signed using our CI system's GPG key. Both binary and signature will be uploaded to the release.

### Release Branches

Each minor release has a corresponding release branch of the form `vX.Y.x`, where `X` and `Y` are the major and minor
release version numbers and the `x` is literal. This branch accepts bug fixes according to our [backport policy][backports].

### Cherry-picking

Once a minor release is complete, bug fixes can be merged into the release branch for the next patch release.
Fixes can be added automatically by posting a `/cherry-pick v1.3.x` comment in the `master` PR, or manually by running:

```sh
git checkout v1.3.x
git checkout -b cherrypick/some-bug
git cherry-pick <commit>
git push upstream cherrypick/some-bug
```

Create and merge a PR from your branch to `v1.3.x`.

### GitHub Release Information

GitHub releases live under the [`Releases` tab][release-page] in the operator-sdk repo.

### Updating OLM Bindata

Prior to an Operator SDK release, add bindata (if required) for a new OLM version by following these steps:

1. Add the new version to the [`OLM_VERSIONS`][olm_version] variable in the Makefile.
2. Remove the _lowest_ version from that variable, as `operator-sdk` only supports 3 versions at a time.
3. Run `make bindata`.
4. Check that all files were correctly updated by running this script from the root directory of the repository:

   - ```sh
     ./hack/check-olm.sh
     ```

     If the check shows that files were missed by the make target, manually edit them to add the new version and remove the obsolete version.
5. Check that the list of supported OLM versions stated in the [`Overview`][overview] section of SDK documentation is updated.
6. Add the changed files to ensure that they will be committed as part of the release commit:

   - ```sh
     git add -u
     ```

### Patch Releases in Parallel

The following should be considered when doing parallel patch releases:

- Releasing in order is nice but not worth the inconvenience. Release order affects the order on GitHub releases, and which
    is labeled "latest release".
- Do not unlock v.* branches while other releases are in progress. Instead, have an admin do the merges.
- Release announcements should be consolidated.

[doc-owners]: https://github.com/graphitehealth/operator-sdk/blob/master/OWNERS
[release-page]:https://github.com/graphitehealth/operator-sdk/releases
[backports]:/docs/upgrading-sdk-version/backport-policy
[of-ggroup]:https://groups.google.com/g/operator-framework
[gh-milestones]:https://github.com/graphitehealth/operator-sdk/milestones
[Makefile]:https://github.com/graphitehealth/operator-sdk/blob/master/Makefile
[olm_version]:https://github.com/graphitehealth/operator-sdk/blob/6002c70fe770cdaba9ba99da72685e0e7b6b69e8/Makefile#L45
[overview]: https://github.com/graphitehealth/operator-sdk/blob/master/website/content/en/docs/overview/_index.md#olm-version-compatibility
