import { useEffect } from 'react';

export default function CalendarConnected() {
  useEffect(() => {
    // Close the window after 5 seconds if it's a popup
    const timer = setTimeout(() => {
      if (window.opener) {
        window.close();
      } else {
        // If not a popup, redirect to home
        window.location.href = '/';
      }
    }, 5000);

    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-4 text-center">
      <div className="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-8 max-w-md">
        <div className="flex justify-center mb-4">
          <div className="rounded-full bg-green-100 p-3">
            <svg className="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
        </div>
        
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
          Google Calendar Connected!
        </h2>
        
        <p className="text-gray-600 dark:text-gray-400 mb-6">
          Your Google Calendar has been successfully connected to Evently. You can now add events to your calendar.
        </p>
        
        <p className="text-sm text-gray-500 dark:text-gray-500">
          This window will close automatically in a few seconds...
        </p>
      </div>
    </div>
  );
}