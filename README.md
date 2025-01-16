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
   
TaskFlowGo/
├── manager/                # Manager service
│   ├── main.go             # Entry point for Manager
│   ├── api/                # API handlers for Manager
│   │   ├── task.go         # Task-related APIs
│   │   └── runner.go       # Runner registration and management APIs
│   ├── db/                 # Database interactions
│   │   └── mongo.go        # MongoDB helper functions
│   └── utils/              # Utilities (shared by Manager)
│       └── logger.go       # Logging utility
├── runner/                 # Runner service
│   ├── main.go             # Entry point for Runner
│   ├── tasks/              # Task execution logic
│   │   ├── shell.go        # Shell task execution
│   │   └── config.go       # Config update execution
│   └── utils/              # Utilities (shared by Runner)
│       └── heartbeat.go    # Heartbeat functionality
├── shared/                 # Shared libraries for Manager and Runner
│   ├── models/             # Shared models (e.g., Task, Runner)
│   │   └── task.go         # Task data model
│   ├── constants/          # Shared constants
│   │   └── types.go        # Task types, states, etc.
│   └── utils/              # Shared utilities
│       └── common.go       # Common helper functions
├── scripts/                # Deployment and setup scripts
│   └── start.sh            # Script to run both services
├── .env                    # Environment variables
├── .gitignore              # Git ignore file
├── go.mod                  # Go modules file
├── go.sum                  # Go modules dependencies
└── README.md               # Project documentation