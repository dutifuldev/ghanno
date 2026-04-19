# Repository Access Grants

`prtags` should support a small local exception system for repository-scoped write access.

The goal is simple.

A user should be able to act as a writer for a repository inside `prtags` even if GitHub itself does not show that user as a repo member.

This is product authorization, not GitHub authorization.

It should not change GitHub org membership, GitHub team membership, or GitHub repository permissions.

## Why This Exists

Today, `prtags` allows writes only when GitHub says the caller has write-level permission on the repository.

That is a good default, but it is too strict for cases where:

- the repository owner wants to delegate `prtags` curation rights without adding a user to the GitHub org
- the repository is public, but `prtags` write actions should still be limited to a small trusted set of people
- `prtags` needs a product-level access model that is separate from GitHub membership

The clean fix is to store local repository access grants in the `prtags` database and check them after the normal GitHub permission check.

## Design Rule

`prtags` should continue to use GitHub for identity.

`prtags` should add its own small repository-scoped grant layer for authorization.

The write decision should be:

1. allow if GitHub says the caller has `admin`, `maintain`, or `push`
2. otherwise allow if `prtags` has a local repository access grant for that user
3. otherwise deny

That keeps the default behavior aligned with GitHub while still supporting clean local exceptions.

## Table Name

The table should be named `repository_access_grants`.

This name is better than `repo_write_grants` because it does not lock the model to one action forever.

It clearly says:

- the grant is repository-scoped
- the row is an access grant
- the system can grow later without renaming the table

## Schema

The first version should stay small.

Suggested columns:

- `id`
- `github_repository_id`
- `github_user_id`
- `github_login`
- `role`
- `granted_by_github_user_id`
- `granted_by_github_login`
- `created_at`
- `updated_at`

## Column Notes

`github_repository_id`

- this is the stable repository identity
- it should be the source of truth, not `owner/name`

`github_user_id`

- this is the stable user identity
- it should be the source of truth, not login

`github_login`

- this is a cached display value
- it is useful for admin views, logs, and debugging
- it should not be the primary identity key

`role`

- use a role string instead of a `can_write` boolean
- the first production value can be `writer`
- this keeps the schema stable if `admin` or `reader` are needed later

`granted_by_github_user_id`
and
`granted_by_github_login`

- these provide a basic audit trail for who created the exception

## Constraints

The important uniqueness rule should be:

- one row per `github_repository_id` plus `github_user_id`

That means one user has at most one effective local role for one repository.

Recommended indexes:

- unique index on `github_repository_id, github_user_id`
- index on `github_user_id`
- index on `github_repository_id`

## Authorization Flow

The service-level write check should become:

1. authenticate the caller with GitHub
2. ask GitHub whether the caller has write-level repository permission
3. if yes, allow
4. if no, look for a row in `repository_access_grants`
5. if a matching row exists with role `writer`, allow
6. otherwise deny

This keeps the exception path explicit and easy to reason about.

## Scope

The first version should apply to all repository-scoped write actions in `prtags`, including:

- creating groups
- updating groups
- adding group members
- removing group members
- setting annotations
- managing field definitions

That is better than creating one special rule only for groups.

## What Not To Do

Do not hardcode repository names in the backend.

Do not add a one-off check like “if repo is `openclaw/openclaw` and login is `dutifulbob`, allow.”

Do not treat `owner/name` as the durable identity key for the grant.

Do not make this change GitHub membership or imply that the user is a GitHub org member.

## Summary

The clean production-ready design is:

- keep GitHub as the identity provider
- keep GitHub write access as the default allow path
- add a local repository-scoped grant table for exceptions
- store those exceptions in `repository_access_grants`
- key them by stable GitHub repository ID and user ID

That is the simplest non-hacky way to let a user act like a repo writer inside `prtags` without adding that user to the GitHub org.
