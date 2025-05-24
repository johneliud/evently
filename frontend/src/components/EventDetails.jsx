import { useState, useEffect } from 'react';
import Notification from './Notification';
import EditEventForm from './EditEventForm';

export default function EventDetails() {
  // Get the event ID from the URL path
  const path = window.location.pathname;
  const id = path.split('/').pop();
  const [event, setEvent] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [notification, setNotification] = useState(null);
  const [rsvpStatus, setRsvpStatus] = useState(null);
  const [rsvpCounts, setRsvpCounts] = useState({
    going: 0,
    maybe: 0,
    not_going: 0,
  });
  const [isRsvpLoading, setIsRsvpLoading] = useState(false);

  // Get the current user ID from localStorage
  const currentUserId = parseInt(localStorage.getItem('userId'), 10);
  const isLoggedIn = !!localStorage.getItem('token');

  useEffect(() => {
    if (id) {
      fetchEventDetails(id);
      if (isLoggedIn) {
        fetchRsvpStatus(id);
        fetchRsvpCounts(id);
      }
    } else {
      setNotification({
        type: 'error',
        message: 'Invalid event ID',
      });
      setIsLoading(false);
    }
  }, [id, isLoggedIn]);

  async function fetchEventDetails(eventId) {
    setIsLoading(true);
    try {
      const response = await fetch(
        `http://localhost:9000/api/events/${eventId}`
      );

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('Event not found');
        }
        const data = await response.json();
        throw new Error(data.message || 'Failed to fetch event details');
      }

      const data = await response.json();
      setEvent(data);
    } catch (error) {
      setNotification({
        type: 'error',
        message:
          error.message || 'An error occurred while fetching event details',
      });
    } finally {
      setIsLoading(false);
    }
  }

  async function fetchRsvpStatus(eventId) {
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      const response = await fetch(
        `http://localhost:9000/api/events/${eventId}/rsvp`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Failed to fetch RSVP status');
      }

      const data = await response.json();
      setRsvpStatus(data ? data.status : null);
    } catch (error) {
      console.error('Error fetching RSVP status:', error);
    }
  }

  async function fetchRsvpCounts(eventId) {
    try {
      const response = await fetch(
        `http://localhost:9000/api/events/${eventId}/rsvp/count`
      );

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Failed to fetch RSVP counts');
      }

      const data = await response.json();
      setRsvpCounts(data);
    } catch (error) {
      console.error('Error fetching RSVP counts:', error);
    }
  }

  async function handleRsvp(status) {
    if (!isLoggedIn) {
      // Use window.location instead of navigate
      window.location.href = '/signin';
      return;
    }

    setIsRsvpLoading(true);
    try {
      const token = localStorage.getItem('token');

      // If clicking the same status again, remove the RSVP
      if (rsvpStatus === status) {
        const response = await fetch(
          `http://localhost:9000/api/events/${id}/rsvp`,
          {
            method: 'DELETE',
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (!response.ok) {
          const data = await response.json();
          throw new Error(data.message || 'Failed to remove RSVP');
        }

        setRsvpStatus(null);
        setNotification({
          type: 'success',
          message: 'RSVP removed successfully',
        });
      } else {
        // Otherwise, update the RSVP
        const response = await fetch(
          `http://localhost:9000/api/events/${id}/rsvp`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              status: status,
            }),
          }
        );

        if (!response.ok) {
          const data = await response.json();
          throw new Error(data.message || 'Failed to update RSVP');
        }

        setRsvpStatus(status);
        setNotification({
          type: 'success',
          message: 'RSVP updated successfully',
        });
      }

      // Refresh RSVP counts
      fetchRsvpCounts(id);
    } catch (error) {
      setNotification({
        type: 'error',
        message: error.message || 'An error occurred while updating RSVP',
      });
    } finally {
      setIsRsvpLoading(false);
    }
  }

  async function handleDeleteEvent() {
    setIsDeleting(true);
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('You must be logged in to delete an event');
      }

      const response = await fetch(
        `http://localhost:9000/api/events/${event.id}`,
        {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Failed to delete event');
      }

      // Show success notification
      setNotification({
        type: 'success',
        message: 'Event deleted successfully!',
      });

      // Set event to null to show the "Event not found" view
      setEvent(null);
    } catch (error) {
      setNotification({
        type: 'error',
        message: error.message || 'An error occurred while deleting the event',
      });
    } finally {
      setIsDeleting(false);
    }
  }

  // Format date for display
  function formatDate(dateString) {
    const options = {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    };
    return new Date(dateString).toLocaleDateString(undefined, options);
  }

  // Calculate days remaining until event
  function getDaysRemaining(dateString) {
    const eventDate = new Date(dateString);
    const today = new Date();

    // Reset time part for accurate day calculation
    today.setHours(0, 0, 0, 0);
    const eventDateNoTime = new Date(eventDate);
    eventDateNoTime.setHours(0, 0, 0, 0);

    const diffTime = eventDateNoTime - today;
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Tomorrow';
    return `${diffDays} days away`;
  }

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto p-4">
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
        </div>
      </div>
    );
  }

  if (!event && !isLoading) {
    return (
      <div className="max-w-4xl mx-auto p-4">
        {notification && (
          <Notification
            type={notification.type}
            message={notification.message}
            onClose={() => setNotification(null)}
          />
        )}
        <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6 text-center">
          <p className="text-gray-600 dark:text-gray-400">Event not found.</p>
          <a
            href="/upcoming-events"
            className="mt-4 inline-block px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700"
          >
            View All Events
          </a>
        </div>
      </div>
    );
  }

  // Check if the current user is the event creator
  const isEventCreator = currentUserId === event.user_id;

  if (isEditing) {
    return (
      <div className="max-w-4xl mx-auto p-4">
        <EditEventForm
          eventId={event.id}
          onCancel={() => setIsEditing(false)}
          onSuccess={() => {
            setIsEditing(false);
            fetchEventDetails(event.id);
          }}
        />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto p-4">
      {notification && (
        <Notification
          type={notification.type}
          message={notification.message}
          onClose={() => setNotification(null)}
        />
      )}

      <div className="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
        <div className="p-6">
          <div className="flex justify-between items-start mb-4">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
              {event.title}
            </h1>
            <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-200">
              {getDaysRemaining(event.date)}
            </span>
          </div>

          <div className="flex items-center text-sm text-gray-500 dark:text-gray-400 mb-4">
            <svg
              className="h-5 w-5 mr-2"
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
            <span>{formatDate(event.date)}</span>
          </div>

          <div className="flex items-center text-sm text-gray-500 dark:text-gray-400 mb-6">
            <svg
              className="h-5 w-5 mr-2"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
            <span>{event.location}</span>
          </div>

          {/* RSVP Section */}
          {!isEventCreator && (
            <div className="mb-8 border-t border-b border-gray-200 dark:border-gray-700 py-4">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">
                Will you attend?
              </h3>
              <div className="flex flex-wrap gap-3">
                <button
                  onClick={() => handleRsvp('going')}
                  disabled={isRsvpLoading}
                  className={`px-4 py-2 rounded-md text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 ${
                    rsvpStatus === 'going'
                      ? 'bg-green-600 text-white hover:bg-green-700'
                      : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-600'
                  }`}
                >
                  Going ({rsvpCounts.going})
                </button>
                <button
                  onClick={() => handleRsvp('maybe')}
                  disabled={isRsvpLoading}
                  className={`px-4 py-2 rounded-md text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 ${
                    rsvpStatus === 'maybe'
                      ? 'bg-yellow-500 text-white hover:bg-yellow-600'
                      : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-600'
                  }`}
                >
                  Maybe ({rsvpCounts.maybe})
                </button>
                <button
                  onClick={() => handleRsvp('not_going')}
                  disabled={isRsvpLoading}
                  className={`px-4 py-2 rounded-md text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 ${
                    rsvpStatus === 'not_going'
                      ? 'bg-red-600 text-white hover:bg-red-700'
                      : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-600'
                  }`}
                >
                  Not Going ({rsvpCounts.not_going})
                </button>
              </div>
              {!isLoggedIn && (
                <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                  <a
                    href="/signin"
                    className="text-primary-600 hover:text-primary-500 dark:text-primary-400"
                  >
                    Sign In
                  </a>{' '}
                  to RSVP to this event.
                </p>
              )}
            </div>
          )}

          <div className="mb-8">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
              About this event
            </h2>
            <p className="text-gray-700 dark:text-gray-300 whitespace-pre-line">
              {event.description || 'No description provided.'}
            </p>
          </div>

          <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
              Organizer
            </h2>
            <p className="text-gray-700 dark:text-gray-300">
              {event.organizer_first_name} {event.organizer_last_name}
            </p>
          </div>

          {isEventCreator && (
            <div className="border-t border-gray-200 dark:border-gray-700 mt-6 pt-6">
              <div className="flex flex-col sm:flex-row sm:justify-end gap-3">
                <button
                  onClick={() => setIsEditing(true)}
                  className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-600"
                >
                  Edit Event
                </button>
                <button
                  onClick={() => {
                    if (
                      window.confirm(
                        'Are you sure you want to delete this event? This action cannot be undone.'
                      )
                    ) {
                      handleDeleteEvent();
                    }
                  }}
                  disabled={isDeleting}
                  className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
                >
                  {isDeleting ? 'Deleting...' : 'Delete Event'}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
