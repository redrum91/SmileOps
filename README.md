# SmileOps

Desktop application for managing dental implants and patient procedures.

## Tech Stack

- **Backend**: Go (Wails v2)
- **Frontend**: Vue 3 + TypeScript
- **UI**: shadcn-vue (Tailwind CSS)
- **Data**: Excel file storage

## Features

- Patient management with implant tracking
- Multiple operation types (Removal, Sinus Lift, Bone Grafting, etc.)
- Date and tooth number tracking per operation
- Patient list view with comprehensive operation history
- Excel data export

## Development

Run the app in development mode with hot reload:

```bash
wails dev
```

Access the dev server at http://localhost:34115

## Building

Build a production package:

```bash
wails build
```

## Configuration

Edit `wails.json` to configure the project.

More info: https://wails.io/docs/reference/project-config
