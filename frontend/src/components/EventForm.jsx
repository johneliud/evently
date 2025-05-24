import { useState } from 'react';
import Notification from './Notification';

export default function EventForm() {
  const [isLoading, setIsLoading] = useState(false);
  const [notification, setNotification] = useState(null);
  
  // Get today's date in YYYY-MM-DD format for min attribute
  const today = new Date().toISOString().split('T')[0];

  async function handleSubmit(e) {
    e.preventDefault();
    setIsLoading(true);
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
        throw new Error('You must be logged in to create an event');
      }

      const response = await fetch('http://localhost:9000/api/events', {
        method: 'POST',
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

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Failed to create event');
      }

      // Show success notification
      setNotification({
        type: 'success',
        message: 'Event created successfully!',
      });

      // Reset form
      e.target.reset();
    } catch (error) {
      setNotification({
        type: 'error',
        message: error.message || 'An error occurred while creating the event',
      });
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <>
      {notification && (
        <Notification
          type={notification.type}
          message={notification.message}
          onClose={() => setNotification(null)}
        />
      )}

      <div className="flex items-center justify-center min-h-[calc(100vh-80px)]">
        <div className="w-full max-w-2xl mx-auto p-6 bg-white dark:bg-gray-800 rounded-lg shadow-md">
          <h2 className="text-2xl font-bold mb-6 text-gray-900 dark:text-white">
            Create New Event
          </h2>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label
                htmlFor="title"
                className="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Event Title *
              </label>
              <input
                type="text"
                id="title"
                name="title"
                required
                className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                placeholder="Enter event title"
              />
            </div>

            <div>
              <label
                htmlFor="description"
                className="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Description
              </label>
              <textarea
                id="description"
                name="description"
                rows="4"
                className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                placeholder="Describe your event"
              ></textarea>
            </div>

            <div>
              <label
                htmlFor="location"
                className="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Location *
              </label>
              <input
                type="text"
                id="location"
                name="location"
                required
                className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                placeholder="Enter event location"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label
                  htmlFor="time"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >
                  Time *
                </label>
                <input
                  type="time"
                  id="time"
                  name="time"
                  required
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                />
              </div>
              <div>
                <label
                  htmlFor="date"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >
                  Date *
                </label>
                <input
                  type="date"
                  id="date"
                  name="date"
                  required
                  min={today}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
                />
              </div>
            </div>

            <div>
              <button
                type="submit"
                disabled={isLoading}
                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isLoading ? 'Creating...' : 'Create Event'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  );
}
