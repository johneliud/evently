import { useState, useEffect } from 'react';
import './App.css';
import Header from './components/Header';
import Sidebar from './components/Sidebar';
import SignupForm from './components/SignupForm';
import SigninForm from './components/SigninForm';
import EventForm from './components/EventForm';
import EventList from './components/EventList';
import UpcomingEvents from './components/UpcomingEvents';
import EventDetails from './components/EventDetails';
import EventSearch from './components/EventSearch';
import CalendarConnected from './components/CalendarConnected';

function App() {
  const [currentPath, setCurrentPath] = useState(window.location.pathname);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  useEffect(() => {
    const handleLocationChange = () => {
      setCurrentPath(window.location.pathname);
      // Close sidebar on navigation (mobile)
      setIsSidebarOpen(false);
    };

    // Listen for popstate events (back/forward navigation)
    window.addEventListener('popstate', handleLocationChange);

    return () => {
      window.removeEventListener('popstate', handleLocationChange);
    };
  }, []);

  // Toggle sidebar
  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

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
      case '/calendar-connected':
        return <CalendarConnected />;
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
      <Header toggleSidebar={toggleSidebar} />
      <div className="flex min-h-[calc(100vh-64px)]">
        <Sidebar isOpen={isSidebarOpen} onClose={() => setIsSidebarOpen(false)} />
        <main className="flex-1 py-5 px-4 lg:px-10 overflow-y-auto">
          {renderContent()}
        </main>
      </div>
    </div>
  );
}

export default App;
