# A good static site generator

Do you like the idea behind [simple-site](https://github.com/wjohnsto/simple-site) but absolutely refuse to build another JavaScript app? Well this static site generator is for you. It is built with:

1. [Golang](https://go.dev/)
1. [templ](https://templ.guide/)
1. Markdown

It provides:

1. Building static pages from markdown files
1. Generating nonsense output to serve to spammers
1. Merging styles into a single `<style>` tag (no external CSS)
1. A two-line Dockerfile for hosting with nginx
1. Scripts for local development
1. Scripts for deployment to fly.io
1. JS if you want it (ugh)

This site is **intended for blog sites**, but it could be used for any static site if you feel so inclined.

## Table of contents

- [A good static site generator](#a-good-static-site-generator)
  - [Table of contents](#table-of-contents)
  - [Dependencies](#dependencies)
  - [Getting started](#getting-started)
    - [Running locally outside docker](#running-locally-outside-docker)
    - [Running locally in docker](#running-locally-in-docker)
    - [Configuring the site](#configuring-the-site)
    - [Creating a page](#creating-a-page)
    - [Creating markdown pages](#creating-markdown-pages)
      - [Writing a markdown post](#writing-a-markdown-post)
    - [Helper components](#helper-components)
  - [Deployment](#deployment)
  - [Nonsense generator](#nonsense-generator)

## Dependencies

- [golang](https://go.dev/doc/install)
  - [templ](https://templ.guide/quick-start/installation)
- [cargo](https://doc.rust-lang.org/cargo/commands/cargo-install.html)
  - [minhtml CLI](https://github.com/wilsonzlin/minify-html)
- [node](https://nodejs.org/)
  - [purgecss CLI](https://purgecss.com/CLI.html)
  - [lightningcss CLI](https://lightningcss.dev/docs.html#from-the-cli)
- [flyctl](https://fly.io/docs/flyctl/install/)
- [docker](https://www.docker.com/get-started/)
- inotifywait (sudo apt install inotify-tools)

## Getting started

The premise of this site builder is you write static pages using [templ](https://templ.guide/) in the `/pkg/ui/pages` directory. They can be static, include CSS, JS, whatever. See the "[creating a page](#creating-a-page)" section for more information. When you build the site, the `templ` files are built into `index.html` pages at the routes you determine. All the styles are bundled together, purged, minified, and included in the `index.html` files.

Run `make` or `make help` to view the available commands in the `Makefile`.

### Running locally outside docker

```bash
make dev
```

Your site will now be served at `https://localhost:8080`. The `posts/` and `pkg/` directories will be watched for changes. When changes happen the site will be rebuilt. Refresh your browser to see the changes.

### Running locally in docker

```bash
make docker
```

Your site will now be served at `https://localhost`. The docker container is set up with a volume that points to your local `/public` directory. This way you can run the following command to develop your site and get automatic changes:

```bash
make watch
```

### Configuring the site

Copy the `config.example.yaml` file to `config.yaml` and edit the variables to your liking. There is a description of what each variable is for. If you prefer json use a `config.json` with the same schema. You can also override variables at runtime with a `.env` file or environment variables. See the `.env.example` for the variable values.

### Creating a page

To create a page create a `.templ` file in the `pkg/ui/pages` directory by copying an existing file or starting from scratch. You can follow the [templ guide](https://templ.guide/) for building out the page, but in order to get it to work with this static site generator it needs to implement the `Page` interface found in `pkg/ui/pages/routes.go`.

Once you have something that implements the `Page` interface, create an `init` function in your `.templ` file that calls the `AddRoute` function found in `pkg/ui/pages/routes.go` to wire the page up to the static generator for the specified route. Below you will find the most basic example of a `.templ` file:

```go
package pages

import (
  "context"
)

type MyPage struct{}

templ (h *MyPage) Render() {
  <!DOCTYPE html><html><body>Hello World!</body></html>
}

func NewMyPage() *MyPage {
  return &MyPage{}
}

func init() {
  AddRoute("/my-page", func() Page {
    return NewMyPage()
  })
}
```

### Creating markdown pages

There are two existing `.templ` files that are set up to show list pages and single pages for markdown files. They are `listMarkdown.templ` and `markdown.templ`. By default, these pages don't publish routes. You can use them to build your own pages for displaying lists/single Markdown files. There is an existing `posts.go` file that creates a blog list and single pages at the `/posts` root URL. The `posts.go` page does a few things:

1. It calls the `markdown` service in `pkg/services` to find all the markdown files under the `/posts/` path within the root markdown directory. It supplies its own method for determining the post slug. By default, it uses `/posts/<year>/<month>/<day>/<slug>`
1. It takes the markdown files and converts them into `MarkdownPost` objects to be used with the `listMarkdown` and `markdown` templ components
1. It adds routes for the `list` and `single` post pages
1. It generates atom, RSS, and JSON feeds that include all the posts

You can use this component to spin up multiple feeds/list/single page combinations for the various content types you want to display.

#### Writing a markdown post

The site is set up to make writing markdown as straightforward as possible while providing some utility if things get screwed up. Don't worry about any meta/front matter to start, just create a `.md` file in the appropriate directory (e.g. `markdown/posts/hello.md`) and write some text. Here's an example:

```md
# Hello world

Test post, please ignore
```

When you build the site, the markdown file will be edited with metadata based on the `h1` tag (i.e. "# Hello World"), the first paragraph, and the site configuration:

```md
---
title: "Hello World"
slug: "hello-world"
author: "Author McAuthorson"
description: "Test post, please ignore."
created_date: "2025-01-07T12:41:04-08:00"
updated_date: "2025-01-27T12:41:04-08:00"
---
# Hello World

Test post, please ignore.
```

You can of course edit this stuff or add it in the beginning to override the defaults.

### Helper components

The `pkg/ui/components` directory contains a few components for helping you build your site. The main one is the `layout.templ` component. You can use it to scaffold all the `doctype`, `html`, `head`, and `body` stuff that you don't want to write for every page. You can see an example of how to call it in any of the template pages.

## Deployment

This site is a static site, so it can be deployed almost anywhere. Included is a `fly.toml` that you can use to deploy to `fly.io`. Use the `make deploy` and `make deploy-remote` commands to handle deployment. The `make deploy` command will build the docker image locally and push it to `registry.fly.io`. Then it will call the `fly` CLI to deploy. This saves you from having to build in the cloud using fly's build tools. You can still use the cloud build with the `make deploy-remote` script.


## Nonsense generator

This site comes bundled with the ability to serve nonsense pages to spammers, bots, etc. (you know what I'm talking about). It's up to you to use this and determine who you want to serve the nonsense to. The nonsense is generated using a Markov-chain trained on your own `public/` HTML files. You can adjust this in `cmd/nonsense/main.go` if you like. By default, 20% of all the text on your site will be converted to nonsense. You can configure this by looking at the `build-nonsense` task in the `Makefile`. The `nginx.conf` has an existing rule for serving the nonsense as a baseline. Edit it for your own needs.
