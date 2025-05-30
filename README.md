# Evently

Evently is a full-stack event management application that allows users to create, manage, and RSVP to events. The application features user authentication, event creation and management, RSVP functionality, and Google Calendar integration.

## Features

- **User Authentication**
  - Email/password registration and login
  - Google OAuth integration
  - JWT-based authentication

- **Event Management**
  - Create, read, update, and delete events
  - View upcoming events
  - Search for events
  - View event details including location, date, and description

- **RSVP System**
  - RSVP to events (Going, Maybe, Not Going)
  - View RSVP counts for events
  - Email notifications for RSVPs

- **Google Calendar Integration**
  - Connect your Google Calendar
  - Add events to your Google Calendar

- **Responsive Design**
  - Mobile-friendly interface
  - Dark mode support

## Tech Stack

### Frontend
- React 19
- Tailwind CSS
- Vite

### Backend
- Go 1.24
- PostgreSQL
- JWT Authentication
- Google OAuth 2.0
- Google Calendar API

## Project Structure

```
evently/
├── backend/
│   ├── cmd/
│   │   └── main.go         # Application entry point
│   ├── controllers/        # HTTP request handlers
│   ├── db/                 # Database connection and migrations
│   ├── models/             # Data models
│   ├── repositories/       # Database operations
│   ├── server/             # Server setup and configuration
│   └── services/           # Business logic
├── frontend/
│   ├── public/             # Static assets
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── App.jsx         # Main application component
│   │   ├── main.jsx        # Application entry point
│   │   └── index.css       # Global styles
│   ├── index.html          # HTML template
│   ├── package.json        # Frontend dependencies
│   ├── tailwind.config.js  # Tailwind CSS configuration
│   └── vite.config.js      # Vite configuration
├── .env                    # Environment variables (not in repo)
├── go.mod                  # Go dependencies
├── go.sum                  # Go dependencies checksums
└── README.md               # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Node.js 18 or higher
- PostgreSQL
- Google Cloud Platform account (for OAuth and Calendar API)

### Environment Setup

1. Create a `.env` file in the root directory with the following variables:

```
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=evently

# JWT
JWT_SECRET=your_jwt_secret

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:9000/api/auth/google/callback

# Google Calendar
GOOGLE_CALENDAR_REDIRECT_URL=http://localhost:9000/api/calendar/callback

# Email Service
EMAIL_FROM=your_email@example.com
EMAIL_PASSWORD=your_email_password
EMAIL_SMTP_HOST=smtp.example.com
EMAIL_SMTP_PORT=587
```

2. Create a `google_client_credentials.json` file for Google Calendar API (download from Google Cloud Console)

### Database Setup

1. Create a PostgreSQL database named `evently`
2. The application will automatically run migrations on startup

### Running the Application

#### Backend

```bash
# From the root directory
go run backend/cmd/main.go
```

The backend server will start on http://localhost:9000

#### Frontend

```bash
# From the frontend directory
cd frontend
npm install
npm run dev
```

The frontend development server will start on http://localhost:5173

## API Endpoints

### Authentication

- `POST /api/signup` - Register a new user
- `POST /api/signin` - Login a user
- `GET /api/auth/google` - Get Google OAuth URL
- `GET /api/auth/google/callback` - Google OAuth callback

### Events

- `GET /api/events/upcoming` - Get upcoming events
- `GET /api/events/user` - Get current user's events
- `POST /api/events` - Create a new event
- `GET /api/events/:id` - Get event by ID
- `PUT /api/events/:id` - Update event
- `DELETE /api/events/:id` - Delete event
- `GET /api/events/search` - Search events

### RSVPs

- `GET /api/events/:id/rsvp` - Get user's RSVP status for an event
- `POST /api/events/:id/rsvp` - Create or update RSVP
- `DELETE /api/events/:id/rsvp` - Delete RSVP
- `GET /api/events/:id/rsvp/count` - Get RSVP counts for an event
- `GET /api/events/:id/rsvps` - Get all RSVPs for an event

### Google Calendar

- `GET /api/calendar/authorize` - Get Google Calendar authorization URL
- `GET /api/calendar/callback` - Google Calendar authorization callback
- `POST /api/calendar/add-event` - Add event to Google Calendar
- `GET /api/calendar/check-connection` - Check if user has connected Google Calendar

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - [see the LICENSE file for details](https://github.com/johneliud/evently/blob/main/LICENSE).

## Acknowledgements

- [React](https://reactjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Google Calendar API](https://developers.google.com/calendar)
