import { useState, useEffect } from 'react';
import Notification from './Notification';
import GoogleAuthButton from './GoogleAuthButton';

export default function SigninForm() {
  const [isLoading, setIsLoading] = useState(false);
  const [notification, setNotification] = useState(null);
  const [redirecting, setRedirecting] = useState(false);

  useEffect(() => {
    if (notification && notification.type === 'success') {
      const timer = setTimeout(() => {
        setRedirecting(true);
        window.location.href = '/';
      }, 2000);

      return () => clearTimeout(timer);
    }
  }, [notification]);

  async function handleSubmit(e) {
    e.preventDefault();
    setIsLoading(true);
    setNotification(null);

    const formData = new FormData(e.target);
    const email = formData.get('email');
    const password = formData.get('password');

    try {
      const response = await fetch('http://localhost:9000/api/signin', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Invalid email or password');
      }

      // Store token in localStorage
      localStorage.setItem('token', data.token);
      localStorage.setItem('userId', data.user_id);

      // Show success notification
      setNotification({
        type: 'success',
        message: 'Sign in successful! Redirecting to dashboard...',
      });
    } catch (error) {
      setNotification({
        type: 'error',
        message: error.message || 'An error occurred during sign in',
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
          showProgress={notification.type === 'success'}
        />
      )}

      <main className="min-h-[calc(100vh-80px)] flex items-center justify-center px-4 sm:px-6 lg:px-8">
        <div className="w-full max-w-2xl md:w-1/2 space-y-8 bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md">
          <div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
              Welcome Back
            </h2>
          </div>

          <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
            <div className="rounded-md shadow-sm -space-y-px">
              <div className="mb-4">
                <label
                  htmlFor="email"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  Email address
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-400 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-700 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
                  placeholder="Email address"
                />
              </div>
              <div className="mt-4">
                <label
                  htmlFor="password"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  Password
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  required
                  className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-400 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-700 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
                  placeholder="Enter your password"
                />
              </div>
            </div>

            <div>
              <button
                type="submit"
                disabled={isLoading || redirecting}
                className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-md font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isLoading
                  ? 'Signing in...'
                  : redirecting
                  ? 'Redirecting...'
                  : 'Sign In'}
              </button>
            </div>

            <div className="mt-6 relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-300 dark:border-gray-600"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400">
                  Or continue with
                </span>
              </div>
            </div>

            {/* Google Auth Button */}
            <div className="mt-6">
              <GoogleAuthButton text="Sign in with Google" />
            </div>

            <div className="text-center">
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Don't have an account?{' '}
                <a
                  href="/signup"
                  className="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400"
                >
                  Sign Up
                </a>
              </p>
            </div>
          </form>
        </div>
      </main>
    </>
  );
}
