"use client"
import React, { useState } from 'react';

const SignUp = () => {
    const [email, setEmail] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    return (
        <div className="flex justify-center items-center h-screen">
            <div className="bg-gray-900 p-10 rounded-lg shadow-lg max-w-sm w-full mx-4">
                <h2 className="text-2xl font-bold text-white text-center mb-8">Create an account</h2>
                <form>
                    <div className="mb-4">
                        <label htmlFor="email" className="block text-white text-sm mb-2">Email *</label>
                        <input
                            type="email"
                            id="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="bg-gray-800 text-white w-full p-3 rounded focus:bg-gray-700 focus:outline-none"
                            required
                        />
                    </div>
                    <div className="mb-4">
                        <label htmlFor="username" className="block text-white text-sm mb-2">Username *</label>
                        <input
                            type="text"
                            id="username"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            className="bg-gray-800 text-white w-full p-3 rounded focus:bg-gray-700 focus:outline-none"
                            required
                        />
                    </div>
                    <div className="mb-6">
                        <label htmlFor="password" className="block text-white text-sm mb-2">Password *</label>
                        <input
                            type="password"
                            id="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className="bg-gray-800 text-white w-full p-3 rounded focus:bg-gray-700 focus:outline-none"
                            required
                        />
                    </div>
                    <button
                        type="submit"
                        className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                    >
                        Continue
                    </button>
                </form>
                <p className="text-gray-400 text-xs text-center mt-4">
                    Need an account? <a href="#" className="text-blue-500 hover:text-blue-400">Register</a>
                </p>
            </div>
        </div>
    );
};

export default SignUp;
