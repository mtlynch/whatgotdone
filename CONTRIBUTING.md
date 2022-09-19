# How to contribute to What Got Done

## Overview

Thanks for your interest in contributing to What Got Done!

## Expectations

[Michael Lynch](https://mtlynch.io) is the sole maintainer of What Got Done. He does it in his spare time on a best-effort basis.

His target response times are as follows:

* Review a pull request: <= 5 business days
* Triage a bug: <= 3 business days
* Respond to an email about the project: <= 3 business days
* Add a new feature: *Best effort*

## Submitting a good pull request

To help get your pull request merged in quickly, keep these guidelines in mind:

* Keep the pull request **narrowly scoped**.
  * A pull request should do just one thing (e.g., fix a single bug, refactor one file)
* Write a descriptive pull request title and commit message.
  * e.g., "Refactoring indexHandler to a separate file"
* [Rebase your changes](https://www.atlassian.com/git/tutorials/rewriting-history/git-rebase) onto `mtlynch:master`.
  * Ensure that your pull request doesn't have any distracting merge commits.

## Build checks

On any pull request, [Circle CI](https://circleci.com/gh/mtlynch/whatgotdone) automatically runs What Got Done's unit tests. Ensure that your pull request passes these checks.

## I want to add a feature

If this is a small feature (e.g., you want to tweak the UI to look prettier), feel free to go ahead and submit a pull request.

For larger features (e.g., you want to support daily updates or email notifications), [create a new issue](https://github.com/mtlynch/whatgotdone/issues/new) to lay out your plans to ensure you don't invest too much time on something that I can't accept.
