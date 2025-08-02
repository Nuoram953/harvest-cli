# Harvest CLI Tool

A simple command-line interface to manage your [Harvest](https://www.getharvest.com/) time tracking entries, tasks, and projects efficiently.

## Features

This CLI is built in phases:

### Phase 1: Time Entry Management

- Create, edit, view, and delete time entries.

### Phase 2: Task Management *(Planned)*

- Create, edit, and delete tasks.

### Phase 3: Project Management *(Planned)*

- Create, edit, and delete projects.

---

## 🚀 Usage

All commands follow the pattern:

```bash
harvest [entity] [action] [options]
```

## Entry Commands

```bash
harvest entry create
```

Create a new time entry.

### Options

- `-p, --project <project>`: Specify the project ID.
- `-t, --task <task>`: Specify the task ID.
- `-d, --date <date>`: Specify the date (default: today).
- `-h, --hours <hours>`: Specify the number of hours.

---

## Global Options

```bash
- `-n, --noconfirm`: Skip confirmation prompts.
```
