import { useState } from 'react';
import { useActionState } from 'react';

export default function SignupForm() {
  const [isLoading, setIsLoading] = useState(false);

  async function handleSignup(prevState, formData) {
    setIsLoading(true);

    const email = formData.get('email');
    const password = formData.get('password');
    const confirmedPassword = formData.get('confirmed_password');
    const firstName = formData.get('firstName');
    const lastName = formData.get('lastName');
    
    try {
      const response = await fetch('http://localhost:9000/api/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          password,
          confirmed_password: confirmedPassword,
          first_name: firstName,
          last_name: lastName,
        }),
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        return errorData.message || 'Signup failed';
      }
      
      return { success: true, message: 'Account created successfully!' };
    } catch (error) {
      return error.message || 'An error occurred during signup';
    } finally {
      setIsLoading(false);
    }
  }
  
  const [result, signupAction] = useActionState(handleSignup, null);
  
  
}