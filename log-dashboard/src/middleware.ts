import { jwtVerify } from "jose";
import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

export async function middleware(request: NextRequest) {
  const session = request.cookies.get("session")?.value;
  const path = request.nextUrl.pathname;

  const isAuthPage = path.startsWith("/login") || path.startsWith("/register");
  const isDashboard = path.startsWith("/dashboard");
  const isRoot = path === "/";

  console.log(session);

  if (!session) {
    if (!isDashboard) {
      return NextResponse.redirect(new URL("/login", request.url));
    }
    return NextResponse.next();
  }

  try {
    const secret = new TextEncoder().encode(
      process.env.JWT_SECRET || "my_super_long_jwt_secret_key_for_log_engine"
    );

    await jwtVerify(session, secret, {
      clockTolerance: 30,
    });

    if (isAuthPage || isRoot) {
      return NextResponse.redirect(new URL("/dashboard", request.url));
    }

    return NextResponse.next();
  } catch (error) {
    console.log(error, "middle------------");
    console.error("🔒 Auth Failed:", error);
    return NextResponse.redirect(new URL("/login", request.url));
  }
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};
