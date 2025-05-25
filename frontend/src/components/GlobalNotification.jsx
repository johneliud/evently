import { useState, useEffect } from 'react';
import Notification from './Notification';

export default function GlobalNotification() {
  const [notification, setNotification] = useState(null);

  useEffect(() => {
    // Listen for custom events to show notifications
    const handleNotification = (event) => {
      setNotification(event.detail);
      
      // Auto-dismiss success notifications after 3 seconds
      if (event.detail.type === 'success') {
        setTimeout(() => {
          setNotification(null);
        }, 3000);
      }
    };

    window.addEventListener('showNotification', handleNotification);
    
    // Check for auth callback parameters in URL
    const checkAuthCallback = () => {
      const urlParams = new URLSearchParams(window.location.search);
      const authSuccess = urlParams.get('auth_success');
      const authError = urlParams.get('auth_error');
      
      if (authSuccess) {
        setNotification({
          type: 'success',
          message: decodeURIComponent(authSuccess)
        });
        
        // Remove the query parameter
        const newUrl = window.location.pathname;
        window.history.replaceState({}, document.title, newUrl);
        
        // Auto-dismiss after 3 seconds
        setTimeout(() => {
          setNotification(null);
        }, 3000);
      } else if (authError) {
        setNotification({
          type: 'error',
          message: decodeURIComponent(authError)
        });
        
        // Remove the query parameter
        const newUrl = window.location.pathname;
        window.history.replaceState({}, document.title, newUrl);
      }
    };
    
    checkAuthCallback();

    return () => {
      window.removeEventListener('showNotification', handleNotification);
    };
  }, []);

  // Helper function to show notifications from anywhere in the app
  window.showNotification = (type, message) => {
    const event = new CustomEvent('showNotification', {
      detail: { type, message }
    });
    window.dispatchEvent(event);
  };

  if (!notification) return null;

  return (
    <Notification
      type={notification.type}
      message={notification.message}
      onClose={() => setNotification(null)}
    />
  );
}