import { NextResponse, NextRequest } from "next/server";
// import  

const publicRoutes = ['/','/contact'];

export function middleware(request: NextRequest) {
    const url = request.nextUrl.clone();
    const cookies = request.cookies;
    // console.log(cookies)
    // if (cookies) {
    //     console.log({ cookies })
    // }
    return NextResponse.next()
}