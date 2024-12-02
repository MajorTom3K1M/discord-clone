"use client"
import React, { ChangeEvent, FormEvent, useState } from 'react';
import Link from 'next/link';

import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { useForm } from 'react-hook-form';
import { useAuth } from '@/components/providers/AuthProvider';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage
} from '@/components/ui/Form';
import { Input } from '../ui/Input';
import { FileUpload } from '../FileUpload';
import { Button } from '../ui/Button';

const formSchema = z.object({
    email: z.string().min(1, 'Email is required').email('Invalid email format'),
    name: z.string().min(1, 'Username is required'),
    password: z.string().min(1, 'Password is required'),
    imageUrl: z.string().min(1, 'Image URL is required'),
});

const SignUp = () => {
    const { signup } = useAuth();

    const onSubmit = async (values: z.infer<typeof formSchema>) => {
        try {
            await signup(values.name, values.imageUrl, values.email, values.password);

            form.reset();
            // router.refresh();
        } catch (error) {
            console.log(error)
        }
    }

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
    });

    return (
        <div className="flex justify-center items-center h-screen flex-1">
            <div className="bg-white dark:bg-[#2B2D31] bg-[#F2F3F5] p-8 rounded-lg shadow-lg max-w-sm w-full flex-1">
                <h2 className="text-2xl font-semibold text-gray-700 dark:text-gray-200 text-center">Create an account</h2>
                <p className="text-gray-500 dark:text-gray-400 text-center mb-6">We&apos;re so excited to see you again!</p>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-8'>
                        <div className='space-y-4'>
                            <FormField
                                control={form.control}
                                name="imageUrl"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel
                                            className='uppercase text-xs font-bold text-gray-700 dark:text-gray-200'
                                        >
                                            Profile Image
                                        </FormLabel>
                                        <div className='flex items-center justify-center text-center'>
                                            <FormControl>
                                                <FileUpload
                                                    className="bg-gray-100 p-4 rounded-lg"
                                                    endpoint='profileImage'
                                                    value={field.value}
                                                    onChange={field.onChange}
                                                />
                                            </FormControl>
                                        </div>
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='email'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel
                                            className='uppercase text-xs font-bold text-gray-700 dark:text-gray-200 px-0'
                                        >
                                            Email
                                        </FormLabel>
                                        <FormControl>
                                            <Input
                                                disabled={false}
                                                className='border-0 focus-visible:ring-0 text-white focus-visible:ring-offset-0'
                                                placeholder='Enter email'
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='name'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel
                                            className='uppercase text-xs font-bold text-gray-700 dark:text-gray-200 px-0'
                                        >
                                            Username
                                        </FormLabel>
                                        <FormControl>
                                            <Input
                                                disabled={false}
                                                className='border-0 focus-visible:ring-0 text-white focus-visible:ring-offset-0'
                                                placeholder='Enter username'
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name='password'
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel
                                            className='uppercase text-xs font-bold text-gray-700 dark:text-gray-200 px-0'
                                        >
                                            Password
                                        </FormLabel>
                                        <FormControl>
                                            <Input
                                                disabled={false}
                                                className='border-0 focus-visible:ring-0 text-white focus-visible:ring-offset-0'
                                                placeholder='Enter password'
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                        <Button className="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
                            Create
                        </Button>
                    </form>
                </Form>
                <p className="text-gray-500 dark:text-gray-400 text-xs text-center mt-4">
                    <Link href="/sign-in">
                        <button className="text-blue-500 hover:text-blue-800 hover:underline">
                            Already have an account?
                        </button>
                    </Link>
                </p>
            </div>
        </div>
    );
};

export default SignUp;
