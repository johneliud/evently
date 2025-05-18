import { useState, useEffect } from 'react';
import './App.css';
import Header from './components/Header';
import SignupForm from './components/SignupForm';
import SigninForm from './components/SigninForm';
import EventForm from './components/EventForm';
import EventList from './components/EventList';

function App() {
  const [currentPath, setCurrentPath] = useState(window.location.pathname);

  useEffect(() => {
    const handleLocationChange = () => {
      setCurrentPath(window.location.pathname);
    };

    // Listen for popstate events (back/forward navigation)
    window.addEventListener('popstate', handleLocationChange);

    return () => {
      window.removeEventListener('popstate', handleLocationChange);
    };
  }, []);

  // Simple routing
  const renderContent = () => {
    switch (currentPath) {
      case '/signin':
        return <SigninForm />;
      case '/signup':
        return <SignupForm />;
      case '/create-event':
        return requireAuth(<EventForm />);
      case '/my-events':
        return requireAuth(<EventList />);
      default: {
        // Check if user is authenticated
        const token = localStorage.getItem('token');
        if (!token) {
          window.location.href = '/signin';
          return null;
        }
        return <EventList />;
      }
    }
  };

  // Helper function to require authentication
  const requireAuth = (component) => {
    const token = localStorage.getItem('token');
    if (!token) {
      window.location.href = '/signin';
      return null;
    }
    return component;
  };

  return (
    <div className="min-h-screen w-full">
      <Header />
      <div className="container mx-auto py-8 px-4">
        {renderContent()}
      </div>
    </div>
  );
}

export default App;
