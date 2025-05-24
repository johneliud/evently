import { useState, useEffect } from 'react';
import './App.css';
import Header from './components/Header';
import SignupForm from './components/SignupForm';
import SigninForm from './components/SigninForm';
import EventForm from './components/EventForm';
import EventList from './components/EventList';
import UpcomingEvents from './components/UpcomingEvents';
import EventDetails from './components/EventDetails';
import EventSearch from './components/EventSearch';

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
    // Check if path matches event details pattern (/event/{id})
    if (currentPath.match(/^\/event\/\d+$/)) {
      return <EventDetails />;
    }

    switch (currentPath) {
      case '/signin':
        return <SigninForm />;
      case '/signup':
        return <SignupForm />;
      case '/create-event':
        return requireAuth(<EventForm />);
      case '/my-events':
        return requireAuth(<EventList />);
      case '/upcoming-events':
        return <UpcomingEvents />;
      case '/search':
        return <EventSearch />;
      default: {
        // Check if user is authenticated
        const token = localStorage.getItem('token');
        if (!token) {
          return <UpcomingEvents />;
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
      <div className="w-full py-5 px-10">
        {renderContent()}
      </div>
    </div>
  );
}

export default App;
