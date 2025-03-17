# Agent OS: Secure Environment for Autonomous Agents (WIP)

## Overview

**Agent OS** is an open-source framework designed to provide a secure, isolated environment for running autonomous agents. By leveraging containerization technology, Agent OS ensures that each agent operates within its own sandbox, minimizing security risks and preserving system integrity. This tool is ideal for developers and researchers building scalable, safe, and efficient agent-based systems. Agent OS currently utilizes the [Browserless](https://github.com/browserless/browserless) image to power its browser instances.

### Key Features

- **Secure Execution**  
  Agents operate in isolated containers, preventing interference and bolstering security.

- **Autonomous Task Management**  
  Supports a wide range of automated tasks, including web browsing, file operations, and system interactions.

- **Seamless Playwright Integration**  
  Easily integrates with Playwright for advanced automation and testing capabilities.

- **Scalable Infrastructure**  
  Efficiently spawns, manages, and monitors multiple agents concurrently.

- **Command Execution**  
  Execute shell commands directly within a container for enhanced control and debugging.

## Getting Started

> **Note:** Agent OS is currently in development. The features below reflect its initial capabilities, with more to come soon!

Currently, Agent OS supports creating a browser instance for AI agents using the [Browserless](https://github.com/browserless/browserless) image. To get started, simply run:

```bash
agentos create --name my-container
```

This command sets up a containerized browser environment ready for agent operations, accessible at `IP-ADDRESS:3000`.

## Web Browser

### With Puppeteer

Puppeteer can connect to the Agent OS browser instance via the `browserWSEndpoint` option. Here's how to adapt your code:

**Before:**

```javascript
const puppeteer = require('puppeteer');
const browser = await puppeteer.launch();
```

**After:**

```javascript
const puppeteer = require('puppeteer');
const browser = await puppeteer.connect({
  browserWSEndpoint: 'ws://IP-ADDRESS:3000',
});
```

### With Playwright

Agent OS supports Playwright's remote connection protocols. Here's an example using Firefox:

**Before:**

```javascript
const pw = require('playwright');
const browser = await pw.firefox.launch();
```

**After:**

```javascript
const pw = require('playwright-core');
const browser = await pw.firefox.connect({
  wsEndpoint: 'ws://IP-ADDRESS:3000/firefox/playwright',
});
```

## Execute Command

Run shell commands directly inside a running container using the `exec` command. This feature allows you to interact with the container environment seamlessly.

### Usage

```bash
agentos exec --containerId <container-id> --cmd "<shell-command>"
```

### Example

To execute a simple command like `echo Hello, World` in a container:

```bash
agentos exec --containerId mycontainerId --cmd "echo Hello, World"
```

### Details

- **`--containerId` (required):** The ID of the container where the command will be executed.
- **`--cmd` (required):** The shell command to run inside the container.
- The container must be in a "running" state for the command to execute successfully.
- Output from the command is logged for easy inspection.

This feature is useful for debugging, testing, or performing one-off tasks within the containerized environment.

---
