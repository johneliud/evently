import { useState, useEffect } from 'react';
import './App.css';
import Header from './components/Header';
import SignupForm from './components/SignupForm';
import SigninForm from './components/SigninForm';

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
      default: {
        // Check if user is authenticated
        const token = localStorage.getItem('token');
        if (!token) {
          window.location.href = '/signin';
          return null;
        }
        return <div className="p-8 text-center">Welcome to Evently! You are logged in.</div>;
      }
    }
  };

  return (
    <div className="min-h-screen w-full">
      <Header />
      {renderContent()}
    </div>
  );
}

export default App;
