# Palworld App (Backend)

A Go-based REST API backend for managing Palworld creatures (Pals) and their passive skills. This application provides endpoints for storing, retrieving, and managing Pal data with a web interface.

## Features

- RESTful API for Pal management
- CORS support for web frontend integration
- Environment-based configuration
- JSON data storage
- Passive skills and combinations management

## Data Sources

Data scraped from:
- https://game8.co/games/Palworld
- https://palworld.fandom.com/wiki
- https://palworkd.wiki.gg

## Frontend Integration

This backend works with the Next.js frontend located in `../palworld_web/dumbcode_palworld_web/`. See the [Deployment Guide](../DEPLOYMENT.md) for setup instructions.

## Configuration

The application now supports environment variable configuration. You can configure the app using:

1. **Environment variables** - Set directly in your system
2. **`.env` file** - Create a `.env` file in the project root

### Available Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `GIN_MODE` | `release` | Gin framework mode (debug/release) |
| `ALLOWED_ORIGINS` | `http://localhost:3000,http://localhost:3001` | CORS allowed origins (comma-separated) |
| `ALLOWED_METHODS` | `GET,POST,PUT,DELETE,OPTIONS` | CORS allowed methods (comma-separated) |
| `ALLOWED_HEADERS` | `Origin,Content-Type,Accept,Authorization` | CORS allowed headers (comma-separated) |
| `DATA_DIR` | `./data` | Directory containing data files |
| `PALS_FILE` | `pals.json` | Pals data file name |
| `STORED_PALS_FILE` | `stored_pals.json` | Stored pals data file name |
| `PASSIVE_SKILLS_FILE` | `passive_skills.json` | Passive skills data file name |
| `PASSIVE_SKILL_COMBOS_FILE` | `passive_skill_combos.json` | Passive skill combos data file name |

### Setup

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Modify `.env` with your desired configuration

3. Run the application:
   ```bash
   go run main.go
   ```

### Example .env file

```env
PORT=9000
GIN_MODE=release
ALLOWED_ORIGINS=https://your-frontend-domain.com,http://localhost:3000
```

## CORS Configuration

The application includes CORS support for web frontend integration. Configure allowed origins in the `ALLOWED_ORIGINS` environment variable:

- **Development**: Include `http://localhost:3000` for local frontend
- **Production**: Add your deployed frontend domain(s)
- **Multiple origins**: Separate with commas

Example:
```env
ALLOWED_ORIGINS=https://myapp.vercel.app,https://myapp.netlify.app,http://localhost:3000
```

## API Endpoints

- `GET /store` - Get all stored Pals
- `POST /add-pal` - Add a new Pal
- `GET /options/passive-skills` - Get available passive skills
- `GET /options/pal-species` - Get available Pal species
- `POST /update-data` - Update data from external sources

## Deployment

See the [Deployment Guide](../DEPLOYMENT.md) for detailed deployment instructions including:
- Platform-specific setup
- Environment variable configuration
- CORS troubleshooting
- Frontend-backend connectivity