// Configuration for different environments
const config = {
  // API base URL - automatically detects environment
  apiBaseUrl: import.meta.env.PROD
    ? 'https://evently-backend-gs5n.onrender.com'
    : 'http://localhost:9000',

  googleAuthEnabled: true,
  calendarIntegrationEnabled: true,
};

export default config;
