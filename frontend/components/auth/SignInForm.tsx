"use client"
import { ChangeEvent, FormEvent, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/components/providers/AuthProvider';

const SignIn = () => {
    const { signin } = useAuth();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const handleEmailChange = (event: ChangeEvent<HTMLInputElement>) => {
        setEmail(event.target.value);
    };

    const handlePasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
        setPassword(event.target.value);
    };

    const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        console.log('Email:', email, 'Password:', password);
    };

    return (
        <div className="flex justify-center items-center h-screen">
            <div className="bg-white dark:bg-[#2B2D31] bg-[#F2F3F5] p-8 rounded-lg shadow-lg max-w-sm w-full">
                <h2 className="text-2xl font-semibold text-gray-700 dark:text-gray-200 text-center">Welcome back!</h2>
                <p className="text-gray-500 dark:text-gray-400 text-center mb-6">We&apos;re so excited to see you again!</p>
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label htmlFor="email" className="block text-gray-700 dark:text-gray-200 text-sm font-semibold mb-2">
                            Email or Phone Number
                        </label>
                        <input
                            type="text"
                            id="email"
                            className="shadow appearance-none border rounded w-full py-2 px-3 text-white leading-tight focus:outline-none focus:shadow-outline"
                            onChange={handleEmailChange}
                            required
                        />
                    </div>
                    <div className="mb-6">
                        <label htmlFor="password" className="block text-gray-700 dark:text-gray-200 text-sm font-semibold mb-2">
                            Password
                        </label>
                        <input
                            type="password"
                            id="password"
                            className="shadow appearance-none border rounded w-full py-2 px-3 text-white mb-3 leading-tight focus:outline-none focus:shadow-outline"
                            onChange={handlePasswordChange}
                            required
                        />
                        <button 
                            className="text-sm text-blue-500 hover:text-blue-800 hover:underline"
                        >
                            Forgot your password?
                        </button>
                    </div>
                    <button
                        type="submit"
                        className="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                        onClick={() => signin(email, password)}
                    >
                        Log In
                    </button>
                </form>
                <p className="text-gray-500 dark:text-gray-400 text-xs text-center mt-4">
                    Need an account? {' '}
                    <Link href="/sign-up">
                        <button className="text-blue-500 hover:text-blue-800 hover:underline">
                            Register
                        </button>
                    </Link>
                </p>
            </div>
        </div>
    );
};

export default SignIn;
