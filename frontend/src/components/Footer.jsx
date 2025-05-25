export default function Footer() {
  return (
    <footer className="bg-white dark:bg-gray-800 shadow-md py-4 mt-auto">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="mb-4 md:mb-0">
            <p className="text-gray-600 dark:text-gray-300 text-sm">
              &copy; {new Date().getFullYear()} Evently.
            </p>
          </div>
          <div className="flex space-x-4">
            <a
              href="https://github.com/johneliud"
              className="text-gray-600 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400 text-sm"
            >
              Made by John
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
