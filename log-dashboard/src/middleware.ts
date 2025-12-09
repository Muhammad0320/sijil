import { jwtVerify } from "jose";
import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

export async function middleware(request: NextRequest) {
  const session = request.cookies.get("session")?.value;

  const isAuthPage =
    request.nextUrl.pathname.startsWith("/login") ||
    request.nextUrl.pathname.startsWith("/register");

  if (!session) {
    if (!isAuthPage) {
      return NextResponse.redirect(new URL("/login", request.url));
    }
    return NextResponse.next();
  }

  try {
    const secret = new TextEncoder().encode(
      process.env.JWT_KEY || "my_super_long_jwt_secret_key_for_log_engine"
    );

    await jwtVerify(session, secret);

    if (isAuthPage) {
      return NextResponse.redirect(new URL("/", request.url));
    }

    return NextResponse.next();
  } catch (error) {
    const response = NextResponse.redirect(new URL("/login", request.url));
    response.cookies.delete("session");
    return response;
  }

  if (session && isAuthPage) {
    return NextResponse.redirect(new URL("/", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};
