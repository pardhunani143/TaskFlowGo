# TaskFlowGo

TaskFlowGo is a Go-based orchestration and task management system designed for learning and experimentation. It consists of:
- **Manager**: Assigns tasks to Runners and tracks their execution.
- **Runner**: Executes tasks assigned by the Manager.

## Features
- Dynamic task assignment
- Shell script execution
- Config updates for multiple groups and applications
- Heartbeat monitoring

## Folder Structure
- `manager/`: Handles task orchestration and Runner registration.
- `runner/`: Processes tasks from the Manager and executes them.
- `shared/`: Contains shared libraries and utilities.

## Setup Instructions
1. Clone the repository:
   ```bash
   git clone https://github.com/pardhunani143/TaskFlowGo.git
   cd TaskFlowGo
