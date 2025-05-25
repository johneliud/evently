import { useState, useEffect } from 'react';
import Notification from './Notification';

export default function GoogleCalendarButton({ eventId }) {
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [notification, setNotification] = useState(null);

  useEffect(() => {
    checkCalendarConnection();
  }, []);

  async function checkCalendarConnection() {
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      const response = await fetch(
        'http://localhost:9000/api/calendar/check-connection',
        {
          method: 'GET',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Failed to check calendar connection');
      }

      const data = await response.json();
      setIsConnected(data.connected);
    } catch (error) {
      console.error('Error checking calendar connection:', error);
    }
  }

  async function handleConnectCalendar() {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setNotification({
          type: 'error',
          message: 'You must be logged in to connect Google Calendar',
        });
        return;
      }

      setIsLoading(true);

      const response = await fetch(
        'http://localhost:9000/api/calendar/authorize',
        {
          method: 'GET',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Failed to get authorization URL');
      }

      const data = await response.json();

      // Open the authorization URL in a new window
      window.open(data.auth_url, '_blank');
    } catch (error) {
      setNotification({
        type: 'error',
        message:
          error.message ||
          'An error occurred while connecting to Google Calendar',
      });
    } finally {
      setIsLoading(false);
    }
  }

  async function handleAddToCalendar() {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setNotification({
          type: 'error',
          message: 'You must be logged in to add events to Google Calendar',
        });
        return;
      }

      setIsLoading(true);

      // Convert eventId to a number if it's a string
      const eventIdNumber = parseInt(eventId, 10);

      if (isNaN(eventIdNumber)) {
        throw new Error('Invalid event ID');
      }

      const response = await fetch(
        'http://localhost:9000/api/calendar/add-event',
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            event_id: eventIdNumber,
          }),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();

        // Check if authorization is required
        if (
          response.status === 401 &&
          errorData.status === 'authorization_required'
        ) {
          // Prompt user to connect their Google Calendar
          setIsConnected(false);
          throw new Error('Please connect your Google Calendar first');
        }

        throw new Error(
          errorData.message || 'Failed to add event to Google Calendar'
        );
      }

      await response.json();

      setNotification({
        type: 'success',
        message: 'Event added to Google Calendar successfully!',
      });
    } catch (error) {
      setNotification({
        type: 'error',
        message:
          error.message ||
          'An error occurred while adding the event to Google Calendar',
      });
    } finally {
      setIsLoading(false);
    }
  }

  // Check for connection status periodically (every 5 seconds) to help detect when the user completes the OAuth flow in another tab
  useEffect(() => {
    const intervalId = setInterval(() => {
      if (!isConnected) {
        checkCalendarConnection();
      }
    }, 5000);

    return () => clearInterval(intervalId);
  }, [isConnected]);

  return (
    <div className="mt-4">
      {notification && (
        <Notification
          type={notification.type}
          message={notification.message}
          onClose={() => setNotification(null)}
        />
      )}

      {isConnected ? (
        <button
          onClick={handleAddToCalendar}
          disabled={isLoading}
          className="flex items-center justify-center w-full px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {isLoading ? (
            <span className="flex items-center">
              <svg
                className="w-5 h-5 mr-2 animate-spin"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              Adding to Calendar...
            </span>
          ) : (
            <span className="flex items-center">
              <svg
                className="w-5 h-5 mr-2"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
              Add to Google Calendar
            </span>
          )}
        </button>
      ) : (
        <button
          onClick={handleConnectCalendar}
          disabled={isLoading}
          className="flex items-center justify-center w-full px-4 py-2 text-sm font-medium text-white bg-red-600 border border-transparent rounded-md shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
        >
          {isLoading ? (
            <span className="flex items-center">
              <svg
                className="w-5 h-5 mr-2 animate-spin"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              Connecting...
            </span>
          ) : (
            <span className="flex items-center">
              <svg
                className="w-5 h-5 mr-2"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
              Connect Google Calendar
            </span>
          )}
        </button>
      )}
    </div>
  );
}
