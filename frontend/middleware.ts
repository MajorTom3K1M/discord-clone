import { NextResponse, NextRequest } from "next/server";
// import  

const publicRoutes = ['/','/contact'];

export function middleware(request: NextRequest) {
    return NextResponse.next()
}