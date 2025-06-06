import { useState, useEffect } from 'react';
import Notification from './Notification';
import config from '../config';

export default function EditEventForm({ eventId, onCancel, onSuccess }) {
  const [event, setEvent] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [notification, setNotification] = useState(null);

  useEffect(() => {
    const fetchEvent = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(`${config.apiBaseUrl}/api/events/${eventId}`);

        if (!response.ok) {
          throw new Error('Failed to fetch event details');
        }

        const data = await response.json();
        setEvent(data);
      } catch (error) {
        setNotification({
          type: 'error',
          message: error.message || 'An error occurred while fetching event details'
        });
      } finally {
        setIsLoading(false);
      }
    };
    fetchEvent();
  }, [eventId]);

  async function handleSubmit(e) {
    e.preventDefault();
    setIsSaving(true);
    setNotification(null);

    const formData = new FormData(e.target);
    const title = formData.get('title');
    const description = formData.get('description');
    const date = formData.get('date');
    const time = formData.get('time');
    const location = formData.get('location');

    // Combine date and time
    const dateTime = new Date(`${date}T${time}`);

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('You must be logged in to update an event');
      }

      const response = await fetch(`${config.apiBaseUrl}/api/events/${eventId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          title,
          description,
          date: dateTime.toISOString(),
          location,
        }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Failed to update event');
      }

      // Show success notification
      setNotification({
        type: 'success',
        message: 'Event updated successfully!',
      });

      // Call the success callback
      if (onSuccess) {
        setTimeout(() => {
          onSuccess();
        }, 1500);
      }
    } catch (error) {
      setNotification({
        type: 'error',
        message: error.message || 'An error occurred while updating the event',
      });
    } finally {
      setIsSaving(false);
    }
  }

  // Format date for input field
  function formatDateForInput(dateString) {
    const date = new Date(dateString);
    return date.toISOString().split('T')[0];
  }

  // Format time for input field
  function formatTimeForInput(dateString) {
    const date = new Date(dateString);
    return date.toISOString().split('T')[1].substring(0, 5);
  }

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-40">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
      </div>
    );
  }

  if (!event) {
    return (
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6 text-center">
        <p className="text-gray-600 dark:text-gray-400">Event not found.</p>
        <button
          onClick={onCancel}
          className="mt-4 inline-block px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700"
        >
          Go Back
        </button>
      </div>
    );
  }

  // Get today's date in YYYY-MM-DD format for min attribute
  const today = new Date().toISOString().split('T')[0];

  return (
    <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Edit Event</h2>
      
      {notification && (
        <Notification
          type={notification.type}
          message={notification.message}
          onClose={() => setNotification(null)}
        />
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="title" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Event Title *
          </label>
          <input
            type="text"
            name="title"
            id="title"
            required
            defaultValue={event.title}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
          />
        </div>

        <div>
          <label htmlFor="description" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Description
          </label>
          <textarea
            name="description"
            id="description"
            rows="4"
            defaultValue={event.description}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
          ></textarea>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label htmlFor="time" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Time *
            </label>
            <input
              type="time"
              name="time"
              id="time"
              required
              defaultValue={formatTimeForInput(event.date)}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>

          <div>
            <label htmlFor="date" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Date *
            </label>
            <input
              type="date"
              name="date"
              id="date"
              required
              min={today}
              defaultValue={formatDateForInput(event.date)}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
            />
          </div>
        </div>

        <div>
          <label htmlFor="location" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Location *
          </label>
          <input
            type="text"
            name="location"
            id="location"
            required
            defaultValue={event.location}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
          />
        </div>

        <div className="flex justify-end space-x-3">
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-600"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isSaving}
            className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
          >
            {isSaving ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </form>
    </div>
  );
}
