// middleware.js
import { NextResponse } from 'next/server';

export function middleware(request) {
  const { nextUrl } = request;
  const { pathname } = request.nextUrl;

  // 如果访问dashboard且未登录，则重定向到登录页面
  if (pathname.startsWith('/dashboard')) {
    const token = request.cookies.get('token');
    console.log(token);
    if (!token) {
      // Assuming you want to redirect to the home page
      const url = nextUrl.clone();
      url.pathname = '/';
      return NextResponse.redirect(url);
    }
  }

  return NextResponse.next();
}
