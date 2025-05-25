import { useEffect, useState } from 'react';
import Notification from './Notification';

export default function AuthCallback() {
  const [notification, setNotification] = useState(null);
  const [redirecting, setRedirecting] = useState(false);

  useEffect(() => {
    // Get token and user_id from URL parameters
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');
    const userId = params.get('user_id');

    if (token && userId) {
      // Store token and user ID in localStorage
      localStorage.setItem('token', token);
      localStorage.setItem('userId', userId);
      
      // Show success notification
      setNotification({
        type: 'success',
        message: 'Authentication successful! Redirecting to dashboard...'
      });
      
      // Redirect to home page after a delay
      const timer = setTimeout(() => {
        setRedirecting(true);
        window.location.href = '/';
      }, 2000);
      
      return () => clearTimeout(timer);
    } else {
      // Show error notification
      setNotification({
        type: 'error',
        message: 'Authentication failed. Please try again.'
      });
    }
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="max-w-md w-full space-y-8 p-10 bg-white dark:bg-gray-800 rounded-xl shadow-md">
        <div className="text-center">
          <h2 className="mt-6 text-3xl font-extrabold text-gray-900 dark:text-white">
            {notification?.type === 'success' ? 'Authentication Successful' : 'Authentication Failed'}
          </h2>
          
          <div className="mt-8">
            {notification?.type === 'success' ? (
              <div className="flex justify-center">
                <div className="rounded-full bg-green-100 p-3">
                  <svg className="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
              </div>
            ) : (
              <div className="flex justify-center">
                <div className="rounded-full bg-red-100 p-3">
                  <svg className="h-8 w-8 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
              </div>
            )}
          </div>
          
          <p className="mt-4 text-md text-gray-600 dark:text-gray-400">
            {notification?.message}
          </p>
          
          {notification?.type === 'error' && (
            <div className="mt-6">
              <a 
                href="/signin" 
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              >
                Return to Sign In
              </a>
            </div>
          )}
          
          {notification?.type === 'success' && redirecting && (
            <div className="mt-6">
              <div className="w-full bg-gray-200 rounded-full h-2.5 dark:bg-gray-700">
                <div className="bg-primary-600 h-2.5 rounded-full animate-pulse" style={{ width: '100%' }}></div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
