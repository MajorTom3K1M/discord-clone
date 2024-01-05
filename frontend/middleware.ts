import { NextResponse, NextRequest } from "next/server";
// import  

const publicRoutes = ['/','/contact'];

export function middleware(request: NextRequest) {
    const url = request.nextUrl.clone();
    // const token = 
    
}